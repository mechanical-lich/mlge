package audio

import "testing"

// newTestMusic builds a director whose players are fakePlayers, returning the
// director and the players in creation order (one per unique resource).
func newTestMusic() (*MusicDirector, *[]*fakePlayer) {
	players := &[]*fakePlayer{}
	m := newMusicDirector(func(*AudioResource) voicePlayer {
		p := &fakePlayer{}
		*players = append(*players, p)
		return p
	})
	m.pick = func(int) int { return 0 } // deterministic: first track of a state
	m.fadeStep = 1                       // instant fades so a single Update settles
	return m, players
}

func TestMusicDirectorPlaysStateTrack(t *testing.T) {
	m, players := newTestMusic()
	m.AddTrack("field", &AudioResource{})

	m.SetState("field")
	if !(*players)[0].playing {
		t.Fatal("field track should be playing after SetState")
	}
	// Same state again is a no-op (no restart).
	rewinds := (*players)[0].rewinds
	m.SetState("field")
	if (*players)[0].rewinds != rewinds {
		t.Fatal("re-entering the same state should not restart the track")
	}
}

func TestMusicDirectorCrossfades(t *testing.T) {
	m, players := newTestMusic()
	field := &AudioResource{}
	combat := &AudioResource{}
	m.AddTrack("field", field)
	m.AddTrack("combat", combat)

	m.SetState("field")
	m.Update() // field ramps to full
	if (*players)[0].volume != 1 {
		t.Fatalf("field vol = %v, want 1", (*players)[0].volume)
	}

	m.SetState("combat") // crossfade: field out, combat in
	if !(*players)[1].playing {
		t.Fatal("combat should start playing")
	}
	m.Update() // outgoing field fades to 0 and stops; combat to full
	if (*players)[0].playing {
		t.Fatal("field should have faded out and stopped")
	}
	if (*players)[1].volume != 1 {
		t.Fatalf("combat vol = %v, want 1", (*players)[1].volume)
	}
}

func TestMusicDirectorLoopsSingleTrackState(t *testing.T) {
	m, players := newTestMusic()
	m.AddTrack("combat", &AudioResource{})
	m.SetState("combat")
	m.Update()

	p := (*players)[0]
	rewinds := p.rewinds
	p.finish() // track reached its end
	m.Update() // single-track state → loop (rewind + play)
	if !p.playing {
		t.Fatal("single-track state should loop")
	}
	if p.rewinds <= rewinds {
		t.Fatal("looping should rewind the track")
	}
}

func TestMusicDirectorMasterVolume(t *testing.T) {
	m, players := newTestMusic()
	m.AddTrack("field", &AudioResource{})
	m.SetVolume(0.5)
	m.SetState("field")
	m.Update()
	if got := (*players)[0].volume; got != 0.5 { // curFade 1 × master 0.5
		t.Fatalf("volume = %v, want 0.5 (master applied)", got)
	}
}
