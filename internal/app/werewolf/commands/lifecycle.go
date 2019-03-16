package commands

type InitializeGame struct {
	Platform  string
	ChannelID, ChannelName string
	CreatorID, CreatorName string
}

type CancelGame struct {
}

type StartGame struct {
}

type EndGame struct {
}

type Advance struct {
}

type JoinGame struct {
}
