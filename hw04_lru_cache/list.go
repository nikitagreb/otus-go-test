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
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := NewListItem()
	item.Value = v

	l.len++

	// [] -> [10]
	if l.front == nil {
		l.back = item
		l.front = item

		return item
	}

	// [10] -> [20, 10]
	item.Next = l.front
	l.front.Prev = item
	l.front = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := NewListItem()
	item.Value = v

	l.len++

	// [] -> [10]
	if l.back == nil && l.front == nil {
		l.back = item
		l.front = item

		return item
	}

	// [10] -> [10, 20]
	item.Prev = l.back
	l.back.Next = item
	l.back = item

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	switch {
	case i.Next == nil && i.Prev == nil:
		l.back = nil
		l.front = nil
	case i.Next != nil && i.Prev != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	case i.Next == nil && i.Prev != nil:
		l.back = i.Prev
		l.back.Next = nil
	case i.Next != nil && i.Prev == nil:
		l.back = i.Next
		i.Next.Prev = nil
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}

func NewListItem() *ListItem {
	return new(ListItem)
}
