package main

import (
	"net/http"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"os"
	"image"
	"time"
	"golang.org/x/image/colornames"
	_ "image/png"
	"strings"
)

const screenWidth = 1024
const screenHeight = 768

var (
	frames            = 0
	second            = time.Tick(time.Second)
	windowTitlePrefix	= "Game client with an API"
	win               *pixelgl.Window
	sprite            *pixel.Sprite
	spriteBatch       *pixel.Batch
	x = float64(screenWidth/2)
	y = float64(screenHeight/2)
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func moveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathBits := strings.Split( r.URL.Path, "/")
	direction := pathBits[len(pathBits)-1]
	fmt.Println("/direction/", direction)

	switch direction {
		case "up": y += 25
		case "down": y -= 25
		case "left": x -= 25
		case "right": x += 25
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

	spriteImage, initError := loadImageFile("sonic.png")
	if initError != nil {
		panic(initError)
	}

	spritePic := pixel.PictureDataFromImage(spriteImage)

	sprite = pixel.NewSprite(spritePic, spritePic.Bounds())

	spriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)

	for !win.Closed() {

		spriteBatch.Clear()

		matrix := pixel.IM.Rotated(pixel.ZV, 0).Scaled(pixel.ZV, 0.2).Moved(pixel.Vec{X: x, Y: y})

		sprite.Draw(spriteBatch, matrix)

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