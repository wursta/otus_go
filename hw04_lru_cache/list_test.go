package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("simple", func(t *testing.T) {
		l := NewList()
		l.PushBack(10) // [10]
		l.PushBack(20) // [10, 20]
		l.PushBack(30) // [10, 20, 30]

		l.Remove(l.Front()) // [20, 30]
		require.Equal(t, 20, l.Front().Value)

		l.Remove(l.Back()) // [20]
		require.Equal(t, 20, l.Back().Value)

		l.Remove(l.Front()) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		l2 := NewList()
		l2.PushFront(10) // [10]
		l2.PushFront(20) // [10, 20]
		l2.PushFront(30) // [10, 20, 30]

		l2.Remove(l2.Back())  // [10, 20]
		l2.Remove(l2.Front()) // [20]
		l2.Remove(l2.Back())  // []
		require.Equal(t, 0, l2.Len())
		require.Nil(t, l2.Front())
		require.Nil(t, l2.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, toValuesList(l))

		// Проверяем перемещение элемента из середины списка
		l.MoveToFront(l.Front().Next.Next) // [60, 70, 80, 40, 10, 30, 50]
		require.Equal(t, []int{60, 70, 80, 40, 10, 30, 50}, toValuesList(l))
	})
}

func toValuesList(list List) []int {
	elems := make([]int, 0, list.Len())
	for i := list.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	return elems
}
