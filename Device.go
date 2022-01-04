package onvif

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"onvif/soap"
	"onvif/types/device"
	"reflect"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

type DeviceParams struct {
	Ipddr    string
	Username string
	Password string
}

type Device struct {
	Params     DeviceParams
	httpClient *http.Client
	endpoints  map[string]string
}

//DeviceType alias for int
type DeviceType int

// Onvif Device Tyoe
const (
	NVD DeviceType = iota
	NVS
	NVA
	NVT
)

func (devType DeviceType) String() string {
	stringRepresentation := []string{
		"NetworkVideoDisplay",
		"NetworkVideoStorage",
		"NetworkVideoAnalytics",
		"NetworkVideoTransmitter",
	}
	i := uint8(devType)
	switch {
	case i <= uint8(NVT):
		return stringRepresentation[i]
	default:
		return strconv.Itoa(int(i))
	}
}

//Xlmns XML Scheam
var Xlmns = map[string]string{
	"onvif":   "http://www.onvif.org/ver10/schema",
	"tds":     "http://www.onvif.org/ver10/device/wsdl",
	"trt":     "http://www.onvif.org/ver10/media/wsdl",
	"tev":     "http://www.onvif.org/ver10/events/wsdl",
	"tptz":    "http://www.onvif.org/ver20/ptz/wsdl",
	"timg":    "http://www.onvif.org/ver20/imaging/wsdl",
	"tan":     "http://www.onvif.org/ver20/analytics/wsdl",
	"xmime":   "http://www.w3.org/2005/05/xmlmime",
	"wsnt":    "http://docs.oasis-open.org/wsn/b-2",
	"xop":     "http://www.w3.org/2004/08/xop/include",
	"wsa":     "http://www.w3.org/2005/08/addressing",
	"wstop":   "http://docs.oasis-open.org/wsn/t-1",
	"wsntw":   "http://docs.oasis-open.org/wsn/bw-2",
	"wsrf-rw": "http://docs.oasis-open.org/wsrf/rw-2",
	"wsaw":    "http://www.w3.org/2006/05/addressing/wsdl",
}

func GetAvailableDevicesAtSpecificEthernetInterface(interfaceName string) []Device {
	/* Call an ws-discovery Probe Message to Discover NVT type Devices */
	devices := SendProbe(interfaceName, nil, []string{"dn:" + NVT.String()}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	nvtDevices := make([]Device, 0)

	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			log.Printf("error:%s", err.Error())
			return nil
		}

		endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
		for _, xaddr := range endpoints {
			xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
			c := 0
			for c = 0; c < len(nvtDevices); c++ {
				if nvtDevices[c].Params.Ipddr == xaddr {
					fmt.Println(nvtDevices[c].Params.Ipddr, "==", xaddr)
					break
				}
			}
			if c < len(nvtDevices) {
				continue
			}
			dev, err := NewDevice(DeviceParams{Ipddr: strings.Split(xaddr, " ")[0]})
			if err != nil {
				log.Printf("error:%s", err.Error())
				continue
			} else {
				nvtDevices = append(nvtDevices, *dev)

			}
		}
	}
	return nvtDevices
}

//NewDevice function construct a ONVIF Device entity
func NewDevice(params DeviceParams) (*Device, error) {
	dev := new(Device)
	dev.Params = params
	dev.endpoints = make(map[string]string)
	dev.addEndpoint("Device", "http://"+dev.Params.Ipddr+"/onvif/device_service")

	if dev.httpClient == nil {
		dev.httpClient = new(http.Client)
	}

	getCapabilities := device.GetCapabilities{Category: "All"}
	resp, err := dev.CallMethod(getCapabilities)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("camera is not available at " + dev.Params.Ipddr + " or it does not support ONVIF services")
	}
	dev.getSupportedServices(resp)
	return dev, nil
}

