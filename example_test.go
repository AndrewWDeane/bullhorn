package bullhorn_test

func ExampleAdd() {
	// Register my inbound channel against an interesting feed
	interestingFeed := make(chan interface{}, 1024)
	Add(Subscription{"Interesting Feed", interestingFeed})
}

func ExampleDelete() {
	// Remove an existingsubscription
	Delete(Subscription{"Interesting Feed", interestingFeed})
}

func ExamplePublish() {
	// Publish an event with asc data to a key
	Delete(Event{"Interesting Feed", "Interesting Data"})
}

func ExampleStart() {
	// Start the service with buffered subscription and event channels
	Start(1024, 2048)
}

func Example() {

	type A struct {
		AField int
	}

	type B struct {
		BField string
	}

	inbound := make(chan interface{}, 16)
	bullhorn.Start(4, 16)

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

			bullhorn.Publish(bullhorn.Event{"A feed", A{i}})
			bullhorn.Publish(bullhorn.Event{"B feed", B{fmt.Sprintf("%v", i)}})
			i++
		}
	}()

	// subscribe to A feed
	bullhorn.Add(bullhorn.Subscription{"A feed", inbound})

	time.Sleep(2e9)

	// subscribe to B feed
	bullhorn.Add(bullhorn.Subscription{"B feed", inbound})

	time.Sleep(2e9)

	// delete A feed subscription
	bullhorn.Delete(bullhorn.Subscription{"A feed", inbound})

	time.Sleep(2e9)

	// Start over
	bullhorn.Start(4, 16)

	time.Sleep(2e9)

	// subscribe to B feed again
	bullhorn.Add(bullhorn.Subscription{"B feed", inbound})

	time.Sleep(2e9)

	// enough
	return

}