package main

import (
	"log"

	"github.com/pltanton/gosd/notifier"
	"github.com/pltanton/gosd/volume"
)

func main() {
	n := notifier.NewNotifier()
	n.Subscribe(volume.NewVolumeListener())
	err := n.Start()
	if err != nil {
		log.Panicln(err)
	}
}
