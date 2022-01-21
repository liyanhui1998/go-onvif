package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/liyanhui1998/go-onvif"
	"github.com/liyanhui1998/go-onvif/types/device"
	"github.com/liyanhui1998/go-onvif/types/hikvision"
	"github.com/liyanhui1998/go-onvif/types/media"
)

func main222() {
	/* 扫描局域网内摄像头 */
	deviceList := onvif.GetAvailableDevicesAtSpecificEthernetInterface("以太网")
	for _, value := range deviceList {
		log.Printf("scan IP:%s\r\n", value.Params.Ipddr)
	}

	/* 连接设备 */
	dev, _ := onvif.NewDevice(onvif.DeviceParams{Ipddr: "10.1.1.200", Username: "admin", Password: "123qweasdZXC"})

	/* 获取能力集合 */
	retServices := dev.GetServices()
	for key, value := range retServices {
		log.Printf("server key : %s value : %s\r\n", key, value)
	}

	/* 获取设备基本信息 */
	deviceInfo := device.GetDeviceInformationResponse{}
	dev.CallMethodInterface(device.GetDeviceInformation{}, &deviceInfo)
	jsonString, _ := json.Marshal(deviceInfo)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备支持网络协议信息 */
	deviceNetProtocols := device.GetNetworkProtocolsResponse{}
	dev.CallMethodInterface(device.GetNetworkProtocols{}, &deviceNetProtocols)
	jsonString, _ = json.Marshal(deviceNetProtocols)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备视频资源信息 */
	mediaVideoSour := media.GetVideoSourcesResponse{}
	dev.CallMethodInterface(media.GetVideoSources{}, &mediaVideoSour)
	jsonString, _ = json.Marshal(mediaVideoSour)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备视频配置信息 */
	mediaProfiles := media.GetProfilesResponse{}
	dev.CallMethodInterface(media.GetProfiles{}, &mediaProfiles)
	jsonString, _ = json.Marshal(mediaProfiles)
	log.Printf("%s\r\n", jsonString)

	/* 获取摄像头抓拍地址 */
	mediaSnapshot := media.GetSnapshotUriResponse{}
	dev.CallMethodInterface(media.GetSnapshotUri{ProfileToken: "Profile_1"}, &mediaSnapshot)
	jsonString, _ = json.Marshal(mediaSnapshot)
	log.Printf("%s\r\n", jsonString)
	/* 获取图片 */
	ddd, _ := hikvision.DowloadHttpSnapshotImage(string(mediaSnapshot.MediaUri.Uri), "admin", "123qweasdZXC")
	f, err := os.Create("123123.jpg")
	if err != nil {
		panic(err)
	}
	f.Write(ddd)
	f.Close()

	/* 获取设备RTSP直播地址 */
	mediaRTSP := media.GetStreamUriResponse{}
	dev.CallMethodInterface(media.GetStreamUri{ProfileToken: "Profile_1"}, &mediaRTSP)
	jsonString, _ = json.Marshal(mediaRTSP)
	log.Printf("%s\r\n", jsonString)
}
