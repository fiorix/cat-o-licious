// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/fiorix/cat-o-licious/game"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conf := game.DefaultConfig
	flag.IntVar(&conf.FPS, "fps", conf.FPS, "game fps")
	flag.IntVar(&conf.Width, "width", conf.Width, "game width")
	flag.IntVar(&conf.Height, "height", conf.Height, "game height")
	flag.IntVar(&conf.PlayerSpeed, "speed", conf.PlayerSpeed, "player speed")
	flag.Parse()
	err := game.Run(&conf)
	if err != nil {
		log.Fatal(err)
	}
}
