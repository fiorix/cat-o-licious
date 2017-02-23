// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"errors"
	"sync/atomic"

	"github.com/veandco/go-sdl2/sdl"
	sdlmix "github.com/veandco/go-sdl2/sdl_mixer"
)

// Direction is the player's lateral movement direction.
type Direction int

// Player's lateral movement.
const (
	Center Direction = iota
	Left
	Right
)

// DefaultPlayerSide configures the default facing side of the player.
var DefaultPlayerSide = Right

// DefaultPlayerHitSquare defines the hit square within the
// player image, in percentages.
var DefaultPlayerHitSquare = struct{ X, Y, W, H, FlipOffset float32 }{
	X:          .4,  // left offset
	Y:          .8,  // top offset
	W:          .3,  // pct of the img width
	H:          .15, // pct of the img height
	FlipOffset: .1,  // X offset for the horizontal flip
}

// Player is a game player.
type Player interface {
	// Move moves the player.
	Move(d Direction, steps int32)

	// Hit returns true when the given drop's area intersects
	// with the player's hit area. This is when you make points.
	Hit(d Drop) bool

	// Draw draws the player.
	Draw(viewport *sdl.Rect)
}

// players need at least 3 images, loaded in the following sequence:
const (
	moving  = 0
	winning = 1
	losing  = 2
)

type player struct {
	r    *sdl.Renderer   // main renderer
	imgs []Image         // available images
	img  Image           // current player image
	x    int32           // lateral movement
	y    int32           // vertical position relative to viewport
	d    Direction       // facing side
	hitP sdl.Rect        // hit area of the player
	hitC int             // hit counter to swap image for N frames
	sfx  []*sdlmix.Chunk // available sfx
}

// NewPlayer creates and initializes a new player.
func NewPlayer(r *sdl.Renderer) (Player, error) {
	imgs, err := NewImageSetFromFiles(r, "assets/img/player_frame_")
	if err != nil {
		return nil, err
	}
	if len(imgs) < 3 {
		return nil, errors.New("need at least 3 player frames")
	}
	sfxwin, err := sdlmix.LoadWAV("assets/snd/score_win_1.wav")
	if err != nil {
		return nil, err
	}
	sfxlose, err := sdlmix.LoadWAV("assets/snd/score_lose_1.wav")
	if err != nil {
		return nil, err
	}
	p := &player{
		r:    r,
		imgs: imgs,
		img:  imgs[moving],
		d:    Center,
		sfx:  []*sdlmix.Chunk{nil, sfxwin, sfxlose},
	}
	return p, nil
}

// Move implements the Player interface.
func (p *player) Move(d Direction, steps int32) {
	p.d = d
	switch d {
	case Left:
		atomic.AddInt32(&p.x, -steps)
	case Right:
		atomic.AddInt32(&p.x, steps)
	}
}

// Draw implements the Player interface.
func (p *player) Draw(viewport *sdl.Rect) {
	x := atomic.LoadInt32(&p.x)
	if p.hitC > 0 {
		// reset player image from hit to moving after 10 frames
		if p.hitC == 10 {
			p.hitC = 0
			p.img = p.imgs[moving]
		} else {
			p.hitC++
		}
	}
	w, h := p.img.Size()
	if p.d == Center {
		x = (viewport.W / 2) - (w / 2)
		p.d = DefaultPlayerSide
	} else {
		// fix player left position in case it has moved
		// past the left border
		if x < 0 {
			x = 0
		}
		// fix player right position in case is has moved
		// past the right border
		lim := viewport.W - w
		if x > lim {
			x = lim
		}
	}
	atomic.StoreInt32(&p.x, x)
	p.y = int32(float32(viewport.H)/1.5) - (h / 2)
	r := &sdl.Rect{X: x, Y: p.y, W: w, H: h}
	hs := DefaultPlayerHitSquare
	// save player's current hit square position
	p.hitP = sdl.Rect{
		X: x + int32(float32(w)*hs.X),
		Y: p.y + int32(float32(h)*hs.Y),
		W: int32(float32(w) * hs.W),
		H: int32(float32(h) * hs.H),
	}
	// render player facing default side
	if p.d == DefaultPlayerSide {
		p.r.Copy(p.img.Texture(), nil, r)
		//p.r.DrawRect(&p.hitP) // debug
		return
	}
	// render player facing the opposide side, then adjust
	// the hit square X position
	p.r.CopyEx(p.img.Texture(), nil, r, 0, nil, sdl.FLIP_HORIZONTAL)
	p.hitP.X = p.hitP.X - int32(float32(w)*hs.FlipOffset)
	//p.r.DrawRect(&p.hitP) // debug
}

// Hit implements the Player interface.
func (p *player) Hit(d Drop) bool {
	hit := p.hitP
	area := d.Pos()
	if (area.Y+area.H) < hit.Y || area.Y > hit.Y {
		return false
	}
	if (area.X+area.W) < hit.X || area.X > (hit.X+hit.W) {
		return false
	}
	p.hitC = 1
	if d.Points() > 0 {
		p.img = p.imgs[winning]
		p.sfx[winning].Play(1, 0)
	} else {
		p.img = p.imgs[losing]
		p.sfx[losing].Play(1, 0)
	}
	return true
}
