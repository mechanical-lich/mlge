package audio

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Bus is a mix group with its own volume knob. Every voice plays on exactly one
// bus, and the Master bus multiplies all of them — so a settings screen can
// expose Master/Music/SFX/UI/Global sliders that all feed the same engine.
type Bus int

const (
	BusMaster Bus = iota // global gain applied to every voice; not a play target
	BusMusic
	BusSFX
	BusUI
	BusGlobal
	numBuses
)

func (b Bus) String() string {
	switch b {
	case BusMaster:
		return "master"
	case BusMusic:
		return "music"
	case BusSFX:
		return "sfx"
	case BusUI:
		return "ui"
	case BusGlobal:
		return "global"
	default:
		return "unknown"
	}
}

// PlayOptions tunes a single one-shot.
type PlayOptions struct {
	Bus    Bus     // target bus; BusMaster (or out of range) falls back to BusSFX
	Volume float64 // per-voice gain in (0,1]; <= 0 means full volume (1.0)
	Loop   bool    // restart automatically when the clip ends (looping ambiences)
}

// voicePlayer is the slice of *ebiten/audio.Player the Mixer relies on. Pulling
// it behind an interface lets tests drive all the bookkeeping — pooling, voice
// caps, coalescing, volume math — with a fake, no audio device required.
type voicePlayer interface {
	Play()
	Pause()
	Rewind() error
	IsPlaying() bool
	SetVolume(float64)
	Close() error
}

type voice struct {
	player voicePlayer
	clip   *Clip // which variation is playing (the pool keys on this)
	key    string
	bus    Bus
	loop   bool
}

// Mixer plays short sound effects through mix buses. A key can back several clip
// variations, one of which Play picks at random (so repeated impacts/footsteps
// don't machine-gun the same sample). It reuses finished players from a
// per-variation pool (so a burst never allocates), caps concurrent voices per
// bus with oldest-first stealing, and coalesces identical plays within a frame
// so a colony sim firing dozens of the same sound at once stays legible.
//
// A Mixer is not safe for concurrent use; drive it from the game-loop goroutine
// (Play during event handling, Update once per frame).
type Mixer struct {
	// newPlayer builds a player for a clip's PCM. nil disables playback (Play
	// becomes a no-op), which keeps the game running on a machine with no audio
	// device and lets tests inject a fake.
	newPlayer func(pcm []byte) voicePlayer
	// pick chooses a variation index in [0,n); overridable for deterministic
	// tests. Defaults to a uniform random pick.
	pick func(n int) int

	clips  map[string][]*Clip // key → variations
	volume [numBuses]float64

	pool      map[*Clip][]voicePlayer // idle players by specific variation
	active    []voice                 // playing voices, oldest first
	liveByBus [numBuses]int

	maxVoices   [numBuses]int
	maxPool     int
	coalesce    int            // max identical plays honoured per frame (0 = unlimited)
	playedFrame map[string]int // per-key play count since the last Update
}

// NewMixer builds a Mixer backed by the shared ebiten audio context (created on
// first use, matching the rest of this package).
func NewMixer() *Mixer {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
	ctx := audioContext
	return newMixer(func(pcm []byte) voicePlayer {
		return ctx.NewPlayerFromBytes(pcm)
	})
}

// newMixer is the shared constructor; tests call it with a fake factory. A nil
// factory yields a silent (no-op) mixer.
func newMixer(factory func(pcm []byte) voicePlayer) *Mixer {
	m := &Mixer{
		newPlayer: factory,
		pick: func(n int) int {
			if n <= 1 {
				return 0
			}
			return rand.Intn(n)
		},
		clips:       make(map[string][]*Clip),
		pool:        make(map[*Clip][]voicePlayer),
		playedFrame: make(map[string]int),
		maxPool:     8,
		coalesce:    3,
	}
	for i := range m.volume {
		m.volume[i] = 1
	}
	// Sensible per-bus ceilings; music is effectively single-voice here (the
	// music director owns crossfades separately).
	m.maxVoices = [numBuses]int{BusMusic: 2, BusSFX: 16, BusUI: 8, BusGlobal: 8}
	return m
}

// Register adds a clip as a variation under key. Calling it repeatedly with the
// same key builds up a variation set that Play chooses from at random.
func (m *Mixer) Register(key string, clip *Clip) {
	if clip == nil {
		return
	}
	m.clips[key] = append(m.clips[key], clip)
}

