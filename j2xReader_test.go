package j2x

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	var data [3][]byte
	data[0] = []byte(`{"here":"is", "the_first":1, "json":"string"}{"here":"is", "the_second":2, "string":true}`)
	data[1] = []byte(`{"here":"is", "the_first":1, "json":"string"}, {"here":"is", "the_second":2, "string":true}`)
	data[2] = []byte(` {"here":"is", "the_first":1, "json":"string"}
							  {"here":"is", "the_second":2, "string":true }
							`)

	fmt.Println("\nj2xReader_test.go ...")

	for i := 0 ; i < len(data) ; i++ {
		fmt.Println("\ndata:",i,"string:",string(data[i]))
		r := bytes.NewReader(data[i])
		for {
			m, err := JsonReaderToMap(r)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Error("data:",i,"err:",err.Error())
				continue
			}
			fmt.Println("data:",i,"map",m)
		}
		fmt.Println("\ndata:",i,"string:",string(data[i]))
		r = bytes.NewReader(data[i])
		for {
			d, err := JsonReaderToDoc(r)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Error("data:",i,"err:",err.Error())
				continue
			}
			fmt.Println("data:",i,"doc",d)
		}
	}
}
