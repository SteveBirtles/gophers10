package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"net/http"
	"os"
	"strings"
	"time"
)

const screenWidth = 1024
const screenHeight = 768

var (
	frames            = 0
	second            = time.Tick(time.Second)
	quarter           = time.Tick(time.Second / 10)
	windowTitlePrefix = "Game client with an API"
	win               *pixelgl.Window
	bombermanSprite   [4][4][4]*pixel.Sprite
	bombermanFrame    = 0
	spriteBatch       *pixel.Batch
	x                 = float64(screenWidth / 2)
	y                 = float64(screenHeight / 2)
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func moveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathBits := strings.Split(r.URL.Path, "/")
	direction := pathBits[len(pathBits)-1]
	fmt.Println("/direction/", direction)

	switch direction {
	case "up":
		y += 25
	case "down":
		y -= 25
	case "left":
		x -= 25
	case "right":
		x += 25
	}

}

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

func client() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
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
				bombermanSprite[i][j][k] = pixel.NewSprite(spritePic, pixel.R(float64(k+j*4)*20, float64(120-30*(i+1)), float64(k+j*4+1)*20, float64(120-30*i)))
			}
			bombermanSprite[i][j][3] = bombermanSprite[i][j][1]
		}
	}

	spriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)

	for !win.Closed() {

		spriteBatch.Clear()

		matrix := pixel.IM.Rotated(pixel.ZV, 0).Scaled(pixel.ZV, 4).Moved(pixel.Vec{X: x, Y: y})

		select {
		case <-quarter:
			bombermanFrame--
			if bombermanFrame < 0 {
				bombermanFrame = 3
			}
		default:
		}

		bombermanSprite[0][0][bombermanFrame].Draw(spriteBatch, matrix)

		win.Clear(colornames.Black)

		win.SetComposeMethod(pixel.ComposeOver)
		spriteBatch.Draw(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}

	}

}

func server() {

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/direction/", moveHandler)

	err := http.ListenAndServe(":8081", http.DefaultServeMux)
	if err != nil {
		fmt.Println("Error:", err)
	}

}

func main() {

	go server()

	pixelgl.Run(client)

}
