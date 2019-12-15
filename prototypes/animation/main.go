package main

import (
	"image/color"
	"log"
	"time"

	"github.com/ArnaudCalmettes/gozelda/graphics"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenW     = 160
	screenH     = 160
	screenScale = 4
	title       = "Prototype: animations"
)

var anims = [4]string{
	"link_walk_default_left",
	"link_walk_default_up",
	"link_walk_default_down",
	"link_walk_default_right",
}

func main() {
	err := graphics.Load("assets/sprites/manifest.json")
	if err != nil {
		log.Fatal(err)
	}

	sprites := make([]*graphics.AnimatedSprite, len(anims))
	for i, anim := range anims {
		sprites[i], err = graphics.NewAnimatedSprite(anim)
		if err != nil {
			log.Fatal(err)
		}
	}

	now := time.Now()
	update := func(screen *ebiten.Image) error {
		// Compute dt & update the clock
		dt := time.Since(now)
		now = now.Add(dt)

		// Update the sprite animations
		for _, sprite := range sprites {
			sprite.Update(dt)
		}

		if ebiten.IsDrawingSkipped() {
			return nil
		}

		screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

		for i, sprite := range sprites {
			sprite.DrawAt(screen, float64(48+16*i), 72)
		}
		return nil
	}

	if err = ebiten.Run(update, screenW, screenH, screenScale, title); err != nil {
		log.Fatal(err)
	}
}
