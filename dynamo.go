package dynamo

import "go/format"

type File struct {
	Package string

	imports []string
	structs map[string]Struct
	funcs   map[string]Func
	methods map[string]Method
}

type Func struct {
	Arguments []Parameter
	Outputs   []Parameter

	Body string
}

type Struct map[string]Feild

type Method struct {
	Struct   string
	Receiver string

	Func Func
}

type Parameter struct {
	Name string
	Type string
}

type Feild struct {
	Type string
	Note string
}

const space = "    "

func (file *File) Import(paths []string) {
	file.imports = append(file.imports, paths...)
}

func (file *File) Struct(name string, body Struct) {
	file.structs[name] = body
}

func (file *File) Func(name string, body Func) {
	file.funcs[name] = body
}

func (file *File) Method(name string, body Method) {
	file.methods[name] = body
}

func (file *File) Compile() (string, error) {
	plain := "package " + file.Package + "\n\n"

	plain = plain + "import (\n"

	for index, _ := range file.imports {
		plain = plain + space + `"` + file.imports[index] + `"` + "\n"
	}

	plain = plain + ")\n\n"

	for name, body := range file.structs {
		plain = plain + "type " + name + " struct {\n"

		for feildName, feildBody := range body {
			plain = plain + space + feildName + " " + feildBody.Type + "\n"
		}

		plain = plain + "}\n\n"
	}

	for name, method := range file.methods {
		plain = plain + "func (" + method.Receiver + " " + method.Struct + ") " + name + "(" + joinParameters(method.Func.Arguments, " ") + ") "

		if len(method.Func.Outputs) > 0 {
			plain = plain + " (" + joinParameters(method.Func.Outputs, ", ") + ")"
		}

		plain = plain + " {\n" + method.Func.Body + "\n}\n\n"
	}

	for name, function := range file.funcs {
		plain = plain + "func " + name + "(" + joinParameters(function.Arguments, " ") + ")"

		if len(function.Outputs) > 0 {
			plain = plain + " (" + joinParameters(function.Outputs, ", ") + ")"
		}

		plain = plain + " {\n" + function.Body + "\n}\n\n"
	}

	formatted, err := format.Source([]byte(plain))

	return string(formatted), err
}

func joinParameters(parameters []Parameter, delimiter string) string {
	var plain string

	for index, parameter := range parameters {
		plain = plain + parameter.Name + " " + parameter.Type

		if index != len(parameters)-1 {
			plain = plain + delimiter
		}
	}

	return plain
}

func New(name string) File {
	return File{
		Package: name,

		imports: []string{},
		structs: make(map[string]Struct),
		funcs:   make(map[string]Func),
		methods: make(map[string]Method),
	}
}
