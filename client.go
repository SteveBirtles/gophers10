package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
	"time"
)

const ScreenWidth = 1024
const ScreenHeight = 768

var (
	frames = 0
	second = time.Tick(time.Second)
	tick   = time.Tick(time.Second / 16)
	title  = "Game startClient with an API"
	win    *pixelgl.Window

	sprite [4][4][4]*pixel.Sprite
	batch  *pixel.Batch
)

func loadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func startClient() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, ScreenWidth, ScreenHeight),
		VSync:  true,
	}

	win, initError = pixelgl.NewWindow(cfg)
	if initError != nil {
		panic(initError)
	}

	spriteImage, initError := loadImageFile("bomberman.png")
	if initError != nil {
		panic(initError)
	}

	spritePic := pixel.PictureDataFromImage(spriteImage)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 3; k++ {
				sprite[i][j][k] = pixel.NewSprite(spritePic,
					pixel.R(float64(k+j*3)*20, float64(120-30*(i+1)), float64(k+j*3+1)*20, float64(120-30*i)))
			}
			sprite[i][j][3] = sprite[i][j][1]
		}
	}

	batch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)

	mainLoop()

}

func mainLoop() {

	for !win.Closed() {

		batch.Clear()

		select {
		case <-tick:
			for i := range player {
				player[i].updatePosition()
			}
		default:
		}

		for i := range player {
			player[i].draw()
		}

		win.Clear(colornames.Black)

		win.SetComposeMethod(pixel.ComposeOver)
		batch.Draw(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
			frames = 0
		default:
		}

	}

}
