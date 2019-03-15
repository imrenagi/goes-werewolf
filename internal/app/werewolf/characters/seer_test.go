package characters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeer_String(t *testing.T) {
	s := NewSeer()
	require.Equal(t, "Seer", s.String())
}

func TestSeer_VisitAnyBodyNothingHappened(t *testing.T) {
	s := NewSeer()

	t.Run("visit wolf nothing happened", func(t *testing.T) {
		w := NewWolf()
		s.VisitWolf(w)

		require.True(t, s.IsAlive())
		require.True(t, w.IsAlive())
	})

	t.Run("visit villager nothing happened", func(t *testing.T) {
		v := NewVillager()
		s.VisitVillager(v)

		require.True(t, s.IsAlive())
		require.True(t, v.IsAlive())
	})

	t.Run("visit doctor nothing happened", func(t *testing.T) {
		d := NewDoctor()
		s.VisitDoctor(d)

		require.True(t, s.IsAlive())
		require.True(t, d.IsAlive())
	})

	t.Run("visit seer is not allowed", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if seer visits seer")
			}
		}()

		s2 := NewSeer()
		s.VisitSeer(s2)
	})
}

func TestSeer_Accept(t *testing.T) {

	t.Run("Seer dies, after accepting wolf", func(t *testing.T) {
		s := NewSeer()
		w := NewWolf()

		s.Accept(w)

		require.False(t, s.IsAlive())
		require.True(t, w.IsAlive())
	})

	t.Run("Seer is revive after visited by doctor", func(t *testing.T) {
		s := NewSeer()
		d := NewDoctor()

		s.Die()
		s.Accept(d)

		require.True(t, s.IsAlive())
		require.True(t, d.IsAlive())

	})

	t.Run("Seer should not be visited by another seer", func(t *testing.T) {
		s1 := NewSeer()
		s2 := NewSeer()

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if seer visits another seer")
			}
		}()
		s1.Accept(s2)
	})

	t.Run("Seer should not be visited by villager", func(t *testing.T) {
		s := NewSeer()
		v := NewVillager()

		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villanger visits seer")
			}
		}()
		s.Accept(v)
	})
}
