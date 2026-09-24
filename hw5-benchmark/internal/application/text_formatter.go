package application

import (
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"
)

type TextFormatter struct{}

func NewTextFormatter() TextFormatter {
	return TextFormatter{}
}

func (formatter TextFormatter) Format(info dto.StructInfo) (string, error) {
	var b strings.Builder

	b.WriteString("Class: ")
	b.WriteString(info.Class)
	b.WriteString("\n")

	b.WriteString("Embeds:")
	if len(info.Embeds) > 0 {
		b.WriteString(" ")
		b.WriteString(strings.Join(info.Embeds, ", "))
	}
	b.WriteString("\n")

	b.WriteString("Fields:\n")
	for _, f := range info.Fields {
		b.WriteString("  - ")
		b.WriteString(f.Name)
		b.WriteString(" (")
		b.WriteString(f.Type)
		b.WriteString(")\n")
	}

	b.WriteString("Methods:\n")
	for _, m := range info.Methods {
		b.WriteString("  - ")
		b.WriteString(m.Name)
		b.WriteString("(")
		b.WriteString(strings.Join(m.Params, ", "))
		b.WriteString(")")
		if len(m.ReturnType) > 0 {
			b.WriteString(" : ")
			b.WriteString(strings.Join(m.ReturnType, ", "))
		}
		b.WriteString("\n")
	}

	b.WriteString("Hierarchy:\n")
	writeHierarchyTree(&b, info.Hierarchy, "", true)

	return b.String(), nil
}

func writeHierarchyTree(b *strings.Builder, node dto.TreeNode, prefix string, isLast bool) {
	if prefix == "" {
		b.WriteString("  ")
		b.WriteString(node.Name)
		b.WriteString("\n")
	} else {
		b.WriteString(prefix)
		if isLast {
			b.WriteString("└── ")
		} else {
			b.WriteString("├── ")
		}
		b.WriteString(node.Name)
		b.WriteString("\n")
	}

	newPrefix := prefix
	if prefix != "" {
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
	} else {
		newPrefix = "    "
	}

	for i, ch := range node.Children {
		writeHierarchyTree(b, ch, newPrefix, i == len(node.Children)-1)
	}
}
