// j2xReader.go - wrap j2X with an io.Reader input option

package j2x

import (
	"encoding/json"
	"errors"
	"io"
)

// JsonReaderToXml implements JsonToXml() by wrapping MapToXml() with an io.Reader.
// Repeated calls will bulk process the stream of anonymous JSON strings.
// The function returns: XML string, pointer to source JSON value, error.
func JsonReaderToXml(rdr io.Reader, rootTag ...string) (string, *[]byte, error) {
	m, jb, err := JsonReaderToMap(rdr)
	if err != nil {
		return "", jb, err
	}
	doc, derr :=  MapToXml(m, rootTag...)
	return doc, jb, derr
}

// JsonReaderToMap wraps json.Unmarshal() with an io.Reader.
// Repeated calls will bulk process the stream of anonymous JSON strings.
// The function returns: map[string]interface{}, pointer to source JSON value, error.
func JsonReaderToMap(rdr io.Reader) (map[string]interface{}, *[]byte, error) {
	jb, err := getJson(rdr)
	if err != nil {
		return nil, jb, err
	}

	// Unmarshal the 'presumed' JSON string
	val := make(map[string]interface{}, 0)
	err = json.Unmarshal(*jb, &val)
	return val, jb, err
}

// JsonReaderToStruct - wraps json.Unmarshal to load instances of a structure.
// The function returns: pointer to source JSON value, error - structPtr holds the data.
func JsonReaderToStruct(rdr io.Reader, structPtr interface{}) (*[]byte, error) {
	jb, err := getJson(rdr)
	if err != nil {
		return jb, err
	}

	err = json.Unmarshal(*jb,structPtr)
	return jb, err
}

func getJson(rdr io.Reader) (*[]byte, error) {
	bval := make([]byte, 1)
	jb := make([]byte, 0)
	var inQuote, inJson bool
	var parenCnt int

	// scan the input for a matched set of {...}
	// json.Unmarshal will handle syntax checking.
	for {
		_, err := rdr.Read(bval)
		if err != nil {
			if err == io.EOF && inJson && parenCnt > 0 {
				return nil, errors.New("no closing } for JSON string: "+string(jb))
			}
			return nil, err
		}
		switch bval[0] {
		case '{':
			if !inQuote {
				parenCnt++
				inJson = true
			}
		case '}':
			if !inQuote {
				parenCnt--
			}
			if parenCnt < 0 {
				return nil, errors.New("closing } without opening {: "+string(jb))
			}
		case '"':
			if inQuote {
				inQuote = false
			} else {
				inQuote = true
			}
		case '\n', '\r', '\t', ' ':
			if !inQuote {
				continue
			}
		}
		if inJson {
			jb = append(jb, bval...)
			if parenCnt == 0 {
				break
			}
		}
	}

	return &jb, nil
}
