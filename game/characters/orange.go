package characters

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Orange struct {
	PosX             float64
	PosY             float64
	Speed            float64
	LastDirection    Direction
	CurrentDirection Direction
	counter          int
	move             bool
	keys             []ebiten.Key

	moveUpIdle *ebiten.Image
	moveUpA    *ebiten.Image
	moveUpB    *ebiten.Image

	moveDownIdle *ebiten.Image
	moveDownA    *ebiten.Image
	moveDownB    *ebiten.Image

	moveLeftIdle *ebiten.Image
	moveLeftA    *ebiten.Image
	moveLeftB    *ebiten.Image

	moveRightIdle *ebiten.Image
	moveRightA    *ebiten.Image
	moveRightB    *ebiten.Image

	moveUpLeftIdle *ebiten.Image
	moveUpLeftA    *ebiten.Image
	moveUpLeftB    *ebiten.Image

	moveUpRightIdle *ebiten.Image
	moveUpRightA    *ebiten.Image
	moveUpRightB    *ebiten.Image

	moveDownLeftIdle *ebiten.Image
	moveDownLeftA    *ebiten.Image
	moveDownLeftB    *ebiten.Image

	moveDownRightIdle *ebiten.Image
	moveDownRightA    *ebiten.Image
	moveDownRightB    *ebiten.Image
}

func (o *Orange) MoveUp(screen *ebiten.Image) {
	o.PosY -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveUpA, opts)
	} else {
		screen.DrawImage(o.moveUpB, opts)
	}
	o.LastDirection = UP
}
func (o *Orange) MoveUpLeft(screen *ebiten.Image) {
	o.PosY -= o.Speed
	o.PosX -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveUpLeftA, opts)
	} else {
		screen.DrawImage(o.moveUpLeftB, opts)
	}
	o.LastDirection = UP_LEFT
}
func (o *Orange) MoveUpRight(screen *ebiten.Image) {
	o.PosY -= o.Speed
	o.PosX += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveUpRightA, opts)
	} else {
		screen.DrawImage(o.moveUpRightB, opts)
	}
	o.LastDirection = UP_RIGHT
}
func (o *Orange) MoveDown(screen *ebiten.Image) {
	o.PosY += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveDownA, opts)
	} else {
		screen.DrawImage(o.moveDownB, opts)
	}
	o.LastDirection = DOWN
}
func (o *Orange) MoveDownLeft(screen *ebiten.Image) {
	o.PosY += o.Speed
	o.PosX -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveDownLeftA, opts)
	} else {
		screen.DrawImage(o.moveDownLeftB, opts)
	}
	o.LastDirection = DOWN_LEFT
}
func (o *Orange) MoveDownRight(screen *ebiten.Image) {
	o.PosY += o.Speed
	o.PosX += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveDownRightA, opts)
	} else {
		screen.DrawImage(o.moveDownRightB, opts)
	}
	o.LastDirection = DOWN_RIGHT
}
func (o *Orange) MoveLeft(screen *ebiten.Image) {
	o.PosX -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveLeftA, opts)
	} else {
		screen.DrawImage(o.moveLeftB, opts)
	}
	o.LastDirection = LEFT
}
func (o *Orange) MoveRight(screen *ebiten.Image) {
	o.PosX += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if o.move {
		screen.DrawImage(o.moveRightA, opts)
	} else {
		screen.DrawImage(o.moveRightB, opts)
	}
	o.LastDirection = RIGHT
}
func (o *Orange) Idle(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	switch o.LastDirection {
	case UP:
		screen.DrawImage(o.moveUpIdle, opts)
	case UP_RIGHT:
		screen.DrawImage(o.moveUpRightIdle, opts)
	case UP_LEFT:
		screen.DrawImage(o.moveUpLeftIdle, opts)
	case DOWN:
		screen.DrawImage(o.moveDownIdle, opts)
	case DOWN_LEFT:
		screen.DrawImage(o.moveDownLeftIdle, opts)
	case DOWN_RIGHT:
		screen.DrawImage(o.moveDownRightIdle, opts)
	case LEFT:
		screen.DrawImage(o.moveLeftIdle, opts)
	case RIGHT:
		screen.DrawImage(o.moveRightIdle, opts)
	}
}
func (o *Orange) Move(screen *ebiten.Image) {
	switch o.CurrentDirection {
	case UP:
		o.MoveUp(screen)
	case UP_RIGHT:
		o.MoveUpRight(screen)
	case UP_LEFT:
		o.MoveUpLeft(screen)
	case DOWN:
		o.MoveDown(screen)
	case DOWN_LEFT:
		o.MoveDownLeft(screen)
	case DOWN_RIGHT:
		o.MoveDownRight(screen)
	case LEFT:
		o.MoveLeft(screen)
	case RIGHT:
		o.MoveRight(screen)
	default:
		o.Idle(screen)
	}
}

