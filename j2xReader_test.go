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

func TestReaderToStruct(t *testing.T) {
	data := []byte(`{"Key1":"value1", "Key2":3.14159625, "Key3":true},
						 {"Key1":"value2", "Key2":31.4159625, "Key3":false}`)

	type tstruct struct {
		Key1 string
		Key2 float64
		Key3 bool
	}

	fmt.Println("\ndata for structs:", string(data))
	r := bytes.NewReader(data)
	for {
		v := new(tstruct)
		err := JsonReaderToStruct(r,v)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error("err:",err.Error())
		}
		fmt.Println("v:",v)
	}
}
