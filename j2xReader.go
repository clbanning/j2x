// j2xReader.go - wrap j2X with an io.Reader input option

package j2x

import (
	"encoding/json"
	"io"
)

// JsonReaderToDoc implements JsonToDoc() by wrapping MapToDoc() with an io.Reader.
// Repeated calls will bulk process the stream of anonymous JSON strings.
func JsonReaderToDoc(rdr io.Reader, rootTag ...string) (string, error) {
	m, err := JsonReaderToMap(rdr)
	if err != nil {
		return "", err
	}
	return MapToDoc(m, rootTag...)
}

// JsonReaderToDoc wraps json.Unmarshal() with an io.Reader.
// Repeated calls will bulk process the stream of anonymous JSON strings.
func JsonReaderToMap(rdr io.Reader) (map[string]interface{}, error) {
	bval := make([]byte, 1)
	jb := make([]byte, 0)
	var inQuote, inJson bool
	var parenCnt int

	// scan the input for a matched set of {...}
	// json.Unmarshal will handle syntax checking.
	for {
		_, err := rdr.Read(bval)
		if err != nil {
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

	// Unmarshal the 'presumed' JSON string
	val := make(map[string]interface{}, 0)
	err := json.Unmarshal(jb, &val)
	return val, err
}
