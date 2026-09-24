package application

import (
	"path/filepath"
	"reflect"
	"testing"

	"go/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/mocks"
	"golang.org/x/tools/go/packages"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"
)

func projectRootAbs(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err)
	return root
}

func examplePkgPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRootAbs(t), "example")
}

func testPkgPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRootAbs(t), "test_package")
}

func TestStartInspect(t *testing.T) {
	tests := []struct {
		name        string
		cfg         domain.InputConfig
		contains    string
		expectedErr bool
	}{
		{
			name: "class Manager",
			cfg: domain.NewInputConfig(
				examplePkgPath(t),
				"Manager",
				"text",
				"out.txt",
			),
			contains:    "Class: Manager",
			expectedErr: false,
		},
		{
			name: "struct not found",
			cfg: domain.NewInputConfig(
				examplePkgPath(t),
				"Nope",
				"text",
				"out.txt",
			),
			expectedErr: true,
		},
		{
			name: "invalid format",
			cfg: domain.NewInputConfig(
				examplePkgPath(t),
				"Manager",
				"xml",
				"out.txt",
			),
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := StartInspect(test.cfg)
			assert.Equal(t, test.expectedErr, err != nil)
			if err != nil {
				return
			}
			assert.Contains(t, out, test.contains)
		})
	}
}

func TestLoadPackage(t *testing.T) {
	tests := []struct {
		name        string
		pkgPath     string
		findName    string
		expectedErr bool
	}{
		{
			name:        "load example",
			pkgPath:     examplePkgPath(t),
			findName:    "Manager",
			expectedErr: false,
		},
		{
			name:        "load test_package",
			pkgPath:     testPkgPath(t),
			findName:    "Person",
			expectedErr: false,
		},
		{
			name:        "invalid dir",
			pkgPath:     filepath.Join(projectRootAbs(t), "nope"),
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadPackage(test.pkgPath)
			assert.Equal(t, test.expectedErr, err != nil)
		})
	}
}

