package dynamo

import "testing"

func compile() (string, error) {
	file := New("main")

	file.Import([]string{"fmt"})

	file.Struct("Person", Struct{
		"Name": Feild{"string", ""},
		"Age":  Feild{"int", "`json:\"age\"`"},
	})

	file.Method("Greet", Method{
		Struct:   "Person",
		Receiver: "person",

		Func: Func{
			Arguments: []Parameter{},
			Outputs: []Parameter{
				Parameter{"response", "string"},
			},

			Body: `return fmt.Sprintf("Hello %s!", person.Name)`,
		},
	})

	file.Func("main", Func{
		Body: `person := Person{Name: "Linden", Age: 99}` + "\n\n" + `fmt.Println(person.Greet())`,
	})

	return file.Compile()
}

func TestCompile(test *testing.T) {
	body, err := compile()

	if err != nil {
		test.Errorf("%v", err)
	} else {
		test.Logf("compiled\n%v", body)
	}
}

func BenchmarkCompile(bench *testing.B) {
	for index := 0; index < bench.N; index++ {
		compile()
	}
}
