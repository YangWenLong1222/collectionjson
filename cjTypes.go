package cj

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type CollectionJsonType struct {
	Collection CollectionType `json:"collection"` // REQUIRED
	Queries    []QueryType    `json:"queries"`    // OPTIONAL ? top-level
	Links      []LinkType     `json:"links"`      // OPTIONAL ? top-level
	Template   TemplateType   `json:"template"`   // REQUIRED when it's a request.
}

type CollectionType struct {
	Version  string       `json:"version"` //TODO: always be 1.0, how can do it?
	Href     URIType      `json:"href"`
	Links    []LinkType   `json:"links"`
	Items    []ItemType   `json:"items"`
	Queries  []QueryType  `json:"queries"`
	Template TemplateType `json:"template"`
	Error    ErrorType    `json:"error"`
}

type LinkType struct {
	Href   URIType `json:"href"`   // REQUIRED
	Rel    string  `json:"rel"`    // REQUIRED
	Name   string  `json:"name"`   // OPTIONAL
	Render string  `json:"render"` // OPTIONAL MUST be "image" or "link"
	Prompt string  `json:"prompt"` // OPTIONAL
}

type ItemType struct {
	Href  URIType    `json:"href"`
	Data  []DataType `json:"data"`
	Links []LinkType `json:"links"` // OPTIONAL
}

type URIType string

type QueryType struct {
	Href   URIType    `json:"href"`   // REQUIRED
	Rel    string     `json:"rel"`    // REQUIRED
	Name   string     `json:"name"`   // OPTIONAL
	Prompt string     `json:"prompt"` // OPTIONAL
	Data   []DataType `json:"data"`   // OPTIONAL
}

type TemplateType struct {
	Data []DataType `json:"data"`
}

type DataType struct {
	Name   string    `json:"name"`   // REQUIRED
	Value  ValueType `json:"value"`  // OPTIONAL
	Prompt string    `json:"prompt"` // OPTIONAL
}

type ValueType interface{}

type ErrorType struct {
	Title   string `json:"title"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

/*
 * A very simple pack of json package.
 * Let callers need not import json package in some case.
 */
func ReadCollectionJson(inputData interface{}) (CollectionJsonType, error) {
	var cj CollectionJsonType
	var err error
	var buf []byte

	switch inputData.(type) {
	case string:
		buf = []byte(inputData.(string))
	case []byte:
		buf = inputData.([]byte)
	case map[string]interface{}:
		buf, err = json.Marshal(inputData.(map[string]interface{}))
	}

	if err == nil {
		err = json.Unmarshal(buf, &cj)
	}
	return cj, err
}

func WriteCollectionJson(cj CollectionJsonType) ([]byte, error) {
	return json.Marshal(cj)
}

/*
 * Set field of me with the value of the same field of cj, if this field of me is empty("" or nil)
 * Append array value of cj into the same field of me.
 * Using JoinMe() instead of using merge() directly, Since there may some special logical in some type.
 */
func (me *CollectionJsonType) JoinMe(cj CollectionJsonType) {
	merge(me, cj)
	return
}

func (me *ItemType) JoinMe(cj ItemType) {
	merge(me, cj)
	return
}

func merge(me interface{}, cj interface{}) {
	//TODO: 判断me是否是一个指针。如果不是报错。
	//TODO: 判断me指向的类型和cj是否是一个类型，如果不是，报错。如果要放开，必须要做这个检查，如果不，则不用了。
	meValue := reflect.ValueOf(me).Elem()

	switch meValue.Kind() {

	case reflect.String:
		value := meValue.String()
		if value == "" {
			meValue.Set(reflect.ValueOf(cj))
		}

	case reflect.Slice:
		//TODO: 不能简单的append，需要去掉重复的，依据name。
		meValue.Set(reflect.AppendSlice(meValue, reflect.ValueOf(cj)))

	case reflect.Struct:
		num := meValue.NumField()
		for i := 0; i < num; i++ {
			pField := meValue.Field(i).Addr().Interface()
			merge(pField, reflect.ValueOf(cj).Field(i).Interface())
		}

	case reflect.Interface:
		value := meValue.Interface()
		if value == nil {
			meValue.Set(reflect.ValueOf(cj))
		}

	default:
		fmt.Println("default")
	}
	return
}
