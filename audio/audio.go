package audio

import (
	"errors"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var audioContext *audio.Context

type AudioResource struct {
	Source    AudioStream
	MusicType MusicType
}

type AudioStream interface {
	io.ReadSeeker
	Length() int64
}

type MusicType int

const sampleRate = 32000

const (
	TypeOgg MusicType = iota
	TypeMP3
	TypeWav
)

func (t MusicType) String() string {
	switch t {
	case TypeOgg:
		return "Ogg"
	case TypeMP3:
		return "MP3"
	case TypeWav:
		return "Wav"
	default:
		return "unsupported"
	}
}

func Init() {
	audioContext = audio.NewContext(sampleRate)
}

func LoadAudioFromFile(path string, musicType MusicType) (*AudioResource, error) {
	fileStream, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}

	var s AudioStream

	switch musicType {
	case TypeOgg:
		var err error

		s, err = vorbis.DecodeWithSampleRate(sampleRate, fileStream)
		if err != nil {
			return nil, err
		}
	case TypeMP3:
		var err error
		s, err = mp3.DecodeWithSampleRate(sampleRate, fileStream)
		if err != nil {
			return nil, err
		}
	case TypeWav:
		var err error
		s, err = wav.DecodeWithSampleRate(sampleRate, fileStream)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid music type")
	}

	aR := &AudioResource{MusicType: musicType}

	aR.Source = s
	return aR, nil
}
