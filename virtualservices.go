package kempclient

import (
	"encoding/xml"
	"fmt"
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
	IPAddress        string `xml:"VSAddress"`
	Port             string `xml:"VSPort"`
	Protocol         string
	Status           string
	NickName         string
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
	CheckUse11       string `xml:"CheckUse1.1"`
	MatchLen         string
	CheckUseGet      string
	SSLRewrite       string
	VStype           string
	FollowVSID       string
	Schedule         string
	CheckType        string
	PersistTimeout   string
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

func (c *Client) ShowVirtualService(id string) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = id

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

func (c *Client) DeleteVirtualService(id string) error {
	parameters := make(map[string]string)
	parameters["vs"] = id

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

func (c *Client) AddVirtualService(vs VirtualService) (VirtualService, error) {
	parameters := make(map[string]string)
	parameters["vs"] = vs.IPAddress
	parameters["port"] = vs.Port
	parameters["prot"] = vs.Protocol

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
