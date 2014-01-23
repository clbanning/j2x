package j2x

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	var data [3][]byte
	data[0] = []byte(`{"here":"is", "the_first":1, "json":"string"}{"here":"is", "the_second":2, "string":true}`)
	data[1] = []byte(`{"here":"is", "the_first":1, "json":"string"}, {"here":"is", "the_second":2, "string":true}`)
	data[2] = []byte(` {"here":"is", "the_first":1, "json":"string"}
							  {"here":"is", "the_second":2, "string":true }
							`)

	fmt.Println("\nj2xwriter_test.go ... TestWriter")
	buf := make([]byte,1024)
	w := bytes.NewBuffer(buf)

	for i := 0 ; i < len(data) ; i++ {
		fmt.Println("\ndata:",i,"string:",string(data[i]))
		r := bytes.NewReader(data[i])
		for {
			m, jb, err := JsonReaderToMap(r)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Error("data:",i,"err:",err.Error())
				continue
			}
			fmt.Println("data:",i,"jb:",string(*jb),"map",m)

			ps, err := MapToXmlWriter(m,w)
			fmt.Println("*ps:",*ps)
			s := make([]byte,w.Len())
			_, err = w.Read(s)
			fmt.Println("s  :",string(s))
		}
	}
}

func TestReadWriter(t *testing.T) {
	var data [3][]byte
	data[0] = []byte(`{"here":"is", "the_first":1, "json":"string"}{"here":"is", "the_second":2, "string":true}`)
	data[1] = []byte(`{"here":"is", "the_first":1, "json":"string"}, {"here":"is", "the_second":2, "string":true}`)
	data[2] = []byte(` {"here":"is", "the_first":1, "json":"string"}
							  {"here":"is", "the_second":2, "string":true }
							`)

	fmt.Println("\nj2xwriter_test.go ... TestReadWriter")
	buf := make([]byte,1024)
	w := bytes.NewBuffer(buf)

	for i := 0 ; i < len(data) ; i++ {
		fmt.Println("\ndata:",i,"string:",string(data[i]))
		r := bytes.NewReader(data[i])
		for {
			j, x, err := JsonReaderToXmlWriter(r,w)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Error("data:",i,"err:",err.Error())
				continue
			}
			fmt.Println("j:",*j)
			fmt.Println("x:",*x)
			s := make([]byte,w.Len())
			_, err = w.Read(s)
			fmt.Println("s:",string(s))
		}
	}
}

