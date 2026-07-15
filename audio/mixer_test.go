package audio

import "testing"

// fakePlayer stands in for *ebiten/audio.Player so the Mixer's bookkeeping can
// be exercised without an audio device. playing flips to false when the test
// calls finish(), simulating a clip reaching its end before the next Update.
type fakePlayer struct {
	pcm     []byte
	volume  float64
	playing bool
	plays   int
	rewinds int
	closed  bool
}

func (p *fakePlayer) Play()               { p.playing = true; p.plays++ }
func (p *fakePlayer) Pause()              { p.playing = false }
func (p *fakePlayer) Rewind() error       { p.rewinds++; return nil }
func (p *fakePlayer) IsPlaying() bool     { return p.playing }
func (p *fakePlayer) SetVolume(v float64) { p.volume = v }
func (p *fakePlayer) Close() error        { p.closed = true; return nil }
func (p *fakePlayer) finish()             { p.playing = false }

// newTestMixer builds a Mixer whose players are fakePlayers, and returns the
// mixer plus a slice recording every player it created (in creation order).
func newTestMixer() (*Mixer, *[]*fakePlayer) {
	created := &[]*fakePlayer{}
	m := newMixer(func(pcm []byte) voicePlayer {
		p := &fakePlayer{pcm: pcm}
		*created = append(*created, p)
		return p
	})
	return m, created
}

func registerSilence(m *Mixer, keys ...string) {
	for _, k := range keys {
		m.Register(k, NewClipFromPCM([]byte{0, 0, 0, 0}))
	}
}

func TestPlayIgnoresUnknownKeyAndSilentMixer(t *testing.T) {
	// Silent mixer (nil factory): Play must be a harmless no-op.
	silent := newMixer(nil)
	registerSilence(silent, "click")
	silent.Play("click", PlayOptions{Bus: BusUI})
	if got := silent.ActiveVoices(); got != 0 {
		t.Fatalf("silent mixer played a voice: %d", got)
	}

	m, created := newTestMixer()
	m.Play("missing", PlayOptions{Bus: BusUI})
	if got := m.ActiveVoices(); got != 0 {
		t.Fatalf("unknown key produced a voice: %d", got)
	}
	if len(*created) != 0 {
		t.Fatalf("unknown key created a player: %d", len(*created))
	}
}

func TestPlayAppliesMasterAndBusGain(t *testing.T) {
	m, created := newTestMixer()
	registerSilence(m, "click")
	m.SetBusVolume(BusUI, 0.5)
	m.SetBusVolume(BusMaster, 0.5)

	m.Play("click", PlayOptions{Bus: BusUI, Volume: 0.8})
	if len(*created) != 1 {
		t.Fatalf("expected 1 player, got %d", len(*created))
	}
	// 0.5 master × 0.5 bus × 0.8 per-voice = 0.2
	if got := (*created)[0].volume; got != 0.2 {
		t.Fatalf("effective volume = %v, want 0.2", got)
	}
}

func TestPlayDefaultsZeroVolumeToFull(t *testing.T) {
	m, created := newTestMixer()
	registerSilence(m, "click")
	m.Play("click", PlayOptions{Bus: BusUI}) // Volume left 0 → full
	if got := (*created)[0].volume; got != 1 {
		t.Fatalf("zero volume should default to 1.0, got %v", got)
	}
}

func TestBusFallbackForMaster(t *testing.T) {
	m, _ := newTestMixer()
	registerSilence(m, "boom")
	// Targeting BusMaster is invalid; it must fall back to SFX and still play.
	m.Play("boom", PlayOptions{Bus: BusMaster})
	if got := m.liveByBus[BusSFX]; got != 1 {
		t.Fatalf("BusMaster target should fall back to SFX, live sfx = %d", got)
	}
}

func TestCoalesceLimitsIdenticalPlaysPerFrame(t *testing.T) {
	m, created := newTestMixer()
	registerSilence(m, "dig")
	m.coalesce = 3
	for i := 0; i < 10; i++ {
		m.Play("dig", PlayOptions{Bus: BusSFX})
	}
	if len(*created) != 3 {
		t.Fatalf("coalesce should cap identical plays at 3/frame, got %d", len(*created))
	}
	// A new frame resets the window.
	m.Update()
	m.Play("dig", PlayOptions{Bus: BusSFX})
	if len(*created) != 4 {
		t.Fatalf("coalesce window should reset each Update, got %d players", len(*created))
	}
}

