package soap

import (
	"encoding/xml"
	"log"

	"github.com/beevik/etree"
)

//SoapMessage type from string
type SoapMessage string

// NewEmptySOAP return new SoapMessage
func NewEmptySOAP() SoapMessage {
	doc := buildSoapRoot()
	res, _ := doc.WriteToString()
	return SoapMessage(res)
}

//NewSOAP Get a new soap message
func NewSOAP(headContent []*etree.Element, bodyContent []*etree.Element, namespaces map[string]string) SoapMessage {
	doc := buildSoapRoot()
	res, _ := doc.WriteToString()
	return SoapMessage(res)
}

func (msg SoapMessage) String() string {
	return string(msg)
}

//StringIndent handle indent
func (msg SoapMessage) StringIndent() string {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	doc.IndentTabs()
	res, _ := doc.WriteToString()
	return res
}

//Body return body from Envelope
func (msg SoapMessage) Body() string {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	bodyTag := doc.Root().SelectElement("Body").ChildElements()[0]
	doc.SetRoot(bodyTag)
	doc.IndentTabs()
	res, _ := doc.WriteToString()
	return res
}

//AddStringBodyContent for Envelope
func (msg *SoapMessage) AddStringBodyContent(data string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(data); err != nil {
		log.Println(err.Error())
	}
	element := doc.Root()
	doc = etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	bodyTag := doc.Root().SelectElement("Body")
	bodyTag.AddChild(element)
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
}

//AddBodyContent for Envelope
func (msg *SoapMessage) AddBodyContent(element *etree.Element) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	bodyTag := doc.Root().SelectElement("Body")
	bodyTag.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	bodyTag.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	bodyTag.AddChild(element)
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
}

//AddBodyContents for Envelope body
func (msg *SoapMessage) AddBodyContents(elements []*etree.Element) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	bodyTag := doc.Root().SelectElement("Body")
	if len(elements) != 0 {
		for _, j := range elements {
			bodyTag.AddChild(j)
		}
	}
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
}

//AddStringHeaderContent for Envelope body
func (msg *SoapMessage) AddStringHeaderContent(data string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(data); err != nil {
		return err
	}
	element := doc.Root()
	doc = etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		return err
	}
	bodyTag := doc.Root().SelectElement("Header")
	bodyTag.AddChild(element)
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
	return nil
}

//AddHeaderContent for Envelope body
func (msg *SoapMessage) AddHeaderContent(element *etree.Element) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	bodyTag := doc.Root().SelectElement("Header")
	bodyTag.AddChild(element)
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
}

//AddHeaderContents for Envelope body
func (msg *SoapMessage) AddHeaderContents(elements []*etree.Element) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	headerTag := doc.Root().SelectElement("Header")
	if len(elements) != 0 {
		for _, j := range elements {
			headerTag.AddChild(j)
		}
	}
	res, _ := doc.WriteToString()
	*msg = SoapMessage(res)
}

//AddRootNamespace for Envelope body
func (msg *SoapMessage) AddRootNamespace(key, value string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
	doc.Root().CreateAttr("xmlns:"+key, value)
	res, _ := doc.WriteToString()

	*msg = SoapMessage(res)
}

//AddRootNamespaces for Envelope body
func (msg *SoapMessage) AddRootNamespaces(namespaces map[string]string) {
	for key, value := range namespaces {
		msg.AddRootNamespace(key, value)
	}
}

func buildSoapRoot() *etree.Document {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	env := doc.CreateElement("s:Envelope")
	env.CreateElement("s:Header")
	env.CreateElement("s:Body")
	env.CreateAttr("xmlns:s", "http://www.w3.org/2003/05/soap-envelope")
	env.CreateAttr("xmlns:a", "http://www.w3.org/2005/08/addressing")
	env.CreateAttr("xmlns:soap-enc", "http://www.w3.org/2003/05/soap-encoding")
	return doc
}

//AddWSSecurity Header for soapMessage
func (msg *SoapMessage) AddWSSecurity(username, password string) {
	/* Getting an WS-Security struct representation */
	auth := NewSecurity(username, password)
	/* Adding WS-Security namespaces to root element of SOAP message */
	soapReq, err := xml.Marshal(auth)
	if err != nil {
		log.Panic(err)
	}
	/*Adding WS-Security struct to SOAP header*/
	msg.AddStringHeaderContent(string(soapReq))
}

//AddAction Header handling for soapMessage
func (msg *SoapMessage) AddAction() {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg.String()); err != nil {
		log.Println(err.Error())
	}
}
