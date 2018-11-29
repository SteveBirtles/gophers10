package main

import (
	"github.com/faiface/pixel"
)

type Location struct {
	x, y int
}

type Player struct {
	lastLocation     Location
	targetQueue      []Location
	progress         float64
	screenX, screenY float64
	frame, direction int
	ai               bool
	spriteRow        int
}

var player [4]Player

func preparePlayers() {

	/* Code goes here */

}

func (p *Player) updatePosition() {

	/* Code goes here */

}

func (p *Player) pickRandomDirection() {

	/* Code goes here */

}

func (p *Player) draw() {

	if len(p.targetQueue) > 0 {
		currentTarget := p.targetQueue[0]
		p.screenX = 80 * (float64(p.lastLocation.x) + float64(currentTarget.x-p.lastLocation.x)*p.progress)
		p.screenY = ScreenHeight - 80*(float64(p.lastLocation.y)+float64(currentTarget.y-p.lastLocation.y)*p.progress)
	} else {
		p.screenX = 80 * (float64(p.lastLocation.x))
		p.screenY = ScreenHeight - 80*(float64(p.lastLocation.y))
	}
	matrix := pixel.IM.Rotated(pixel.ZV, 0).Scaled(pixel.ZV, 4).Moved(pixel.Vec{X: p.screenX - 8, Y: p.screenY - 16})
	sprite[p.spriteRow][p.direction][p.frame].Draw(batch, matrix)

}
