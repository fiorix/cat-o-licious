package game

import (
	"fmt"
	"sync/atomic"

	"github.com/fiorix/cat-o-licious/wasm/media"
)

// Scoreboard tracks the player's score and draws the scoreboard.
type Scoreboard interface {
	// Add adds the given delta to the player's score.
	Add(delta int64)

	// Points returns the current player's points.
	Points() int64

	// Draw draws the scoreboard.
	Draw(canvas media.Canvas)
}

type scoreboard struct {
	points int64
}

// NewScoreboard creates and initializes a new scoreboard.
func NewScoreboard() Scoreboard {
	return &scoreboard{}
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
func (sb *scoreboard) Draw(canvas media.Canvas) {
	p := atomic.LoadInt64(&sb.points)
	text := fmt.Sprintf("%d", p)
	canvas.DrawText(text, 100, 100)
}
