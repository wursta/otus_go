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
	itemsCount int
	firstEl    *ListItem
	lastEl     *ListItem
}

func (l *list) Len() int {
	return l.itemsCount
}

func (l *list) Front() *ListItem {
	return l.firstEl
}

func (l *list) Back() *ListItem {
	return l.lastEl
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Next: l.firstEl}
	l.firstEl = newListItem

	if l.firstEl.Next != nil {
		l.firstEl.Next.Prev = newListItem
	}

	// Если это первый элемент в цепочке, то он же является последним
	if l.lastEl == nil {
		l.lastEl = l.firstEl
	}

	l.itemsCount++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Prev: l.lastEl}

	if l.lastEl == nil {
		l.firstEl = newListItem
		l.lastEl = newListItem
	} else {
		l.lastEl.Next = newListItem
		l.lastEl = newListItem
	}
	l.itemsCount++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	// Если удаляется последний элемент, то нужно обновить lastEl
	if i.Next == nil {
		l.lastEl = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	l.itemsCount--

	if l.itemsCount == 0 {
		l.firstEl = nil
		l.lastEl = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// Если перемещается первый элемент, то ничего делать не нужно
	if i.Prev == nil {
		return
	}

	// Если перемещается последний элемент, то нужно обновить lastEl
	if i.Next == nil {
		l.lastEl = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = l.firstEl
	i.Prev = nil
	l.firstEl = i
}

func NewList() List {
	return new(list)
}
