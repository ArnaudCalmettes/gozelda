package graphics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

///////////////////////// Global graphics assets manifest

// A Manifest is a json file describing the whole game's graphic assets
type Manifest struct {
	Collections []*AssetCollection `json:"collections"`
}

// An AssetCollection gathers related spritesheets and animations
type AssetCollection struct {
	Name        string   `json:"name"`
	SpriteSheet string   `json:"spritesheet"`
	Animations  []string `json:"animations"`
}

// Check sanity of the Manifest
func (m *Manifest) Check() error {
	for i, c := range m.Collections {
		if c.Name == "" {
			return fmt.Errorf("Collection #%d doesn't have a name", i)
		}
		if c.SpriteSheet == "" {
			return fmt.Errorf(
				"collection #%d (%s) doesn't have an associated spritesheet",
				i, c.Name,
			)
		}
		if len(c.Animations) == 0 {
			return fmt.Errorf(
				"collection #%d (%s) doesn't have any associated animations",
				i, c.Name,
			)
		}
	}
	return nil
}

// LoadManifest loads a manifest file
func LoadManifest(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	m := &Manifest{}
	if err = d.Decode(m); err != nil {
		return nil, err
	}
	return m, m.Check()
}

//////////////////////// Spritesheet manifest

// A SpriteSheetManifest is a structure tied to an image packing multiple sprites.
type SpriteSheetManifest struct {
	Name string // Name of the spritesheet
	Path string // Full path of the spritesheet
	// Individual frames
	Frames []*SSFrame `json:"frames"`

	// Metadata associated with the sprite sheet
	Meta SSMeta `json:"meta"`
}

// A SSFrame is a single sprite packed in a spritesheet
type SSFrame struct {
	// Name of the sprite
	Name string `json:"filename"`

	// Area within the spritesheet that contains the sprite
	ROI Rect `json:"frame"`

	// Pivot point of the sprite (anchor to align animation frames)
	Pivot PointF64 `json:"pivot"`
}

// A Rect describes a rectangular area within an image.
type Rect struct {
	X int `json:"x"` // X position
	Y int `json:"y"` // Y position
	W int `json:"w"` // Width
	H int `json:"h"` // Height
}

// XYWH unpacks a Rect into its X, Y, W, H int components.
func (r Rect) XYWH() (int, int, int, int) {
	return r.X, r.Y, r.W, r.H
}

// A PointF64 is a point whose coordinates are float64 values.
type PointF64 struct {
	X float64 `json:"x"` // X position
	Y float64 `json:"y"` // Y position
}

func (p PointF64) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// SSMeta contains various informations about the spritesheet
type SSMeta struct {
	// The image file this spritesheed is tied to.
	// Note: the file has to be located in the same folder as the spritesheet.
	Image string `json:"image"`

	// Image size
	Size Size `json:"size"`
}

// A Size is the combination of width and height.
type Size struct {
	W int `json:"w"` // Width
	H int `json:"h"` // Height
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.W, s.H)
}

// Check sanity of the SpriteSheetManifest
func (s *SpriteSheetManifest) Check() error {
	if s.Meta.Image == "" {
		return fmt.Errorf("spritesheet '%s' isn't associated to an image file", s.Name)
	}
	if s.Meta.Size.W <= 0 || s.Meta.Size.H <= 0 {
		return fmt.Errorf("spritesheet '%s' has an invalid image size: %s", s.Name, s.Meta.Size)
	}
	for i, f := range s.Frames {
		if f.Name == "" {
			return fmt.Errorf("frame #%d doesn't have a name", i)
		}
		r := &f.ROI
		if r.X+r.W > s.Meta.Size.W || r.Y+r.H > s.Meta.Size.H {
			return fmt.Errorf("frame #%d (%s) is out of image boundaries", i, f.Name)
		}
		if f.Pivot.X < 0 || f.Pivot.X > 1 || f.Pivot.Y < 0 || f.Pivot.Y > 1 {
			return fmt.Errorf("frame #%d (%s) has an invalid pivot point: %s", i, f.Name, f.Pivot)
		}
	}
	return nil
}

// SetPath sets both the path and the name of the spritesheet.
func (s *SpriteSheetManifest) SetPath(path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	s.Path = p
	s.Name = filepath.Base(p)
	return err
}

// LoadSpriteSheetManifest loads a json spritesheet file.
func LoadSpriteSheetManifest(path string) (*SpriteSheetManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	s := &SpriteSheetManifest{}
	if err = d.Decode(s); err != nil {
		return nil, err
	}
	if err = s.SetPath(path); err != nil {
		return nil, err
	}
	return s, s.Check()
}

//////////////////////// Animation manifest

// An AnimationManifest describes all animations for given sprite sheet
type AnimationManifest struct {
	// Spritesheet where the animation frames are located
	SpriteSheet string `json:"spritesheet"`

	// Animation descriptions
	Anims []*AnimationDesc `json:"animations"`
}

// An AnimationDesc describes an animation (name, speed and frames)
type AnimationDesc struct {
	// Unique name of the animation
	Name string `json:"name"`
	// Animation speed in frames per second
	FPS int `json:"fps"`
	// Specification for the animation frames
	Frames []*FrameSpec `json:"frames"`
}

// A FrameSpec describes a unique animation frame
type FrameSpec struct {
	// Key of the frame in the global frameAtlas
	Key string `json:"key"`

	// Flip the frame horizontally
	FlipH bool `json:"flipH"`

	// Flip the frame vertically
	FlipV bool `json:"flipV"`
}

// Check sanity of the AnimationManifest
func (a *AnimationManifest) Check() error {
	for i, anim := range a.Anims {
		if anim.Name == "" {
			return fmt.Errorf("anim #%d doesn't have a name", i)
		}
		if len(anim.Frames) == 0 {
			return fmt.Errorf("anim #%d (%s) has no frames", i, anim.Name)
		}
		if anim.FPS == 0 && len(anim.Frames) > 1 {
			return fmt.Errorf("anim #%d (%s) should has a null FPS", i, anim.Name)
		}
		for j, frame := range anim.Frames {
			if frame.Key == "" {
				return fmt.Errorf("frame #%d from anim #%d (%s) has no key", j, i, anim.Name)
			}
		}
	}
	return nil
}

// LoadAnimationManifest loads a json manifest file that describes animations
func LoadAnimationManifest(path string) (*AnimationManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	a := &AnimationManifest{}
	if err = d.Decode(a); err != nil {
		return nil, err
	}
	return a, a.Check()
}
