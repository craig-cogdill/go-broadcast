package broadcast

import "sync"

type broadcastQueue chan<- interface{}
type subscriberQueue <-chan interface{}

type Broadcaster interface {
	Subscribe() *Subscription
	Broadcast(msg interface{})
	Close()
}

type broadcaster struct {
	m           sync.Mutex
	subscribers map[int]broadcastQueue
}

type Subscription struct {
	id    int
	queue subscriberQueue
}

func (s *Subscription) ID() int {
	return s.id
}

func (s *Subscription) Queue() subscriberQueue {
	return s.queue
}

func New() Broadcaster {
	return &broadcaster{
		subscribers: make(map[int]broadcastQueue),
	}
}

func (b *broadcaster) Subscribe() *Subscription {
	newSubscriberChan := make(chan interface{})
	newId := len(b.subscribers)
	b.subscribers[newId] = newSubscriberChan
	return &Subscription{
		id:    newId,
		queue: newSubscriberChan,
	}
}

func (b *broadcaster) Broadcast(msg interface{}) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.subscribers == nil || len(b.subscribers) == 0 {
		return
	}
	for _, subscriber := range b.subscribers {
		subscriber <- msg
	}
}

func (b *broadcaster) Close() {
	for _, channel := range b.subscribers {
		close(channel)
	}
	b.subscribers = nil
}

func (b *broadcaster) Unsubscribe(id int) {
	b.m.Lock()
	defer b.m.Unlock()
	channel, ok := b.subscribers[id]
	if ok {
		close(channel)
		delete(b.subscribers, id)
	}
}
