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
