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
	perRow      = 8
)

var anims = []string{
	"link_walk_default_left",
	"link_walk_default_up",
	"link_walk_default_down",
	"link_walk_default_right",
	"link_walk_shield_off_left",
	"link_walk_shield_off_up",
	"link_walk_shield_off_down",
	"link_walk_shield_off_right",
	"link_walk_shield_on_left",
	"link_walk_shield_on_up",
	"link_walk_shield_on_down",
	"link_walk_shield_on_right",
	"link_walk_mirror_off_left",
	"link_walk_mirror_off_up",
	"link_walk_mirror_off_down",
	"link_walk_mirror_off_right",
	"link_walk_mirror_on_left",
	"link_walk_mirror_on_up",
	"link_walk_mirror_on_down",
	"link_walk_mirror_on_right",
	"link_swim_left",
	"link_swim_up",
	"link_swim_down",
	"link_swim_right",
	"link_jump_left",
	"link_jump_down",
	"link_jump_up",
	"link_jump_right",
	"link_swim_full_left",
	"link_dive",
	"link_fall",
	"link_swim_full_right",
	"link_push_left",
	"link_push_up",
	"link_push_down",
	"link_push_right",
	"link_pull_left",
	"link_pull_up",
	"link_pull_down",
	"link_pull_right",
	"link_spin",
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
			log.Println("Dropping 1 frame")
			return nil
		}

		screen.Fill(color.RGBA{0x31, 0x8b, 0x6a, 0xff})

		for i, sprite := range sprites {
			row := i / perRow
			col := i % perRow
			sprite.DrawAt(screen, float64(16+16*col), float64(24+16*row))
		}

		return nil
	}

	if err = ebiten.Run(update, screenW, screenH, screenScale, title); err != nil {
		log.Fatal(err)
	}
}
