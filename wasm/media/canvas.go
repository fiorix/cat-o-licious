package media

import (
	"syscall/js"
)

// Rect represents a rectangle.
type Rect struct {
	X, Y,
	W, H int
}

// Canvas represents an HTML5 canvas.
type Canvas struct {
	Value js.Value
	ctx2d js.Value
}

// GetCanvas gets the canvas with given id from HTML5.
func GetCanvas(id string) Canvas {
	canvasEl := js.Global().Get("document").Call("getElementById", id)
	canvas2d := canvasEl.Call("getContext", "2d")

	c := Canvas{
		Value: canvasEl,
		ctx2d: canvas2d,
	}

	canvasEl.Set("width", c.ClientW())
	canvasEl.Set("height", c.ClientH())

	return c
}

// ClientW returns the client width.
func (c Canvas) ClientW() int {
	return c.Value.Get("clientWidth").Int()
}

// ClientH returns the client height.
func (c Canvas) ClientH() int {
	return c.Value.Get("clientHeight").Int()
}

// ClearRect clears rectangle r.
func (c Canvas) ClearRect(r Rect) {
	c.ctx2d.Call("clearRect", r.X, r.Y, r.W, r.H)
}

// DrawImage draws an image into the canvas.
func (c Canvas) DrawImage(img Image, x, y, w, h int) {
	c.ctx2d.Call("drawImage", img.Value, x, y, w, h)
}

// DrawImageFlipHorizontal ...
func (c Canvas) DrawImageFlipHorizontal(img Image, x, y, w, h int) {
	c.ctx2d.Call("save")
	c.ctx2d.Call("translate", x+w/2, y+h/2)
	c.ctx2d.Call("scale", -1, 1)
	c.ctx2d.Call("drawImage", img.Value, -w/2, -h/2, w, h)
	c.ctx2d.Call("restore")
}

// SetFont ...
func (c Canvas) SetFont(name, style string) {
	c.ctx2d.Set("font", "50px Score")
	c.ctx2d.Set("fillStyle", style)
}

// DrawText ...
func (c Canvas) DrawText(s string, x, y int) {
	c.ctx2d.Call("fillText", s, x, y)
}
