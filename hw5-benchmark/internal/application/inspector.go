package application

import (
	"errors"
	"go/types"
	"reflect"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/domain/dto"
	"golang.org/x/tools/go/packages"
)

const CurrentDir = "."

var (
	ErrStructNotFound   = errors.New("struct not found")
	ErrKindNotSupported = errors.New("kind not supported")
)

func StartInspect(cfg domain.InputConfig) (string, error) {
	formatter := NewFormatter(cfg.Format)

	inspectCfg, err := adapter.ValidateArgs(cfg)
	if err != nil {
		return "", err
	}

	pkgs, err := LoadPackage(inspectCfg.PkgPath)
	if err != nil {
		return "", err
	}

	str, err := FindStruct(inspectCfg.StructName, pkgs)
	if err != nil {
		return "", err
	}

	info, err := GetStructInfo(str)
	if err != nil {
		return "", err
	}

	return formatter.Format(info)
}

func GetStructInfo(obj types.Object) (dto.StructInfo, error) {
	named := obj.Type().(*types.Named)
	str := named.Underlying().(*types.Struct)

	embeds, fields := ExtractEmbedsAndFields(str)

	return dto.StructInfo{
		Class:     obj.Name(),
		Embeds:    embeds,
		Fields:    fields,
		Methods:   ExtractMethods(named),
		Hierarchy: BuildHierarchy(named),
	}, nil
}

func ExtractEmbedsAndFields(st *types.Struct) ([]string, []dto.Field) {
	var embeds []string
	var fields []dto.Field

	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Embedded() {
			fields = append(fields, dto.Field{
				Name: f.Name(),
				Type: f.Type().String(),
			})
			continue
		}

		embeds = append(embeds, EmbeddedName(f.Type()))
	}

	return embeds, fields
}

func EmbeddedName(t types.Type) string {
	switch x := t.(type) {
	case *types.Named:
		return x.Obj().Name()
	case *types.Pointer:
		return EmbeddedName(x.Elem())
	default:
		return t.String()
	}
}

func BuildHierarchy(named *types.Named) dto.TreeNode {
	node := dto.TreeNode{
		Name: named.Obj().Name(),
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return node
	}

	numFields := st.NumFields()
	for i := 0; i < numFields; i++ {
		f := st.Field(i)
		if !f.Embedded() {
			continue
		}

		emb := UnwrapNamedType(f.Type())
		if emb != nil {
			node.Children = append(node.Children, BuildHierarchy(emb))
		}
	}

	return node
}

func UnwrapNamedType(t types.Type) *types.Named {
	switch x := t.(type) {
	case *types.Named:
		return x
	case *types.Pointer:
		return UnwrapNamedType(x.Elem())
	default:
		return nil
	}
}

func ExtractMethods(named *types.Named) []dto.Method {
	var methods []dto.Method

	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		sig, ok := m.Type().(*types.Signature)
		if !ok {
			continue
		}

		params := make([]string, 0)
		lenParams := sig.Params().Len()
		for j := 0; j < lenParams; j++ {
			params = append(params, sig.Params().At(j).Type().String())
		}

		results := make([]string, 0)
		lenResults := sig.Results().Len()
		for j := 0; j < lenResults; j++ {
			results = append(results, sig.Results().At(j).Type().String())
		}

		methods = append(methods, dto.Method{
			Name:       m.Name(),
			Params:     params,
			ReturnType: results,
		})
	}

	return methods
}

func FindStruct(structName string, pkgs []*packages.Package) (types.Object, error) {
	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(structName)
		if obj == nil {
			continue
		}
		t, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}
		_, ok = t.Underlying().(*types.Struct)
		if !ok {
			continue
		}
		return obj, nil
	}
	return nil, ErrStructNotFound
}

func LoadPackage(pkgPath string) ([]*packages.Package, error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes |
			packages.NeedTypesInfo | packages.NeedSyntax,
		Dir: pkgPath,
	}
	pkgs, err := packages.Load(&cfg, CurrentDir)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// CreateInstance creates and fills a new instance of type T.
// T must be a compile-time known type.
func CreateInstance[T any]() (T, error) {
	var obj T
	t := reflect.TypeOf(obj)
	v := reflect.New(t)
	rnd := NewInspectorRand()

	err := FillObject(v.Elem(), rnd)
	if err != nil {
		return obj, err
	}

	obj, ok := v.Elem().Interface().(T)
	if !ok {
		return obj, ErrKindNotSupported
	}
	return obj, nil
}

func FillObject(v reflect.Value, rnd domain.Random) error {
	t := v.Type()

	switch t.Kind() {
	case reflect.Bool:
		v.SetBool(rnd.RandBoll())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(rnd.RandInt64())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(rnd.RandUint64())

	case reflect.Float32, reflect.Float64:
		v.SetFloat(rnd.RandFloat64())

	case reflect.String:
		v.SetString(rnd.RandString())

	case reflect.Array:
		array := reflect.New(t).Elem()
		for i := 0; i < t.Len(); i++ {
			err := FillObject(array.Index(i), rnd)
			if err != nil {
				return err
			}
		}
		v.Set(array)

	case reflect.Slice:
		vLen := rnd.RandIntPositive()
		slice := reflect.MakeSlice(t, vLen, vLen)
		for i := 0; i < vLen; i++ {
			err := FillObject(slice.Index(i), rnd)
			if err != nil {
				return err
			}
		}
		v.Set(slice)

	case reflect.Map:
		vLen := rnd.RandIntPositive()
		m := reflect.MakeMap(t)
		mKeyType := t.Key()
		mValType := t.Elem()

		for i := 0; i < vLen; i++ {
			randKey := reflect.New(mKeyType).Elem()
			randVal := reflect.New(mValType).Elem()

			err := FillObject(randKey, rnd)
			if err != nil {
				return err
			}

			err = FillObject(randVal, rnd)
			if err != nil {
				return err
			}

			m.SetMapIndex(randKey, randVal)
		}
		v.Set(m)

	case reflect.Struct:
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			if !v.Field(i).CanSet() {
				continue
			}
			err := FillObject(v.Field(i), rnd)
			if err != nil {
				return err
			}
		}

	case reflect.Pointer:
		ptr := reflect.New(t.Elem())
		err := FillObject(ptr.Elem(), rnd)
		if err != nil {
			return err
		}
		v.Set(ptr)

	default:
		return ErrKindNotSupported
	}
	return nil
}
