package dto

type Hierarchy map[string]Hierarchy

type JsonStructInfo struct {
	Class       string       `json:"class"`
	Superclass  string       `json:"superclass"`
	Interfaces  []string     `json:"interfaces"`
	Fields      []JsonField  `json:"fields"`
	Methods     []JsonMethod `json:"methods"`
	Annotations []string     `json:"annotations"`
	Hierarchy   Hierarchy    `json:"hierarchy"`
}

type JsonField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type JsonMethod struct {
	Name       string   `json:"name"`
	Params     []string `json:"params"`
	ReturnType []string `json:"returnType"`
}
