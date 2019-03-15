package characters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVillager_String(t *testing.T) {
	v := NewVillager()
	require.Equal(t, "Villager", v.String())
}

func TestVillager_CantVisitAnybody(t *testing.T) {
	v := NewVillager()

	t.Run("Visit wolf should panic", func(t *testing.T) {
		w := NewWolf()

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visits wolf")
			}
		}()
		v.VisitWolf(w)
	})

	t.Run("Visit doctor should panic", func(t *testing.T) {
		d := NewDoctor()

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visits doctor")
			}
		}()
		v.VisitDoctor(d)
	})

	t.Run("Visit seer should panic", func(t *testing.T) {
		s := NewSeer()

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visits seer")
			}
		}()
		v.VisitSeer(s)
	})

	t.Run("Visit other villager should panic", func(t *testing.T) {
		v2 := NewVillager()
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visits other villager")
			}
		}()
		v.VisitVillager(v2)

	})
}

func TestVillager_Accept(t *testing.T) {

	t.Run("Accept wolf, villager dies", func(t *testing.T) {
		v := NewVillager()
		w := NewWolf()

		v.Accept(w)

		require.False(t, v.IsAlive())
		require.True(t, w.IsAlive())
	})

	t.Run("Accept doctor, villager is revived", func(t *testing.T) {
		v := NewVillager()
		d := NewDoctor()

		require.True(t, v.IsAlive())
		require.True(t, d.IsAlive())

		v.Accept(d)

		require.True(t, v.IsAlive())
		require.True(t, d.IsAlive())

	})

	t.Run("Accept doctor after villager dies, villager is revived", func(t *testing.T) {
		v := NewVillager()
		d := NewDoctor()

		v.Die()

		require.False(t, v.IsAlive())

		v.Accept(d)

		require.True(t, v.IsAlive())
		require.True(t, d.IsAlive())
	})

	t.Run("Accept seer, villager got nothing", func(t *testing.T) {
		v := NewVillager()
		s := NewSeer()

		v.Accept(s)

		require.True(t, v.IsAlive())
		require.True(t, s.IsAlive())
	})
}
