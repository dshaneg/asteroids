package assets

import (
	"bytes"
	"embed"
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *
var assets embed.FS

var PlayerSprite = mustLoadImage("images/playerShip1_blue.png")
var PlumeSprite = mustLoadImage("images/fire10.png")
var LaserSprite = mustLoadImage("images/laserRed07.png")
var ScoreFont = mustLoadFont("font.ttf")

var AsteroidBigBrownSprites = mustLoadImages("images/meteors/*Brown_big*.png")
var AsteroidMediumBrownSprites = mustLoadImages("images/meteors/*Brown_med*.png")
var AsteroidSmallBrownSprites = mustLoadImages("images/meteors/*Brown_small*.png")
var AsteroidTinyBrownSprites = mustLoadImages("images/meteors/*Brown_tiny*.png")

var AsteroidBigGraySprites = mustLoadImages("images/meteors/*Grey_big*.png")
var AsteroidMediumGraySprites = mustLoadImages("images/meteors/*Grey_med*.png")
var AsteroidSmallGraySprites = mustLoadImages("images/meteors/*Grey_small*.png")
var AsteroidTinyGraySprites = mustLoadImages("images/meteors/*Grey_tiny*.png")

var BoomSprite = mustLoadImage("images/boom.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string) *text.GoTextFaceSource {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	ff, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		panic(err)
	}

	return ff
}
