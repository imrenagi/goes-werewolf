package models

type Poll struct {
	ID        string
	PollID    string
	GameID    string
	GameState string
	GameDay   int
	Voter     Player
	Choice    Player
}
