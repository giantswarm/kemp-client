package kempclient

import (
	"encoding/xml"
	"fmt"

	"github.com/juju/errgo"
)

// Information about the content rules can be found in https://support.kemptechnologies.com/hc/en-us/articles/203863435-RESTful-API

const (
	ContentRuleAddHeader    = "1"
	ContentRuleUpdateHeader = "3"
	ContentRuleDeleteHeader = "2"
	DeleteHeaderProtoName   = "DeleteHeaderProto"
	DeleteHeaderPortName    = "DeleteHeaderPort"
	DeleteHeaderPortValue   = "X-Forwarded-Port"
	DeleteHeaderProtoValue  = "X-Forwarded-Proto"
)

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

func (c *Client) AddHeaderContentRule(name, headerKey, headerValue string) error {
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

func (c *Client) AddProtoPortHeaderRequestRules() error {
	data := ContentRuleResponse{}
	ruleParameters := make(map[string]string)

	ruleParameters["name"] = DeleteHeaderProtoName
	ruleParameters["type"] = ContentRuleDeleteHeader

	if err := c.Request("showrule", ruleParameters, &data); err != nil {
		// The Rule doesn't exists
		ruleParameters["pattern"] = DeleteHeaderProtoValue

		err := c.Request("addrule", ruleParameters, &data)
		if err != nil {
			return errgo.NoteMask(err, fmt.Sprintf("kemp unable to add content rule %s with value %s for virtual service %s '%#v'", DeleteHeaderProtoName, DeleteHeaderProtoValue, ruleParameters), errgo.Any)
		}
	}

	ruleParameters["name"] = DeleteHeaderPortName

	if err := c.Request("showrule", ruleParameters, &data); err != nil {
		ruleParameters["pattern"] = DeleteHeaderPortValue

		err = c.Request("addrule", ruleParameters, &data)
		if err != nil {
			return errgo.NoteMask(err, fmt.Sprintf("kemp unable to add content rule %s with value %s for virtual service %s '%#v'", DeleteHeaderPortName, DeleteHeaderPortValue, ruleParameters), errgo.Any)
		}
	}

	return nil
}

func (c *Client) UpdateHeaderContentRule(name, headerKey, headerValue string) error {
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

func (c *Client) DeleteHeaderContentRule(name string) error {
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
