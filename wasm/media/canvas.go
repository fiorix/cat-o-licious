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
func (c Canvas) DrawImage(img Image, r Rect) {
	c.ctx2d.Call("drawImage", img.Value, r.X, r.Y, r.W, r.H)
}

// DrawImageFlipHorizontal ...
func (c Canvas) DrawImageFlipHorizontal(img Image, r Rect) {
	c.ctx2d.Call("save")
	c.ctx2d.Call("translate", r.X+r.W/2, r.Y+r.H/2)
	c.ctx2d.Call("scale", -1, 1)
	c.ctx2d.Call("drawImage", img.Value, -r.W/2, -r.H/2, r.W, r.H)
	c.ctx2d.Call("restore")
}

// SetFont ...
func (c Canvas) SetFont(name, style string) {
	c.ctx2d.Set("font", name)
	c.ctx2d.Set("fillStyle", style)
}

// DrawText ...
func (c Canvas) DrawText(s string, x, y int) {
	c.ctx2d.Call("fillText", s, x, y)
}

// MouseClick ...
type MouseClick int8

// Mouse clicks ...
const (
	MouseUp MouseClick = iota
	MouseDown
)

var mouseClickNames = map[MouseClick]string{
	MouseUp:   "mouseup",
	MouseDown: "mousedown",
}

func (m MouseClick) String() string {
	return mouseClickNames[m]
}

// OnMouse ...
func (c Canvas) OnMouse(mapTouch bool, f func(click MouseClick, x, y int)) {
	pos := func(ev js.Value) (int, int) {
		x := ev.Get("clientX").Int()
		y := ev.Get("clientY").Int()
		return x, y
	}
	for _, ev := range []MouseClick{MouseUp, MouseDown} {
		ev := ev
		c.Value.Call("addEventListener", ev.String(), js.NewCallback(func(args []js.Value) {
			x, y := pos(args[0])
			f(ev, x, y)
		}))
	}
	if !mapTouch {
		return
	}
	for _, ev := range []struct {
		touchEv string
		mouseEv MouseClick
	}{
		{"touchstart", MouseDown},
		{"touchend", MouseUp},
		{"touchcancel", MouseUp},
	} {
		ev := ev
		c.Value.Call("addEventListener", ev.touchEv, js.NewCallback(func(args []js.Value) {
			if args[0].Get("touches").Get("length").Int() > 1 {
				return // ignore multi-touch events
			}
			x, y := pos(args[0].Get("changedTouches").Call("item", 0))
			f(ev.mouseEv, x, y)
		}))
	}
}
