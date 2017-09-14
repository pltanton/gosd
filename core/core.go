package core

import (
	"strings"
)

const barChar string = "━"
const halfBarChar string = "╸"

type Listener interface {
	StartMonitor()
	Chan() chan NotificationMessage
}

type NotificationMessage struct {
	Title   string
	Message string
	Icon    string
}

func RenderBar(volume int) string {
	// TODO: get value from configuration file
	maxLength := 30
	actualLength := int(float64(maxLength*2)*float64(volume)/100 + .5)
	return strings.Repeat(barChar, actualLength/2) + strings.Repeat(halfBarChar, actualLength%2)
}
