package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len    int
	queue  *ListItem
	lastEl *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.queue
}

func (l *list) Back() *ListItem {
	return l.lastEl
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Next: l.queue}
	l.queue = newListItem

	if l.queue.Next != nil {
		l.queue.Next.Prev = newListItem
	}

	// Если это первый элемент в очереди, то устанавливаем lastEl
	if l.lastEl == nil {
		l.lastEl = l.queue
	}

	l.len++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Prev: l.lastEl}

	if l.lastEl == nil {
		l.queue = newListItem
		l.lastEl = newListItem
	} else {
		l.lastEl.Next = newListItem
		l.lastEl = newListItem
	}
	l.len++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	// Если это последний элемент, то нужно обновить lastEl
	if i.Next == nil {
		l.lastEl = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	l.len--

	// Если элементов в очереди не осталось, то обнуляем lastEl
	if l.len == 0 {
		l.lastEl = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// Если это первый элемент, то ничего делать не нужно
	if i.Prev == nil {
		return
	}

	// Если это последний элемент, то нужно обновить lastEl
	if i.Next == nil {
		l.lastEl = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = l.queue
	i.Prev = nil
	l.queue = i
}

func NewList() List {
	return new(list)
}
