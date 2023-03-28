package game

import "github.com/jodios/samplemovement/game/characters"

type GameState struct {
	Characters     map[string]*characters.Orange `json:"characters,omitempty"`
	AddedCharacter *string                       `json:"character,omitempty"`
}
