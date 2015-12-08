package kempclient

import (
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/juju/errgo"
)

// The type of the virutalservice.
//
// Doc from webui:
// > Select the type of service that will be run over this Virtual Service.
// > The LoadMaster automatically tries to determine the type of the service
// > but the user can override this.
const (
	VStypeHTTP           = "http"
	VStypeGeneric        = "gen"
	VStypeSTARTTLS       = "tls"
	VStypeRemoteTerminal = "ts"
	VStypeLogInsight     = "log"
)

// These are the options for the AddVia field.
// In the kemp webui it relates to the realserver option under Advanced Properties > Add HTTP Headers
//
// Doc from the webui
// > Select which headers are to be added to HTTP requests.
// > X-ClientSide and X-Forwarded-For are only added to Non-Transparent Connections.
const (
	VSAddViaLegacyXClientSide    = "0"
	VSAddViaNone                 = "2"
	VSAddViaXClientSideWithVia   = "3"
	VSAddViaXClientSideNoVia     = "4"
	VSAddViaXForwardedForWithVia = "1"
	VSAddViaXForwardedForNoVia   = "5"
	VSAddViaViaOnly              = "6"
	ContentRuleAddHeader         = "1"
	ContentRuleUpdateHeader      = "3"
	ContentRuleDeleteHeader      = "2"
)

type VirtualServiceListResponse struct {
	Debug   string             `xml:",innerxml"`
	XMLName xml.Name           `xml:"Response"`
	Data    VirtualServiceList `xml:"Success>Data"`
}

type VirtualServiceList struct {
	VS []VirtualService `xml:",any"`
}

type Header struct {
	Key   string
	Value string
}

type VirtualServiceParams struct {
	Name                    string
	IPAddress               string
	Port                    string
	Protocol                string
	CheckType               string
	CheckURL                string
	CheckPort               string
	SSLAcceleration         bool
	Transparent             bool
	AddVia                  string
	VStype                  string
	ExtraRequestHeaderKey   string
	ExtraRequestHeaderValue string
	Headers                 map[string]string
}

type VirtualServiceResponse struct {
	Debug   string         `xml:",innerxml"`
	XMLName xml.Name       `xml:"Response"`
	VS      VirtualService `xml:"Success>Data"`
}

type ContentRuleResponse struct {
	Debug   string      `xml:",innerxml"`
	XMLName xml.Name    `xml:"Response"`
	CR      ContentRule `xml:"Success>Data"`
}

type ContentRule struct {
	Name        string `xml:"Name"`
	Header      string `xml:"Header"`
	HeaderValue string `xml:"HeaderValue"`
}

type VirtualService struct {
	ID               int    `xml:"Index"`
	Name             string `xml:"NickName"`
	IPAddress        string `xml:"VSAddress"`
	Port             string `xml:"VSPort"`
	Protocol         string
	Status           string
	Enable           string
	SSLReverse       string
	SSLReencrypt     string
	Intercept        string
	InterceptOpts    []string `xml:"InterceptOpts>Opt"`
	AlertThreshold   string
	Transactionlimit string
	Transparent      string
	ServerInit       string
	StartTLSMode     string
	Idletime         string
	Cache            string
	Compress         string
	Verify           string
	UseforSnat       string
	ForceL7          string
	ClientCert       string
	ErrorCode        string
	CertFile         string
	CheckURL         string `xml:"CheckUrl"`
	CheckUse11       string `xml:"CheckUse1.1"`
	MatchLen         string
	CheckUseGet      string
	SSLRewrite       string
	VStype           string
	FollowVSID       int
	Schedule         string
	CheckType        string
	PersistTimeout   string
	SSLAcceleration  string
	CheckPort        string
	NRules           string
	NRequestRules    string
	NResponseRules   string
	NPreProcessRules string
	EspEnabled       string
	InputAuthMode    string
	OutputAuthMode   string
	MasterVS         string
	MasterVSID       int
	AddVia           string
	TlsType          string
	NeedHostName     string
	OCSPVerify       string
	NumberOfRSs      string
	Rs               []RealServer `xml:"Rs"`
	ExtraHdrKey      string
	ExtraHdrValue    string
}

