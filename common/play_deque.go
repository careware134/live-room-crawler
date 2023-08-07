package common

type PlayDeque struct {
	maxSize int
	deque   []UserActionEvent
}

func NewFixedSizeDeque(maxSize int) *PlayDeque {
	return &PlayDeque{
		maxSize: maxSize,
		deque:   make([]UserActionEvent, 0, maxSize),
	}
}

func (d *PlayDeque) PushFront(item UserActionEvent) {
	if len(d.deque) >= d.maxSize {
		d.deque = d.deque[:d.maxSize-1]
	}
	d.deque = append([]UserActionEvent{item}, d.deque...)
}

func (d *PlayDeque) PushBack(item UserActionEvent) {
	if d.deque == nil {

	}
	if len(d.deque) >= d.maxSize {
		d.deque = d.deque[1:]
	}
	d.deque = append(d.deque, item)
}

func (d *PlayDeque) PopFront() *UserActionEvent {
	if len(d.deque) == 0 {
		return nil
	}
	item := d.deque[0]
	d.deque = d.deque[1:]
	return &item
}

func (d *PlayDeque) PopBack() *UserActionEvent {
	if len(d.deque) == 0 {
		return nil
	}
	lastIndex := len(d.deque) - 1
	item := d.deque[lastIndex]
	d.deque = d.deque[:lastIndex]
	return &item
}

func (d *PlayDeque) Size() int {
	return len(d.deque)
}

func (d *PlayDeque) IsEmpty() bool {
	return len(d.deque) == 0
}

func (d *PlayDeque) IsFull() bool {
	return len(d.deque) == d.maxSize
}

func (d *PlayDeque) Clear() {
	d.deque = make([]UserActionEvent, 0, d.maxSize)
}
