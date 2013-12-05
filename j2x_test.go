package j2x

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	var s = `{ "key":"value" }`

	fmt.Println("\nTestSimple ... JsonToDoc:",s)
	v, err := JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	fmt.Println("\nTestSimple ... JsonToDoc, rootTag: zoom")
	v, err = JsonToDoc(s,"zoom")
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	s = `{ "one":1, "two":1.999, "3":"three", "four":false }`

	fmt.Println("\nTestSimple ... JsonToDoc:",s)
	v, err = JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)
}

func TestNotSoSimple(t *testing.T) {
	var s = `{ "json":{ "one":1, "pi":3.1415962535, "bool":true, "jsonJR":{ "key":"value" } } }`

	fmt.Println("\nTestNotSoSimple ... JsonToDoc:",s)
	v, err := JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	s = `{ "json":[ "one", 3.1415962535, true, { "key":"value" } ] }`

	fmt.Println("\nTestNotSoSimple ... JsonToDoc:",s)
	v, err = JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)
}

func TestAttr(t *testing.T) {
	var s = `{ "json":{ "-one":1, "-pi":3.1415962535, "-bool":true, "jsonJR":{ "-key":"value" } } }`

	fmt.Println("\nTestAttr ... JsonToDoc:",s)
	v, err := JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	s = `{ "json":{ "-one":1, "-pi":3.1415962535, "-bool":true, "jsonJR":{ "-attr":"value", "#text":"value" } } }`

	fmt.Println("\nTestAttr ... #test:",s)
	v, err = JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	s = `{ "json":[ "one", 3.1415962535, true, { "-key":"value" } ] }`

	fmt.Println("\nTestAttr ... list:",s)
	v, err = JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)

	s = `{ "json":[ "one", 3.1415962535, true, { "-key":"value", "#text":"Now is the time..." } ] }`

	fmt.Println("\nTestAttr ... #text:",s)
	v, err = JsonToDoc(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",v)
}

func TestGoofy(t *testing.T) {
	var s = `{ "json":{ "-one":1, "-pi":3.1415962535, "-bool":true, "jsonJR":{ "-key":"value" } } }`
	type goofy struct {
		S string
		Sp *string
	}
	g := new(goofy)
	g.S = "Now is the time for"
	tmp := "all good men to come to"
	g.Sp = &tmp

	m := make(map[string]interface{},0)
	_ = json.Unmarshal([]byte(s),&m)

	m["goofyVal"] = interface{}(g)
	m["byteVal"] = interface{}([]byte(`the aid of their country`))
	m["nilVal"] = interface{}(nil)

	fmt.Println("\nTestGoofy ... MapToDoc:",m)
	v, _ := MapToDoc(m)
	fmt.Println("v:",v)

	type goofier struct {
		G *goofy
		B []byte
		N *string
	}
	gg := new(goofier)
	gg.G = g
	gg.B = []byte(`the tree of freedom must periodically be`)
	gg.N = nil
	m["goofierVal"] = interface{}(gg)

	fmt.Println("\nTestGoofier ... MapToDoc:",m)
	v, _ = MapToDoc(m)
	fmt.Println("v:",v)
}

func TestMarshal(t *testing.T) {
	var s = `{ "json":{ "-one":1, "-pi":3.1415962535, "-bool":true, "jsonJR":{ "-key":"value" } } }`
	type goofy struct {
		S string
		Sp *string
	}
	g := new(goofy)
	g.S = "Now is the time for"
	tmp := "all good men to come to"
	g.Sp = &tmp

	m := make(map[string]interface{},0)
	_ = json.Unmarshal([]byte(s),&m)

	m["goofyVal"] = interface{}(g)

	fmt.Println("\nTestMarshal ... :",s)
	v, err := Marshal(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	fmt.Println("\nTestMarshal ... :",g)
	v, err = Marshal(g)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	fmt.Println("\nTestMarshal ... :",g.Sp)
	v, err = Marshal(g.Sp)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	fmt.Println("\nTestMarshal ... :",[]byte(g.S))
	v, err = Marshal([]byte(g.S))
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	fmt.Println("\nTestMarshal ... :",m)
	v, err = Marshal(m)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))
}

func TestSingleRootKey(t *testing.T) {
	var s = `{ "head":[ "one", 2, true, { "key":"value" } ] }`

	fmt.Println("\nTestSingleRootKey ... list :", s)
	v, err := Marshal(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	s = `{ "head":{ "line":[ "one", 2, true, { "key":"value" } ] } }`
	fmt.Println("\nTestSingleRootKey ... JSON:", s)
	v, err = Marshal(s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))
}

func TestBangTextError(t *testing.T) {
	var s = `{ "-attr":"value", "#text":true }`

	m := make(map[string]interface{},0)
	_ = json.Unmarshal([]byte(s),&m)
	m["something"] = interface{}("else")

	fmt.Println("\nTestBangTextError ... map :", m)
	v, err := Marshal(m)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))
}
