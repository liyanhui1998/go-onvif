package main

import (
	"log"

	"github.com/liyanhui1998/go-onvif"
)

func main() {
	/* 扫描局域网内摄像头 */
	deviceList := onvif.GetAvailableDevicesAtSpecificEthernetInterface("以太网 2")
	for _, value := range deviceList {
		log.Printf("ip : %s\r\n", value.Params.Ipddr)
		log.Printf("mac : %s\r\n", value.Params.MAC)
		log.Printf("uuid : %s\r\n", value.Params.Uuid)
		log.Printf("name : %s\r\n", value.Params.Name)
		log.Printf("types : %s\r\n", value.Params.Types)
		log.Printf("model : %s\r\n", value.Params.Model)
	}

}
