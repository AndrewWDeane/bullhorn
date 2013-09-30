// The bullhorn package provides a lightweight internal pub/sub mechanism for goroutines
package bullhorn

import (
	"errors"
)

var (
	subscriptionTable map[interface{}]map[chan interface{}]chan interface{}
	subscriptions     chan admin
	events            chan Event
)

// internal type to handle the addition, deletion, and querying of subscriptions
type admin struct {
	action       operation
	subscription Subscription
	respondOn    chan bool
}

// internal admin operations
type operation int

const (
	add operation = iota
	del
	query
)

// The Event struct holds the key and data to be published
type Event struct {
	Key  interface{}
	Data interface{}
}

// The Subscription struct holds the key to subscribe to and the channel
// the subscriber will be listening on
type Subscription struct {
	Key      interface{}
	Receiver chan interface{}
}

// Start the publish and subscribe service
func Start(subscriptionBuffer, eventBuffer int64) {

	// set up the subscription map
	subscriptionTable = make(map[interface{}]map[chan interface{}]chan interface{})

	// and the internal channels
	subscriptions = make(chan admin, subscriptionBuffer)
	events = make(chan Event, eventBuffer)

	// start main processing
	go process()

	return
}

// process selects on the subscription and event channels,
// administering the subscriptions and forwarding published events
func process() {

	var subMap map[chan interface{}]chan interface{}
	var ok bool

	for {

		select {

		case a := <-subscriptions:

			switch a.action {

			case add:
				subMap, ok = subscriptionTable[a.subscription.Key]

				if !ok {
					subMap = make(map[chan interface{}]chan interface{})
				}

				subMap[a.subscription.Receiver] = a.subscription.Receiver

				subscriptionTable[a.subscription.Key] = subMap

			case del:
				delete(subscriptionTable[a.subscription.Key], a.subscription.Receiver)

			case query:
				subMap, ok = subscriptionTable[a.subscription.Key]

				if !ok {
					a.respondOn <- false
				}

				_, ok = subMap[a.subscription.Receiver]
				a.respondOn <- ok

			}

		case e := <-events:

			if subs, ok := subscriptionTable[e.Key]; ok {

				for _, sub := range subs {

					// in a seperate routine so that slow consumers don't block others
					go func(ch chan interface{}, i interface{}) {
						ch <- i
					}(sub, e.Data)

				}

			}

		}
	}

}

// Add a new subscription
func Add(s Subscription) error {

	if subscriptions == nil {
		return errors.New("Error Adding. Have you Started()?")
	}

	subscriptions <- admin{action: add, subscription: s}
	return nil
}

// Delete an existing subscription
func Delete(s Subscription) error {

	if subscriptions == nil {
		return errors.New("Error Deleting. Have you Started()?")
	}

	subscriptions <- admin{action: del, subscription: s}
	return nil
}

// Publish an event
func Publish(e Event) error {

	if events == nil {
		return errors.New("Error Publishing. Have you Started()?")
	}

	events <- e
	return nil

}

// Make a blocking query returning true if a channel is subscribed to a key
func Subscribed(key interface{}, receiver chan interface{}) bool {

	if subscriptions == nil {
		return false
	}

	response := make(chan bool, 1)

	subscriptions <- admin{action: query, subscription: Subscription{Key: key, Receiver: receiver}, respondOn: response}
	return <-response

}
