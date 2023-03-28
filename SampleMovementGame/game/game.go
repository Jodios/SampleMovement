package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/samplemovement/game/characters"
)

type Game struct {
	InnerWidth       int
	InnerHeight      int
	Character        *characters.Orange
	MoveUp           bool
	MoveDown         bool
	MoveLeft         bool
	MoveRight        bool
	CharacterChannel chan characters.Orange
	GameState        *GameState
}

// Called every tick
func (g *Game) Update() error {
	g.Character.Update()
	if len(g.Character.Keys) > 0 {
		g.CharacterChannel <- *g.Character
	}
	return nil
}

// called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{213, 218, 241, 255})
	g.Character.Move(screen)
	for _, v := range g.GameState.Characters {
		v.Move(screen)
	}
}

// accepts the window dimensions and returns the inside
// size autoscaled (PRETTY COOL)
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.InnerWidth, g.InnerHeight
}
