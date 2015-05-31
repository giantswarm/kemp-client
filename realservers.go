package kempclient

import (
	"encoding/xml"
	"fmt"
	"net"
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

func (c *Client) AddRealServerById(id string, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = id
	parameters["rs"] = rs.IPAddress
	parameters["rsport"] = rs.Port

	return c.addRealServer(parameters)
}

func (c *Client) AddRealServerByData(ip, port, protocol string, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = ip
	parameters["port"] = port
	parameters["prot"] = protocol
	parameters["rs"] = rs.IPAddress
	parameters["rsport"] = rs.Port

	return c.addRealServer(parameters)
}

func (c *Client) addRealServer(parameters map[string]string) error {
	if net.ParseIP(parameters["rs"]) == nil {
		return fmt.Errorf("%s is not a valid ip address", parameters["rs"])
	}
	if parameters["rs"] == "" {
		return fmt.Errorf("A virtual service needs a port")
	}

	data := RealServerResponse{}
	err := c.Request("addrs", parameters, &data)
	if err != nil {
		return err
	}

	return nil
}
