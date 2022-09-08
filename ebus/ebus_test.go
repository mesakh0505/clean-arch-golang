package ebus_test

import (
	"context"
	"fmt"
	"time"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
)

func ExampleBus() {
	bus := new(ebus.Bus)
	handleFunc := func(ctx context.Context, e ebus.Event) {
		fmt.Printf("name: '%s', body: '%v'\n", e.Name, e.Data)
	}
	bus.Subscribe(ebus.HandlerFunc(handleFunc))
	ctx := context.Background()
	bus.Publish(ctx, ebus.Event{
		Name:        "Greet",
		Data:        "Hello World!",
		OccuredTime: time.Now(),
	})

	bus.Publish(ctx, ebus.Event{
		Name:        "Asked",
		Data:        "How are you?",
		OccuredTime: time.Now(),
	})

	// output:
	// name: 'Greet', body: 'Hello World!'
	// name: 'Asked', body: 'How are you?'
}
