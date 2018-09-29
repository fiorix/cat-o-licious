package game

import (
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
	for {
		e.s.Draw(time.Now(), e.c)
		time.Sleep(1 * time.Second / fps)
	}
}
