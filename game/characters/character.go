package characters

import "github.com/hajimehoshi/ebiten/v2"

type Direction string

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
