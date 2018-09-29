package main

import (
	"log"

	"github.com/fiorix/cat-o-licious/wasm/game"
)

func main() {
	engine, err := game.NewEngine("canvas")
	if err != nil {
		log.Fatal(err)
	}

	engine.Run()
}
