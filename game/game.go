package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jodios/samplemovement/game/characters"
)

type Game struct {
	InnerWidth  int
	InnerHeight int
	Character   *characters.Orange
	MoveUp      bool
	MoveDown    bool
	MoveLeft    bool
	MoveRight   bool
}

// Called every tick
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.MoveUp = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		g.MoveUp = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.MoveRight = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.MoveRight = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.MoveDown = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		g.MoveDown = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.MoveLeft = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		g.MoveLeft = false
	}
	return nil
}

// called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	if g.MoveUp {
		g.Character.MoveUp(screen)
	} else if g.MoveDown {
		g.Character.MoveDown(screen)
	} else if g.MoveLeft {
		g.Character.MoveLeft(screen)
	} else if g.MoveRight {
		g.Character.MoveRight(screen)
	}
	g.Character.Idle(screen)
}

// accepts the window dimensions and returns the inside
// size autoscaled (PRETTY COOL)
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.InnerWidth, g.InnerHeight
}
