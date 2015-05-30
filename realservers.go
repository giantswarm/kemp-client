package kempclient

import (
	"encoding/xml"
	"fmt"
)

type RealServerResponse struct {
	Debug   string     `xml:",innerxml"`
	XMLName xml.Name   `xml:"Response"`
	Data    RealServer `xml:"Success>Data"`
}

type RealServerList struct {
	Rs []RealServer `xml:",any"`
}

type RealServer struct {
	ID             string `xml:"RsIndex"`
	Status         string
	VirtualService string `xml:"VsIndex"`
	IPAddress      string `xml:"Addr"`
	Port           string
	Forward        string
	Weight         string
	Limit          string
	Enable         string
}

func (c *Client) AddRealServer(vs VirtualService, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = vs.IPAddress
	parameters["port"] = vs.Port
	parameters["prot"] = vs.Protocol
	parameters["rs"] = rs.IPAddress
	parameters["rsport"] = rs.Port

	data := RealServerResponse{}
	err := c.Request("addrs", parameters, data)
	if err != nil {
		return err
	}

	fmt.Println(data)
	return nil
}
