package broadcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

		testSubscriber := Subscriber{
			id:    testId,
			queue: nil,
		}
		assert.Equal(testId, testSubscriber.ID())
	})

	t.Run("can fetch input channel", func(t *testing.T) {
		assert := assert.New(t)

		testChan := make(subscriberQueue)
		testSubscriber := Subscriber{
			id:    testId,
			queue: testChan,
		}
		assert.Equal(testChan, testSubscriber.Queue())
	})
}

func Test_AddSubscriber(t *testing.T) {

	// Need to initialize the map, as using the New() constructor
	//	would not allow access to 'subscribers' member data
	getDefaultBroadcaster := func() *broadcaster {
		return &broadcaster{
			subscribers: make(map[int]broadcastQueue),
		}
	}

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
