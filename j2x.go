// j2x package - mirror of x2j package
//	Marshal dynamic / arbitrary XML docs from arbitrary JSON string and map[string]interface{} variables.
// Copyright 2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
/*	Marshal dynamic / arbitrary XML docs from arbitrary JSON string and map[string]interface{} variables.

	Compliments the x2j package functions.

	Uses x2j conventions:
		- Keys that begin with a hyphen, '-', are treated as attributes.
		- The "#text" key is treated as the value for a simple element.

	Map values that are not standard JSON types - can be a structure, etc. - are marshal'd using xml.Marshal().

	There are numerous examples in j2x_test.go.
*/
package j2x

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

const (
	DefaultRootTag = "doc"
)

// Extends xml.Marshal() to handle JSON and map[string]interface{} types.
//	This is the inverse of x2j.Unmarshal().
//	Strings are interpreted as JSON strings; use xml.Marshal() to marshal
//	a string as "<string>...</string>" - the standard package handling.
//	See MapToDoc() for encoding rules.
func Marshal(v interface{}, rootTag ...string) ([]byte, error) {
	switch v.(type) {
	case string:
		xmlString, err := JsonToDoc(v.(string), rootTag...)
		return []byte(xmlString), err
	case map[string]interface{}:
		xmlString := MapToDoc(v.(map[string]interface{}), rootTag...)
		return []byte(xmlString), nil
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
	return MapToDoc(m, rootTag...), nil
}

// Encode a map[string]interface{} variable as XML.  The inverse of x2j.DocToMap().
// The following rules apply.
//    - The key label "#text" is treated as the value for a simple element with attributes.
//    - Map keys that begin with a hyphen, '-', are interpreted as attributes.
//    - Map value type encoding:
//          > string, bool, float64, int, int32, int64, float32: per "%v" formating
//          > []bool: by casting to string (which is how xml.Marshal handles such structure members)
//          > structures, etc.: handed to xml.Marshal() - if there is an error, the element
//            value is "UNKNOWN"
//    - Elements with only attribute values or are null are terminated using "/>".
//    - If len(m) == 1 and no rootTag is provided, then the map key is used as the root tag.
//      Thus, `{ "key":"value" }` encodes as `<key>value</key>`.
func MapToDoc(m map[string]interface{}, rootTag ...string) string {
	s := new(string)

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			if _, ok := value.([]interface{}); ok {
				mapToDoc(s, DefaultRootTag, m)
			} else {
				mapToDoc(s, key, value)
			}
		}
		return *s
	} else if len(rootTag) == 1 {
		mapToDoc(s, rootTag[0], m)
		return *s
	}
	mapToDoc(s, DefaultRootTag, m)
	return *s
}

// where the work actually happens
// no errors - we coerce everything using "%v"
func mapToDoc(s *string, key string, value interface{}) {
	var endTag bool
	var isList bool

	if _, isList = value.([]interface{}); !isList {
		*s += `<` + key
	}
	switch value.(type) {
	case map[string]interface{}:
		vv := value.(map[string]interface{})
		// scan out attributes - keys have prepended hyphen, '-'
		var cntAttr int
		for k, v := range vv {
			if k[:1] == "-" {
				*s += ` ` + k[1:] + `="` + fmt.Sprintf("%v", v) + `"`
				cntAttr++
			}
		}
		// only attributes?
		if cntAttr == len(vv) {
			break
		}
		// simple element?
		if v, ok := vv["#text"]; ok {
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
		endTag = true
	case nil:
		// terminate the tag
		break
	default: // handle anything - even goofy stuff
		var tmp string
		switch value.(type) {
		case float64, string, bool, int, int32, int64, float32:
			tmp = fmt.Sprintf("%v", value)
		case []byte:
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

	if isList {
		return
	}
	if endTag {
		*s += "</" + key + ">"
	} else {
		*s += "/>"
	}
}

