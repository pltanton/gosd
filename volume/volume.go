package volume

import (
	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"

	"github.com/pltanton/gosd/core"
)

type volume struct {
	out chan core.NotificationMessage
}

func (v volume) Chan() chan core.NotificationMessage {
	return v.out
}

func NewVolumeListener() volume {
	return volume{
		out: make(chan core.NotificationMessage, 1),
	}
}

func (v volume) StartMonitor() {
	client := newPulseClient()

	// TODO: put sink inside some config file
	sink := "/org/pulseaudio/core1/sink0"
	for range client.event {
		message := formatVolume(client.getVolume(sink))
		v.out <- message
	}
}

func formatVolume(volume int, muted bool) core.NotificationMessage {
	var icon string
	if muted {
		icon = "audio-off"
	} else {
		icon = "audio-on"
	}
	message := core.RenderBar(volume)
	return core.NotificationMessage{
		Title:   "",
		Message: message,
		Icon:    icon,
	}
}

type pulseClient struct {
	*pulseaudio.Client
	event chan bool
}

func (pc *pulseClient) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	pc.event <- true
}

func (pc *pulseClient) DeviceMuteUpdated(path dbus.ObjectPath, values bool) {
	pc.event <- true
}

func (pc *pulseClient) getVolume(sink string) (vol int, mute bool) {
	dev := pc.Device(dbus.ObjectPath(sink))

	mute, _ = dev.Bool("Mute")
	vols, _ := dev.ListUint32("Volume")
	vol = int(float64(volAvg(vols)*100)/65535. + .5)

	return
}

func volAvg(vols []uint32) (vol uint32) {
	if l := len(vols); l > 0 {
		for _, v := range vols {
			vol += v
		}
		vol /= uint32(l)
	}
	return vol
}

func newPulseClient() *pulseClient {
	pulse, err := pulseaudio.New()
	if err != nil {
		panic(err.Error())
	}

	client := &pulseClient{pulse, make(chan bool, 1)}
	pulse.Register(client)

	go pulse.Listen()

	return client
}