func TestVoiceCapStealsOldest(t *testing.T) {
	m, created := newTestMixer()
	m.coalesce = 0 // don't let coalescing mask the cap in this test
	m.maxVoices[BusUI] = 2
	registerSilence(m, "a", "b", "c")

	m.Play("a", PlayOptions{Bus: BusUI})
	m.Play("b", PlayOptions{Bus: BusUI})
	m.Play("c", PlayOptions{Bus: BusUI}) // exceeds cap → steal oldest ("a")

	if got := m.liveByBus[BusUI]; got != 2 {
		t.Fatalf("live UI voices = %d, want 2 (capped)", got)
	}
	if got := m.ActiveVoices(); got != 2 {
		t.Fatalf("active voices = %d, want 2", got)
	}
	// "a" was the oldest → its player should have been paused (stolen).
	if (*created)[0].playing {
		t.Fatalf("oldest voice should have been stolen (paused)")
	}
	if !(*created)[2].playing {
		t.Fatalf("newest voice should be playing")
	}
}

func TestUpdateReclaimsAndPoolReuses(t *testing.T) {
	m, created := newTestMixer()
	registerSilence(m, "click")

	m.Play("click", PlayOptions{Bus: BusUI})
	first := (*created)[0]
	first.finish() // clip ended
	m.Update()     // reclaim into pool

	if got := m.liveByBus[BusUI]; got != 0 {
		t.Fatalf("finished voice not reclaimed, live = %d", got)
	}
	if got := m.ActiveVoices(); got != 0 {
		t.Fatalf("active voices after reclaim = %d, want 0", got)
	}

	m.Play("click", PlayOptions{Bus: BusUI}) // should reuse pooled player
	if len(*created) != 1 {
		t.Fatalf("expected pool reuse (1 player total), got %d", len(*created))
	}
	if first.rewinds == 0 {
		t.Fatalf("reused player should have been rewound")
	}
}

func TestPlayPicksRegisteredVariation(t *testing.T) {
	m, created := newTestMixer()
	m.Register("impact", NewClipFromPCM([]byte{1}))
	m.Register("impact", NewClipFromPCM([]byte{2}))
	m.Register("impact", NewClipFromPCM([]byte{3}))
	if got := len(m.clips["impact"]); got != 3 {
		t.Fatalf("expected 3 variations registered, got %d", got)
	}

	m.pick = func(int) int { return 1 } // deterministic: middle variation
	m.Play("impact", PlayOptions{Bus: BusSFX})
	if len(*created) != 1 {
		t.Fatalf("expected 1 player, got %d", len(*created))
	}
	if string((*created)[0].pcm) != string([]byte{2}) {
		t.Fatalf("played variation %v, want the index-1 clip {2}", (*created)[0].pcm)
	}
}

func TestVariationPoolIsPerClip(t *testing.T) {
	m, created := newTestMixer()
	m.Register("imp", NewClipFromPCM([]byte{1}))
	m.Register("imp", NewClipFromPCM([]byte{2}))
	idx := 0
	m.pick = func(int) int { return idx }

	// Play variation 0, let it finish, reclaim it.
	m.Play("imp", PlayOptions{Bus: BusSFX})
	(*created)[0].finish()
	m.Update()

	// Same variation again → reuse the pooled player (still 1 created).
	m.Play("imp", PlayOptions{Bus: BusSFX})
	if len(*created) != 1 {
		t.Fatalf("variation 0 should reuse its pooled player, got %d players", len(*created))
	}

	// A different variation can't reuse variation 0's player → a new one.
	(*created)[len(*created)-1].finish()
	m.Update()
	idx = 1
	m.Play("imp", PlayOptions{Bus: BusSFX})
	if len(*created) != 2 {
		t.Fatalf("variation 1 needs its own player, got %d players", len(*created))
	}
}

func TestLoopingVoiceRestartsOnUpdate(t *testing.T) {
	m, created := newTestMixer()
	registerSilence(m, "hum")
	m.Play("hum", PlayOptions{Bus: BusGlobal, Loop: true})
	p := (*created)[0]

	p.finish() // reached the end
	m.Update() // loop → restart, stays active
	if got := m.ActiveVoices(); got != 1 {
		t.Fatalf("looping voice should stay active, got %d", got)
	}
	if !p.playing {
		t.Fatalf("looping voice should have been restarted")
	}
	if len(*created) != 1 {
		t.Fatalf("looping restart should reuse the same player, got %d", len(*created))
	}
}
