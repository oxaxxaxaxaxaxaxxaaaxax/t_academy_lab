package domain

type InputConfig struct {
	PkgName    string
	StructName string
	Format     string
	Output     string
}

func NewInputConfig(pkgName, structName, format, output string) InputConfig {
	return InputConfig{
		PkgName:    pkgName,
		StructName: structName,
		Format:     format,
		Output:     output,
	}
}

type InspectConfig struct {
	PkgPath    string
	StructName string
	Format     string
}

func NewInspectConfig(pkgPath, structName, format string) InspectConfig {
	return InspectConfig{
		PkgPath:    pkgPath,
		StructName: structName,
		Format:     format,
	}
}
