package sense

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestWidth = 10
const TestHeight = 10

func TestNewSenseScapeCreatesAValidSoundScape(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	assert.Equal(t, TestWidth, s.width)
	assert.Equal(t, TestHeight, s.height)
	assert.Len(t, s.data, TestWidth)
	assert.Len(t, s.data[0], TestHeight)
}

func TestGetStimuliAtWithNothingThere(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	stimuli, _ := s.GetStimuliAt(0, 0)

	assert.Nil(t, stimuli)
}
func TestGetStimuliAtOutOfBounds(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	_, err := s.GetStimuliAt(TestWidth+1, 0)

	assert.NotNil(t, err)
}
func TestAddStimulus(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)
	//Test applying a brand new stimulus
	err := s.ApplyStimulus(0, 0, Stimulus{Type: SoundStimuli, Intensity: 1})
	assert.Nil(t, err)
	stimuli, _ := s.GetStimuliAt(0, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 1, stimuli[0].Intensity)

	//Test applying the same stimulus again
	err = s.ApplyStimulus(0, 0, Stimulus{Type: SoundStimuli, Intensity: 1})
	assert.Nil(t, err)
	stimuli, _ = s.GetStimuliAt(0, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 2, stimuli[0].Intensity)
}
func TestMakeSound(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)
	//Test applying a brand new stimulus
	s.MakeSound(5, 5, "TEST", 5)

	for x := 0; x < TestWidth; x++ {
		for y := 0; y < TestHeight; y++ {
			if len(s.data[x][y].Stimuli) > 0 {
				fmt.Print(s.data[x][y].Stimuli[0].Intensity)
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println("")
	}

	stimuli, _ := s.GetStimuliAt(5, 5)

	assert.NotNil(t, stimuli)
	assert.Equal(t, 5, stimuli[0].Intensity)
	stimuli, _ = s.GetStimuliAt(4, 4)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 5, stimuli[0].Intensity)
}
func TestAddStimulusOutOfBounds(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	err := s.ApplyStimulus(TestWidth+1, 0, Stimulus{Type: SoundStimuli, Intensity: 1})

	assert.NotNil(t, err)
}
func TestUpdateDoesNotError(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)
	err := s.Update()

	assert.Nil(t, err)
}
func TestUpdateDispersesScentCorrectly(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	err := s.ApplyStimulus(0, 0, Stimulus{Type: ScentStimuli, Intensity: 4, Decay: 1})
	assert.Nil(t, err)
	err = s.Update()
	assert.Nil(t, err)

	stimuli, _ := s.GetStimuliAt(0, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 3, stimuli[0].Intensity)
	stimuli, _ = s.GetStimuliAt(1, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 2, stimuli[0].Intensity)
	stimuli, _ = s.GetStimuliAt(2, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 1, stimuli[0].Intensity)
}
func TestUpdateDecaysPheremones(t *testing.T) {
	s := NewSenseScape(TestWidth, TestHeight)

	err := s.ApplyStimulus(0, 0, Stimulus{Type: PheremoneStimuli, Intensity: 5, Decay: 1})
	assert.Nil(t, err)
	err = s.Update()
	assert.Nil(t, err)

	stimuli, _ := s.GetStimuliAt(0, 0)
	assert.NotNil(t, stimuli)
	assert.Equal(t, 4, stimuli[0].Intensity)
}
