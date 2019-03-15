package characters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWolf_String(t *testing.T) {
	w := Wolf{}
	require.Equal(t, "Wolf", w.String())
}

func TestWolf_VisitVillager(t *testing.T) {
	w := NewWolf()

	t.Run("Villager dies after it is visited by wolf", func(t *testing.T) {
		v := NewVillager()
		require.True(t, v.IsAlive())
		w.VisitVillager(v)
		require.False(t, v.IsAlive())
		require.True(t, w.IsAlive())
	})
}

func TestWolf_VisitDoctor(t *testing.T) {
	w := NewWolf()

	t.Run("Doctor dies after it is visited by wolf", func(t *testing.T) {
		d := NewDoctor()
		require.True(t, d.IsAlive())
		w.VisitDoctor(d)
		require.False(t, d.IsAlive())
		require.True(t, w.IsAlive())
	})
}

func TestWolf_VisitSeer(t *testing.T) {
	w := NewWolf()

	t.Run("Seer dies after being visited by wolf", func(t *testing.T) {
		s := NewSeer()
		require.True(t, s.IsAlive())
		w.VisitSeer(s)
		require.False(t, s.IsAlive())
		require.True(t, w.IsAlive())
	})

}

func TestWolf_Accept(t *testing.T) {

	t.Run("Panic if wolf is visited by other wolf", func(t *testing.T) {
		w := NewWolf()
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if wolf visit wolf")
			}
		}()
		w2 := NewWolf()
		w.Accept(w2)
	})

	t.Run("Panic if wolf is visited by villager", func(t *testing.T) {
		w := NewWolf()
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visit wolf")
			}
		}()
		v := NewVillager()
		w.Accept(v)
	})

	t.Run("Seer should be alive if wolf accept him", func(t *testing.T) {
		w := NewWolf()
		s := NewSeer()

		w.Accept(s)

		require.True(t, w.IsAlive())
		require.True(t, s.IsAlive())
	})

	t.Run("Doctor dies if wolf accepts him", func(t *testing.T) {
		w := NewWolf()
		d := NewDoctor()

		w.Accept(d)

		require.True(t, w.IsAlive())
		require.False(t, d.IsAlive())
	})
}
