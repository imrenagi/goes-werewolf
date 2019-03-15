package services

import "github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"

type PlayerDAO interface {
	GetPlayers(gameID string) ([]models.Player, error)
	GetPlayerWithID(gameID string, id string) (*models.Player, error)
	SavePlayer(gameID string, player models.Player) error

	//Only return alive player
	GetWolves(gameID string) ([]models.Player, error)
	GetWolf(gameID, playerID string) ([]models.Player, error)
	GetDoctor(gameID string) (*models.Player, error)
	GetSeer(gameID string) (*models.Player, error)
	GetVillager(gameID, playerID string) (*models.Player, error)
}
