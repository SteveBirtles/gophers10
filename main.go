package main

import "github.com/faiface/pixel/pixelgl"

func main() {

	preparePlayers()

	go startServer()

	pixelgl.Run(startClient)

}
