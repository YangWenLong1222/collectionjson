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
		var tm TemplateTypeExt
		map2struct(cj.Template, &tm)
		cj.Template = tm
	case map[string]interface{}:
		isExt = false
		var ts TemplateTypeStandard
		map2struct(cj.Template, &ts)
		cj.Template = ts
	default:
		fmt.Println("Err template.") //TODO: find a better logger.
	}

	return cj, isExt, err
}

func WriteCollectionJson(cj CollectionJsonType) ([]byte, error) {
	replaceStructNilWithOne(&cj, false)
	return json.Marshal(cj)
}

/*
 * Transfer cj.template content into local struct.
 */
func (me CollectionJsonType) AbstractTo(outputData interface{}) {
	outputDataValue := reflect.ValueOf(outputData).Elem()

	switch me.Template.(type) {
	case TemplateTypeStandard:
		nv2Struct(me.Template.(TemplateTypeStandard).Data, outputDataValue)
	case TemplateTypeExt:
		sliceType := outputDataValue.Type()

		sliceValue := reflect.MakeSlice(sliceType, 1, 1) // TODO: the len and cap should inc auto.
		elemType := sliceValue.Index(0).Type()           //TODO: Is there a better way to get the type of element in slice?

		for _, item := range me.Template.(TemplateTypeExt) {
			for _, data := range item.Data {
				fmt.Println("every data", data)
			}
			dataValue := reflect.New(elemType).Elem()
			nv2Struct(item.Data, dataValue)
			sliceValue = reflect.Append(sliceValue, dataValue)
		}
		outputDataValue.Set(sliceValue.Slice(1, sliceValue.Len())) // rm the 1st empty elem.
	}
}

func nv2Struct(dataArr []DataType, outputDataValue reflect.Value) {
	for _, data := range dataArr {
		func() {
			fieldName := strings.Title(data.Name)
			field := outputDataValue.FieldByName(fieldName)
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(e)
					fmt.Println(fieldName, data.Value)
				}
			}()
			if field.IsValid() {
				var dataValue reflect.Value
				switch data.Value.(type) {
				case float64:
					dataValue = reflect.ValueOf(data.Value)
					dataValue = dataValue.Convert(field.Type())
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
		}()
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
	tmp := reflect.New(reflect.TypeOf(inputData)).Elem()
	tmp.Set(inputDataValue) // make a copy of input, so that it's accessable.
	//TODO: find a better way without making a copy to access the inputData parameter,
	//      as in C, the input param should be a local copy, and could be accessed.
	replaceStructNilWithOne(tmp.Addr().Interface(), false)
	for i := 0; i < inputDataValue.NumField(); i++ {
		var data DataType
		// elem := inputDataValue.Field(i)
		// data.Name = strings.ToLower(inputDataValue.Type().Field(i).Name)
		elem := tmp.Field(i)
		data.Name = strings.ToLower(tmp.Type().Field(i).Name)
		data.Value = ValueType(elem.Interface())
		me.Data = append(me.Data, data)
	}
	return me
}

/*
 * Automaticly generate a template data module according to the dataModule struct type.
 */
func TemplateMaker(dataModule interface{}) TemplateTypeStandard {
	// assert dataModule is a struct type.
	tmpNew := reflect.New(reflect.TypeOf(dataModule)).Elem()
	tmpNewAddr := tmpNew.Addr().Interface()
	replaceStructNilWithOne(tmpNewAddr, true)
	item := ConcreteFrom(tmpNew.Interface(), URIType("not a real one"))
	var ret TemplateTypeStandard
	ret.Data = item.Data
	return ret
}

/*
 * p is the pointer of a struct which should not have nil as some fields' value.
 * if isSliceExt is true, the nil slice would have one element which won't have a nil too.
 * and if not, there would be an empty slice to replace the nil slice.
 */
func replaceStructNilWithOne(p interface{}, isSliceExt bool) {
	pValue := reflect.ValueOf(p).Elem()
	if pValue.Kind() == reflect.Struct {
		num := pValue.NumField()
		for i := 0; i < num; i++ {
			pNext := pValue.Field(i).Addr().Interface()
			replaceStructNilWithOne(pNext, isSliceExt)
		}
		return
	}

	var valueToSet reflect.Value
	var couldNotBeNil bool
	switch pValue.Kind() {
	case reflect.Chan:
		valueToSet = reflect.MakeChan(pValue.Type(), 1)
	case reflect.Func:
		valueToSet = reflect.MakeFunc(pValue.Type(), func(args []reflect.Value) (results []reflect.Value) {
			var ret []reflect.Value
			return ret
		})
	case reflect.Map:
		valueToSet = reflect.MakeMap(pValue.Type())

	case reflect.Slice:
		if isSliceExt {
			valueToSet = reflect.MakeSlice(pValue.Type(), 1, 1)
			replaceStructNilWithOne(valueToSet.Index(0).Addr().Interface(), isSliceExt)
		} else {
			valueToSet = reflect.MakeSlice(pValue.Type(), 0, 1)
		}

	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		valueToSet = reflect.New(pValue.Type()).Elem()

	default:
		couldNotBeNil = true
	}

	if !couldNotBeNil && pValue.IsNil() {
		pValue.Set(valueToSet)
		return
	}
	return
}
