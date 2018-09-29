package game

import (
	"math/rand"
	"time"

	"github.com/fiorix/cat-o-licious/wasm/media"
)

// Rain ...
type Rain interface {
	SetRate(n int)
	Drops() []Drop
	Draw(now time.Time, canvas media.Canvas)
}

// Drop ...
type Drop interface {
	Pos() media.Rect
	Points() int64
}

// NewRain ...
func NewRain() (Rain, error) {
	var rd []*raindrop
	imgs, err := media.NewImageSet("assets/img/drop_good_")
	if err != nil {
		return nil, err
	}
	for i, img := range imgs {
		rd = append(rd, &raindrop{
			img:    img,
			points: int64(i+1) * 5,
		})
	}
	imgs, err = media.NewImageSet("assets/img/drop_bad_")
	if err != nil {
		return nil, err
	}
	for i, img := range imgs {
		rd = append(rd, &raindrop{
			img:    img,
			points: -(int64(i+1) * 10),
		})
	}
	ra := &rain{
		lastdrop:  time.Now(),
		available: rd,
		delay:     time.Second,
	}
	return ra, nil
}

type rain struct {
	canvas    media.Canvas
	lastdrop  time.Time
	available []*raindrop
	drops     []*drop
	delay     time.Duration
}

type raindrop struct {
	img    media.Image
	points int64
}

func (r *rain) SetRate(n int) {
	const min = 100 * time.Millisecond
	r.delay = time.Duration(1000-(n-1)*100) * time.Millisecond
	if r.delay < min {
		r.delay = min
	}
}

func (r *rain) Draw(now time.Time, canvas media.Canvas) {
	if now.Sub(r.lastdrop) >= r.delay {
		r.newDrop(canvas)
		r.lastdrop = now
	}
	r.drawAndDrain(canvas)
}

func (r *rain) newDrop(canvas media.Canvas) {
	border := int(float64(canvas.ClientW()) * .05)
	idx := rand.Intn(len(r.available))
	rd := r.available[idx]
	w, h := rd.img.W(), rd.img.H()
	lim := canvas.ClientW() - w - (border * 2)
	pos := media.Rect{
		X: border + rand.Intn(lim),
		Y: -h,
		W: w,
		H: h,
	}
	d := &drop{
		src:   rd,
		pos:   pos,
		speed: 5 + rand.Intn(10),
	}
	r.drops = append(r.drops, d)
}

func (r *rain) drawAndDrain(canvas media.Canvas) {
	rmset := make(map[int]struct{})
	for i, d := range r.drops {
		if d.pos.Y > canvas.ClientH() {
			rmset[i] = struct{}{}
			continue
		}
		d.Draw(canvas)
	}
	if len(rmset) == 0 {
		return
	}
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

func (r *rain) Drops() []Drop {
	drops := make([]Drop, len(r.drops))
	for i, d := range r.drops {
		drops[i] = d
	}
	return drops
}

type drop struct {
	src   *raindrop
	pos   media.Rect
	speed int
}

func (d *drop) Draw(canvas media.Canvas) {
	d.pos.Y = d.pos.Y + d.speed
	w, h := d.src.img.W(), d.src.img.H()
	canvas.DrawImage(d.src.img, d.pos.X, d.pos.Y, w, h)
}

func (d *drop) Pos() media.Rect {
	return d.pos
}

func (d *drop) Points() int64 {
	return d.src.points
}