func (o *Orange) Update() {
	if o.counter%(20/int(o.Speed+1)) == 0 {
		o.move = !o.move
	}
	o.counter++
	o.CurrentDirection = NONE
	o.keys = inpututil.AppendPressedKeys(o.keys[:0])
	if contains(o.keys, ebiten.KeyW) && contains(o.keys, ebiten.KeyD) {
		o.CurrentDirection = UP_RIGHT
	} else if contains(o.keys, ebiten.KeyW) && contains(o.keys, ebiten.KeyA) {
		o.CurrentDirection = UP_LEFT
	} else if contains(o.keys, ebiten.KeyS) && contains(o.keys, ebiten.KeyD) {
		o.CurrentDirection = DOWN_RIGHT
	} else if contains(o.keys, ebiten.KeyS) && contains(o.keys, ebiten.KeyA) {
		o.CurrentDirection = DOWN_LEFT
	} else if contains(o.keys, ebiten.KeyW) {
		o.CurrentDirection = UP
	} else if contains(o.keys, ebiten.KeyS) {
		o.CurrentDirection = DOWN
	} else if contains(o.keys, ebiten.KeyA) {
		o.CurrentDirection = LEFT
	} else if contains(o.keys, ebiten.KeyD) {
		o.CurrentDirection = RIGHT
	}
}

func contains(s []ebiten.Key, k ebiten.Key) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}
	return false
}

func NewOrange(loaded map[string]*ebiten.Image) *Orange {
	return &Orange{
		LastDirection: DOWN,
		Speed:         .5,
		PosX:          50,
		PosY:          50,

		moveUpIdle: loaded["orange_up_idle.png"],
		moveUpA:    loaded["orange_up_move_a.png"],
		moveUpB:    loaded["orange_up_move_b.png"],

		moveDownIdle: loaded["orange_down_idle.png"],
		moveDownA:    loaded["orange_down_move_a.png"],
		moveDownB:    loaded["orange_down_move_b.png"],

		moveLeftIdle: loaded["orange_left_idle.png"],
		moveLeftA:    loaded["orange_left_move_a.png"],
		moveLeftB:    loaded["orange_left_move_b.png"],

		moveRightIdle: loaded["orange_right_idle.png"],
		moveRightA:    loaded["orange_right_move_a.png"],
		moveRightB:    loaded["orange_right_move_b.png"],

		moveUpLeftIdle: loaded["orange_up_left_idle.png"],
		moveUpLeftA:    loaded["orange_up_left_move_a.png"],
		moveUpLeftB:    loaded["orange_up_left_move_b.png"],

		moveUpRightIdle: loaded["orange_up_right_idle.png"],
		moveUpRightA:    loaded["orange_up_right_move_a.png"],
		moveUpRightB:    loaded["orange_up_right_move_b.png"],

		moveDownLeftIdle: loaded["orange_down_left_idle.png"],
		moveDownLeftA:    loaded["orange_down_left_move_a.png"],
		moveDownLeftB:    loaded["orange_down_left_move_b.png"],

		moveDownRightIdle: loaded["orange_down_right_idle.png"],
		moveDownRightA:    loaded["orange_down_right_move_a.png"],
		moveDownRightB:    loaded["orange_down_right_move_b.png"],
	}
}
