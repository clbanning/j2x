// j2x_writer.go - wrap j2X with an io.Reader input option

package j2x

import (
	"io"
)

// JsonToXmlWriter decodes JSON string and writes it using io.Writer
// Returns pointer to encoded XML, error.
func JsonToXmlWriter(s string, wtr io.Writer) (*string, error) {
	x, err := JsonToXml(s)
	if err != nil {
		return nil, err
	}

	_, err = wtr.Write([]byte(x))
	return &x, err
}

// MapToXmlWriter encodes the map as XML and writes in on the io.Writer.
// The function returns: pointer to encoded XML, error.
func MapToXmlWriter(m map[string]interface{}, wtr io.Writer) (*string, error) {
	x, err := MapToXml(m)
	if err != nil {
		return nil, err
	}

	_, err = wtr.Write([]byte(x))
	return &x, err
}

// Decodes next value from a JSON io.Reader and writes it using io.Writer
// Returns: pointer to JSON, pointer to encoded XML, error.
func JsonReaderToXmlWriter(rtr io.Reader, wtr io.Writer, rootTag ...string) (*string, *string, error) {
	rt := DefaultRootTag
	if len(rootTag) == 1 {
		rt = rootTag[0]
	}

	doc, jval, err := JsonReaderToXml(rtr,rt)
	if err != nil {
		return nil, nil, err
	}

	j := string(*jval)
	if err != nil {
		return &j, &doc, err
	}

	_, err = wtr.Write([]byte(doc))
	return &j, &doc, err
}

