package media

import (
	"fmt"
	"syscall/js"
)

// Audio represents an HTML5 audio element.
type Audio struct {
	Value js.Value
}

// NewAudio creates and initializes an Audio element that loads the given URI.
func NewAudio(uri string) (Audio, error) {
	audio := js.Global().Get("Audio").New(uri)

	errc := make(chan error, 1)

	audio.Set("oncanplaythrough", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(errc)
		return nil
	}))

	audio.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("audio failed: %s", uri)
		close(errc)
		return nil
	}))

	audio.Call("load")

	if err := <-errc; err != nil {
		return Audio{}, err
	}

	return Audio{audio}, nil
}

// Play plays the audio from the beginning.
func (a Audio) Play() {
	a.Value.Set("currentTime", 0)
	a.Value.Call("play")
}

// PlayLoop plays the audio in a loop.
func (a Audio) PlayLoop() {
	a.Value.Set("loop", true)
	a.Value.Call("play")
}

// Playing returns true if the audio is currently playing.
func (a Audio) Playing() bool {
	return !a.Value.Get("paused").Bool()
}
