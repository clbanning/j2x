// j2x package - mirror of x2j package
//	Marshal XML docs from arbitrary JSON string and map[string]interface{} variables.
// Copyright 2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
/*	Marshal dynamic / arbitrary XML docs from arbitrary JSON string and map[string]interface{} variables.

	Compliments the x2j package functions.

	Uses x2j conventions:
		- Keys that begin with a hyphen, '-', are treated as attributes.
		- The "#text" key is treated as the value for a simple element.

	Map values that are not standard JSON types - can be a structure, etc. - are marshal'd using xml.Marshal().
	However, attribute keys are restricted to string, numeric, or boolean types.

	If the map[string]interface{} has a single key, it is used as the XML root tag.  If it doesn't have
	a single key, then a root tag - rootTag - must be provided or the default root tag value is used.
*/

package j2x

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

type pretty struct {
	indent string
	cnt int
	padding string
	inList bool
	inMap bool
}

func (p *pretty)Indent() {
	p.padding += p.indent
	p.cnt++
}

func (p *pretty)Dedent() {
	if p.cnt > 0 {
		p.padding = p.padding[:len(p.padding)-len(p.indent)]
		p.cnt--
	}
}

// Extends xml.MarshalIndent() to handle JSON and map[string]interface{} types.
// See Marshal().
func MarshalIndent(v interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	switch v.(type) {
	case string:
		xmlString, err := JsonToXmlIndent([]byte(v.(string)), prefix, indent, rootTag...)
		return xmlString, err
	case map[string]interface{}:
		xmlString, err := MapToXmlIndent(v.(map[string]interface{}), prefix, indent, rootTag...)
		return xmlString, err
	}
	return xml.MarshalIndent(v, prefix, indent)
}

// Encode a JSON string as pretty XML string.
//	See JsonToXml().
func JsonToXmlIndent(jsonString []byte, prefix, indent string, rootTag ...string) ([]byte, error) {
	m := make(map[string]interface{}, 0)
	if err := json.Unmarshal(jsonString, &m); err != nil {
		return nil, err
	}
	return MapToXmlIndent(m, prefix, indent, rootTag...)
}

// Encode a map[string]interface{} variable as a pretty XML string.
// See MapToXml().
func MapToXmlIndent(m map[string]interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			if _, ok := value.([]interface{}); ok {
				err = p.mapToXmlIndent(s, DefaultRootTag, m)
			} else {
				err = p.mapToXmlIndent(s, key, value)
			}
		}
	} else if len(rootTag) == 1 {
		err = p.mapToXmlIndent(s, rootTag[0], m)
	} else {
		err = p.mapToXmlIndent(s, DefaultRootTag, m)
	}
	return []byte(*s), err
}

// where the work actually happens
// returns an error if an attribute is not atomic
func (p *pretty)mapToXmlIndent(s *string, key string, value interface{}) error {
	var endTag bool
	var isSimple bool

	switch value.(type) {
	case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32:
		*s += p.padding + `<` + key
	}
	switch value.(type) {
	case map[string]interface{}:
		vv := value.(map[string]interface{})
		lenvv := len(vv)
		// scan out attributes - keys have prepended hyphen, '-'
		var cntAttr int
		for k, v := range vv {
			if k[:1] == "-" {
				switch v.(type) {
				case string, float64, bool, int, int32, int64, float32:
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", v) + `"`
					cntAttr++
				case []byte: // allow standard xml pkg []byte transform, as below
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", string(v.([]byte))) + `"`
					cntAttr++
				default:
					return errors.New("invalid attribute value for: " + k)
				}
			}
		}
		// only attributes?
		if cntAttr == lenvv {
			break
		}
		// simple element? Note: '#text" is an invalid XML tag.
		if v, ok := vv["#text"]; ok {
			if cntAttr+1 < lenvv {
				return errors.New("#text key occurs with other non-attribute keys")
			}
			*s += ">" + fmt.Sprintf("%v", v)
			endTag = true
			break
		}
		// close tag with possible attributes
		*s += ">"
		*s += "\n"
		// something more complex
		p.inMap = true
		for k, v := range vv {
			if k[:1] == "-" {
				continue
			}
			switch v.(type) {
			case []interface{}:
			default:
				p.Indent()
			}
			p.mapToXmlIndent(s, k, v)
			switch v.(type) {
			case []interface{}:	// handled in []interface{} case
			default:
				if !p.inList { p.Dedent() }
			}
		}
		p.inMap = false
		endTag = true
	case []interface{}:
		p.inList = true
		for _, v := range value.([]interface{}) {
			p.Indent()
			p.mapToXmlIndent(s, key, v)
			p.Dedent()
		}
		p.inList = false
		return nil
	case nil:
		// terminate the tag
		break
	default: // handle anything - even goofy stuff
		switch value.(type) {
		case string, float64, bool, int, int32, int64, float32:
			*s += ">" + fmt.Sprintf("%v", value)
		case []byte: // NOTE: byte is just an alias for uint8
			// similar to how xml.Marshal handles []byte structure members
			*s += ">" + fmt.Sprintf("%v", string(value.([]byte)))
		default:
			var v []byte
			var err error
				v, err = xml.MarshalIndent(value,p.padding,p.indent)
			if err != nil {
				*s += ">UNKNOWN"
			} else {
				*s += string(v)
			}
		}
		isSimple = true
		endTag = true
	}

	if endTag {
		if !isSimple {
			if p.inList { p.Dedent() }
			*s += p.padding
		}
		switch value.(type) {
		case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32:
			*s += `</` + key + ">"
		}
		// *s += "</" + key + ">"
	} else {
		*s += "/>"
	}
	*s += "\n"
	if !p.inList && !p.inMap {
		p.Dedent()
	}

	return nil
}
