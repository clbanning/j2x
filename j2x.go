// Mirror of x2j package. Marshal XML docs from JSON string and map[string]interface{} variables.
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

	EMPTY ELEMENT ENCODING

	Empty (nil) elements or elements with only attributes are encoded as "<tag .../>".  The standard library
	encoding/xml package encodes them as "<tag ...></tag>".  If you're marshaling a map with structure values
	and want a consistent syntax, use the xml_marshal hack of the standard library that conforms encoding/xml
	to the j2x convention.
*/
package j2x

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

const (
	DefaultRootTag = "doc"
)

// Extends xml.Marshal() to handle JSON and map[string]interface{} types.
//	This is the inverse of x2j.Unmarshal().
//	Strings are interpreted as JSON strings; use xml.Marshal() to marshal
//	a string as "<string>...</string>" - the standard package handling.
//	Follows xml.Marshal handling of types except for string and map[string]interface{}
//	values. For more generalized marshal'ing use MapToDoc().
//	See MapToDoc() for encoding rules.
func Marshal(v interface{}, rootTag ...string) ([]byte, error) {
	switch v.(type) {
	case string:
		xmlString, err := JsonToDoc(v.(string), rootTag...)
		return []byte(xmlString), err
	case map[string]interface{}:
		xmlString, err := MapToDoc(v.(map[string]interface{}), rootTag...)
		return []byte(xmlString), err
	}
	return xml.Marshal(v)
}

// Encode a JSON string as XML.  The inverse of x2j.DocToJson().
//	See MapToDoc() for encoding rules.
func JsonToDoc(jsonString string, rootTag ...string) (string, error) {
	m := make(map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(jsonString), &m); err != nil {
		return "", err
	}
	return MapToDoc(m, rootTag...)
}

// Encode a map[string]interface{} variable as XML.  The inverse of x2j.DocToMap().
// The following rules apply.
//    - The key label "#text" is treated as the value for a simple element with attributes.
//    - Map keys that begin with a hyphen, '-', are interpreted as attributes.
//      It is an error if the attribute doesn't have a []byte, string, number, or boolean value.
//    - Map value type encoding:
//          > string, bool, float64, int, int32, int64, float32: per "%v" formating
//          > []bool, []uint8: by casting to string
//          > structures, etc.: handed to xml.Marshal() - if there is an error, the element
//            value is "UNKNOWN"
//    - Elements with only attribute values or are null are terminated using "/>".
//    - If len(m) == 1 and no rootTag is provided, then the map key is used as the root tag.
//      Thus, `{ "key":"value" }` encodes as `<key>value</key>`.
func MapToDoc(m map[string]interface{}, rootTag ...string) (string, error) {
	var err error
	s := new(string)

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			if _, ok := value.([]interface{}); ok {
				err = mapToDoc(s, DefaultRootTag, m)
			} else {
				err = mapToDoc(s, key, value)
			}
		}
	} else if len(rootTag) == 1 {
		err = mapToDoc(s, rootTag[0], m)
	} else {
		err = mapToDoc(s, DefaultRootTag, m)
	}
	return *s, err
}

// where the work actually happens
// returns an error if an attribute is not atomic
func mapToDoc(s *string, key string, value interface{}) error {
	var endTag bool

	if _, ok := value.([]interface{}); !ok {
		*s += `<` + key
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
				case []byte:		// allow standard xml pkg []byte transform, as below
					*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v",string(v.([]byte))) + `"`
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
		// something more complex
		for k, v := range vv {
			if k[:1] == "-" {
				continue
			}
			mapToDoc(s, k, v)
		}
		endTag = true
	case []interface{}:
		for _, v := range value.([]interface{}) {
			mapToDoc(s, key, v)
		}
		return nil
	case nil:
		// terminate the tag
		break
	default: // handle anything - even goofy stuff
		var tmp string
		switch value.(type) {
		case string, float64, bool, int, int32, int64, float32:
			tmp = fmt.Sprintf("%v", value)
		case []byte:			// NOTE: byte is just an alias for uint8
			// similar to how xml.Marshal handles []byte structure members
			tmp = fmt.Sprintf("%v", string(value.([]byte)))
		default:
			v, err := xml.Marshal(value)
			if err != nil {
				tmp = "UNKNOWN"
			} else {
				tmp = string(v)
			}
		}
		*s += ">" + tmp
		endTag = true
	}

	if endTag {
		*s += "</" + key + ">"
	} else {
		*s += "/>"
	}
	return nil
}
