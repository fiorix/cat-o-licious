// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Version is the version of the game engine.
var Version = "tip"

// Engine is the game engine.
type Engine interface {
	// Run runs the engine. Blocks until Q is pressed.
	Run()
}

type engine struct {
	c *Config
	w *sdl.Window
	r *sdl.Renderer
	s Scene
}

// NewEngine creates and initializes a new game engine.
func NewEngine(c *Config) (Engine, error) {
	w, r, err := sdl.CreateWindowAndRenderer(
		c.Width,
		c.Height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}
	w.SetTitle("cat-o-licious " + Version)
	s, err := NewScene(r)
	if err != nil {
		return nil, err
	}
	return &engine{c, w, r, s}, nil
}

// Run implements the Engine interface.
func (e *engine) Run() {
	stop := make(chan struct{})
	go e.Draw(stop)
	e.HandleInput()
	close(stop)
}

// Draw draws frames at a given FPS.
func (e *engine) Draw(stop chan struct{}) {
	var viewport sdl.Rect
	fps := uint32(e.c.FPS)
loop:
	for {
		select {
		case <-stop:
			break loop
		default:
		}
		e.r.GetViewport(&viewport)
		e.s.Draw(time.Now(), &viewport)
		e.r.Present()
		sdl.Delay(1000 / fps)
	}
}

// HandleInput handles keyboard input.
func (e *engine) HandleInput() {
	p := e.s.Player()
	playerSpeed := int32(e.c.PlayerSpeed)
	fullscreen := false
loop:
	for {
		ev := sdl.WaitEvent()
		if ev == nil { // cgo barfed here before
			sdl.Delay(1000)
			continue
		}
		switch ev.(type) {
		case *sdl.QuitEvent:
			break loop
		case *sdl.KeyDownEvent:
			switch ev.(*sdl.KeyDownEvent).Keysym.Sym {
			case sdl.K_q:
				break loop
			case sdl.K_f:
				if fullscreen {
					e.w.SetFullscreen(0)
				} else {
					e.w.SetFullscreen(sdl.WINDOW_FULLSCREEN)
				}
				fullscreen = !fullscreen
			case sdl.K_LEFT, sdl.K_a:
				p.Move(Left, playerSpeed)
			case sdl.K_RIGHT, sdl.K_d:
				p.Move(Right, playerSpeed)
			}
		}
	}
}
