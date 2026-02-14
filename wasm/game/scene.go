package game

import (
	"time"

	"github.com/fiorix/cat-o-licious/wasm/media"
)

// Scene ...
type Scene interface {
	Player() Player
	Draw(canvas media.Canvas)
	EnableAudio()
}

func (s *scene) EnableAudio() {
	s.audioEnabled = true
	s.player.EnableAudio()
}

func (s *scene) handleMusic() {
	if !s.audioEnabled {
		return
	}
	if s.musicStarted {
		if s.mus.Playing() {
			return
		}
		s.musicStarted = false
	}
	s.mus.PlayLoop()
	s.musicStarted = true
}

func (s *scene) drawAudioPrompt(canvas media.Canvas) {
	canvas.SetFont("32px Score", "#ffd700")
	canvas.DrawText("Click or press to enable sound", 40, 80)
}

type scene struct {
	lastUpdate   time.Time
	bg           media.Image
	rain         Rain
	player       Player
	score        Scoreboard
	mus          media.Audio
	audioEnabled bool
	musicStarted bool
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
	mus, err := media.NewAudio("assets/snd/music_1.wav")
	if err != nil {
		return nil, err
	}
	s := &scene{
		bg:     bg,
		rain:   rain,
		player: player,
		score:  NewScoreboard(),
		mus:    mus,
	}
	return s, nil
}

func (s *scene) Player() Player {
	return s.player
}

func (s *scene) Draw(canvas media.Canvas) {
	r := media.Rect{
		X: 0,
		Y: 0,
		W: canvas.ClientW(),
		H: canvas.ClientH(),
	}
	canvas.ClearRect(r)
	canvas.DrawImage(s.bg, r)
	s.handleMusic()
	if !s.audioEnabled {
		s.drawAudioPrompt(canvas)
	}
	s.player.Draw(canvas)
	s.rain.Draw(canvas)
	canvas.SetFont("80px Score", "red")
	s.score.Draw(canvas)
	hit := false
	for _, drop := range s.rain.Drops() {
		if s.player.Hit(drop) {
			s.score.Add(drop.Points())
			hit = true
		}
	}
	t := time.Now()
	if !hit || t.Sub(s.lastUpdate) < time.Second {
		return
	}
	s.lastUpdate = t
	p := s.score.Points()
	rate := (int(p) / 1000) + 1
	s.rain.SetRate(rate)
}
