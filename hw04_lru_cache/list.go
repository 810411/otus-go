package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *Item
	Back() *Item
	PushFront(v interface{}) *Item
	PushBack(v interface{}) *Item
	Remove(i *Item)
	MoveToFront(i *Item)
}

type Item struct {
	Value      interface{}
	Next, Prev *Item
}

type list struct {
	front, back *Item
	len         int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *Item {
	return l.front
}

func (l list) Back() *Item {
	return l.back
}

func (l *list) PushFront(v interface{}) *Item {
	i := &Item{Value: v}

	if l.front == nil {
		l.back = i
	} else {
		i.Prev = l.front
		l.front.Next = i
	}

	l.front = i
	l.len++

	return i
}

func (l *list) PushBack(v interface{}) *Item {
	i := &Item{Value: v}

	if l.back == nil {
		l.front = i
	} else {
		i.Next = l.back
		l.back.Prev = i
	}

	l.back = i
	l.len++

	return i
}

func (l *list) Remove(i *Item) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.front = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}

	l.len--
}

func (l *list) MoveToFront(i *Item) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
