// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"log"
	"os"
	"runtime"
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
	renderDriver := os.Getenv("SDL_RENDER_DRIVER")
	userSelectedRenderer := renderDriver != ""
	if renderDriver == "" {
		renderDriver = "software"
		if runtime.GOOS == "linux" {
			renderDriver = "opengl"
		}
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, renderDriver)
	}
	videoDriver, err := sdl.GetCurrentVideoDriver()
	if err != nil {
		videoDriver = "unknown"
	}
	log.Printf("SDL video driver=%q render driver=%q", videoDriver, renderDriver)

	w, r, err := sdl.CreateWindowAndRenderer(
		int32(c.Width),
		int32(c.Height),
		sdl.WINDOW_SHOWN,
	)
	if err != nil && !userSelectedRenderer && runtime.GOOS == "linux" && renderDriver == "opengl" {
		log.Printf("failed to create opengl renderer, falling back to software: %v", err)
		renderDriver = "software"
		sdl.SetHint(sdl.HINT_RENDER_DRIVER, renderDriver)
		w, r, err = sdl.CreateWindowAndRenderer(
			int32(c.Width),
			int32(c.Height),
			sdl.WINDOW_SHOWN,
		)
	}
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
	p := e.s.Player()
	playerSpeed := int32(e.c.PlayerSpeed)
	fullscreen := false
	var viewport sdl.Rect
	fps := uint32(e.c.FPS)
	frameDelay := uint32(1000 / fps)

	running := true
	for {
		if !running {
			break
		}
		frameStart := sdl.GetTicks()

		// 1. Handle Input (Main Thread)
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch t := ev.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.State != sdl.PRESSED {
					continue
				}
				switch t.Keysym.Sym {
				case sdl.K_q:
					running = false
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

		// 2. Draw (Main Thread)
		viewport = e.r.GetViewport()
		e.s.Draw(time.Now(), &viewport)
		e.r.Present()

		// 3. Frame Rate Control
		frameElapsed := sdl.GetTicks() - frameStart
		if frameElapsed < frameDelay {
			sdl.Delay(frameDelay - frameElapsed)
		}
	}
}
