package characters

import "github.com/hajimehoshi/ebiten/v2"

type Direction string

const (
	UP       Direction = "UP"
	UP_LEFT  Direction = "UP_LEFT"
	UP_RIGHT Direction = "UP_RIGHT"

	DOWN       Direction = "DOWN"
	DOWN_LEFT  Direction = "DOWN_LEFT"
	DOWN_RIGHT Direction = "DOWN_RIGHT"

	LEFT Direction = "LEFT"

	RIGHT Direction = "RIGHT"
	NONE  Direction = "NONE"
)

type MoveVector struct {
	Xs int
	Ys int
}

var MoveVectorMapping map[ebiten.Key]*MoveVector = map[ebiten.Key]*MoveVector{
	ebiten.KeyW: {
		Xs: 0,
		Ys: -1,
	},
	ebiten.KeyA: {
		Xs: -1,
		Ys: 0,
	},
	ebiten.KeyS: {
		Xs: 0,
		Ys: 1,
	},
	ebiten.KeyD: {
		Xs: 1,
		Ys: 0,
	},
}
