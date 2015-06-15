package kempclient

import (
	"encoding/xml"
	"fmt"
	"net"
	"strconv"

	"github.com/juju/errgo"
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
	ID             int `xml:"RsIndex"`
	Status         string
	VirtualService int    `xml:"VsIndex"`
	IPAddress      string `xml:"Addr"`
	Port           string
	Forward        string
	Weight         string
	Limit          string
	Enable         string
}

func (c *Client) AddRealServerByID(id int, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = strconv.Itoa(id)
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
		return errgo.Newf("%s is not a valid ip address", parameters["rs"])
	}
	if parameters["rsport"] == "" {
		return errgo.New("A real server needs a port")
	}
	if parameters["vs"] == "" {
		return errgo.New("The virtual service for the real server is missing")
	}

	data := RealServerResponse{}
	err := c.Request("addrs", parameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to delete real server '%#v'", parameters), errgo.Any)
	}

	return nil
}

func (c *Client) DeleteRealServerByID(id int, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = strconv.Itoa(id)
	parameters["rs"] = rs.IPAddress
	parameters["rsport"] = rs.Port

	return c.deleteRealServer(parameters)
}

func (c *Client) DeleteRealServerByData(ip, port, protocol string, rs RealServer) error {
	parameters := make(map[string]string)
	parameters["vs"] = ip
	parameters["port"] = port
	parameters["prot"] = protocol
	parameters["rs"] = rs.IPAddress
	parameters["rsport"] = rs.Port

	return c.deleteRealServer(parameters)
}

func (c *Client) deleteRealServer(parameters map[string]string) error {
	if net.ParseIP(parameters["rs"]) == nil {
		return errgo.Newf("%s is not a valid ip address", parameters["rs"])
	}
	if parameters["rsport"] == "" {
		return errgo.New("A real server needs a port")
	}
	if parameters["vs"] == "" {
		return errgo.New("The virtual service for the real server is missing")
	}

	data := RealServerResponse{}
	err := c.Request("delrs", parameters, &data)
	if err != nil {
		return errgo.NoteMask(err, fmt.Sprintf("kemp unable to delete real server '%#v'", parameters), errgo.Any)
	}

	return nil
}
