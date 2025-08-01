package kgo

import (
	"testing"
)

func TestRing(t *testing.T) {
	t.Run("push multiple elements and then drop them", func(t *testing.T) {
		r := &ring[int]{}

		assertRingPush(t, r, 1, true, false)
		assertRingPush(t, r, 2, false, false)
		assertRingPush(t, r, 3, false, false)

		assertRingDropPeek(t, r, 2, true, false)
		assertRingDropPeek(t, r, 3, true, false)
		assertRingDropPeek(t, r, 0, false, false)
	})

	t.Run("push and drop elements iteratively", func(t *testing.T) {
		r := &ring[int]{}

		// Push an initial element.
		assertRingPush(t, r, 1, true, false)

		// Push an element and them drop the previous one, multiple times.
		for i := 2; i < 10; i++ {
			assertRingPush(t, r, i, false, false)
			assertRingDropPeek(t, r, i, true, false)

			if len(r.overflow) > 0 {
				t.Error("unexpected overflow usage")
			}
		}

		// Finally, drop the last element.
		assertRingDropPeek(t, r, 0, false, false)
	})

	t.Run("push elements above the ring capacity and get them stored in the overflow", func(t *testing.T) {
		r := &ring[int]{}

		for i := 1; i <= 10; i++ {
			assertRingPush(t, r, i, i == 1, false)
		}

		// We expect the overflow has been used.
		if len(r.overflow) == 0 {
			t.Error("unexpected empty overflow")
		}

		for i := 1; i <= 9; i++ {
			assertRingDropPeek(t, r, i+1, true, false)
		}
		assertRingDropPeek(t, r, 0, false, false)

		// At this point the overflow should have been cleared.
		if len(r.overflow) > 0 {
			t.Error("unexpected overflow usage")
		}
	})

	t.Run("interrupt a non-full ring", func(t *testing.T) {
		r := &ring[int]{}

		assertRingPush(t, r, 1, true, false)
		assertRingPush(t, r, 2, false, false)
		assertRingPush(t, r, 3, false, false)

		r.die()

		assertRingPush(t, r, 4, false, true)
		assertRingDropPeek(t, r, 2, true, true)
	})

	t.Run("continuously keeping the items in the ring above the fixed size limit should not grow the overflow slice indefinitely", func(t *testing.T) {
		r := &ring[int]{}

		// Push an initial number of elements above the fixed size length.
		for i := 1; i <= 10; i++ {
			assertRingPush(t, r, i, i == 1, false)
		}

		// Now keep pushing and dropping elements continuously.
		for i := 11; i <= 1000; i++ {
			assertRingPush(t, r, i, false, false)
			assertRingDropPeek(t, r, i-9, true, false)
		}

		if cap(r.overflow) > 20 {
			t.Errorf("unexpected high overflow slice capacity, got: %d", cap(r.overflow))
		}
	})

	t.Run("having a temporarily high number of items in the ring should not keep the overflow slice capacity high indefinitely", func(t *testing.T) {
		r := &ring[int]{}

		// Push a large number of elements.
		for i := 1; i <= 1000; i++ {
			assertRingPush(t, r, i, i == 1, false)
		}

		if cap(r.overflow) < 1000 {
			t.Errorf("unexpected low overflow slice capacity, got: %d, expected >= 1000", cap(r.overflow))
		}

		// Drop most of them, but keep it above the fixed size limit.
		for i := 1; i <= 990; i++ {
			assertRingDropPeek(t, r, i+1, true, false)
		}

		// Push few more items and then drop all the remaining ones.
		for i := 1001; i <= 1010; i++ {
			assertRingPush(t, r, i, false, false)
		}

		for i := 991; i < 1010; i++ {
			assertRingDropPeek(t, r, i+1, true, false)
		}
		assertRingDropPeek(t, r, 0, false, false)

		if cap(r.overflow) > 500 {
			t.Errorf("unexpected high overflow slice capacity, got: %d", cap(r.overflow))
		}
	})
}

func assertRingPush(t *testing.T, r *ring[int], elem int, expectedFirst, expectedDead bool) {
	t.Helper()

	first, dead := r.push(elem)
	if first != expectedFirst {
		t.Errorf("unexpected first: got %t, want %t", first, expectedFirst)
	}
	if dead != expectedDead {
		t.Errorf("unexpected dead: got %t, want %t", dead, expectedDead)
	}
}

func assertRingDropPeek(t *testing.T, r *ring[int], expectedNext int, expectedMore, expectedDead bool) {
	t.Helper()

	next, more, dead := r.dropPeek()
	if next != expectedNext {
		t.Errorf("unexpected next element: got %d, want %d", next, expectedNext)
	}
	if more != expectedMore {
		t.Errorf("unexpected more: got %t, want %t", more, expectedMore)
	}
	if dead != expectedDead {
		t.Errorf("unexpected dead: got %t, want %t", dead, expectedDead)
	}
}
