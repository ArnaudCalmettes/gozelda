package graphics

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

// An AnimatedSprite is a drawable sprite that gets updated with time.
type AnimatedSprite struct {
	data      *AnimationData
	frameTime time.Duration
	frameNum  int
}

// NewAnimatedSprite creates an animated sprite from the related animation
// name. Multiple independent AnimatedSprites can exist at the same time.
func NewAnimatedSprite(animation string) (*AnimatedSprite, error) {
	anim, err := getAnim(animation)
	if err != nil {
		return nil, err
	}
	return &AnimatedSprite{data: anim}, nil
}

// Update the animation: compute the new frame number.
func (a *AnimatedSprite) Update(dt time.Duration) {
	a.frameTime += dt
	if a.data.FPS > 0 {
		a.frameNum = int(a.frameTime.Seconds()*float64(a.data.FPS)) % len(a.data.Frames)
	}
}

// DrawAt draws the current frame at given (x, y) coordinates on an image.
func (a *AnimatedSprite) DrawAt(img *ebiten.Image, x, y float64) error {
	return a.data.Frames[a.frameNum].DrawAt(img, x, y)
}

// AnimationData is the generic data of an animation
type AnimationData struct {
	Frames []*AnimationFrame
	FPS    int
}

func newAnimationData(d *AnimationDesc) (*AnimationData, error) {
	data := &AnimationData{
		Frames: make([]*AnimationFrame, 0, len(d.Frames)),
		FPS:    d.FPS,
	}
	for _, f := range d.Frames {
		// Get frame from global atlas
		img, err := getFrame(f.Key)
		if err != nil {
			return nil, err
		}
		data.Frames = append(data.Frames,
			&AnimationFrame{
				Image: img,
				FlipH: f.FlipH,
				FlipV: f.FlipV,
			},
		)
	}
	return data, nil
}

// An AnimationFrame is an image with optional transformations.
type AnimationFrame struct {
	Image *ebiten.Image
	FlipH bool
	FlipV bool
}

// DrawAt draws the frame at given coordinates on an image,
// flipping it horizontally and/or vertically if needed.
func (af *AnimationFrame) DrawAt(img *ebiten.Image, x, y float64) error {
	op := &ebiten.DrawImageOptions{}
	w, h := af.Image.Size()
	if af.FlipH {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(w), 0)
	}
	if af.FlipV {
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(h))
	}
	op.GeoM.Translate(float64(x), float64(y))
	return img.DrawImage(af.Image, op)
}