// Load decodes a file and registers it as a variation under key.
func (m *Mixer) Load(key, path string, t MusicType) error {
	clip, err := LoadClip(path, t)
	if err != nil {
		return err
	}
	m.Register(key, clip)
	return nil
}

// Has reports whether at least one clip variation is registered under key.
func (m *Mixer) Has(key string) bool {
	return len(m.clips[key]) > 0
}

// SetBusVolume sets a bus's gain (clamped to [0,1]).
func (m *Mixer) SetBusVolume(bus Bus, v float64) {
	if bus < 0 || bus >= numBuses {
		return
	}
	m.volume[bus] = clampVol(v)
}

// BusVolume returns a bus's current gain.
func (m *Mixer) BusVolume(bus Bus) float64 {
	if bus < 0 || bus >= numBuses {
		return 0
	}
	return m.volume[bus]
}

// gain is the fixed multiplier for a bus: master × bus.
func (m *Mixer) gain(bus Bus) float64 {
	return m.volume[BusMaster] * m.volume[bus]
}

// Play starts a one-shot. Unknown keys, a silent mixer, a coalesced repeat, or a
// full bus with nothing to steal are all silently ignored — playing audio is
// best-effort and must never disrupt the game.
func (m *Mixer) Play(key string, opts PlayOptions) {
	if m.newPlayer == nil {
		return
	}
	variations := m.clips[key]
	if len(variations) == 0 {
		return
	}
	bus := opts.Bus
	if bus <= BusMaster || bus >= numBuses {
		bus = BusSFX
	}
	if m.coalesce > 0 && m.playedFrame[key] >= m.coalesce {
		return
	}
	if m.liveByBus[bus] >= m.maxVoices[bus] && !m.stealOldest(bus) {
		return
	}

	clip := variations[m.pick(len(variations))]
	p := m.acquire(clip)
	vol := opts.Volume
	if vol <= 0 {
		vol = 1
	}
	p.SetVolume(clampVol(m.gain(bus) * vol))
	p.Play()

	m.active = append(m.active, voice{player: p, clip: clip, key: key, bus: bus, loop: opts.Loop})
	m.liveByBus[bus]++
	m.playedFrame[key]++
}

// Update reclaims finished voices back into the pool, restarts looping ones, and
// resets the per-frame coalescing window. Call once per frame.
func (m *Mixer) Update() {
	kept := m.active[:0]
	for _, v := range m.active {
		if v.player.IsPlaying() {
			kept = append(kept, v)
			continue
		}
		if v.loop {
			_ = v.player.Rewind()
			v.player.Play()
			kept = append(kept, v)
			continue
		}
		m.liveByBus[v.bus]--
		m.recycle(v.clip, v.player)
	}
	m.active = kept
	for k := range m.playedFrame {
		delete(m.playedFrame, k)
	}
}

// ActiveVoices is the number of currently playing voices (telemetry/tests).
func (m *Mixer) ActiveVoices() int { return len(m.active) }

// acquire returns a ready-to-play player for a specific variation, reusing a
// pooled one when available (rewound to the start) or minting a fresh one.
func (m *Mixer) acquire(clip *Clip) voicePlayer {
	if pl := m.pool[clip]; len(pl) > 0 {
		p := pl[len(pl)-1]
		m.pool[clip] = pl[:len(pl)-1]
		_ = p.Rewind()
		return p
	}
	return m.newPlayer(clip.pcm)
}

// stealOldest stops the oldest voice on a bus to free a slot, returning false if
// the bus has none (shouldn't happen when called at the cap).
func (m *Mixer) stealOldest(bus Bus) bool {
	for i, v := range m.active {
		if v.bus != bus {
			continue
		}
		v.player.Pause()
		m.recycle(v.clip, v.player)
		m.active = append(m.active[:i], m.active[i+1:]...)
		m.liveByBus[bus]--
		return true
	}
	return false
}

// recycle parks a finished player in its variation's pool for reuse, or closes
// it if the pool is already full (bounding idle players per variation).
func (m *Mixer) recycle(clip *Clip, p voicePlayer) {
	p.Pause()
	if len(m.pool[clip]) < m.maxPool {
		m.pool[clip] = append(m.pool[clip], p)
		return
	}
	_ = p.Close()
}

func clampVol(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
