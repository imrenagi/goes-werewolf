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

type Subscriber interface {
	HandleEvent(interface{}) error
	SubscribeToEventType() reflect.Type
}

//NewSubscriber creates new defaultSubscriber
func newSubscriber(ctx context.Context, subscriber Subscriber) *defaultSubscriber {
	sub := defaultSubscriber{
		Ctx:        ctx,
		Subscriber: subscriber,
	}
	return &sub
}

// listener interface which will be used to directly handle the data coming from the observer.
type listener interface {
	// Listen accepts observer.Stream and the handler function which will be used to pull and process
	// the incoming stream
	Listen(stream observer.Stream) error

	// Unsubscribe stops the subscription.
	Unsubscribe() error
}

// defaultSubscriber is the default listener which can be used to capture any event published by the publisher.
// This struct requires EventName and Fn to be supplied. EventName is the name of the event that listener will
// listen to. Fn is the function to handle every message from the listener.
type defaultSubscriber struct {
	listener

	// Handler contains the data processor and the event name this listener should care about.
	Subscriber

	// Context
	Ctx context.Context

	// channel used to notify the end of subscription.
	exit chan bool
}

// Listen listen directly to the go-observer stream.
func (s *defaultSubscriber) Listen(stream observer.Stream) error {
breakLoop:
	for {
		if err := s.HandleEvent(stream.Value()); err != nil {
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
func (s *defaultSubscriber) Unsubscribe() error {
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

//Attach creates a goroutine for every listener and
func Attach(subscriber Subscriber) {
	s := newSubscriber(context.Background(), subscriber)
	subscribeToEventType := s.SubscribeToEventType().String()

	mu.Lock()
	defer mu.Unlock()

	if observers[subscribeToEventType] == nil {
		//Initialize the property will empty struct
		observers[subscribeToEventType] = observer.NewProperty(struct{}{})
	}

	go s.Listen(observers[subscribeToEventType].Observe())
}

// Detach detaches the listener from the publisher.
func Detach(subs listener) {
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
