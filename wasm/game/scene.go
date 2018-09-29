package game

import (
	"time"

	"github.com/fiorix/cat-o-licious/wasm/media"
)

// Scene ...
type Scene interface {
	Player() Player
	Draw(now time.Time, canvas media.Canvas)
}

type scene struct {
	lastUpdate time.Time
	bg         media.Image
	rain       Rain
	player     Player
	score      Scoreboard
}

// NewScene ...
func NewScene() (Scene, error) {
	bg, err := media.NewImage("assets/img/background.png", "background")
	if err != nil {
		return nil, err
	}
	rain, err := NewRain()
	if err != nil {
		return nil, err
	}
	player, err := NewPlayer()
	if err != nil {
		return nil, err
	}
	s := &scene{
		bg:     bg,
		rain:   rain,
		player: player,
		score:  NewScoreboard(),
	}
	return s, nil
}

func (s *scene) Player() Player {
	return s.player
}

func (s *scene) Draw(now time.Time, canvas media.Canvas) {
	canvas.ClearRect(media.Rect{
		X: 0,
		Y: 0,
		W: canvas.ClientW(),
		H: canvas.ClientH(),
	})
	canvas.DrawImage(s.bg, 0, 0, canvas.ClientW(), canvas.ClientH())
	s.player.Draw(canvas)
	s.rain.Draw(now, canvas)
	s.score.Draw(canvas)
	hit := false
	for _, drop := range s.rain.Drops() {
		if s.player.Hit(drop) {
			s.score.Add(drop.Points())
			hit = true
		}
	}
	if !hit || now.Sub(s.lastUpdate) < time.Second {
		return
	}
	s.lastUpdate = now
	p := s.score.Points()
	rate := (int(p) / 1000) + 1
	s.rain.SetRate(rate)
}
