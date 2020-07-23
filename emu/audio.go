package emu

import (
	"math"

	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
	frequency  = 440
)

func NewAudioPlayer() *audio.Player {
	var audioContext *audio.Context
	if currentContext := audio.CurrentContext(); currentContext != nil {
		audioContext = currentContext
	} else {
		var err error
		audioContext, err = audio.NewContext(sampleRate)
		if err != nil {
			panic(err)
		}
	}

	// Pass the (infinite) stream to audio.NewPlayer.
	audioPlayer, err := audio.NewPlayer(audioContext, &stream{})
	if err != nil {
		panic(err)
	}
	audioPlayer.SetVolume(0)

	// After calling Play, the stream never ends as long as the player object lives.
	if err := audioPlayer.Play(); err != nil {
		panic(err)
	}

	return audioPlayer
}

// stream is an infinite stream of 440 Hz sine wave.
type stream struct {
	position  int64
	remaining []byte
}

// Read from io.Reader
// Read fills the data with sine wave samples.
func (s *stream) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]

		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	const length = int64(sampleRate / frequency)
	const max = 32767

	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * max)
		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)

		p++
	}

	s.position += int64(len(buf))
	s.position %= length * 4

	if origBuf != nil {
		n := copy(origBuf, buf)
		s.remaining = buf[n:]

		return n, nil
	}

	return len(buf), nil
}

// Close from io.Closer
func (s *stream) Close() error {
	return nil
}
