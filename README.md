# GoZelda

An implementation of a Game Boy Zelda game using Golang and ebiten.

This project is my way of teaching myself game development. Its purpose
is to achieve a good enough copy of a simple (yet legendary and really fun)
game in order to learn good practices in 2D game dev.

I'm organizing my work around "baby step prototypes", each step consists in a
very simple demo ("prototype") showcasing some mechanism of the game.

## Current step: animated sprites using an "animation atlas"

This is the very first prototype. It can be built & launched using:

```bash
$ cd $GOROOT/github.com/ArnaudCalmettes/gozelda
$ go build -o anim ./prototypes/animation
$ ./anim
```

The purpose of this prototype was to develop:

* A graphics module that loads spritesheets and animations from JSON manifests (see
  `assets/sprites/manifest.json`)
* A way to create independently animated sprites using simply the key to an animation (see
  `prototypes/animation/main.go`)

## Next step: an animated character that responds to user input

Even though the game will probably adopt the entity/component pattern in the future,
the next step is to have Link moving around with his shield, colliding to the screen borders,
and respond to the shield button.

I'm keeping sword-related actions for a later step, because on GameBoy those seem to coordinate
two independent objects & animations (link and the sword).

