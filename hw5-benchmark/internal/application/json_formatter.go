package application

import (
	"encoding/json"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"
)

type JsonFormatter struct{}

func NewJsonFormatter() JsonFormatter {
	return JsonFormatter{}
}

func (formatter JsonFormatter) Format(info dto.StructInfo) (string, error) {
	out := ToJson(info)

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ToJson(info dto.StructInfo) dto.JsonStructInfo {
	out := dto.JsonStructInfo{
		Class:       info.Class,
		Superclass:  "",
		Interfaces:  []string{},
		Fields:      make([]dto.JsonField, 0, len(info.Fields)),
		Methods:     make([]dto.JsonMethod, 0, len(info.Methods)),
		Annotations: []string{},
		Hierarchy:   dto.Hierarchy{},
	}

	if len(info.Embeds) > 0 {
		out.Superclass = info.Embeds[0]
	}

	for _, f := range info.Fields {
		out.Fields = append(out.Fields, dto.JsonField{
			Name: f.Name,
			Type: f.Type,
		})
	}

	for _, m := range info.Methods {
		out.Methods = append(out.Methods, dto.JsonMethod{
			Name:       m.Name,
			Params:     m.Params,
			ReturnType: m.ReturnType,
		})
	}

	if info.Hierarchy.Name != "" {
		out.Hierarchy = treeNodeToHierarchy(info.Hierarchy)
	}

	return out
}

func treeNodeToHierarchy(node dto.TreeNode) dto.Hierarchy {
	if len(node.Children) == 0 {
		return dto.Hierarchy{
			node.Name: dto.Hierarchy{},
		}
	}

	children := dto.Hierarchy{}
	for _, ch := range node.Children {
		chH := treeNodeToHierarchy(ch)
		for k, v := range chH {
			children[k] = v
		}
	}

	return dto.Hierarchy{
		node.Name: children,
	}
}
