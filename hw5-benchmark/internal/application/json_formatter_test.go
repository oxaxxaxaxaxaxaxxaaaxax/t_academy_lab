package application

import (
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"

	"github.com/stretchr/testify/assert"
)

func TestJsonFormatter_Format(t *testing.T) {
	tests := []struct {
		name         string
		info         dto.StructInfo
		expectedJSON string
		expectedErr  error
	}{
		{
			name: "valid StructInfo",
			info: dto.StructInfo{
				Class:  "Manager",
				Embeds: []string{"Employee"},
				Fields: []dto.Field{
					{Name: "Name", Type: "string"},
					{Name: "Age", Type: "int"},
				},
				Methods: []dto.Method{
					{Name: "Work", Params: []string{"int"}, ReturnType: []string{"error"}},
				},
				Hierarchy: dto.TreeNode{
					Name: "Manager",
					Children: []dto.TreeNode{
						{Name: "Lead", Children: nil},
					},
				},
			},
			expectedErr: nil,
			expectedJSON: `{
  "class": "Manager",
  "superclass": "Employee",
  "interfaces": [],
  "fields": [
    {
      "name": "Name",
      "type": "string"
    },
    {
      "name": "Age",
      "type": "int"
    }
  ],
  "methods": [
    {
      "name": "Work",
      "params": [
        "int"
      ],
      "returnType": [
        "error"
      ]
    }
  ],
  "annotations": [],
  "hierarchy": {
    "Manager": {
      "Lead": {}
    }
  }
}`,
		},
		{
			name:        "empty StructInfo",
			info:        dto.StructInfo{},
			expectedErr: nil,
			expectedJSON: `{
  "class": "",
  "superclass": "",
  "interfaces": [],
  "fields": [],
  "methods": [],
  "annotations": [],
  "hierarchy": {}
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewJsonFormatter()

			out, err := f.Format(test.info)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedJSON, out)

		})
	}
}

func TestToJson(t *testing.T) {
	tests := []struct {
		name     string
		info     dto.StructInfo
		expected dto.JsonStructInfo
	}{
		{
			name: "simple valid test",
			info: dto.StructInfo{Class: "User"},
			expected: dto.JsonStructInfo{
				Class:       "User",
				Superclass:  "",
				Interfaces:  []string{},
				Fields:      []dto.JsonField{},
				Methods:     []dto.JsonMethod{},
				Annotations: []string{},
				Hierarchy:   dto.Hierarchy{},
			},
		},
		{
			name: "valid test",
			info: dto.StructInfo{
				Class:  "Manager",
				Embeds: []string{"Employee", "Person"},
				Fields: []dto.Field{
					{Name: "A", Type: "int"},
					{Name: "B", Type: "string"},
				},
				Methods: []dto.Method{
					{Name: "Foo", Params: []string{"int"}, ReturnType: []string{"bool"}},
					{Name: "Bar", Params: nil, ReturnType: []string{"error"}},
				},
			},
			expected: dto.JsonStructInfo{
				Class:       "Manager",
				Superclass:  "Employee",
				Interfaces:  []string{},
				Annotations: []string{},
				Fields: []dto.JsonField{
					{Name: "A", Type: "int"},
					{Name: "B", Type: "string"},
				},
				Methods: []dto.JsonMethod{
					{Name: "Foo", Params: []string{"int"}, ReturnType: []string{"bool"}},
					{Name: "Bar", Params: nil, ReturnType: []string{"error"}},
				},
				Hierarchy: dto.Hierarchy{},
			},
		},
		{
			name: "valid with hierarchy",
			info: dto.StructInfo{
				Class: "Root",
				Hierarchy: dto.TreeNode{
					Name: "Root",
					Children: []dto.TreeNode{
						{
							Name: "Child1",
							Children: []dto.TreeNode{
								{Name: "GrandChild"},
							},
						},
						{Name: "Child2"},
					},
				},
			},
			expected: dto.JsonStructInfo{
				Class:       "Root",
				Superclass:  "",
				Interfaces:  []string{},
				Fields:      []dto.JsonField{},
				Methods:     []dto.JsonMethod{},
				Annotations: []string{},
				Hierarchy: dto.Hierarchy{
					"Root": dto.Hierarchy{
						"Child1": dto.Hierarchy{
							"GrandChild": dto.Hierarchy{},
						},
						"Child2": dto.Hierarchy{},
					},
				},
			},
		},
		{
			name: "hierarchy with empty name",
			info: dto.StructInfo{
				Class:     "NoHierarchy",
				Hierarchy: dto.TreeNode{Name: "", Children: []dto.TreeNode{{Name: "X"}}},
			},
			expected: dto.JsonStructInfo{
				Class:       "NoHierarchy",
				Superclass:  "",
				Interfaces:  []string{},
				Fields:      []dto.JsonField{},
				Methods:     []dto.JsonMethod{},
				Annotations: []string{},
				Hierarchy:   dto.Hierarchy{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToJson(tt.info)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTreeNodeToHierarchy(t *testing.T) {
	tests := []struct {
		name     string
		node     dto.TreeNode
		expected dto.Hierarchy
	}{
		{
			name: "node without child",
			node: dto.TreeNode{
				Name:     "Leaf",
				Children: nil,
			},
			expected: dto.Hierarchy{
				"Leaf": dto.Hierarchy{},
			},
		},
		{
			name: "node with child",
			node: dto.TreeNode{
				Name: "A",
				Children: []dto.TreeNode{
					{Name: "B", Children: nil},
					{
						Name: "C",
						Children: []dto.TreeNode{
							{Name: "D", Children: nil},
						},
					},
				},
			},
			expected: dto.Hierarchy{
				"A": dto.Hierarchy{
					"B": dto.Hierarchy{},
					"C": dto.Hierarchy{
						"D": dto.Hierarchy{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := treeNodeToHierarchy(tt.node)
			assert.Equal(t, tt.expected, got)
		})
	}
}
