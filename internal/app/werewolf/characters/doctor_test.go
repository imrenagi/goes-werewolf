package characters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoctor_String(t *testing.T) {
	d := NewDoctor()
	require.Equal(t, "Doctor", d.String())
}

func TestDoctor_VisitWolf(t *testing.T) {
	d := NewDoctor()
	w := NewWolf()

	d.VisitWolf(w)

	require.True(t, w.IsAlive())
	require.False(t, d.IsAlive())
}

func TestDoctor_VisitDoctor(t *testing.T) {

	t.Run("should revive its self", func(t *testing.T) {
		d := NewDoctor()
		d.Die()
		d.VisitDoctor(d)
		require.True(t, d.IsAlive())
	})

	//t.Run("cant revive another doctor", func(t *testing.T) {
	//	d1 := NewDoctor()
	//	d2 := NewDoctor()
	//
	//	defer func() {
	//		if r := recover(); r != nil {
	//			t.Log("Recovered: ", r)
	//		} else {
	//			t.Fail()
	//			t.Logf("Should panic if doctor visits another doctor")
	//		}
	//	}()
	//
	//	d1.VisitDoctor(d2)
	//})
}

func TestDoctor_VisitNonHarmfulVillager(t *testing.T) {

	t.Run("visit seer", func(t *testing.T) {
		d := NewDoctor()
		s := NewSeer()

		s.Die()
		require.False(t, s.IsAlive())

		d.VisitSeer(s)

		require.True(t, d.IsAlive())
		require.True(t, s.IsAlive())
	})

	t.Run("visit villager", func(t *testing.T) {
		d := NewDoctor()
		v := NewVillager()

		v.Die()
		require.False(t, v.IsAlive())
		d.VisitVillager(v)

		require.True(t, d.IsAlive())
		require.True(t, v.IsAlive())
	})
}

func TestDoctor_Accept(t *testing.T) {

	t.Run("Accept himself when he is alive, nothing happened", func(t *testing.T) {
		d := NewDoctor()

		d.Accept(d)

		require.True(t, d.IsAlive())
	})

	t.Run("Accept himself when he dies, then he is revived", func(t *testing.T) {
		d := NewDoctor()

		d.Die()
		require.False(t, d.IsAlive())

		d.Accept(d)

		require.True(t, d.IsAlive())
	})

	t.Run("Accept wolf, then he dies", func(t *testing.T) {
		d := NewDoctor()
		w := NewWolf()

		d.Accept(w)

		require.False(t, d.IsAlive())
		require.True(t, w.IsAlive())
	})

	t.Run("Accept seer, then nothing happened", func(t *testing.T) {
		d := NewDoctor()
		s := NewSeer()

		d.Accept(s)

		require.True(t, d.IsAlive())
		require.True(t, s.IsAlive())
	})

	t.Run("Accept villager, should panic", func(t *testing.T) {
		d := NewDoctor()
		defer func() {
			if r := recover(); r != nil {
				t.Log("Recovered: ", r)
			} else {
				t.Fail()
				t.Logf("Should panic if villager visits doctor")
			}
		}()
		v := NewVillager()
		d.Accept(v)
	})
}
