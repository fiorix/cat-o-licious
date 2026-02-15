package media

import (
	"fmt"
	"sync"
	"syscall/js"
)

// Audio represents an HTML5 audio element.
type Audio struct {
	Value js.Value
}

// NewAudio creates and initializes an Audio element that loads the given URI.
func NewAudio(uri string) (Audio, error) {
	audio := js.Global().Get("Audio").New(uri)

	result := make(chan error, 1)

	var once sync.Once
	var canPlayFunc, errorFunc js.Func

	finalize := func(err error) {
		once.Do(func() {
			// Detach handlers BEFORE releasing funcs (avoids JS calling a released func).
			audio.Set("oncanplaythrough", js.Null())
			audio.Set("onerror", js.Null())

			if canPlayFunc.Type() == js.TypeFunction {
				canPlayFunc.Release()
			}
			if errorFunc.Type() == js.TypeFunction {
				errorFunc.Release()
			}

			result <- err
		})
	}

	canPlayFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		finalize(nil)
		return nil
	})
	errorFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		finalize(fmt.Errorf("audio failed: %s", uri))
		return nil
	})

	audio.Set("oncanplaythrough", canPlayFunc)
	audio.Set("onerror", errorFunc)
	audio.Call("load")

	if err := <-result; err != nil {
		return Audio{}, err
	}

	return Audio{audio}, nil
}

// Play plays the audio from the beginning.
func (a Audio) Play() {
	a.Value.Set("currentTime", 0)
	ensurePromiseHandled(a.Value.Call("play"))
}

// PlayLoop plays the audio in a loop.
func (a Audio) PlayLoop() {
	a.Value.Set("loop", true)
	ensurePromiseHandled(a.Value.Call("play"))
}

// Playing returns true if the audio is currently playing.
func (a Audio) Playing() bool {
	return !a.Value.Get("paused").Bool()
}

// ensurePromiseHandled attaches handlers to Promise-like values returned by play().
// This prevents unhandled promise rejections from destabilizing the WASM runtime,
// and ensures js.Func resources are always released on resolve/reject.
func ensurePromiseHandled(v js.Value) {
	if v.Type() != js.TypeObject {
		return
	}

	then := v.Get("then")
	catch := v.Get("catch")
	if then.Type() != js.TypeFunction || catch.Type() != js.TypeFunction {
		return
	}

	var onFulfilled js.Func
	var onRejected js.Func

	cleanup := func() {
		if onFulfilled.Type() == js.TypeFunction {
			onFulfilled.Release()
		}
		if onRejected.Type() == js.TypeFunction {
			onRejected.Release()
		}
	}

	onFulfilled = js.FuncOf(func(this js.Value, args []js.Value) any {
		// Preserve value (best-effort), but we don't rely on it.
		defer cleanup()
		if len(args) > 0 {
			return args[0]
		}
		return nil
	})

	onRejected = js.FuncOf(func(this js.Value, args []js.Value) any {
		// Swallow rejection and cleanup; avoids “Unhandled Promise rejection”.
		defer cleanup()
		return nil
	})

	// Attach both handlers so cleanup runs whether it resolves or rejects.
	v.Call("then", onFulfilled)
	v.Call("catch", onRejected)
}
