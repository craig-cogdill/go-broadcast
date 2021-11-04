package broadcast

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Need to initialize the map, as using the New() constructor
//	would not allow access to 'subscribers' member data
func getDefaultBroadcaster() *broadcaster {
	return &broadcaster{
		subscribers: make(map[int]broadcastQueue),
	}
}

func Test_Constructor(t *testing.T) {
	t.Run("verify constructor returns proper type", func(t *testing.T) {
		assert := assert.New(t)
		testBroadcaster := New()
		assert.IsType(&broadcaster{}, testBroadcaster)
	})
}

func Test_Subscriber(t *testing.T) {
	testId := 2112

	t.Run("can query ID", func(t *testing.T) {
		assert := assert.New(t)

		testSubscriber := Subscription{
			id:    testId,
			queue: nil,
		}
		assert.Equal(testId, testSubscriber.ID())
	})

	t.Run("can fetch input channel", func(t *testing.T) {
		assert := assert.New(t)

		testChan := make(subscriberQueue)
		testSubscriber := Subscription{
			id:    testId,
			queue: testChan,
		}
		assert.Equal(testChan, testSubscriber.Queue())
	})
}

func Test_Subscribe(t *testing.T) {
	t.Run("subscribing creates new channel", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.Empty(testBroadcaster.subscribers)

		_ = testBroadcaster.Subscribe()
		assert.Equal(1, len(testBroadcaster.subscribers))
	})

	t.Run("subscribing assigns id", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.Empty(testBroadcaster.subscribers)

		subscriber := testBroadcaster.Subscribe()
		assert.Equal(0, subscriber.ID())
	})

	t.Run("sequential subscribes get sequential ids", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.Empty(testBroadcaster.subscribers)

		firstSubscriber := testBroadcaster.Subscribe()
		assert.Equal(0, firstSubscriber.ID())

		secondSubscriber := testBroadcaster.Subscribe()
		assert.Equal(1, secondSubscriber.ID())
	})
}

func Test_Close(t *testing.T) {
	t.Run("closing nils out the subscriber queues - no entries", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.NotNil(testBroadcaster.subscribers)
		assert.Len(testBroadcaster.subscribers, 0)

		testBroadcaster.Close()
		assert.Nil(testBroadcaster.subscribers)
	})

	t.Run("closing nils out the subscriber queues - one entry", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.NotNil(testBroadcaster.subscribers)

		_ = testBroadcaster.Subscribe()
		assert.Len(testBroadcaster.subscribers, 1)

		testBroadcaster.Close()
		assert.Nil(testBroadcaster.subscribers)
	})

	t.Run("closing nils out the subscriber queues - many entries", func(t *testing.T) {
		assert := assert.New(t)

		numberOfSubscribers := 100

		testBroadcaster := getDefaultBroadcaster()
		assert.NotNil(testBroadcaster.subscribers)

		for i := 0; i < numberOfSubscribers; i++ {
			_ = testBroadcaster.Subscribe()
		}
		assert.Len(testBroadcaster.subscribers, numberOfSubscribers)

		testBroadcaster.Close()
		assert.Nil(testBroadcaster.subscribers)
	})
}

func Test_Broadcast(t *testing.T) {
	t.Run("does not panic when no subscribers", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := New()
		defer testBroadcaster.Close()

		assert.NotPanics(func() { testBroadcaster.Broadcast("don't panic!") })
	})

	t.Run("sends a message to one subscriber", func(t *testing.T) {
		assert := assert.New(t)

		const expectedMsg = "Hello World"

		testBroadcaster := New()
		defer testBroadcaster.Close()

		var threadFinished sync.WaitGroup
		threadFinished.Add(1)

		var threadReady sync.WaitGroup
		threadReady.Add(1)

		var receivedMsg interface{}
		go func() {
			subscription := testBroadcaster.Subscribe()
			threadReady.Done()
			receivedMsg = <-subscription.Queue()
			threadFinished.Done()
		}()

		threadReady.Wait()
		testBroadcaster.Broadcast(expectedMsg)
		threadFinished.Wait()

		assert.Equal(expectedMsg, receivedMsg.(string))
	})

	t.Run("sends a message to many subscribers", func(t *testing.T) {
		assert := assert.New(t)

		const expectedMsg = "Hello World"
		const numberOfSubscribers = 10

		testBroadcaster := New()
		defer testBroadcaster.Close()

		var threadsFinished sync.WaitGroup
		threadsFinished.Add(numberOfSubscribers)

		var threadsReady sync.WaitGroup
		threadsReady.Add(numberOfSubscribers)

		// normal maps are not thread-safe
		var receivedMsgs sync.Map
		for i := 0; i < numberOfSubscribers; i++ {
			go func() {
				subscription := testBroadcaster.Subscribe()
				threadsReady.Done()
				msg := <-subscription.Queue()
				receivedMsgs.Store(subscription.ID(), msg)
				threadsFinished.Done()
			}()
		}

		threadsReady.Wait()
		testBroadcaster.Broadcast(expectedMsg)
		threadsFinished.Wait()

		receivedMsgsCount := 0
		allSuccess := false
		receivedMsgs.Range(func(key, value interface{}) bool {
			receivedMsgsCount += 1
			messageMatchesExpected := value.(string) == expectedMsg
			allSuccess = allSuccess || messageMatchesExpected
			return messageMatchesExpected
		})

		assert.True(allSuccess)
		assert.Equal(numberOfSubscribers, receivedMsgsCount)
	})
}

func Test_Unsubscribe(t *testing.T) {
	t.Run("unsubscribing closes a channel", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.Empty(testBroadcaster.subscribers)

		testId := 2112
		testChannel := make(chan interface{})
		testBroadcaster.m.Lock()
		testBroadcaster.subscribers[testId] = testChannel
		testBroadcaster.m.Unlock()

		testBroadcaster.Unsubscribe(testId)
		assert.Empty(testBroadcaster.subscribers)
	})

	t.Run("unsubscribing leaves other subscriptions untouched", func(t *testing.T) {
		assert := assert.New(t)

		testBroadcaster := getDefaultBroadcaster()
		assert.Empty(testBroadcaster.subscribers)

		testId1 := 1
		testId2 := 2
		testChannel1 := make(chan interface{})
		testChannel2 := make(chan interface{})
		testBroadcaster.m.Lock()
		testBroadcaster.subscribers[testId1] = testChannel1
		testBroadcaster.subscribers[testId2] = testChannel2
		testBroadcaster.m.Unlock()
		assert.Equal(2, len(testBroadcaster.subscribers))

		testBroadcaster.Unsubscribe(testId1)
		assert.Equal(1, len(testBroadcaster.subscribers))
		_, ok := testBroadcaster.subscribers[testId2]
		assert.True(ok)
	})

}
