/*
 * @Author: YanHui Li
 * @Date: 2022-02-10 15:40:55
 * @LastEditTime: 2022-02-15 15:49:33
 * @LastEditors: YanHui Li
 * @Description:
 * @FilePath: \go-onvif\examples\SetDevice\test.go
 *
 */
package main

import (
	"log"

	"github.com/liyanhui1998/go-onvif"
	"github.com/liyanhui1998/go-onvif/types/device"
)

func main() {
	/* connect device */
	dev, _ := onvif.NewDevice(onvif.DeviceParams{Ipddr: "192.168.1.188:8000", Username: "admin", Password: "admin"})

	/* sync time device */
	err := dev.CallMethodInterface(
		device.SetSystemDateAndTime{
			DateTimeType: "Manual", DaylightSavings: false,
			TimeZone:    device.TimeZone{TZ: "CST-8"},
			UTCDateTime: device.DateTime{Time: device.Time{Hour: 12, Minute: 0, Second: 0}, Date: device.Date{Year: 2020, Month: 2, Day: 1}},
		}, nil, "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	/* reboot device */
	reboot := device.SystemRebootResponse{}
	err = dev.CallMethodInterface(device.SystemReboot{}, &reboot, "")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(reboot)
}