func TestFindStruct(t *testing.T) {
	examplePkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	testPkgs, err := LoadPackage(testPkgPath(t))
	require.NoError(t, err)

	allPkgs := append(examplePkgs, testPkgs...)

	tests := []struct {
		name        string
		pkgs        []*packages.Package
		structName  string
		expectedErr error
	}{
		{
			name:        "found Manager in example",
			pkgs:        examplePkgs,
			structName:  "Manager",
			expectedErr: nil,
		},
		{
			name:        "found Person in test_package",
			pkgs:        testPkgs,
			structName:  "Person",
			expectedErr: nil,
		},
		{
			name:        "not found",
			pkgs:        examplePkgs,
			structName:  "Nope",
			expectedErr: ErrStructNotFound,
		},
		{
			name:        "found by searching across both pkgs",
			pkgs:        allPkgs,
			structName:  "Manager",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = FindStruct(test.structName, test.pkgs)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestGetStructInfo(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	structName := "Manager"
	expected := dto.StructInfo{
		Class:  "Manager",
		Embeds: []string{"Employee"},
		Fields: []dto.Field{
			{Name: "Department", Type: "string"},
		},
	}

	obj, err := FindStruct(structName, pkgs)
	assert.NoError(t, err)

	info, err := GetStructInfo(obj)

	assert.Equal(t, expected.Class, info.Class)
	assert.ElementsMatch(t, expected.Embeds, info.Embeds)
	assert.ElementsMatch(t, expected.Fields, info.Fields)

	assert.Equal(t, "Manager", info.Hierarchy.Name)
	require.Len(t, info.Hierarchy.Children, 1)
	assert.Equal(t, "Employee", info.Hierarchy.Children[0].Name)
}

func TestExtractEmbedsAndFields(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	tests := []struct {
		name          string
		structName    string
		expectedEmb   []string
		expectedField []dto.Field
	}{
		{
			name:        "Person embeds Entity",
			structName:  "Person",
			expectedEmb: []string{"Entity"},
			expectedField: []dto.Field{
				{Name: "Name", Type: "string"},
				{Name: "Age", Type: "int"},
			},
		},
		{
			name:        "Employee embeds Person",
			structName:  "Employee",
			expectedEmb: []string{"Person"},
			expectedField: []dto.Field{
				{Name: "Position", Type: "string"},
				{Name: "Salary", Type: "int"},
			},
		},
		{
			name:        "Manager embeds Employee",
			structName:  "Manager",
			expectedEmb: []string{"Employee"},
			expectedField: []dto.Field{
				{Name: "Department", Type: "string"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj, err := FindStruct(test.structName, pkgs)
			require.NoError(t, err)

			named := obj.Type().(*types.Named)
			st := named.Underlying().(*types.Struct)

			embeds, fields := ExtractEmbedsAndFields(st)

			assert.ElementsMatch(t, test.expectedEmb, embeds)
			assert.ElementsMatch(t, test.expectedField, fields)
		})
	}
}

func TestEmbeddedName(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	tests := []struct {
		name       string
		structName string
		fieldIndex int
		expected   string
	}{
		{name: "Person -> Entity", structName: "Person", fieldIndex: 0, expected: "Entity"},
		{name: "Employee -> Person", structName: "Employee", fieldIndex: 0, expected: "Person"},
		{name: "Manager -> Employee", structName: "Manager", fieldIndex: 0, expected: "Employee"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj, err := FindStruct(test.structName, pkgs)
			require.NoError(t, err)

			named := obj.Type().(*types.Named)
			st := named.Underlying().(*types.Struct)

			f := st.Field(test.fieldIndex)
			require.True(t, f.Embedded())

			assert.Equal(t, test.expected, EmbeddedName(f.Type()))
		})
	}
}

func TestBuildHierarchy(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	structName := "Manager"
	expectedRoot := "Manager"
	expectedPath := []string{"Employee", "Person", "Entity"}

	obj, err := FindStruct(structName, pkgs)
	require.NoError(t, err)

	named := obj.Type().(*types.Named)
	h := BuildHierarchy(named)

	assert.Equal(t, expectedRoot, h.Name)

	cur := h
	for _, ch := range expectedPath {
		require.Len(t, cur.Children, 1)
		assert.Equal(t, ch, cur.Children[0].Name)
		cur = cur.Children[0]
	}
}

func TestUnwrapNamedType(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	tests := []struct {
		name         string
		structName   string
		fieldIndex   int
		expectedNil  bool
		expectedName string
	}{
		{name: "Person -> Entity", structName: "Person", fieldIndex: 0, expectedNil: false, expectedName: "Entity"},
		{name: "Person -> Name string -> nil", structName: "Person", fieldIndex: 1, expectedNil: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj, err := FindStruct(test.structName, pkgs)
			require.NoError(t, err)

			named := obj.Type().(*types.Named)
			st := named.Underlying().(*types.Struct)

			f := st.Field(test.fieldIndex)
			got := UnwrapNamedType(f.Type())

			if test.expectedNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, test.expectedName, got.Obj().Name())
		})
	}
}

func TestExtractMethods(t *testing.T) {
	pkgs, err := LoadPackage(examplePkgPath(t))
	require.NoError(t, err)

	tests := []struct {
		name            string
		structName      string
		expectedMethods []dto.Method
	}{
		{
			name:       "Person methods",
			structName: "Person",
			expectedMethods: []dto.Method{
				{Name: "GetName", Params: []string{}, ReturnType: []string{"string"}},
				{Name: "SetName", Params: []string{"string"}, ReturnType: []string{}},
			},
		},
		{
			name:       "Employee methods",
			structName: "Employee",
			expectedMethods: []dto.Method{
				{Name: "GetSalary", Params: []string{}, ReturnType: []string{"int"}},
				{Name: "Promote", Params: []string{"string", "int"}, ReturnType: []string{}},
			},
		},
		{
			name:       "Manager methods",
			structName: "Manager",
			expectedMethods: []dto.Method{
				{Name: "ChangeDepartment", Params: []string{"string"}, ReturnType: []string{}},
				{Name: "GetDepartment", Params: []string{}, ReturnType: []string{"string"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj, err := FindStruct(test.structName, pkgs)
			require.NoError(t, err)

			named := obj.Type().(*types.Named)
			got := ExtractMethods(named)

			assert.ElementsMatch(t, test.expectedMethods, got)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "int64"},
		{name: "string"},
		{name: "bool"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "int64":
				_, err := CreateInstance[int64]()
				assert.NoError(t, err)
			case "string":
				_, err := CreateInstance[string]()
				assert.NoError(t, err)
			case "bool":
				_, err := CreateInstance[bool]()
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateInstance_CustomStruct(t *testing.T) {
	type Inner struct {
		ID   int64
		Name string
	}

	type Custom struct {
		Ok    bool
		Count int64
		Note  string
		Rate  float64

		Array [3]int64
		Slice []int64

		Map map[string]int

		In   *Inner
		List []*Inner
	}

	obj, err := CreateInstance[Custom]()
	assert.NoError(t, err)
	assert.NotNil(t, obj.In)
	assert.NotNil(t, obj.Map)
	assert.NotNil(t, obj.Slice)
	assert.Len(t, obj.Array, 3)
	assert.NotNil(t, obj.List)
}

func TestFillObject_Primitives_WithGoMock(t *testing.T) {
	tests := []struct {
		name     string
		ptr      interface{}
		setup    func(r *mocks.MockRandom)
		expected interface{}
	}{
		{
			name: "int64",
			ptr:  new(int64),
			setup: func(r *mocks.MockRandom) {
				r.EXPECT().RandInt64().Return(int64(-123)).Times(1)
			},
			expected: int64(-123),
		},
		{
			name: "string",
			ptr:  new(string),
			setup: func(r *mocks.MockRandom) {
				r.EXPECT().RandString().Return("hello").Times(1)
			},
			expected: "hello",
		},
		{
			name: "bool",
			ptr:  new(bool),
			setup: func(r *mocks.MockRandom) {
				r.EXPECT().RandBoll().Return(true).Times(1)
			},
			expected: true,
		},
		{
			name: "float64",
			ptr:  new(float64),
			setup: func(r *mocks.MockRandom) {
				r.EXPECT().RandFloat64().Return(3.14).Times(1)
			},
			expected: 3.14,
		},
		{
			name: "uint64",
			ptr:  new(uint64),
			setup: func(r *mocks.MockRandom) {
				r.EXPECT().RandUint64().Return(uint64(999)).Times(1)
			},
			expected: uint64(999),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rnd := mocks.NewMockRandom(ctrl)
			test.setup(rnd)

			v := reflect.ValueOf(test.ptr).Elem()
			err := FillObject(v, rnd)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, v.Interface())
		})
	}
}

func TestFillObject_CustomStruct(t *testing.T) {
	type Inner struct {
		ID   int64
		Name string
	}

	type Custom struct {
		Ok    bool
		Count int64
		U     uint64
		R     float64
		Note  string

		Arr [2]int64

		Slice []string
		Map   map[string]int64

		Ptr  *Inner
		Ptrs []*Inner
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rnd := mocks.NewMockRandom(ctrl)

	gomock.InOrder(
		rnd.EXPECT().RandBoll().Return(true),

		rnd.EXPECT().RandInt64().Return(int64(10)),

		rnd.EXPECT().RandUint64().Return(uint64(999)),

		rnd.EXPECT().RandFloat64().Return(3.14),

		rnd.EXPECT().RandString().Return("note"),

		rnd.EXPECT().RandInt64().Return(int64(1)),
		rnd.EXPECT().RandInt64().Return(int64(2)),

		rnd.EXPECT().RandIntPositive().Return(2),
		rnd.EXPECT().RandString().Return("s1"),
		rnd.EXPECT().RandString().Return("s2"),

		rnd.EXPECT().RandIntPositive().Return(2),
		rnd.EXPECT().RandString().Return("k1"),
		rnd.EXPECT().RandInt64().Return(int64(200)),
		rnd.EXPECT().RandString().Return("k2"),
		rnd.EXPECT().RandInt64().Return(int64(300)),

		rnd.EXPECT().RandInt64().Return(int64(100)),
		rnd.EXPECT().RandString().Return("ptr"),

		rnd.EXPECT().RandIntPositive().Return(2),
		rnd.EXPECT().RandInt64().Return(int64(400)),
		rnd.EXPECT().RandString().Return("p1"),
		rnd.EXPECT().RandInt64().Return(int64(500)),
		rnd.EXPECT().RandString().Return("p2"),
	)

	var cust Custom
	v := reflect.ValueOf(&cust).Elem()

	err := FillObject(v, rnd)
	require.NoError(t, err)

	assert.Equal(t, true, cust.Ok)
	assert.Equal(t, int64(10), cust.Count)
	assert.Equal(t, uint64(999), cust.U)
	assert.Equal(t, 3.14, cust.R)
	assert.Equal(t, "note", cust.Note)

	assert.Equal(t, [2]int64{1, 2}, cust.Arr)
	assert.Equal(t, []string{"s1", "s2"}, cust.Slice)
	assert.Equal(t, map[string]int64{"k1": 200, "k2": 300}, cust.Map)

	require.NotNil(t, cust.Ptr)
	assert.Equal(t, int64(100), cust.Ptr.ID)
	assert.Equal(t, "ptr", cust.Ptr.Name)

	require.Len(t, cust.Ptrs, 2)
	require.NotNil(t, cust.Ptrs[0])
	require.NotNil(t, cust.Ptrs[1])
	assert.Equal(t, int64(400), cust.Ptrs[0].ID)
	assert.Equal(t, "p1", cust.Ptrs[0].Name)
	assert.Equal(t, int64(500), cust.Ptrs[1].ID)
	assert.Equal(t, "p2", cust.Ptrs[1].Name)

}
