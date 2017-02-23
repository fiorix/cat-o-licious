// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Rain makes objects (raindrops) fall from the top of the viewport
// all the way to the bottom.
type Rain interface {
	// SetRate sets the rate of drops per second.
	SetRate(n int)

	// Drops returns all raindrops from the rain.
	Drops() []Drop

	// Draw draws the rain.
	Draw(now time.Time, viewport *sdl.Rect)
}

// Drop is a single drop of rain.
type Drop interface {
	// Pos returns the position of the drop in the viewport.
	Pos() sdl.Rect

	// Points returns delta points for the drop. Good drops
	// return positive numbers while bad drops return negative.
	Points() int64
}

// rain implements the Rain interface.
type rain struct {
	r         *sdl.Renderer
	lastdrop  time.Time
	available []*raindrop
	drops     []*drop
	delay     time.Duration
}

// raindrop is a storage for drop images and the points associated
// to them: good drops have positive points, bad drops have negative.
type raindrop struct {
	r      *sdl.Renderer
	img    Image
	points int64
}

// NewRain creates and initializes a Rain object.
func NewRain(r *sdl.Renderer) (Rain, error) {
	var rd []*raindrop
	imgs, err := NewImageSetFromFiles(r, "assets/img/drop_good_")
	if err != nil {
		return nil, err
	}
	for i, img := range imgs {
		rd = append(rd, &raindrop{
			// you get 5*frameidx(e.g. bacon) points per hit
			r: r, img: img, points: int64(i+1) * 5,
		})
	}
	imgs, err = NewImageSetFromFiles(r, "assets/img/drop_bad_")
	if err != nil {
		return nil, err
	}
	for i, img := range imgs {
		rd = append(rd, &raindrop{
			// lose 20*frameidx(e.g. pineapple) points per hit
			r: r, img: img, points: -(int64(i+1) * 20),
		})
	}
	ra := &rain{
		r:         r,
		lastdrop:  time.Now(), // also initial delay
		available: rd,
		delay:     time.Second, // start at 1 new drop per second
	}
	return ra, nil
}

// SetRate implements the Rain interface.
func (r *rain) SetRate(n int) {
	const min = 100 * time.Millisecond
	r.delay = time.Duration(1000-(n-1)*100) * time.Millisecond
	if r.delay < min {
		r.delay = min
	}
}

// Draw implements the Rain interface.
func (r *rain) Draw(now time.Time, viewport *sdl.Rect) {
	if now.Sub(r.lastdrop) >= r.delay {
		r.newDrop(viewport)
		r.lastdrop = now
	}
	r.drawAndDrain(viewport)
}

// newDrop adds a new drop to the rain.
func (r *rain) newDrop(viewport *sdl.Rect) {
	// border is the pct of the viewport that should not
	// have rain drops
	border := int32(float32(viewport.W) * .05)
	// roll the dice to pick a preloaded raindrop
	idx := rand.Intn(len(r.available))
	rd := r.available[idx]
	// roll the dice to place the raindrop horizontally
	w, h := rd.img.Size()
	lim := viewport.W - w - (border * 2)
	pos := &sdl.Rect{
		X: border + rand.Int31n(lim),
		Y: -h,
		W: w,
		H: h,
	}
	// roll the dice for the drop speed
	d := &drop{
		src:   rd,
		pos:   pos,
		speed: 5 + rand.Int31n(10),
	}
	r.drops = append(r.drops, d)
}

// drawAndDrain draws drops that are within the viewport and drains
// drops that placed past the height of the viewport.
func (r *rain) drawAndDrain(viewport *sdl.Rect) {
	rmset := make(map[int]struct{})
	for i, d := range r.drops {
		if d.pos.Y > viewport.H {
			rmset[i] = struct{}{}
			continue
		}
		d.Draw(viewport)
	}
	if len(rmset) == 0 {
		return
	}
	// drain
	drops := make([]*drop, len(r.drops)-len(rmset))
	j := 0
	for i, d := range r.drops {
		if _, skip := rmset[i]; skip {
			continue
		}
		drops[j] = d
		j++
	}
	r.drops = drops
}

// Drops implement the Rain interface.
func (r *rain) Drops() []Drop {
	drops := make([]Drop, len(r.drops))
	for i, d := range r.drops {
		drops[i] = d
	}
	return drops
}

// drop is a single drop of rain, that falls from top to bottom.
type drop struct {
	src   *raindrop
	pos   *sdl.Rect
	speed int32
}

// Draw draws the drop incrementing its Y position at a given speed.
func (d *drop) Draw(viewport *sdl.Rect) {
	d.pos.Y = d.pos.Y + d.speed
	d.src.r.Copy(d.src.img.Texture(), nil, d.pos)
}

// Pos implements the Drop interface.
func (d *drop) Pos() sdl.Rect {
	return *d.pos
}

// Points implements the Drop interface.
func (d *drop) Points() int64 {
	return d.src.points
}
