package bullhorn_test

import (
	. "bullhorn"
	"fmt"
	"time"
)

func ExampleAdd() {
	// Register my inbound channel against an interesting feed
	interestingFeed := make(chan interface{}, 1024)
	Add(Subscription{"Interesting Feed", interestingFeed})
}

func ExampleDelete() {
	interestingFeed := make(chan interface{}, 1)

	// Remove an existingsubscription
	Delete(Subscription{"Interesting Feed", interestingFeed})
}

func ExamplePublish() {
	// Publish an event with asc data to a key
	Publish(Event{"Interesting Feed", "Interesting Data"})
}

func ExampleStart() {
	// Start the service with buffered subscription and event channels
	Start(1024, 2048)
}

func ExampleSubscribed() {
	ch := make(chan interface{}, 1)

	Add(Subscription{"order book prices", ch})
	if Subscribed("order book prices", ch) {
		// all set, off we go
	}

}

func Example() {

	type A struct {
		AField int
	}

	type B struct {
		BField string
	}

	inbound := make(chan interface{}, 16)
	Start(4, 16)

	// start the consumer
	go func() {

		for data := range inbound {

			fmt.Printf("data:%+v ", data)

			switch data := data.(type) {
			case A:
				fmt.Printf("A:%+v\n", data.AField)
			case B:
				fmt.Printf("B:%+v\n", data.BField)
			}

		}

	}()

	// start the producer
	go func() {
		i := 0
		for {
			time.Sleep(1e8)

			Publish(Event{"A feed", A{i}})
			Publish(Event{"B feed", B{fmt.Sprintf("%v", i)}})
			i++
		}
	}()

	// subscribe to A feed
	Add(Subscription{"A feed", inbound})

	time.Sleep(2e9)

	// subscribe to B feed
	Add(Subscription{"B feed", inbound})

	time.Sleep(2e9)

	// delete A feed subscription
	Delete(Subscription{"A feed", inbound})

	time.Sleep(2e9)

	// Start over
	Start(4, 16)

	time.Sleep(2e9)

	// subscribe to B feed again
	Add(Subscription{"B feed", inbound})

	time.Sleep(2e9)

	// enough
	return

}
