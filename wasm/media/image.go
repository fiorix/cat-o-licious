package media

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"
)

// Image represents a js image.
type Image struct {
	Value js.Value
}

// NewImage creates and initializes an Image detached from the DOM.
func NewImage(uri, alt string) (Image, error) {
	img := js.Global().Get("Image").New()
	img.Set("alt", alt)

	errc := make(chan error, 1)

	var once sync.Once
	var onLoad, onAbort, onStalled, onSuspend, onError js.Func

	finalize := func(err error) {
		once.Do(func() {
			// Detach handlers BEFORE releasing funcs to prevent JS from calling a released func.
			img.Set("onload", js.Null())
			img.Set("onabort", js.Null())
			img.Set("onstalled", js.Null())
			img.Set("onsuspend", js.Null())
			img.Set("onerror", js.Null())

			// Release all js.Func resources.
			onLoad.Release()
			onAbort.Release()
			onStalled.Release()
			onSuspend.Release()
			onError.Release()

			if err != nil {
				errc <- err
			}
			close(errc)
		})
	}

	onLoad = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		finalize(nil)
		return nil
	})

	onAbort = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		finalize(fmt.Errorf("aborted: %s", uri))
		return nil
	})

	onStalled = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		finalize(fmt.Errorf("stalled: %s", uri))
		return nil
	})

	onSuspend = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		finalize(fmt.Errorf("suspended: %s", uri))
		return nil
	})

	onError = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		finalize(fmt.Errorf("failed: %s", uri))
		return nil
	})

	img.Set("onload", onLoad)
	img.Set("onabort", onAbort)
	img.Set("onstalled", onStalled)
	img.Set("onsuspend", onSuspend)
	img.Set("onerror", onError)

	img.Set("src", uri)

	if err := <-errc; err != nil {
		return Image{}, err
	}

	return Image{img}, nil
}

// W returns the image width.
func (img Image) W() int {
	return img.Value.Get("width").Int()
}

// H returns the image height.
func (img Image) H() int {
	return img.Value.Get("height").Int()
}

// NewImageSet ...
func NewImageSet(uriprefix string) ([]Image, error) {
	i := 0
	var imgs []Image
	for {
		idx := strconv.Itoa(i + 1)
		img, err := NewImage(uriprefix+idx+".png", idx)
		if err != nil {
			if i == 0 {
				return nil, err
			}
			break
		}
		imgs = append(imgs, img)
		i++
	}
	return imgs, nil
}
