package characters

import "github.com/hajimehoshi/ebiten/v2"

type Orange struct {
	PosX          float64
	PosY          float64
	Speed         float64
	LastDirection Direction

	MoveUpIdle *ebiten.Image
	MoveUpA    *ebiten.Image
	MoveUpB    *ebiten.Image

	MoveDownIdle *ebiten.Image
	MoveDownA    *ebiten.Image
	MoveDownB    *ebiten.Image

	MoveLeftIdle *ebiten.Image
	MoveLeftA    *ebiten.Image
	MoveLeftB    *ebiten.Image

	MoveRightIdle *ebiten.Image
	MoveRightA    *ebiten.Image
	MoveRightB    *ebiten.Image

	MoveUpLeftIdle *ebiten.Image
	MoveUpLeftA    *ebiten.Image
	MoveUpLeftB    *ebiten.Image

	MoveUpRightIdle *ebiten.Image
	MoveUpRightA    *ebiten.Image
	MoveUpRightB    *ebiten.Image

	MoveDownLeftIdle *ebiten.Image
	MoveDownLeftA    *ebiten.Image
	MoveDownLeftB    *ebiten.Image

	MoveDownRightIdle *ebiten.Image
	MoveDownRightA    *ebiten.Image
	MoveDownRightB    *ebiten.Image
}

func (o *Orange) MoveUp(screen *ebiten.Image) {
	o.PosY -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	screen.DrawImage(o.MoveUpIdle, opts)
	o.LastDirection = UP
}
func (o *Orange) MoveUpLeft(screen *ebiten.Image) {
}
func (o *Orange) MoveUpRight(screen *ebiten.Image) {
}
func (o *Orange) MoveDown(screen *ebiten.Image) {
	o.PosY += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	screen.DrawImage(o.MoveDownIdle, opts)
	o.LastDirection = DOWN
}
func (o *Orange) MoveDownLeft(screen *ebiten.Image) {

}
func (o *Orange) MoveDownRight(screen *ebiten.Image) {

}
func (o *Orange) MoveLeft(screen *ebiten.Image) {
	o.PosX -= o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	screen.DrawImage(o.MoveLeftIdle, opts)
	o.LastDirection = LEFT
}
func (o *Orange) MoveRight(screen *ebiten.Image) {
	o.PosX += o.Speed
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	screen.DrawImage(o.MoveRightIdle, opts)
	o.LastDirection = RIGHT
}
func (o *Orange) Idle(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	switch o.LastDirection {
	case UP:
		screen.DrawImage(o.MoveUpIdle, opts)
	case DOWN:
		screen.DrawImage(o.MoveDownIdle, opts)
	case LEFT:
		screen.DrawImage(o.MoveLeftIdle, opts)
	case RIGHT:
		screen.DrawImage(o.MoveRightIdle, opts)
	}
}

func NewOrange(loaded map[string]*ebiten.Image) *Orange {
	return &Orange{
		LastDirection: DOWN,
		Speed:         .5,
		PosX:          50,
		PosY:          50,

		MoveUpIdle: loaded["orange_up_idle.png"],
		MoveUpA:    loaded["orange_up_move_a.png"],
		MoveUpB:    loaded["orange_up_move_b.png"],

		MoveDownIdle: loaded["orange_down_idle.png"],
		MoveDownA:    loaded["orange_down_move_a.png"],
		MoveDownB:    loaded["orange_down_move_b.png"],

		MoveLeftIdle: loaded["orange_left_idle.png"],
		MoveLeftA:    loaded["orange_left_move_a.png"],
		MoveLeftB:    loaded["orange_left_move_b.png"],

		MoveRightIdle: loaded["orange_right_idle.png"],
		MoveRightA:    loaded["orange_right_move_a.png"],
		MoveRightB:    loaded["orange_right_move_b.png"],

		MoveUpLeftIdle: loaded["orange_up_left_idle.png"],
		MoveUpLeftA:    loaded["orange_up_left_move_a.png"],
		MoveUpLeftB:    loaded["orange_up_left_move_b.png"],

		MoveUpRightIdle: loaded["orange_up_right_idle.png"],
		MoveUpRightA:    loaded["orange_up_right_move_a.png"],
		MoveUpRightB:    loaded["orange_up_right_move_b.png"],

		MoveDownLeftIdle: loaded["orange_down_left_idle.png"],
		MoveDownLeftA:    loaded["orange_down_left_move_a.png"],
		MoveDownLeftB:    loaded["orange_down_left_move_b.png"],

		MoveDownRightIdle: loaded["orange_down_right_idle.png"],
		MoveDownRightA:    loaded["orange_down_right_move_a.png"],
		MoveDownRightB:    loaded["orange_down_right_move_b.png"],
	}
}