func (c *Client) ListVirtualServices() ([]VirtualService, error) {
	parameters := make(map[string]string)

	data := VirtualServiceListResponse{}
	err := c.Request("listvs", parameters, &data)
	if err != nil {
		return []VirtualService{}, errgo.NoteMask(err, "kemp could not list virtual services", errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.Data.VS, nil
}

func (c *Client) FindVirtualServiceByName(name string) (VirtualService, error) {
	list, err := c.ListVirtualServices()
	if err != nil {
		return VirtualService{}, errgo.Mask(err)
	}

	for _, vs := range list {
		if vs.Name == name {
			return vs, nil
		}
	}

	return VirtualService{}, nil
}

func (c *Client) ShowVirtualServiceByData(ip, port, protocol string) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = ip
	parameters["port"] = port
	parameters["prot"] = protocol

	return c.showVirtualService(parameters)
}

func (c *Client) ShowVirtualServiceByID(id int) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = strconv.Itoa(id)

	return c.showVirtualService(parameters)
}

func (c *Client) showVirtualService(parameters map[string]string) (VirtualService, error) {
	data := VirtualServiceResponse{}
	err := c.Request("showvs", parameters, &data)
	if err != nil {
		return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to show virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.VS, nil
}

func (c *Client) DeleteVirtualServiceByID(id int) error {
	parameters := make(map[string]string)
	parameters["vs"] = strconv.Itoa(id)

	return c.deleteVirtualService(parameters)
}

func (c *Client) DeleteVirtualServiceByData(ip, port, protocol string) error {
	parameters := make(map[string]string)
	parameters["vs"] = ip
	parameters["port"] = port
	parameters["prot"] = protocol

	return c.deleteVirtualService(parameters)
}

func (c *Client) deleteVirtualService(parameters map[string]string) error {
	data := VirtualServiceResponse{}
	err := c.Request("delvs", parameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to delete virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return nil
}

func (c *Client) addHeaderContentRule(name, headerKey, headerValue string) error {
	data := ContentRuleResponse{}

	ruleParameters := make(map[string]string)
	// Only accepts alphanumeric names
	ruleParameters["name"] = name
	// Add Header to HTTP Request
	ruleParameters["type"] = ContentRuleAddHeader
	ruleParameters["header"] = headerKey
	ruleParameters["replacement"] = headerValue

	err := c.Request("addrule", ruleParameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to add content rule %s with header %s and value %s for virtual service %s '%#v'", name, headerKey, headerValue, ruleParameters), errgo.Any)
	}

	return nil
}

func (c *Client) updateHeaderContentRule(name, headerKey, headerValue string) error {
	data := ContentRuleResponse{}

	ruleParameters := make(map[string]string)
	// Only accepts alphanumeric names
	ruleParameters["name"] = name
	// Update Header to HTTP Request
	ruleParameters["type"] = ContentRuleUpdateHeader
	ruleParameters["header"] = headerKey
	ruleParameters["replacement"] = headerValue

	err := c.Request("modrule", ruleParameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to update content rule %s with header %s and value %s for virtual service %s '%#v'", name, headerKey, headerValue, ruleParameters), errgo.Any)
	}

	return nil
}

func (c *Client) deleteHeaderContentRule(name string) error {
	data := ContentRuleResponse{}

	ruleParameters := make(map[string]string)
	// Only accepts alphanumeric names
	ruleParameters["name"] = name
	// Delete Header to HTTP Request
	ruleParameters["type"] = ContentRuleDeleteHeader

	err := c.Request("delrule", ruleParameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to delete  content rule %s header for virtual service '%#v'", name, ruleParameters), errgo.Any)
	}
	return nil
}

func (c *Client) UpdateVirtualService(id int, vs VirtualServiceParams) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = strconv.Itoa(id)

	if vs.IPAddress != "" {
		parameters["vsaddress"] = vs.IPAddress
	}
	if vs.Port != "" {
		parameters["vsport"] = vs.Port
	}
	if vs.Protocol != "" {
		parameters["prot"] = vs.Protocol
	}

	c.mapVirtualServiceParamsToRequestParams(vs, parameters)

	for key, value := range vs.Headers {
		// Deleting the content rule http header as there isn't a truly update operation
		if err := c.deleteHeaderContentRule(strings.Replace(vs.Name+key, "-", "", -1)); err != nil {
			fmt.Println(err)
		}
		if err := c.addHeaderContentRule(strings.Replace(vs.Name+key, "-", "", -1), key, value); err != nil {
			return VirtualService{}, err
		}
	}

	data := VirtualServiceResponse{}
	err := c.Request("modvs", parameters, &data)
	if err != nil {
		return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to update virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	for key := range vs.Headers {
		parameters["rule"] = strings.Replace(vs.Name+key, "-", "", -1)
		err := c.Request("addrequestrule", parameters, &data)
		if err != nil {
			return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to add rule to the virtual service '%#v'", parameters), errgo.Any)
		}
	}

	return data.VS, nil
}

func (c *Client) AddVirtualService(vs VirtualServiceParams) (VirtualService, error) {
	parameters := make(map[string]string)
	if net.ParseIP(vs.IPAddress) == nil {
		return VirtualService{}, errgo.Newf("%s is not a valid ip address", vs.IPAddress)
	}
	if vs.Port == "" {
		return VirtualService{}, errgo.New("A virtual service needs a port")
	}
	if vs.Protocol != "tcp" && vs.Protocol != "udp" {
		return VirtualService{}, errgo.New("The protocol of a virtual service is either tcp or udp")
	}

	parameters["vs"] = vs.IPAddress
	parameters["port"] = vs.Port
	parameters["prot"] = vs.Protocol

	c.mapVirtualServiceParamsToRequestParams(vs, parameters)

	for key, value := range vs.Headers {
		// Deleting the content rule http header as there isn't a truly update operation
		if err := c.deleteHeaderContentRule(strings.Replace(vs.Name+key, "-", "", -1)); err != nil {
			fmt.Println(err)
		}
		if err := c.addHeaderContentRule(strings.Replace(vs.Name+key, "-", "", -1), key, value); err != nil {
			return VirtualService{}, err
		}
	}

	data := VirtualServiceResponse{}
	err := c.Request("addvs", parameters, &data)
	if err != nil {
		return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to add virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	for key := range vs.Headers {
		parameters["rule"] = strings.Replace(vs.Name+key, "-", "", -1)
		err := c.Request("addrequestrule", parameters, &data)
		if err != nil {
			return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to add rule to the virtual service '%#v'", parameters), errgo.Any)
		}
	}

	return data.VS, nil
}

func (c *Client) mapVirtualServiceParamsToRequestParams(vs VirtualServiceParams, parameters map[string]string) {
	if vs.Name != "" {
		parameters["nickname"] = vs.Name
	}

	if vs.Transparent {
		parameters["transparent"] = "Y"
	} else {
		parameters["transparent"] = "N"
	}

	if vs.CheckType != "" {
		parameters["checktype"] = vs.CheckType
	}
	if vs.CheckURL != "" {
		parameters["checkurl"] = vs.CheckURL
	}
	if vs.CheckPort != "" {
		parameters["checkport"] = vs.CheckPort
	}

	if vs.SSLAcceleration {
		parameters["sslacceleration"] = "Y"
	} else {
		parameters["sslacceleration"] = "N"
	}

	if vs.AddVia != "" {
		parameters["addvia"] = vs.AddVia
	}

	if vs.ExtraRequestHeaderKey != "" && vs.ExtraRequestHeaderValue != "" {
		parameters["extrahdrkey"] = vs.ExtraRequestHeaderKey
		parameters["extrahdrvalue"] = vs.ExtraRequestHeaderValue
	}

	if vs.VStype != "" {
		parameters["vstype"] = vs.VStype
	}

}
