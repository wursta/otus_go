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
	firstItem  *ListItem
	lastItem   *ListItem
}

func (l *list) Len() int {
	return l.itemsCount
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Next: l.firstItem}
	l.firstItem = newListItem

	if l.firstItem.Next != nil {
		l.firstItem.Next.Prev = newListItem
	}

	// Если это первый элемент в цепочке, то он же является последним
	if l.lastItem == nil {
		l.lastItem = l.firstItem
	}

	l.itemsCount++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v, Prev: l.lastItem}

	if l.lastItem == nil {
		l.firstItem = newListItem
		l.lastItem = newListItem
	} else {
		l.lastItem.Next = newListItem
		l.lastItem = newListItem
	}
	l.itemsCount++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	// Если удаляется первый элемент, то нужно обновить firstItem
	if i.Prev == nil {
		l.firstItem = i.Next
	}
	// Если удаляется последний элемент, то нужно обновить lastItem
	if i.Next == nil {
		l.lastItem = i.Prev
	}

	// Если удаляется промежуточный элемент, то нужно изменить Next и Prev его соседей.
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.itemsCount--

	if l.itemsCount == 0 {
		l.firstItem = nil
		l.lastItem = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// Если перемещается первый элемент, то ничего делать не нужно
	if i.Prev == nil {
		return
	}

	// Если перемещается последний элемент, то нужно обновить lastItem
	if i.Next == nil {
		l.lastItem = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = l.firstItem
	i.Prev = nil
	l.firstItem = i
}

func NewList() List {
	return new(list)
}
