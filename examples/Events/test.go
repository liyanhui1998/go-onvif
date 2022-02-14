package main

import (
	"log"
	"time"

	"github.com/liyanhui1998/go-onvif"
	event "github.com/liyanhui1998/go-onvif/types/events"
)

func main() {
	dev, err := onvif.NewDevice(onvif.DeviceParams{Ipddr: "10.1.1.200", Username: "admin", Password: "123qweasdZXC"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	subres := event.CreatePullPointSubscriptionResponse{}
	err = dev.CallMethodInterface(event.CreatePullPointSubscription{}, &subres, "")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Print(subres)

	for i := 0; i < 10; i++ {
		pull := event.PullMessagesResponse{}
		err := dev.CallMethodInterface(event.PullMessages{Timeout: "PT10S", MessageLimit: 2}, &pull, string(subres.SubscriptionReference.Address))
		if err != nil {
			log.Fatalf(err.Error())
		}
		if pull.NotificationMessage.Topic.TopicKinds != "" {
			log.Printf("%v", pull.NotificationMessage)
		}
		time.Sleep(100 * time.Millisecond)
	}

	err = dev.CallMethodInterface(event.Unsubscribe{}, event.UnsubscribeResponse{}, string(subres.SubscriptionReference.Address))
	if err != nil {
		log.Fatalf(err.Error())
	}
}
