package main

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

const MaxQueueLength = 3

type Location struct {
	x, y int
}

type Player struct {
	lastLocation               Location
	targetQueue                []Location
	screenX, screenY, progress float64
	frame, direction           int
	ai                         bool
	spriteRow                  int
}

var player [4]Player

func preparePlayers() {

	for i := range player {
		player[i].lastLocation = Location{int((i%2)*4 + 2), int(math.Floor(float64(i/2))*4 + 2)}
		player[i].targetQueue = make([]Location, 0)
		player[i].frame = 3
		player[i].spriteRow = i

		if i != 0 {
			player[i].ai = true
			player[i].pickRandomDirection()
		}
	}

}

func (p *Player) pickRandomDirection() {

	pos := p.lastLocation

pick:
	for {

		p.direction = rand.Intn(4)
		switch p.direction {
		case 0:
			if pos.x > 1 {
				p.targetQueue = append(p.targetQueue, Location{pos.x - 1, pos.y})
				p.progress = 0
				break pick
			}
		case 1:
			if pos.y < ScreenHeight/80-1 {
				p.targetQueue = append(p.targetQueue, Location{pos.x, pos.y + 1})
				p.progress = 0
				break pick
			}
		case 2:
			if pos.y > 1 {
				p.targetQueue = append(p.targetQueue, Location{pos.x, pos.y - 1})
				p.progress = 0
				break pick
			}
		case 3:
			if pos.x < ScreenWidth/80 {
				p.targetQueue = append(p.targetQueue, Location{pos.x + 1, pos.y})
				p.progress = 0
				break pick
			}
		}

	}

	p.progress = 0
	p.frame = 3

}

func (p *Player) updatePosition() {

	if len(p.targetQueue) > 0 {

		currentTarget := p.targetQueue[0]

		if p.progress < 1 {

			switch {
			case currentTarget.x < p.lastLocation.x:
				p.direction = 0
			case currentTarget.y > p.lastLocation.y:
				p.direction = 1
			case currentTarget.y < p.lastLocation.y:
				p.direction = 2
			case currentTarget.x > p.lastLocation.x:
				p.direction = 3
			}

			p.progress += 0.125
			p.frame--
			if p.frame < 0 {
				p.frame = 3
			}
		}

		if p.progress >= 1 {

			p.lastLocation = currentTarget
			p.targetQueue = p.targetQueue[1:]
			p.progress = 0

			if p.ai {
				p.pickRandomDirection()
			} else {
				p.frame = 3
			}

		}

	}

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

	matrix := pixel.IM.Rotated(pixel.ZV, 0).Scaled(pixel.ZV, 4).Moved(pixel.Vec{X: p.screenX, Y: p.screenY})
	sprite[p.spriteRow][p.direction][p.frame].Draw(batch, matrix)

}
