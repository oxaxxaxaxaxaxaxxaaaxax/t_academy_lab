package dto

type StructInfo struct {
	Class     string
	Embeds    []string
	Fields    []Field
	Methods   []Method
	Hierarchy TreeNode
}

type Field struct {
	Name string
	Type string
}

type Method struct {
	Name       string
	Params     []string
	ReturnType []string
}

type TreeNode struct {
	Name     string
	Children []TreeNode
}
