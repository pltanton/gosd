package notifier

import (
	"errors"
	"reflect"

	"github.com/TheCreeper/go-notify"

	"github.com/pltanton/gosd/core"
)

type notifier struct {
	subscriptions []core.Listener
	started       bool
}

func NewNotifier() *notifier {
	return &notifier{
		subscriptions: make([]core.Listener, 0),
		started:       false,
	}
}

func (n *notifier) Start() error {
	if n.started {
		return errors.New("Can't start notifier, because it already started")
	}
	n.started = true
	cases := make([]reflect.SelectCase, len(n.subscriptions))
	for i, listener := range n.subscriptions {
		go listener.StartMonitor()
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(listener.Chan())}
	}
	var lastId uint32 = 0
	for {
		_, value, _ := reflect.Select(cases)
		message := value.Interface().(core.NotificationMessage)
		notification := notify.NewNotification(message.Title, message.Message)
		notification.ReplacesID = lastId
		notification.AppIcon = message.Icon
		lastId, _ = notification.Show()
	}
	return nil
}

func (n *notifier) Subscribe(listener core.Listener) error {
	if n.started {
		return errors.New("Can't subscribe anything to notifier because it already started")
	}
	n.subscriptions = append(n.subscriptions, listener)
	return nil
}
