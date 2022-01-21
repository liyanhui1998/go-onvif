package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/liyanhui1998/go-onvif"
	"github.com/liyanhui1998/go-onvif/types/device"
	"github.com/liyanhui1998/go-onvif/types/media"
)

func readResponse(resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func main111() {
	deviceList := onvif.GetAvailableDevicesAtSpecificEthernetInterface("以太网")
	for _, value := range deviceList {
		log.Printf("scan IP:%s\r\n", value.Params.Ipddr)
	}

	dev, _ := onvif.NewDevice(onvif.DeviceParams{Ipddr: "10.1.1.200", Username: "admin", Password: "123qweasdZXC"})
	/*
		2021/12/31 10:04:54 server key:events value:http://10.1.1.200/onvif/Events
		2021/12/31 10:04:54 server key:imaging value:http://10.1.1.200/onvif/Imaging
		2021/12/31 10:04:54 server key:media value:http://10.1.1.200/onvif/Media
		2021/12/31 10:04:54 server key:device value:http://10.1.1.200/onvif/device_service
		2021/12/31 10:04:54 server key:analytics value:http://10.1.1.200/onvif/Analytics
	*/
	retServices := dev.GetServices()
	for key, value := range retServices {
		log.Printf("server key:%s value:%s\r\n", key, value)
	}
	/*
		<?xml version="1.0" encoding="UTF-8"?>
		<env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope" xmlns:soapenc="http://www.w3.org/2003/05/soap-encoding" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:tt="http://www.onvif.org/ver10/schema" xmlns:tds="http://www.onvif.org/ver10/device/wsdl" xmlns:trt="http://www.onvif.org/ver10/media/wsdl" xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl" xmlns:tev="http://www.onvif.org/ver10/events/wsdl" xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl" xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl" xmlns:tst="http://www.onvif.org/ver10/storage/wsdl" xmlns:ter="http://www.onvif.org/ver10/error" xmlns:dn="http://www.onvif.org/ver10/network/wsdl" xmlns:tns1="http://www.onvif.org/ver10/topics" xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl" xmlns:wsoap12="http://schemas.xmlsoap.org/wsdl/soap12" xmlns:http="http://schemas.xmlsoap.org/wsdl/http" xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:wsadis="http://schemas.xmlsoap.org/ws/2004/08/addressing" xmlns:wsnt="http://docs.oasis-open.org/wsn/b-2" xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:wstop="http://docs.oasis-open.org/wsn/t-1" xmlns:wsrf-bf="http://docs.oasis-open.org/wsrf/bf-2" xmlns:wsntw="http://docs.oasis-open.org/wsn/bw-2" xmlns:wsrf-rw="http://docs.oasis-open.org/wsrf/rw-2" xmlns:wsaw="http://www.w3.org/2006/05/addressing/wsdl" xmlns:wsrf-r="http://docs.oasis-open.org/wsrf/r-2" xmlns:trc="http://www.onvif.org/ver10/recording/wsdl" xmlns:tse="http://www.onvif.org/ver10/search/wsdl" xmlns:trp="http://www.onvif.org/ver10/replay/wsdl" xmlns:tnshik="http://www.hikvision.com/2011/event/topics" xmlns:hikwsd="http://www.onvifext.com/onvif/ext/ver10/wsdl" xmlns:hikxsd="http://www.onvifext.com/onvif/ext/ver10/schema" xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl" xmlns:tr2="http://www.onvif.org/ver20/media/wsdl" xmlns:axt="http://www.onvif.org/ver20/analytics">
			<env:Body>
				<tds:GetDeviceInformationResponse>
					<tds:Manufacturer>HIKVISION</tds:Manufacturer>
					<tds:Model>DS-2CD3T45P1-I</tds:Model>
					<tds:FirmwareVersion>V5.5.31 build 180903</tds:FirmwareVersion>
					<tds:SerialNumber>DS-2CD3T45P1-I20210803AACH248538890</tds:SerialNumber>
					<tds:HardwareId>88</tds:HardwareId>
				</tds:GetDeviceInformationResponse>
			</env:Body>
		</env:Envelope>
	*/
	Response, _ := dev.CallMethod(device.GetDeviceInformation{})
	log.Printf("%s\r\n", readResponse(Response))
	/*
		<tds:NetworkProtocols>
			<tt:Name>RTSP</tt:Name>
			<tt:Enabled>true</tt:Enabled>
			<tt:Port>554</tt:Port>
		</tds:NetworkProtocols>
	*/
	Response, _ = dev.CallMethod(device.GetNetworkProtocols{})
	log.Printf("%s\r\n", readResponse(Response))
	/*
		<tt:Resolution>
				<tt:Width>2560</tt:Width>
				<tt:Height>1440</tt:Height>
		</tt:Resolution>
	*/
	Response, _ = dev.CallMethod(media.GetVideoSources{})
	log.Printf("%s\r\n", readResponse(Response))

	/*
		<trt:Profiles token="Profile_1" fixed="true">
			....
		</try>
	*/
	Response, _ = dev.CallMethod(media.GetProfiles{})
	log.Printf("%s\r\n", readResponse(Response))
	/*
		http://10.1.1.200/onvif-http/snapshot?Profile_1
		注:海康摄像头获取快照需要鉴权,onvif获取到的协议无法直接使用

		鉴权方法
			在http请求头中添加
				Authorization  值为 "Basic " + Base64("name:passwd")
			例如:
				admin:123qweasdZXC Base64 后 YWRtaW46MTIzcXdlYXNkWlhD
				则 Authorization 的值为 "Base64 YWRtaW46MTIzcXdlYXNkWlhD"
	*/
	Response, _ = dev.CallMethod(media.GetSnapshotUri{ProfileToken: "Profile_1"})
	log.Printf("%s\r\n", readResponse(Response))

	/*
		注获取地址为
			rtsp://10.1.1.200:554/Streaming/Channels/101?transportmode=unicast&amp;profile=Profile_1
		需添加用户名和密码
			rtsp://admin:123qweasdZXC@10.1.1.200:554/Streaming/Channels/101?transportmode=unicast&amp;profile=Profile_1
	*/
	Response, _ = dev.CallMethod(media.GetStreamUri{ProfileToken: "Profile_1"})
	log.Printf("%s\r\n", readResponse(Response))
}
