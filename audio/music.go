package audio

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// musicTrack is one streaming track (its player).
type musicTrack struct {
	player voicePlayer
}

// MusicDirector streams background music organised into named states ("menu",
// "field", "combat", …), each with a list of tracks. SetState crossfades to a
// random track of that state; when a track ends it rotates to another random
// track of the current state (a single-track state just loops). One master
// volume (the Music slider) scales everything.
//
// Not safe for concurrent use — drive it from the game loop: SetState/SetVolume
// as things change, Update once per frame.
type MusicDirector struct {
	// newPlayer builds a streaming player for a track; nil disables playback.
	newPlayer func(*AudioResource) voicePlayer
	pick      func(n int) int // random track chooser; overridable for tests

	states map[string][]*musicTrack
	byRes  map[*AudioResource]*musicTrack // dedup a file shared across states

	state    string
	current  *musicTrack
	outgoing *musicTrack // fading out during a crossfade
	curFade  float64     // 0..1 fade level of current
	outFade  float64     // 0..1 fade level of outgoing
	fadeStep float64     // fade delta per Update (~60/sec)
	volume   float64     // master music volume
}

// NewMusicDirector builds a director backed by the shared ebiten audio context.
func NewMusicDirector() *MusicDirector {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
	ctx := audioContext
	return newMusicDirector(func(res *AudioResource) voicePlayer {
		p, err := audio.NewPlayer(ctx, res.Source)
		if err != nil {
			return nil
		}
		return p
	})
}

func newMusicDirector(factory func(*AudioResource) voicePlayer) *MusicDirector {
	return &MusicDirector{
		newPlayer: factory,
		pick: func(n int) int {
			if n <= 1 {
				return 0
			}
			return rand.Intn(n)
		},
		states:   make(map[string][]*musicTrack),
		byRes:    make(map[*AudioResource]*musicTrack),
		fadeStep: 1.0 / 120.0, // ~2s at 60 fps
		volume:   1,
	}
}

// SetFadeSeconds sets the crossfade duration (assuming ~60 Update calls/sec).
func (m *MusicDirector) SetFadeSeconds(s float64) {
	if s <= 0 {
		m.fadeStep = 1
		return
	}
	m.fadeStep = 1.0 / (s * 60.0)
}

// AddTrack registers a resource under a state. A resource shared across states
// reuses one player, so a track can carry across a state change.
func (m *MusicDirector) AddTrack(state string, res *AudioResource) {
	if m.newPlayer == nil || res == nil {
		return
	}
	track, ok := m.byRes[res]
	if !ok {
		p := m.newPlayer(res)
		if p == nil {
			return
		}
		track = &musicTrack{player: p}
		m.byRes[res] = track
	}
	m.states[state] = append(m.states[state], track)
}

// SetVolume sets the master music volume (0..1).
func (m *MusicDirector) SetVolume(v float64) { m.volume = clampVol(v) }

// SetState switches the playing state, crossfading to a random track of it.
// A no-op when already in that state; an unknown/empty state fades music out.
func (m *MusicDirector) SetState(name string) {
	if name == m.state {
		return
	}
	m.state = name
	tracks := m.states[name]
	if len(tracks) == 0 {
		m.fadeOutCurrent()
		return
	}
	m.transitionTo(tracks[m.pick(len(tracks))])
}

// State returns the current state name.
func (m *MusicDirector) State() string { return m.state }

// Update advances fades and rotates to a new track when the current one ends.
// Call once per frame.
func (m *MusicDirector) Update() {
	if m.current != nil {
		if m.curFade < 1 {
			m.curFade += m.fadeStep
			if m.curFade > 1 {
				m.curFade = 1
			}
		}
		m.current.player.SetVolume(m.curFade * m.volume)
	}
	if m.outgoing != nil {
		m.outFade -= m.fadeStep
		if m.outFade <= 0 {
			m.stop(m.outgoing)
			m.outgoing = nil
		} else {
			m.outgoing.player.SetVolume(m.outFade * m.volume)
		}
	}
	// Rotate once the current track finishes (and no crossfade is in flight).
	if m.current != nil && m.outgoing == nil && !m.current.player.IsPlaying() {
		if tracks := m.states[m.state]; len(tracks) > 0 {
			m.transitionTo(m.rotationPick(tracks))
		}
	}
}

// transitionTo makes track the new current, crossfading from the old one. The
// same track (single-track state, or a rotation that lands on itself) just loops.
func (m *MusicDirector) transitionTo(track *musicTrack) {
	if track == m.current {
		_ = track.player.Rewind()
		track.player.Play()
		m.curFade = 1
		return
	}
	if m.outgoing != nil {
		m.stop(m.outgoing) // rapid switch: drop the previous outgoing
	}
	m.outgoing = m.current
	m.outFade = m.curFade
	m.current = track
	m.curFade = 0
	_ = track.player.Rewind()
	track.player.SetVolume(0)
	track.player.Play()
}

func (m *MusicDirector) fadeOutCurrent() {
	if m.current == nil {
		return
	}
	if m.outgoing != nil {
		m.stop(m.outgoing)
	}
	m.outgoing = m.current
	m.outFade = m.curFade
	m.current = nil
}

// rotationPick picks a random track of the state, avoiding an immediate repeat
// of the current track when the state has more than one.
func (m *MusicDirector) rotationPick(tracks []*musicTrack) *musicTrack {
	if len(tracks) == 1 {
		return tracks[0]
	}
	for i := 0; i < 4; i++ {
		if t := tracks[m.pick(len(tracks))]; t != m.current {
			return t
		}
	}
	return tracks[m.pick(len(tracks))]
}

func (m *MusicDirector) stop(track *musicTrack) {
	track.player.Pause()
	_ = track.player.Rewind()
}
