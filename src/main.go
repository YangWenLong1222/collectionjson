package main

import (
	"./cj"
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

	fmt.Println(src)

	mapSrc := map[string]interface{}{
		"collection": map[string]interface{}{
			"version": "1.0",
		},
	}

	inputData, _ = cj.ReadCollectionJson(buf)
	inputData, _ = cj.ReadCollectionJson(src)
	inputData, _ = cj.ReadCollectionJson(mapSrc)

	fmt.Println(inputData)
	fmt.Println(inputData.Collection)
	fmt.Println("=================END=================")
}

func testAppendCollectionJson() {
	fmt.Println("================START================")
	fmt.Println("testAppendCollectionJson")
	inputData.Collection.Template.Data = append(
		inputData.Collection.Template.Data,
		cj.DataType{
			"age",
			34,
			"",
		},
	)
	fmt.Println(inputData)
	fmt.Println("=================END=================")
}

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

func testWriteCollectionJson() {
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

func main() {
	testReadCollectionJson()
	testAppendCollectionJson()
	testJoinAnotherCJ()
	testWriteCollectionJson()
}
