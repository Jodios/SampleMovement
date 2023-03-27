package game

type Character struct {
	PosX  float64 `json:"posx,omitempty"`
	PosY  float64 `json:"posy,omitempty"`
	Speed float64 `json:"speed,omitempty"`
}