func readResponse(resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

//GetServices return available endpoints
func (dev *Device) GetServices() map[string]string {
	return dev.endpoints
}

func (dev *Device) getSupportedServices(resp *http.Response) {
	doc := etree.NewDocument()
	data, _ := ioutil.ReadAll(resp.Body)
	if err := doc.ReadFromBytes(data); err != nil {
		return
	}
	services := doc.FindElements("./Envelope/Body/GetCapabilitiesResponse/Capabilities/*/XAddr")
	for _, j := range services {
		dev.addEndpoint(j.Parent().Tag, j.Text())
	}
}

func (dev *Device) addEndpoint(Key, Value string) {
	//use lowCaseKey
	//make key having ability to handle Mixed Case for Different vendor devcie (e.g. Events EVENTS, events)
	lowCaseKey := strings.ToLower(Key)
	// Replace host with host from device params.
	if u, err := url.Parse(Value); err == nil {
		u.Host = dev.Params.Ipddr
		Value = u.String()
	}
	dev.endpoints[lowCaseKey] = Value
}

//getEndpoint functions get the target service endpoint in a better way
func (dev Device) getEndpoint(endpoint string) (string, error) {

	// common condition, endpointMark in map we use this.
	if endpointURL, bFound := dev.endpoints[endpoint]; bFound {
		return endpointURL, nil
	}

	//but ,if we have endpoint like event、analytic
	//and sametime the Targetkey like : events、analytics
	//we use fuzzy way to find the best match url
	var endpointURL string
	for targetKey := range dev.endpoints {
		if strings.Contains(targetKey, endpoint) {
			endpointURL = dev.endpoints[targetKey]
			return endpointURL, nil
		}
	}
	return endpointURL, errors.New("target endpoint service not found")
}

func (dev Device) CallMethodInterface(method interface{}, outStruct interface{}) error {
	pkgPath := strings.Split(reflect.TypeOf(method).PkgPath(), "/")
	pkg := strings.ToLower(pkgPath[len(pkgPath)-1])

	endpoint, err := dev.getEndpoint(pkg)
	if err != nil {
		return err
	}
	retResponse, err := dev.callMethodDo(endpoint, method)
	if err != nil {
		return err
	}
	resString := readResponse(retResponse)
	typ := reflect.TypeOf(outStruct)
	val := reflect.ValueOf(outStruct)
	if typ.Kind() == reflect.Ptr {
		// 传入的interface是指针，需要.Elem()取得指针指向的value
		typ = typ.Elem()
		val = val.Elem()
	} else {
		return errors.New("outStruct must be ptr to struct")
	}
	//获取到该结构体有几个字段
	num := val.NumField()
	//遍历结构体的所有字段
	for i := 0; i < num; i++ {
		//获取到struct标签，需要通过reflect.Type来获取tag标签的值
		value := val.Field(i)
		tagVal := typ.Field(i).Tag.Get("xml")
		if tagVal != "" {
			if strings.Contains(resString, tagVal) {
				startString := "<" + tagVal + ">"
				endString := "</" + tagVal + ">"
				subString := resString[strings.Index(resString, startString)+len(startString) : strings.Index(resString, endString)]
				value.Set(reflect.ValueOf(subString))
			}
		}
	}
	return nil
}

//CallMethod functions call an method, defined <method> struct.
//You should use Authenticate method to call authorized requests.
func (dev Device) CallMethod(method interface{}) (*http.Response, error) {
	pkgPath := strings.Split(reflect.TypeOf(method).PkgPath(), "/")
	pkg := strings.ToLower(pkgPath[len(pkgPath)-1])

	endpoint, err := dev.getEndpoint(pkg)
	if err != nil {
		return nil, err
	}
	return dev.callMethodDo(endpoint, method)
}

//CallMethod functions call an method, defined <method> struct with authentication data
func (dev Device) callMethodDo(endpoint string, method interface{}) (*http.Response, error) {
	output, err := xml.Marshal(method)
	if err != nil {
		return nil, err
	}

	soap, err := dev.buildMethodSOAP(string(output))
	if err != nil {
		return nil, err
	}
	soap.AddRootNamespaces(Xlmns)
	soap.AddAction()
	if dev.Params.Username != "" && dev.Params.Password != "" {
		soap.AddWSSecurity(dev.Params.Username, dev.Params.Password)
	}

	return SendSoap(dev.httpClient, endpoint, soap.String())
}

func (dev Device) buildMethodSOAP(msg string) (soap.SoapMessage, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg); err != nil {
		return "", err
	}
	element := doc.Root()
	soap := soap.NewEmptySOAP()
	soap.AddBodyContent(element)
	return soap, nil
}

// SendSoap send soap message
func SendSoap(httpClient *http.Client, endpoint, message string) (*http.Response, error) {
	resp, err := httpClient.Post(endpoint, "application/soap+xml; charset=utf-8", bytes.NewBufferString(message))
	if err != nil {
		return resp, err
	}

	return resp, nil
}
