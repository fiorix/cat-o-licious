package game

import (
	"errors"
	"sync/atomic"

	"github.com/fiorix/cat-o-licious/wasm/media"
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
	Draw(canvas media.Canvas)
}

// players need at least 3 images, loaded in the following sequence:
const (
	moving  = 0
	winning = 1
	losing  = 2
)

type player struct {
	imgs []media.Image // available images
	img  media.Image   // current player image
	x    int32         // lateral movement
	y    int32         // vertical position relative to viewport
	d    Direction     // facing side
	hitP media.Rect    // hit area of the player
	hitC int           // hit counter to swap image for N frames
}

// NewPlayer creates and initializes a new player.
func NewPlayer() (Player, error) {
	imgs, err := media.NewImageSet("assets/img/player_frame_")
	if err != nil {
		return nil, err
	}
	if len(imgs) < 3 {
		return nil, errors.New("need at least 3 player frames")
	}
	p := &player{
		imgs: imgs,
		img:  imgs[moving],
		d:    Center,
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
func (p *player) Draw(canvas media.Canvas) {
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
	w, h := p.img.W(), p.img.H()
	if p.d == Center {
		x = int32((canvas.ClientW() / 2) - (w / 2))
		p.d = DefaultPlayerSide
	} else {
		// fix player left position in case it has moved
		// past the left border
		if x < 0 {
			x = 0
		}
		// fix player right position in case is has moved
		// past the right border
		lim := int32(canvas.ClientW() - w)
		if x > lim {
			x = lim
		}
	}
	atomic.StoreInt32(&p.x, x)
	p.y = int32(float32(canvas.ClientH())/1.5) - int32(h/2)
	r := &media.Rect{X: int(x), Y: int(p.y), W: w, H: h}
	hs := DefaultPlayerHitSquare
	// save player's current hit square position
	p.hitP = media.Rect{
		X: int(x) + int(float32(w)*hs.X),
		Y: int(p.y) + int(float32(h)*hs.Y),
		W: int(float32(w) * hs.W),
		H: int(float32(h) * hs.H),
	}
	// render player facing default side
	if p.d == DefaultPlayerSide {
		canvas.DrawImage(p.img, r.X, r.Y, r.W, r.H)
		return
	}
	// render player facing the opposide side, then adjust
	// the hit square X position
	canvas.DrawImageFlipHorizontal(p.img, r.X, r.Y, r.W, r.H)
	p.hitP.X = p.hitP.X - int(float32(w)*hs.FlipOffset)
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
		//		p.sfx[winning].Play(1, 0)
	} else {
		p.img = p.imgs[losing]
		//		p.sfx[losing].Play(1, 0)
	}
	return true
}
