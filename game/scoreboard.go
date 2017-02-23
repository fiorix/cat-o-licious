// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/veandco/go-sdl2/sdl"
	sdlttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// Scoreboard tracks the player's score and draws the scoreboard.
type Scoreboard interface {
	// Add adds the given delta to the player's score.
	Add(delta int64)

	// Points returns the current player's points.
	Points() int64

	// Draw draws the scoreboard.
	Draw(viewport *sdl.Rect)
}

type scoreboard struct {
	r      *sdl.Renderer
	f      *sdlttf.Font
	points int64
}

// NewScoreboard creates and initializes a new scoreboard.
func NewScoreboard(r *sdl.Renderer) (Scoreboard, error) {
	f, err := sdlttf.OpenFont("assets/fonts/score.ttf", 50)
	if err != nil {
		return nil, err
	}
	return &scoreboard{r: r, f: f}, nil
}

// Add implements the Scoreboard interface.
func (sb *scoreboard) Add(delta int64) {
	atomic.AddInt64(&sb.points, delta)
}

// Points implements the Scoreboard interface.
func (sb *scoreboard) Points() int64 {
	return atomic.LoadInt64(&sb.points)
}

// Draw implements the Scoreboard interface.
func (sb *scoreboard) Draw(viewport *sdl.Rect) {
	p := atomic.LoadInt64(&sb.points)
	text := fmt.Sprintf("%d", p)
	s, err := sb.f.RenderUTF8_Solid(text, sdl.Color{R: 255})
	if err != nil {
		log.Println("failed to create font surface:", err)
		return
	}
	defer s.Free()
	var clip sdl.Rect
	s.GetClipRect(&clip)
	t, err := sb.r.CreateTextureFromSurface(s)
	if err != nil {
		log.Println("failed to create font texture:", err)
		return
	}
	woff := float32(viewport.W) * .1
	hoff := float32(viewport.H) * .1
	r := &sdl.Rect{
		X: viewport.W - clip.W - int32(woff),
		Y: int32(hoff),
		W: clip.W,
		H: clip.H,
	}
	sb.r.Copy(t, nil, r)
}
