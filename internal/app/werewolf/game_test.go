package werewolf

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/states"
)

func TestCreateNewGame(t *testing.T) {

	t.Run("start new game should start at initial state", func(t *testing.T) {
		game := NewGame()
		init := states.Initial{}
		require.Equal(t, init.String(), game.State.String())
	})

}
