/*
 * Process all the boring things of cj for you.
 * Keep away from collection-json, save your life.
 * I love life, but I hate collection-json.
 * Something more than CJ standard:
 *     1, The value could be Array or Map type which could NOT according to CJ standard.
 *     2, Template could accept an Array to create multi-items together.
 */
package cj

/*
 * The struct which is used to replay the client GET request by server.
 */
type CollectionJsonType struct {
	Collection CollectionType `json:"collection"` // REQUIRED
	// Queries    []QueryType    `json:"queries"`    // OPTIONAL ? top-level
	//TODO: Queries should use another structure.
}

/*
 * The struct which is used to describe a request from client.
 */
type CollectionJsonTemplateType struct {
	Template TemplateType `json:"template"`
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
