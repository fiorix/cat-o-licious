package game

import (
	"sync/atomic"
	"time"

	"github.com/fiorix/cat-o-licious/wasm/media"
)

// Engine ...
type Engine interface {
	Run()
}

type engine struct {
	c media.Canvas
	s Scene
}

// NewEngine ...
func NewEngine(canvasID string) (Engine, error) {
	canvas := media.GetCanvas(canvasID)

	canvas.SetFont("50px Score", "red")

	scene, err := NewScene()
	if err != nil {
		return nil, err
	}

	e := &engine{
		c: canvas,
		s: scene,
	}

	return e, nil
}

func (e *engine) Run() {
	const fps = 30
	const playerSpeed = 20

	media.OnKey(media.KeyDown, func(key string) {
		switch key {
		case "a", "A", "ArrowLeft":
			e.s.Player().Move(Left, playerSpeed)
		case "d", "D", "ArrowRight":
			e.s.Player().Move(Right, playerSpeed)
		}
	})

	var clicking int32

	e.c.OnMouse(func(click media.MouseClick, x, y int) {
		switch click {
		case media.MouseDown:
			side := Left
			if x > e.c.ClientW()/2 {
				side = Right
			}
			atomic.StoreInt32(&clicking, int32(side))
		case media.MouseUp:
			atomic.StoreInt32(&clicking, int32(Center))
		}
	})

	for {
		movingSide := Direction(atomic.LoadInt32(&clicking))
		switch movingSide {
		case Left, Right:
			e.s.Player().Move(movingSide, playerSpeed)
		}

		e.s.Draw(time.Now(), e.c)
		time.Sleep(1 * time.Second / fps)
	}
}
