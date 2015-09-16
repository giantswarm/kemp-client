package kempclient

import (
	"encoding/xml"
	"fmt"
	"net"
	"strconv"

	"github.com/juju/errgo"
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
)

type VirtualServiceListResponse struct {
	Debug   string             `xml:",innerxml"`
	XMLName xml.Name           `xml:"Response"`
	Data    VirtualServiceList `xml:"Success>Data"`
}

type VirtualServiceList struct {
	VS []VirtualService `xml:",any"`
}

type VirtualServiceParams struct {
	Name            string
	IPAddress       string
	Port            string
	Protocol        string
	CheckType       string
	CheckURL        string
	CheckPort       string
	SSLAcceleration bool
	Transparent     bool
	AddVia          string
}

type VirtualServiceResponse struct {
	Debug   string         `xml:",innerxml"`
	XMLName xml.Name       `xml:"Response"`
	VS      VirtualService `xml:"Success>Data"`
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

	data := VirtualServiceResponse{}
	err := c.Request("modvs", parameters, &data)
	if err != nil {
		return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to update virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
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

	data := VirtualServiceResponse{}
	err := c.Request("addvs", parameters, &data)
	if err != nil {
		return VirtualService{}, errgo.NoteMask(err, fmt.Sprintf("kemp unable to add virtual service '%#v'", parameters), errgo.Any)
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
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

}
