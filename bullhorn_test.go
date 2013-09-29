package bullhorn

import (
	"testing"
)

// if you Add before starting we should get an error
func TestAddNotStarted(t *testing.T) {
	checkError(t, Add(Subscription{"A", make(chan interface{})}))
}

// if you Delete before starting we should get an error
func TestDeleteNotStarted(t *testing.T) {
	checkError(t, Delete(Subscription{"A", make(chan interface{})}))
}

// if you Publish before starting we should get an error
func TestPublisheNotStarted(t *testing.T) {
	checkError(t, Publish(Event{"A", "data"}))
}

func checkError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Failed to raise error")
	}
}

func TestSubscribedNotStarted(t *testing.T) {
	// before we start, Subscribed will always return false
	if Subscribed("A", make(chan interface{})) {
		t.Errorf("Error, subscribed")
	}
}

// Test that Add does just that
func TestAdd(t *testing.T) {
	ch := make(chan interface{}, 2)
	Start(2, 2)

	sub := Subscription{"A", ch}
	Add(sub)

	// this will block for a response, meaning any Adds will have been actioned
	// we wont use the result but check independantly
	Subscribed("A", ch)

	subMap, ok := subscriptionTable[sub.Key]

	if !ok {
		t.Errorf("Failed to add")
	}

	_, ok = subMap[sub.Receiver]

	if !ok {
		t.Errorf("Failed to add")
	}

}

// Test that Delete does just that
func TestDelete(t *testing.T) {
	ch := make(chan interface{}, 2)
	Start(2, 2)

	// add it first
	sub := Subscription{"A", ch}
	Add(sub)
	Subscribed("A", ch)

	subMap, ok := subscriptionTable[sub.Key]

	if !ok {
		t.Errorf("Failed to add (del)")
	}

	_, ok = subMap[sub.Receiver]

	if !ok {
		t.Errorf("Failed to add (del)")
	}

	Delete(sub)
	Subscribed("A", ch)
	subMap, ok = subscriptionTable[sub.Key]

	if ok {
		_, ok = subMap[sub.Receiver]
		if ok {
			t.Errorf("Failed to delete")
		}
	}

}

// Test that Delete caters for not present keys
func TestDeleteNotPresentKey(t *testing.T) {
	ch := make(chan interface{}, 2)
	Start(2, 2)

	// add it first
	sub := Subscription{"A", ch}
	sub1 := Subscription{"B", ch}
	Add(sub)
	Subscribed("A", ch)

	// delete B
	Delete(sub1)
	Subscribed("A", ch)

	// make sure A is still present
	subMap, ok := subscriptionTable[sub.Key]

	if !ok {
		t.Errorf("Deleted wrong key")
	}

	_, ok = subMap[sub.Receiver]

	if !ok {
		t.Errorf("Deleted wrong chan")
	}

	// and that B isnt
	subMap, ok = subscriptionTable[sub1.Key]

	if ok {
		t.Errorf("Key is present")
	}

}

// Test that Delete caters for not present chans
func TestDeleteNotChanKey(t *testing.T) {
	ch := make(chan interface{}, 2)
	ch1 := make(chan interface{}, 2)
	Start(2, 2)

	// add it first
	sub := Subscription{"A", ch}
	sub1 := Subscription{"A", ch1}
	Add(sub)
	Subscribed("A", ch)

	// delete ch1
	Delete(sub1)
	Subscribed("A", ch)

	// make sure sub is still present
	subMap, ok := subscriptionTable[sub.Key]

	if !ok {
		t.Errorf("Deleted wrong key")
	}

	_, ok = subMap[sub.Receiver]

	if !ok {
		t.Errorf("Deleted wrong chan")
	}

	// and that sub1 isnt
	_, ok = subMap[sub1.Receiver]

	if ok {
		t.Errorf("Chan still present")
	}

}

// test the subscribed returns the correct value
func TestSubscribed(t *testing.T) {

	ch := make(chan interface{}, 2)
	ch1 := make(chan interface{}, 2)
	Start(2, 2)

	// add it first
	sub := Subscription{"A", ch}

	Add(sub)

	if !Subscribed("A", ch) {
		t.Errorf("Not subscribed")
	}

	if Subscribed("B", ch) {
		t.Errorf("Subscribed")
	}

	if Subscribed("B", ch1) {
		t.Errorf("Subscribed")
	}

}

// Test that once subscribed, we get published messages
func TestPublish(t *testing.T) {
	ch := make(chan interface{}, 2)
	Start(2, 2)

	sub := Subscription{"A", ch}
	Add(sub)
	Subscribed("A", ch)

	Publish(Event{"A", "data"})

	data := <-ch
	if data.(string) != "data" {
		t.Errorf("Publish error")
	}
}
