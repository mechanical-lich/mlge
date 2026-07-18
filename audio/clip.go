package audio

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Clip is a short sound effect decoded fully into memory as raw PCM. Unlike the
// streaming AudioResource (a single io.ReadSeeker, one playback at a time), a
// Clip's bytes can back many independent players at once, so rapid or
// overlapping one-shots — UI clicks, footsteps, a dozen workers digging — never
// block or stutter against each other. Load long music with LoadAudioFromFile
// instead; buffering a whole soundtrack into memory is wasteful.
type Clip struct {
	pcm []byte
}

// LoadClip decodes an ogg or mp3 file fully into memory at the engine sample
// rate, ready to be registered with a Mixer.
func LoadClip(path string, musicType MusicType) (*Clip, error) {
	f, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}

	var stream AudioStream
	switch musicType {
	case TypeOgg:
		stream, err = vorbis.DecodeWithSampleRate(sampleRate, f)
	case TypeMP3:
		stream, err = mp3.DecodeWithSampleRate(sampleRate, f)
	default:
		return nil, errors.New("invalid music type")
	}
	if err != nil {
		return nil, err
	}

	pcm, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	return &Clip{pcm: pcm}, nil
}

// MusicTypeFromExt infers a decoder from a file extension, so callers can load
// a mix of ".ogg" and ".mp3" assets without hard-coding the format per file.
func MusicTypeFromExt(path string) (MusicType, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ogg":
		return TypeOgg, nil
	case ".mp3":
		return TypeMP3, nil
	case ".wav":
		return TypeWav, nil
	default:
		return 0, fmt.Errorf("audio: unsupported extension %q (want .ogg, .mp3 or .wav)", filepath.Ext(path))
	}
}

// LoadClipAuto is LoadClip with the format inferred from the file extension.
func LoadClipAuto(path string) (*Clip, error) {
	t, err := MusicTypeFromExt(path)
	if err != nil {
		return nil, err
	}
	return LoadClip(path, t)
}

// NewClipFromPCM wraps already-decoded 16-bit little-endian stereo PCM at the
// engine sample rate. Handy for synthesised or test audio that never touches a
// file.
func NewClipFromPCM(pcm []byte) *Clip {
	return &Clip{pcm: pcm}
}

// Len reports the clip's decoded size in bytes.
func (c *Clip) Len() int { return len(c.pcm) }
