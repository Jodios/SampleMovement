package characters

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Orange struct {
	PosX      float64
	PosY      float64
	Speed     float64
	counter   int
	move      bool
	keys      []ebiten.Key
	direction [2]int
	movements [3][3][3]*ebiten.Image
}

func (o *Orange) Move(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(o.PosX, o.PosY)
	if len(o.keys) == 0 {
		screen.DrawImage(o.movements[o.direction[1]][o.direction[0]][0], opts)
	} else {
		if o.move {
			screen.DrawImage(o.movements[o.direction[1]][o.direction[0]][1], opts)
		} else {
			screen.DrawImage(o.movements[o.direction[1]][o.direction[0]][2], opts)
		}
	}
}

func (o *Orange) Update() {
	if o.counter%(20/int(o.Speed+1)) == 0 {
		o.move = !o.move
	}
	o.counter++
	o.keys = inpututil.AppendPressedKeys(o.keys[:0])
	x, y := 0, 0
	if len(o.keys) > 0 {
		for _, v := range o.keys {
			if mv, ok := MoveVectorMapping[v]; ok {
				o.PosX += float64(mv.Xs) * o.Speed
				o.PosY += float64(mv.Ys) * o.Speed
				x += int(mv.Xs)
				y += (mv.Ys)
			}
		}
		o.direction[0] = x + 1
		o.direction[1] = y + 1
	}
}

func NewOrange(loaded map[string]*ebiten.Image) *Orange {
	var movements [3][3][3]*ebiten.Image = [3][3][3]*ebiten.Image{
		{
			{loaded["orange_up_left_idle.png"], loaded["orange_up_left_move_a.png"], loaded["orange_up_left_move_b.png"]},
			{loaded["orange_up_idle.png"], loaded["orange_up_move_a.png"], loaded["orange_up_move_b.png"]},
			{loaded["orange_up_right_idle.png"], loaded["orange_up_right_move_a.png"], loaded["orange_up_right_move_b.png"]},
		},
		{
			{loaded["orange_left_idle.png"], loaded["orange_left_move_a.png"], loaded["orange_left_move_b.png"]},
			{loaded["orange_up_idle.png"], loaded["orange_up_move_a.png"], loaded["orange_up_move_b.png"]}, // TODO: when all keys are pressed.... what do?
			{loaded["orange_right_idle.png"], loaded["orange_right_move_a.png"], loaded["orange_right_move_b.png"]},
		},
		{
			{loaded["orange_down_left_idle.png"], loaded["orange_down_left_move_a.png"], loaded["orange_down_left_move_b.png"]},
			{loaded["orange_down_idle.png"], loaded["orange_down_move_a.png"], loaded["orange_down_move_b.png"]},
			{loaded["orange_down_right_idle.png"], loaded["orange_down_right_move_a.png"], loaded["orange_down_right_move_b.png"]},
		},
	}
	return &Orange{
		Speed:     .5,
		PosX:      50,
		PosY:      50,
		movements: movements,
	}
}
