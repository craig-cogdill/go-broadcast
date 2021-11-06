package broadcast

import "sync"

type broadcastQueue chan<- interface{}
type subscriberQueue <-chan interface{}

type Broadcaster interface {
	Subscribe() Subscription
	Broadcast(msg interface{})
	Close()
}

type Subscription interface {
	ID() int
	Queue() subscriberQueue
	Unsubscribe()
}

type broadcaster struct {
	m           sync.Mutex
	subscribers map[int]broadcastQueue
}

type subscription struct {
	id          int
	queue       subscriberQueue
	unsubscribe func(int)
	once        sync.Once
}

func (s *subscription) ID() int {
	return s.id
}

func (s *subscription) Queue() subscriberQueue {
	return s.queue
}

func (s *subscription) Unsubscribe() {
	s.once.Do(func() {
		s.unsubscribe(s.id)
	})
}

func New() Broadcaster {
	return &broadcaster{
		subscribers: make(map[int]broadcastQueue),
	}
}

func (b *broadcaster) Subscribe() Subscription {
	b.m.Lock()
	defer b.m.Unlock()
	newSubscriberChan := make(chan interface{})
	newId := len(b.subscribers)
	b.subscribers[newId] = newSubscriberChan
	return &subscription{
		id:          newId,
		queue:       newSubscriberChan,
		unsubscribe: b.unsubscribe,
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

func (b *broadcaster) unsubscribe(id int) {
	b.m.Lock()
	defer b.m.Unlock()
	channel, ok := b.subscribers[id]
	if ok {
		close(channel)
		delete(b.subscribers, id)
	}
}
