/*
 * @Author: YanHui Li
 * @Date: 2022-02-08 09:37:30
 * @LastEditTime: 2022-02-23 17:03:06
 * @LastEditors: YanHui Li
 * @Description:
 * @FilePath: \go-onvif\examples\GetDeviceInfo\test.go
 *
 */
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/liyanhui1998/go-onvif"
	"github.com/liyanhui1998/go-onvif/types/device"
	"github.com/liyanhui1998/go-onvif/types/media"
)

func main() {
	/* 连接设备 */
	dev, _ := onvif.NewDevice(onvif.DeviceParams{Ipddr: "10.1.1.210", Username: "admin", Password: "123qweasdZXC"})

	/* 获取能力集合 */
	retServices := dev.GetServices()
	for key, value := range retServices {
		log.Printf("server key : %s value : %s\r\n", key, value)
	}

	/* 获取设备基本信息 */
	deviceInfo := device.GetDeviceInformationResponse{}
	dev.CallMethodInterface(device.GetDeviceInformation{}, &deviceInfo, "")
	jsonString, _ := json.Marshal(deviceInfo)
	log.Printf("%s\r\n", jsonString)

	// /* 获取设备网卡信息 */
	deviceInterfaces := device.GetNetworkInterfacesResponse{}
	dev.CallMethodInterface(device.GetNetworkInterfaces{}, &deviceInterfaces, "")
	jsonString, _ = json.Marshal(deviceInterfaces)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备支持网络协议信息 */
	deviceNetProtocols := device.GetNetworkProtocolsResponse{}
	dev.CallMethodInterface(device.GetNetworkProtocols{}, &deviceNetProtocols, "")
	jsonString, _ = json.Marshal(deviceNetProtocols)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备视频资源信息 */
	mediaVideoSour := media.GetVideoSourcesResponse{}
	dev.CallMethodInterface(media.GetVideoSources{}, &mediaVideoSour, "")
	jsonString, _ = json.Marshal(mediaVideoSour)
	log.Printf("%s\r\n", jsonString)

	/* 获取设备视频配置信息 */
	mediaProfiles := media.GetProfilesResponse{}
	dev.CallMethodInterface(media.GetProfiles{}, &mediaProfiles, "")
	jsonString, _ = json.Marshal(mediaProfiles)
	log.Printf("%s\r\n", jsonString)

	// /* 获取摄像头抓拍地址 */
	mediaSnapshot := media.GetSnapshotUriResponse{}
	dev.CallMethodInterface(media.GetSnapshotUri{ProfileToken: mediaProfiles.Profiles[0].Token}, &mediaSnapshot, "")
	jsonString, _ = json.Marshal(mediaSnapshot)
	log.Printf("%s\r\n", jsonString)
	/* 获取图片 */
	ddd, _ := onvif.HttpDigestAuthGetSnapshotImage(string(mediaSnapshot.MediaUri.Uri), "admin", "123qweasdZXC")
	f, err := os.Create("123123.jpg")
	if err != nil {
		panic(err)
	}
	f.Write(ddd)
	f.Close()

	/* 获取设备RTSP直播地址 */
	mediaRTSP := media.GetStreamUriResponse{}
	/* 大华摄像头获取视频流地址需要指定rstp信息 */
	if err := dev.CallMethodInterface(media.GetStreamUri{ProfileToken: mediaProfiles.Profiles[0].Token,
		StreamSetup: device.StreamSetup{Stream: "RTP-Unicast", Transport: device.Transport{Protocol: "UDP"}}}, &mediaRTSP, ""); err != nil {
		log.Fatalln(err)
	}
	jsonString, _ = json.Marshal(mediaRTSP)
	log.Printf("%s\r\n", jsonString)
}
