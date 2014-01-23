package j2x

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMashalIndent(t *testing.T) {
	var s = `{ "head":[ "one", 2, true, { "key":"value" } ] }`

	fmt.Println("\nTestMashalIndent ... list :", s)
	v, err := MarshalIndent(s,"  ","   ")
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Printf("v:\n%s",string(v))

	s = `{ "head":{ "line":[ "one", 2, true, { "key":"value" } ] } }`
	m := make(map[string]interface{},0)
	err = json.Unmarshal([]byte(s), &m)
	type mystruct struct {
		S string
		F float64
	}
	ms := mystruct{ S:"now's the time", F:3.14159625 }
	m["mystruct"] = interface{}(ms)
	fmt.Println("\nTestMarshalIndent ... mystruct", m)
	v, err = MarshalIndent(m,"   ","  ")
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Printf("v:\n%s",string(v))
}
