package broadcast

type broadcastQueue chan<- interface{}
type subscriberQueue <-chan interface{}

type Broadcaster interface {
	Subscribe() *Subscriber
	Close()
}

type broadcaster struct {
	subscribers map[int]broadcastQueue
}

type Subscriber struct {
	id    int
	queue subscriberQueue
}

func (s *Subscriber) ID() int {
	return s.id
}

func (s *Subscriber) Queue() subscriberQueue {
	return s.queue
}

func New() Broadcaster {
	return &broadcaster{
		subscribers: make(map[int]broadcastQueue),
	}
}

func (b *broadcaster) Subscribe() *Subscriber {
	newSubscriberChan := make(chan interface{})
	newId := len(b.subscribers)
	b.subscribers[newId] = newSubscriberChan
	return &Subscriber{
		id:    newId,
		queue: newSubscriberChan,
	}
}

func (b *broadcaster) Close() {
	for _, channel := range b.subscribers {
		close(channel)
	}
	b.subscribers = nil
}
