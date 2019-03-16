package services

import (
	"context"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/commands"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf"
	"github.com/imrenagi/goes-werewolf/pkg/eventsource/sourcing"
	"github.com/imrenagi/goes-werewolf/pkg/eventsource/storage"
)

type EventRepository interface {
	Store(context.Context, sourcing.EventSource) error
	Load(context.Context, storage.AggregateID) []sourcing.Event
}

type GameService struct {
	EventRepository EventRepository
}

func (g GameService) ProcessInitializeGame(ctx context.Context, command commands.InitializeGame) error {
	game, err := werewolf.NewGame(
		command.ChannelID,
		command.ChannelName,
		command.CreatorID,
		command.CreatorName,
		command.Platform,
	)
	if err != nil {
		return err
	}

	err = g.storeAndPublishEvent(ctx, *game)
	if err != nil {
		return err
	}
	return nil
}

func (g GameService) loadGame(ctx context.Context, aggregateID storage.AggregateID) (*werewolf.Game, error) {
	events := g.EventRepository.Load(ctx, aggregateID)
	game, err := werewolf.NewGameFromHistory(sourcing.NewEventSourceId(), events)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g GameService) storeAndPublishEvent(ctx context.Context, game werewolf.Game) error {
	err := g.EventRepository.Store(ctx, game)
	if err != nil {
		return err
	}

	//todo dispatch events to all listener
	//todo the listener can listen to the event and store the current state to database. use goroutine

	//remove all stored events in memory
	game.Accept()

	return nil
}
