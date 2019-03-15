package eventbus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/imkira/go-observer"
)

type Publisher interface {
	Publish(evt interface{})
}

type EventHandler interface {
	// Stringer interface used to store the name of the event this subscriber is listening to.
	fmt.Stringer

	// ProcessData is the handler used to processed the data.
	ProcessData(data interface{}) error
}

// Subscriber interface which will be used to directly handle the data coming from the observer.
type Subscriber interface {
	// Stringer interface used to store the name of the event this subscriber is listening to.
	fmt.Stringer

	// Subscribe accepts observer.Stream and the handler function which will be used to pull and process
	// the incoming stream
	Subscribe(stream observer.Stream) error

	// Unsubscribe stops the subscription.
	Unsubscribe() error
}

// NewSubscriber creates new StandardSubscriber
func NewSubscriber(ctx context.Context, handler EventHandler) *StandardSubscriber {
	return &StandardSubscriber{
		Ctx:     ctx,
		Handler: handler,
	}
}

// StandardSubscriber is the default subscriber which can be used to capture any event published by the publisher.
// This struct requires EventName and Fn to be supplied. EventName is the name of the event that subscriber will
// listen to. Fn is the function to handle every message from the subscriber.
type StandardSubscriber struct {
	Subscriber
	// Context
	Ctx context.Context
	// Handler contains the data processor and the event name this subscriber should care about.
	Handler EventHandler
	// channel used to notify the end of subscription.
	exit chan bool
}

// String return the name of the event
func (s *StandardSubscriber) String() string {
	return s.Handler.String()
}

// Subscribe listen directly to the go-observer stream.
func (s *StandardSubscriber) Subscribe(stream observer.Stream) error {
breakLoop:
	for {
		if err := s.Handler.ProcessData(stream.Value()); err != nil {
			fmt.Printf("Got error : %s\n", err)
		}
		select {
		case _, isOpen := <-s.exit:
			if isOpen {
				close(s.exit)
			}
			break breakLoop
		case <-stream.Changes():
			stream.Next()
		}
	}
	return nil
}

// Unsubscribe will stop the subscription.
func (s *StandardSubscriber) Unsubscribe() error {
	if s.exit == nil {
		s.exit = make(chan bool)
	}
	s.exit <- true
	return nil
}

var (
	// map of the eventname and the channel listening to the event
	observers map[string]observer.Property

	// RWmutex for making sure there is no multiple routine accessing the observers map
	mu *sync.RWMutex
)

func init() {
	mu = &sync.RWMutex{}
	observers = make(map[string]observer.Property)
}

// Attach creates a goroutine for every subscriber and
func Attach(subs Subscriber) {
	mu.Lock()
	defer mu.Unlock()

	eventName := subs.String()
	if observers[eventName] == nil {
		//Initialize the property will empty struct
		observers[eventName] = observer.NewProperty(struct{}{})
	}

	go subs.Subscribe(observers[eventName].Observe())
}

// Detach detaches the subscriber from the publisher.
func Detach(subs Subscriber) {
	subs.Unsubscribe()
}

// Notify send evt to all channels/subscribers registered to a specific event
func Publish(evt interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	eventName := getType(evt)
	prop := observers[eventName]
	if prop != nil {
		prop.Update(evt)
	}
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
