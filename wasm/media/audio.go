package media

import (
	"fmt"
	"syscall/js"
)

var audioCtx js.Value

func getAudioContext() js.Value {
	if audioCtx.IsUndefined() {
		ac := js.Global().Get("AudioContext")
		if ac.IsUndefined() {
			ac = js.Global().Get("webkitAudioContext")
		}
		audioCtx = ac.New()
	}
	return audioCtx
}

// ResumeAudioContext resumes the audio context if it is suspended.
// Must be called from a user gesture handler (e.g., touch/click)
// to satisfy iOS Safari's autoplay restrictions.
func ResumeAudioContext() {
	ctx := getAudioContext()
	if ctx.Get("state").String() == "suspended" {
		ctx.Call("resume")
	}
}

// Audio represents an audio source backed by the Web Audio API.
type Audio struct {
	state *audioState
}

type audioState struct {
	buffer  js.Value // decoded AudioBuffer
	started bool     // whether PlayLoop has been initiated
}

// NewAudio creates and initializes an Audio source that loads the given URI.
func NewAudio(uri string) (Audio, error) {
	ctx := getAudioContext()

	errc := make(chan error, 1)
	var buffer js.Value

	onDecode := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		buffer = args[0]
		close(errc)
		return nil
	})

	onDecodeError := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("audio decode failed: %s", uri)
		close(errc)
		return nil
	})

	onFetchError := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("audio fetch failed: %s", uri)
		close(errc)
		return nil
	})

	onArrayBuffer := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("decodeAudioData", args[0], onDecode, onDecodeError)
		return nil
	})

	onResponse := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return args[0].Call("arrayBuffer")
	})

	js.Global().Call("fetch", uri).
		Call("then", onResponse).
		Call("then", onArrayBuffer).
		Call("catch", onFetchError)

	if err := <-errc; err != nil {
		return Audio{}, err
	}

	return Audio{state: &audioState{buffer: buffer}}, nil
}

// Play plays the audio from the beginning.
func (a Audio) Play() {
	if a.state == nil {
		return
	}
	ctx := getAudioContext()
	source := ctx.Call("createBufferSource")
	source.Set("buffer", a.state.buffer)
	source.Call("connect", ctx.Get("destination"))
	source.Call("start", 0)
}

// PlayLoop plays the audio in a loop.
func (a Audio) PlayLoop() {
	if a.state == nil || a.state.started {
		return
	}
	ctx := getAudioContext()
	source := ctx.Call("createBufferSource")
	source.Set("buffer", a.state.buffer)
	source.Set("loop", true)
	source.Call("connect", ctx.Get("destination"))
	source.Call("start", 0)
	a.state.started = true
}

// Playing returns true if the audio is currently playing.
func (a Audio) Playing() bool {
	if a.state == nil {
		return false
	}
	return a.state.started
}
