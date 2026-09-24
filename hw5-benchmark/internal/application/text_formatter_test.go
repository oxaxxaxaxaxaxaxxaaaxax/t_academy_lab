package application

import (
	"strings"
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatter_Format(t *testing.T) {
	tests := []struct {
		name         string
		info         dto.StructInfo
		expectedText string
		expectedErr  error
	}{
		{
			name: "valid StructInfo",
			info: dto.StructInfo{
				Class:  "Manager",
				Embeds: []string{"Employee", "Person"},
				Fields: []dto.Field{
					{Name: "Name", Type: "string"},
					{Name: "Age", Type: "int"},
				},
				Methods: []dto.Method{
					{Name: "Work", Params: []string{"int"}, ReturnType: []string{"error"}},
					{Name: "Ping", Params: nil, ReturnType: nil},
				},
				Hierarchy: dto.TreeNode{
					Name: "Manager",
					Children: []dto.TreeNode{
						{
							Name: "Lead",
							Children: []dto.TreeNode{
								{Name: "Junior"},
							},
						},
						{Name: "Intern"},
					},
				},
			},
			expectedText: "" +
				"Class: Manager\n" +
				"Embeds: Employee, Person\n" +
				"Fields:\n" +
				"  - Name (string)\n" +
				"  - Age (int)\n" +
				"Methods:\n" +
				"  - Work(int) : error\n" +
				"  - Ping()\n" +
				"Hierarchy:\n" +
				"  Manager\n" +
				"    ├── Lead\n" +
				"    │   └── Junior\n" +
				"    └── Intern\n",
			expectedErr: nil,
		},
		{
			name: "empty StructInfo",
			info: dto.StructInfo{
				Class:   "Empty",
				Embeds:  nil,
				Fields:  nil,
				Methods: nil,
				Hierarchy: dto.TreeNode{
					Name:     "Empty",
					Children: nil,
				},
			},
			expectedText: "" +
				"Class: Empty\n" +
				"Embeds:\n" +
				"Fields:\n" +
				"Methods:\n" +
				"Hierarchy:\n" +
				"  Empty\n",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewTextFormatter()

			out, err := f.Format(tt.info)

			assert.ErrorIs(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedText, out)
		})
	}
}

func TestWriteHierarchyTree(t *testing.T) {
	tests := []struct {
		name     string
		node     dto.TreeNode
		expected string
	}{
		{
			name: "single node",
			node: dto.TreeNode{Name: "Root"},
			expected: "" +
				"  Root\n",
		},
		{
			name: "two children",
			node: dto.TreeNode{
				Name: "Root",
				Children: []dto.TreeNode{
					{Name: "A"},
					{Name: "B"},
				},
			},
			expected: "" +
				"  Root\n" +
				"    ├── A\n" +
				"    └── B\n",
		},
		{
			name: "many children",
			node: dto.TreeNode{
				Name: "Root",
				Children: []dto.TreeNode{
					{
						Name: "A",
						Children: []dto.TreeNode{
							{Name: "A1"},
							{Name: "A2"},
						},
					},
					{Name: "B"},
				},
			},
			expected: "" +
				"  Root\n" +
				"    ├── A\n" +
				"    │   ├── A1\n" +
				"    │   └── A2\n" +
				"    └── B\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			writeHierarchyTree(&b, tt.node, "", true)
			assert.Equal(t, tt.expected, b.String())
		})
	}
}
