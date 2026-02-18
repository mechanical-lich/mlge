---
layout: default
title: Audio
nav_order: 6
---

# Audio

`github.com/mechanical-lich/mlge/audio`

Wraps Ebitengine's audio context for loading and playing background music in MP3 and OGG formats.

## Initialization

Call `Init()` once at startup to initialize the audio context:

```go
import "github.com/mechanical-lich/mlge/audio"

audio.Init()
```

## Loading Audio

```go
resource, err := audio.LoadAudioFromFile("assets/music/bgm.ogg", audio.TypeOgg)
resource2, err := audio.LoadAudioFromFile("assets/music/title.mp3", audio.TypeMP3)
```

### Music Types

| Constant | Format |
|----------|--------|
| `audio.TypeOgg` | OGG Vorbis |
| `audio.TypeMP3` | MP3 |

## Background Audio Player

The `BackgroundAudioPlayer` manages a playlist of audio resources and automatically advances to the next song when the current one finishes.

```go
player, err := audio.NewBackgroundAudioPlayer([]*audio.AudioResource{
    resource,
    resource2,
})
```

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Update` | `()` | Call each frame to manage playback |
| `SetActiveSong` | `(index int) error` | Switch to a specific song by index |
| `SetVolume` | `(volume float64)` | Set volume (0.0 to 1.0) |

### Usage in Game Loop

```go
func (g *Game) Update() error {
    g.bgPlayer.Update()
    // ...
    return nil
}
```

## Types

### AudioResource

```go
type AudioResource struct {
    Source    AudioStream
    MusicType MusicType
}
```

### AudioStream

```go
type AudioStream interface {
    io.ReadSeeker
    Length() int64
}
```
