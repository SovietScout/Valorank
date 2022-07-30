package models

type Match struct {
	State     State
	Players   []*Player
	GamePodID string
	Err       error
}
