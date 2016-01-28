package main

import (
	".."
	"encoding/json"
	"fmt"
	"os"
)

const (
	INPUT_FILE_NAME  = "input.txt"
	OUTPUT_FILE_NAME = "output.txt"
	FILE_SIZE_LIMIT  = 1024 * 1024 // 1M
)

var inputData cj.CollectionJsonType

func testReadCollectionJson() {
	fmt.Println("================START================")
	fmt.Println("testReadCollectionJson")
	file, err := os.Open(INPUT_FILE_NAME) // This kind Open is just for reading.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	var n int
	buf := make([]byte, FILE_SIZE_LIMIT)
	n, err = file.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(n)
	buf = buf[:n]

	src := string(buf)

	fmt.Println(src)

	mapSrc := map[string]interface{}{
		"collection": map[string]interface{}{
			"version": "1.0",
		},
	}

	inputData, _, _ = cj.ReadCollectionJson(buf)
	inputData, _, _ = cj.ReadCollectionJson(src)
	inputData, _, _ = cj.ReadCollectionJson(mapSrc)

	fmt.Println(inputData)
	fmt.Println(inputData.Collection)
	fmt.Println("=================END=================")
}

func testAppendCollectionJson() {
	fmt.Println("================START================")
	fmt.Println("testAppendCollectionJson")
	templ := inputData.Collection.Template.(cj.TemplateTypeStandard)
	templ.Data = append(
		templ.Data,
		cj.DataType{
			"age",
			34,
			"",
		},
	)
	inputData.Collection.Template = templ
	fmt.Println(inputData)
	fmt.Println("=================END=================")
}

/*
func testJoinAnotherCJ() {
	fmt.Println("================START================")
	fmt.Println("testJoinAnotherCJ")
	var data1 cj.CollectionJsonType
	data1.Collection.Template.Data = []cj.DataType{
		{
			"company",
			"zmeng",
			"",
		},
	}

	data1.Collection.Href = "www.sina.com.cn"
	data1.Links = []cj.LinkType{
		{
			"http://www.163.com",
			"rel string",
			"netease",
			"what is render.",
			"",
		},
	}

	fmt.Println(data1)
	inputData.JoinMe(data1)
	fmt.Println(inputData)
	fmt.Println("=================END=================")
}
*/

func testWriteCollectionJson(inputData cj.CollectionJsonType) {
	fmt.Println("================START================")
	fmt.Println("testWriteCollectionJson")
	retBytes, _ := cj.WriteCollectionJson(inputData)

	fmt.Println(string(retBytes))

	file, err := os.Create(OUTPUT_FILE_NAME)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	// I don't know if this is nessary.
	// There is also no this byte in real response now.
	// However, you'd better add this byte if you want to write it into a file, to avoid a bad shown of a text file.
	retBytes = append(retBytes, 0x0A)

	_, err = file.Write(retBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("=================END=================")
}

func testTemplateArray() {
	fmt.Println("================START================")
	fmt.Println("testWriteCollectionJson")
	file, err := os.Open(INPUT_FILE_NAME)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	var n int
	buf := make([]byte, FILE_SIZE_LIMIT)
	n, err = file.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(n)
	buf = buf[:n]

	src := string(buf)
	// fmt.Println(src)

	inputData, _, err = cj.ReadCollectionJson(src)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(inputData)

	fmt.Println("=================END=================")
}

func testAbstractTo() {
	fmt.Println("================START================")
	fmt.Println("testAbstractTo")
	file, err := os.Open(INPUT_FILE_NAME) // This kind Open is just for reading.
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	buf := make([]byte, FILE_SIZE_LIMIT)
	var n int
	n, err = file.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf = buf[:n]
	srcData, isExt, err1 := cj.ReadCollectionJson(buf)
	fmt.Println(srcData)
	fmt.Println(isExt)
	fmt.Println(err1)

	type tarType struct {
		Name    string
		Address string
		age     int
	}
	var tar tarType
	tar.age = 33
	var tarArr []tarType
	if isExt {
		fmt.Println("Extension")
		srcData.AbstractTo(&tarArr)
		if err != nil {
			fmt.Println("ab err")
		}
		fmt.Println(tarArr)
	} else {
		fmt.Println("Standard")
		srcData.AbstractTo(&tar)
		if err != nil {
			fmt.Println("ab err")
		}
		fmt.Println(tar)
	}

	testWriteCollectionJson(srcData)
	fmt.Println("=================END=================")
}

func testConcreteFrom() {
	fmt.Println("================START================")
	fmt.Println("testConcreteFrom")
	src := struct {
		Merchant int
		Name     string
	}{62, "RD_TEST"}
	fmt.Println(src)
	dest := cj.ConcreteFrom(src, "llala")
	fmt.Println(dest)
	var ccjj cj.CollectionJsonType
	ccjj.Collection.Items = []cj.ItemType{}
	ccjj.Collection.Items = append(ccjj.Collection.Items, dest)
	buf, _ := cj.WriteCollectionJson(ccjj)
	fmt.Println(string(buf))

	fmt.Println("=================END=================")
}

func testOther() {
	src := []map[string]interface{}{
		{
			"title":   "first rule",
			"checked": 1,
		},
		{
			"title":   "second rule",
			"checked": 0,
		},
	}

	type RuleType struct {
		Title   string
		Checked int
		Promote string
	}

	fmt.Println(src)
	buf, _ := json.Marshal(src)
	fmt.Println(src)
	fmt.Println(buf)
	fmt.Println(string(buf))

	res := []RuleType{}
	json.Unmarshal(buf, &res)
	fmt.Println(res)
}

func testTemplateMaker() {
	fmt.Println("================START================")
	fmt.Println("testTemplateMaker")
	src := struct {
		Merchant int
		Name     string
	}{62, "RD_TEST"}

	ret := cj.TemplateMaker(src)
	var cjRet cj.CollectionJsonType
	cjRet

	fmt.Println(ret)

	fmt.Println("=================END=================")
}

func main() {
	// testReadCollectionJson()
	// testAppendCollectionJson()
	// testJoinAnotherCJ()
	// testTemplateArray()
	// testAbstractTo()
	// testWriteCollectionJson()
	// testConcreteFrom()
	// testOther()
	testTemplateMaker()
}
