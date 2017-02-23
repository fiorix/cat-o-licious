// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

// Package game provides the cat-o-licious game.
package game

import (
	"github.com/veandco/go-sdl2/sdl"
	sdlimg "github.com/veandco/go-sdl2/sdl_image"
	sdlmix "github.com/veandco/go-sdl2/sdl_mixer"
	sdlttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// Config is the game configuration.
type Config struct {
	// FPS is the rendering rate of the game. It cannot be changed.
	FPS int

	// Width is the initial width for the game window.
	Width int

	// Height is the initial height of the game window.
	Height int

	// PlayerSpeed is the speed of the player's lateral movement.
	PlayerSpeed int
}

// DefaultConfig is the game's default configuration.
var DefaultConfig = Config{
	FPS:         30,
	Width:       800,
	Height:      600,
	PlayerSpeed: 20,
}

// Run runs the game with the given config, which is optional.
// The nil config uses DefaultConfig.
func Run(c *Config) error {
	if c == nil {
		nc := DefaultConfig
		c = &nc
	}
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()
	// init fonts
	if err := sdlttf.Init(); err != nil {
		return err
	}
	defer sdlttf.Quit()
	// init audio
	if err := sdlmix.Init(0); err != nil {
		return err
	}
	defer sdlmix.Quit()
	sndfmt := uint16(sdlmix.DEFAULT_FORMAT)
	if err := sdlmix.OpenAudio(44100, sndfmt, 2, 1024); err != nil {
		return err
	}
	// init images
	sdlimg.Init(sdlimg.INIT_PNG)
	defer sdlimg.Quit()
	// create and run the game engine
	engine, err := NewEngine(c)
	if err != nil {
		return err
	}
	engine.Run()
	return nil
}
