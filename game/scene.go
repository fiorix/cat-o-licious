// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
	sdlmix "github.com/veandco/go-sdl2/mix"
)

// Scene is the game scene.
// It contains the background image, scoreboard, player, rain, and sound.
type Scene interface {
	// Player returns the game player so callers can control
	// its lateral movement.
	Player() Player

	// Draw draws the scene.
	Draw(now time.Time, viewport *sdl.Rect)
}

type scene struct {
	r          *sdl.Renderer
	lastupdate time.Time

	bg     Image
	score  Scoreboard
	rain   Rain
	player Player
	mus    *sdlmix.Music
}

// NewScene creates and initializes the game scene.
func NewScene(r *sdl.Renderer) (Scene, error) {
	bg, err := NewImageFromFile(r, "assets/img/background.png")
	if err != nil {
		return nil, err
	}
	score, err := NewScoreboard(r)
	if err != nil {
		return nil, err
	}
	rain, err := NewRain(r)
	if err != nil {
		return nil, err
	}
	player, err := NewPlayer(r)
	if err != nil {
		return nil, err
	}
	mus, err := sdlmix.LoadMUS("assets/snd/music_1.wav")
	if err != nil {
		return nil, err
	}
	s := &scene{
		r:      r,
		bg:     bg,
		score:  score,
		rain:   rain,
		player: player,
		mus:    mus,
	}
	return s, nil
}

// Draw implements the Scene interface.
func (s *scene) Draw(now time.Time, viewport *sdl.Rect) {
	if !sdlmix.PlayingMusic() {
		s.mus.Play(0)
	}
	s.r.Copy(s.bg.Texture(), nil, nil)
	s.score.Draw(viewport)
	s.player.Draw(viewport)
	s.rain.Draw(now, viewport)
	hit := false
	for _, drop := range s.rain.Drops() {
		if s.player.Hit(drop) {
			s.score.Add(drop.Points())
			hit = true
		}
	}
	// update rain delay at most 1/s, only on player hits.
	if !hit || now.Sub(s.lastupdate) < time.Second {
		return
	}
	s.lastupdate = now
	// TODO: do this less frequent, not on all hits
	// TODO: switch music every 1000 points?
	p := s.score.Points()
	rate := (int(p) / 1000) + 1
	s.rain.SetRate(rate)
}

// Player implements the Scene interface.
func (s *scene) Player() Player {
	return s.player
}
