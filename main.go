package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/audio"
	"github.com/mechanical-lich/mlge/event"
)

const (
	Test event.EventType = iota
)

var m *event.QueuedEventManager

type TestEventData struct {
	x int
	y int
}

func (t TestEventData) GetType() event.EventType {
	return Test
}

type TestListener struct {
}

func (t *TestListener) HandleEvent(data event.EventData) error {
	fmt.Println("Handling event: ", data)
	m.QueueEvent(data)
	return nil
}

type Game struct {
	b *audio.BackgroundAudioPlayer
}

func (g *Game) Update() error {
	g.b.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 100, 100
}
func main() {
	m = &event.QueuedEventManager{}

	testListener := &TestListener{}

	m.RegisterListener(testListener, Test)
	m.RegisterListener(testListener, Test)
	m.RegisterListener(testListener, Test)

	//Test event send
	testEventData := TestEventData{x: 5, y: 6}
	m.QueueEvent(testEventData)
	m.HandleQueue()
	m.HandleQueue()
	//m.UnregisterListener(testListener, Test)
	m.UnregisterListenerFromAll(testListener)

	audio.Init()
	a, err := audio.LoadAudioFromFile("./assets/audio/ding.mp3", audio.TypeMP3)
	if err != nil {
		log.Fatal(err)
	}

	b, err := audio.NewBackgroundAudioPlayer([]*audio.AudioResource{a})
	if err != nil {
		log.Fatal(err)
	}
	b.SetActiveSong(0)
	b.SetVolume(.4)

	g := &Game{b: b}

	ebiten.SetWindowSize(100, 100)
	ebiten.SetWindowTitle("Test")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
