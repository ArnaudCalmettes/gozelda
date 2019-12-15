package graphics

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	frameAtlas = make(map[string]*ebiten.Image)
	animAtlas  = make(map[string]*AnimationData)
)

func getFrame(key string) (*ebiten.Image, error) {
	frame, ok := frameAtlas[key]
	if !ok {
		return nil, fmt.Errorf("unknown frame '%s'", key)
	}
	return frame, nil
}

func getAnim(key string) (*AnimationData, error) {
	anim, ok := animAtlas[key]
	if !ok {
		return nil, fmt.Errorf("unknown animation '%s'", key)
	}
	return anim, nil
}

// Load graphics assets from given manifest file
func Load(path string) error {
	m, err := LoadManifest(path)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(path)
	for _, c := range m.Collections {
		log.Printf("Loading collection '%s'", c.Name)
		if err = loadSprites(filepath.Join(baseDir, c.SpriteSheet)); err != nil {
			return err
		}

		for _, a := range c.Animations {
			path = filepath.Join(baseDir, a)
			if err = loadAnimations(filepath.Join(baseDir, a)); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadSprites(path string) error {
	log.Printf("Loading sprites from %s", path)
	sheet, err := LoadSpriteSheetManifest(path)
	if err != nil {
		return err
	}

	log.Printf("Loading spritesheet '%s' (%d frames)", sheet.Meta.Image, len(sheet.Frames))

	imgPath := filepath.Join(filepath.Dir(sheet.Path), sheet.Meta.Image)
	img, _, err := ebitenutil.NewImageFromFile(imgPath, ebiten.FilterDefault)
	if err != nil {
		return err
	}

	for _, frame := range sheet.Frames {
		if _, exists := frameAtlas[frame.Name]; exists {
			return fmt.Errorf("frame '%s' already exists", frame.Name)
		}
		x, y, w, h := frame.ROI.XYWH()
		sub, ok := img.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
		if !ok {
			// Unreachable (according to the doc)
			panic("SubImage() didn't return an *ebiten.Image")
		}
		frameAtlas[frame.Name] = sub
	}

	return nil
}

func loadAnimations(path string) error {
	log.Printf("Loading animations from %s", path)
	m, err := LoadAnimationManifest(path)
	if err != nil {
		return err
	}
	for _, a := range m.Anims {
		if _, exists := animAtlas[a.Name]; exists {
			return fmt.Errorf("animation '%s' already exists", a.Name)
		}
		animData, err := newAnimationData(a)
		if err != nil {
			return fmt.Errorf("while loading animation '%s': '%s'", a.Name, err)
		}
		animAtlas[a.Name] = animData
	}
	return nil
}
