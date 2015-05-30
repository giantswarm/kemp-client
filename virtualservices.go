package kempclient

import (
	"encoding/xml"
	"fmt"
	"net"
)

type VirtualServiceListResponse struct {
	Debug   string             `xml:",innerxml"`
	XMLName xml.Name           `xml:"Response"`
	Data    VirtualServiceList `xml:"Success>Data"`
}

type VirtualServiceList struct {
	VS []VirtualService `xml:",any"`
}

type VirtualServiceResponse struct {
	Debug   string         `xml:",innerxml"`
	XMLName xml.Name       `xml:"Response"`
	VS      VirtualService `xml:"Success>Data"`
}

type VirtualService struct {
	ID               string `xml:"Index"`
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
	CheckUrl         string
	CheckUse11       string `xml:"CheckUse1.1"`
	MatchLen         string
	CheckUseGet      string
	SSLRewrite       string
	VStype           string
	FollowVSID       string
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
	MasterVSID       string
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
		return []VirtualService{}, err
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.Data.VS, nil
}

func (c *Client) ShowVirtualServiceByData(ip, port, protocol string) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = ip
	parameters["port"] = port
	parameters["prot"] = protocol

	return c.showVirtualService(parameters)
}

func (c *Client) ShowVirtualServiceById(id string) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = id

	return c.showVirtualService(parameters)
}

func (c *Client) showVirtualService(parameters map[string]string) (VirtualService, error) {
	data := VirtualServiceResponse{}
	err := c.Request("showvs", parameters, &data)
	if err != nil {
		return VirtualService{}, err
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.VS, nil
}

func (c *Client) DeleteVirtualServiceById(id string) error {
	parameters := make(map[string]string)
	parameters["vs"] = id

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
		return err
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return nil
}

func (c *Client) UpdateVirtualService(id string, vs VirtualService) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = id

	if vs.Name != "" {
		parameters["nickname"] = vs.Name
	}
	if vs.IPAddress != "" {
		parameters["vsaddress"] = vs.IPAddress
	}
	if vs.Port != "" {
		parameters["vsport"] = vs.Port
	}
	if vs.Protocol != "" {
		parameters["prot"] = vs.Protocol
	}
	if vs.Transparent != "" {
		parameters["transparent"] = vs.Transparent
	}
	if vs.CheckType != "" {
		parameters["checktype"] = vs.CheckType
	}
	if vs.CheckUrl != "" {
		parameters["checkurl"] = vs.CheckUrl
	}
	if vs.CheckPort != "" {
		parameters["checkport"] = vs.CheckPort
	}
	if vs.SSLAcceleration != "" {
		parameters["sslacceleration"] = vs.SSLAcceleration
	}

	data := VirtualServiceResponse{}
	err := c.Request("modvs", parameters, &data)
	if err != nil {
		return VirtualService{}, err
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.VS, nil
}

func (c *Client) AddVirtualService(vs VirtualService) (VirtualService, error) {
	parameters := make(map[string]string)
	if net.ParseIP(vs.IPAddress) == nil {
		return VirtualService{}, fmt.Errorf("%s is not a valid ip address", vs.IPAddress)
	}
	if vs.Port == "" {
		return VirtualService{}, fmt.Errorf("A virtual service needs a port")
	}
	if vs.Protocol != "tcp" && vs.Protocol != "udp" {
		return VirtualService{}, fmt.Errorf("The protocol of a virtual service is either tcp or udp")
	}

	parameters["vs"] = vs.IPAddress
	parameters["port"] = vs.Port
	parameters["prot"] = vs.Protocol

	if vs.Name != "" {
		parameters["nickname"] = vs.Name
	}
	if vs.Transparent != "" {
		parameters["transparent"] = vs.Transparent
	}
	if vs.CheckType != "" {
		parameters["checktype"] = vs.CheckType
	}
	if vs.CheckUrl != "" {
		parameters["checkurl"] = vs.CheckUrl
	}
	if vs.CheckPort != "" {
		parameters["checkport"] = vs.CheckPort
	}
	if vs.SSLAcceleration != "" {
		parameters["sslacceleration"] = vs.SSLAcceleration
	}

	data := VirtualServiceResponse{}
	err := c.Request("addvs", parameters, &data)
	if err != nil {
		return VirtualService{}, err
	}

	if c.debug {
		fmt.Println("DEBUG:", data.Debug)
	}

	return data.VS, nil
}
