package media

import (
	"fmt"
	"strconv"
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

	img.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(errc)
		return nil
	}))

	img.Set("onabort", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("aborted: %s", uri)
		close(errc)
		return nil
	}))

	img.Set("onstalled", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("stalled: %s", uri)
		close(errc)
		return nil
	}))

	img.Set("onsuspend", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("suspended: %s", uri)
		close(errc)
		return nil
	}))

	img.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errc <- fmt.Errorf("failed: %s", uri)
		defer close(errc)
		return nil
	}))

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
