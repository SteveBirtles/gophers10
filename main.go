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
	"math"
	"math/rand"
)

const screenWidth = 1024
const screenHeight = 768

type Player struct {
	lastP, lastQ, targetP, targetQ int
	x, y, progress                 float64
	frame, direction               int
}

var (
	frames = 0
	second = time.Tick(time.Second)
	tick   = time.Tick(time.Second / 16)
	title  = "Game client with an API"
	win    *pixelgl.Window
	sprite [4][4][4]*pixel.Sprite
	batch  *pixel.Batch
	player [4]Player
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
		player[0].y += 25
	case "down":
		player[0].y -= 25
	case "left":
		player[0].x -= 25
	case "right":
		player[0].x += 25
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
				sprite[i][j][k] = pixel.NewSprite(spritePic,
					pixel.R(float64(k+j*3)*20, float64(120-30*(i+1)), float64(k+j*3+1)*20, float64(120-30*i)))
			}
			sprite[i][j][3] = sprite[i][j][1]
		}
	}


	for i := range player {
		player[i].lastP = int((i % 2) * 4 + 2)
		player[i].lastQ = int(math.Floor(float64(i / 2)) * 4 + 2)
		player[i].targetP = player[i].lastP
		player[i].targetQ = player[i].lastQ
		player[i].direction = i
		player[i].progress = 1
	}

	batch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)

	for !win.Closed() {

		batch.Clear()

		select {
		case <-tick:
			for i := range player {
				if player[i].progress < 1 {
					player[i].progress += 0.125
					player[i].frame--
					if player[i].frame < 0 {
						player[i].frame = 3
					}
				}
				if player[i].progress >= 1 {
					player[i].progress = 1
					player[i].frame = 3
					player[i].lastP = player[i].targetP
					player[i].lastQ = player[i].targetQ

					player[i].direction = rand.Intn(4)
					switch player[i].direction {
					case 0:
						if player[i].lastP > 1 {
							player[i].targetP--
							player[i].progress = 0
						}
					case 1:
						if player[i].lastQ < screenHeight/80-1 {
							player[i].targetQ++
							player[i].progress = 0
						}
					case 2:
						if player[i].lastQ > 1 {
							player[i].targetQ--
							player[i].progress = 0
						}
					case 3:
						if player[i].lastP < screenWidth/80 {
							player[i].targetP++
							player[i].progress = 0
						}
					}

				}
				player[i].x = 80 * (float64(player[i].lastP) + float64(player[i].targetP - player[i].lastP) * player[i].progress)
				player[i].y = screenHeight - 80 * (float64(player[i].lastQ) + float64(player[i].targetQ - player[i].lastQ) * player[i].progress)
			}
		default:
		}

		for i, p := range player {
			matrix := pixel.IM.Rotated(pixel.ZV, 0).Scaled(pixel.ZV, 4).Moved(pixel.Vec{X: p.x, Y: p.y})
			sprite[i][p.direction][p.frame].Draw(batch, matrix)
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
