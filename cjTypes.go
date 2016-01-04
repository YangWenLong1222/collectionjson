/*
 * Process all the boring things of cj for you.
 * Keep away from collection-json, save your life.
 * I love life, but I hate collection-json.
 * Something more than CJ standard:
 *     1, The value could be Array or Map type which could NOT according to CJ standard.
 *     2, Template could accept an Array to create multi-items together.
 */
package cj

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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

type TemplateType interface{}

type TemplateTypeStandard struct {
	Data []DataType `json:"data"`
}

type TemplateTypeExt []struct {
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
 * ret val:
 * cj    -- The result cj struct val.
 * isExt -- Tell the caller this cj is an Extension of cj or not.
 * err   -- Any err when reading cj content.
 */
func ReadCollectionJson(inputData interface{}) (CollectionJsonType, bool, error) {
	var cj CollectionJsonType
	var isExt bool
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

	if err != nil {
		return cj, false, err
	}

	switch cj.Template.(type) {
	case []interface{}:
		isExt = true
	case map[string]interface{}:
		isExt = false
	default:
		fmt.Println("Err template.") //TODO: find a better logger.
	}

	return cj, isExt, err
}

func WriteCollectionJson(cj CollectionJsonType) ([]byte, error) {
	return json.Marshal(cj)
}

/*
 * Transfer cj.template content into local struct.
 */
func (me CollectionJsonType) AbstractTo(outputData interface{}) {
	outputDataValue := reflect.ValueOf(outputData).Elem()

	switch me.Template.(type) {
	case map[string]interface{}:
		fmt.Println("runs here.map")
		var ts TemplateTypeStandard
		map2struct(me.Template, &ts)
		fmt.Println(ts)
		nv2Struct(ts.Data, outputDataValue)
	case []interface{}:
		fmt.Println("runs here.multi")
		sliceType := outputDataValue.Type()
		fmt.Println(sliceType)

		sliceValue := reflect.MakeSlice(sliceType, 1, 1) // TODO: the len and cap should inc auto.
		elemType := sliceValue.Index(0).Type()
		fmt.Println(elemType)

		var tm TemplateTypeExt
		map2struct(me.Template, &tm)
		me.Template = tm
		for _, item := range tm {
			for _, data := range item.Data {
				fmt.Println("every data", data)
			}
			dataValue := reflect.New(elemType).Elem()
			nv2Struct(item.Data, dataValue)
			sliceValue = reflect.Append(sliceValue, dataValue)
		}
		outputDataValue.Set(sliceValue.Slice(1, sliceValue.Len()))
	}
}

func nv2Struct(dataArr []DataType, outputDataValue reflect.Value) {
	for _, data := range dataArr {
		fieldName := strings.Title(data.Name)
		field := outputDataValue.FieldByName(fieldName)
		if field.IsValid() {
			var dataValue reflect.Value
			switch data.Value.(type) {
			case []interface{}:
				dataValue = reflect.New(field.Type())
				map2struct(data.Value, dataValue.Interface())
				dataValue = dataValue.Elem()
			default:
				dataValue = reflect.ValueOf(data.Value)
			}
			field.Set(dataValue)
		} else {
			fmt.Println("no field: " + fieldName) //TODO: use a better logger.
		}
	}
}

/*
 * Maybe there would be a better way using reflect but not json.
 * Actually json use reflect tech too.
 */
func map2struct(src interface{}, destPointer interface{}) {
	buf, _ := json.Marshal(src) // TODO: err process.
	json.Unmarshal(buf, destPointer)
}

/*
 * Change the item according to the input local struct and href string.
 * Ignore the links in cj, since it is only OPTIONAL.
 * panic if there is any err.
 */
func ConcreteFrom(inputData interface{}, href URIType) ItemType {
	var me ItemType
	me.Href = URIType(href)
	inputDataValue := reflect.ValueOf(inputData)
	for i := 0; i < inputDataValue.NumField(); i++ {
		var data DataType
		elem := inputDataValue.Field(i)
		data.Name = strings.ToLower(inputDataValue.Type().Field(i).Name)
		data.Value = ValueType(elem.Interface())
		me.Data = append(me.Data, data)
	}
	return me
}

/*
 * Set field of me with the value of the same field of cj, if this field of me is empty("" or nil)
 * Append array value of cj into the same field of me.
 */
func (me *CollectionJsonType) JoinMe(cj CollectionJsonType) {
	merge(me, cj)
	return
}

func merge(me interface{}, cj interface{}) {
	//TODO: 判断me是否是一个指针。
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
