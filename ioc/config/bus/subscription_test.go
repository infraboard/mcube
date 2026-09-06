package bus

import (
	"errors"
	"testing"
)

func TestFuncSubscriptionUnsubscribe(t *testing.T) {
	called := false
	sub := FuncSubscription(func() error {
		called = true
		return nil
	})
	if err := sub.Unsubscribe(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("expected unsubscribe fn to be called")
	}
}

func TestFuncSubscriptionNil(t *testing.T) {
	var sub FuncSubscription
	if err := sub.Unsubscribe(); err != nil {
		t.Fatal(err)
	}
}

func TestFuncSubscriptionError(t *testing.T) {
	want := errors.New("closed")
	sub := FuncSubscription(func() error { return want })
	if err := sub.Unsubscribe(); err != want {
		t.Fatalf("got %v want %v", err, want)
	}
}
