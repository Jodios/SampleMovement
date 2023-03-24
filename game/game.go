package game

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	g.Character.Update()
	return nil
}

// called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.Character.Move(screen)
}

// accepts the window dimensions and returns the inside
// size autoscaled (PRETTY COOL)
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.InnerWidth, g.InnerHeight
}
