# cat-o-licious

This is a simple cat game written in Go and SDL, inspired by [flappy](https://github.com/campoy/flappy/). My kids were too excited to see flappy's source code and wanted to design and code their own game, so that's what we did together in a rainy Sunday.

### Install

It's been only tested on MacOS, and requires the [SDL](https://www.libsdl.org) library and Go bindings:

```
brew install sdl2 sdl2_image sdl2_mixer sdl2_ttf
go get github.com/fiorix/cat-o-licious
```

The assets directory must be relative to the path of the binary. Assets include fonts, images, and sounds used by the game. The font was copied from flappy, images randomly downloaded from the Internet (by them), and the game soundtrack is my daughter's composition in Garage Band. Go figure.

Run:

```
cd $GOPATH/src/github.com/fiorix/cat-o-licious
./cat-o-licious
```

There's a minimal set of command line flags for things like screen resolution, player speed, and FPS.

![cat-o-licious](assets/screenshot.png)

### Keys

Arrows left and right, as well as A and D for lateral movement.
F for full screen and Q to quit.

### Playing

You're the cat, and food falls from the top of the screen. The more good stuff you lick the more points you make. The more points you make the more food drops, and it gets really hard to get out of the way of the broccoli, tomatos and pineapples.

My kids love veggies btw, but they say that cats don't.
