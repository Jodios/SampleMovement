package game

/*
Game stores the state of the
Connected user.
*/
type Game struct {
	Characters map[string]*Character `json:"characters,omitempty"`
	Added      *string               `json:"character,omitempty"`
}
