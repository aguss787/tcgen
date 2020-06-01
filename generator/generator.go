package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type Runner interface {
	AddSample(input []interface{})
	Add(input []interface{})
	GenerateTC(solution interface{})
}

type runner struct {
	sampleInputs [][]reflect.Value
	inputs [][]reflect.Value
}

func (r *runner) Add(rawInput []interface{}) {
	var input []reflect.Value
	for _, i := range rawInput {
		input = append(input, reflect.ValueOf(i))
	}
	r.inputs = append(r.inputs, input)
}

func (r *runner) AddSample(rawInput []interface{}) {
	var input []reflect.Value
	for _, i := range rawInput {
		input = append(input, reflect.ValueOf(i))
	}
	r.sampleInputs = append(r.sampleInputs, input)
}

func (r *runner) GenerateTC(solution interface{}) {
	err := os.RemoveAll("tc")
	if err != nil {
		log.Fatal("Error removing old tc", err)
	}

	err = os.Mkdir("tc", os.ModeDir | 0777)
	if err != nil {
		log.Fatal("Error making tc directory", err)
	}

	err = os.Chdir("tc")
	if err != nil {
		log.Fatal("Error chdir to tc directory", err)
	}

	defer func() {
		err := os.Chdir("..")
		if err != nil {
			log.Fatal("Error returning to base directory", err)
		}
	}()

	r.generateSample(solution)
	r.generateOther(solution)
}

func (r *runner) generateSample(solution interface{}) {
	for i, input := range r.sampleInputs {
		generateSingleTestcase(solution, input, "sample", i)
	}
}

func (r *runner) generateOther(solution interface{}) {
	for i, input := range r.inputs {
		generateSingleTestcase(solution, input, "tc", i)
	}
}

func generateSingleTestcase(solution interface{}, input []reflect.Value, prefix string, id int) {
	returnValue := reflect.ValueOf(solution).Call(input)

	res := []interface{}{}
	for _, value := range returnValue {
		res = append(res, value.Interface())
	}

	output, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Error marshalling output file", err)
	}

	var rawInput []interface{}
	for _, value := range input {
		rawInput = append(rawInput, value.Interface())
	}

	inputByte, err := json.Marshal(rawInput)
	if err != nil {
		log.Fatal("Error marshalling input file", err)
	}

	err = save(prefix, id, inputByte, INPUT)
	if err != nil {
		log.Fatal("Error saving output file", err)
	}

	err = save(prefix, id, output, OUTPUT)
	if err != nil {
		log.Fatal("Error saving output file", err)
	}
}

func NewRunner() Runner {
	return &runner{}
}

func Test(t *testing.T, solver interface{}) {
	t.Log("Running test suite")

	solverValue := reflect.ValueOf(solver)
	solverType := reflect.TypeOf(solver)

	err := os.Chdir("tc")
	if err != nil {
		t.Fatalf("TC directory not found")
	}

	defer func() {
		err := os.Chdir("..")
		if err != nil {
			log.Fatalf("Cannot return to base directory")
		}
	}()

	dir, err := os.Open(".")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = dir.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	dirInfos, err := dir.Readdir(-1)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(dirInfos, func(i, j int) bool {
		return dirInfos[i].Name() < dirInfos[j].Name()
	})

	for _, fileInfo := range dirInfos {
		filename := fileInfo.Name()
		if !strings.HasSuffix(filename, string(INPUT)) {
			continue
		}

		filename = strings.TrimSuffix(filename, string(INPUT))

		t.Logf("Running %s", filename)

		rawInput, err := read(fmt.Sprintf("%s%s", filename, INPUT))
		if err != nil {
			t.Fatal(err)
		}

		rawExpectedOutput, err := read(fmt.Sprintf("%s%s", filename, OUTPUT))
		if err != nil {
			t.Fatal(err)
		}

		var expectedOutput []interface{}
		for i, value := range rawExpectedOutput {
			outputType := solverType.Out(i)

			val := reflect.New(outputType)
			valP := val.Interface()
			json.Unmarshal(value, valP)

			expectedOutput = append(expectedOutput, val.Elem().Interface())
		}

		var input []reflect.Value
		for i, value := range rawInput {
			inputType := solverType.In(i)

			val := reflect.New(inputType)
			valP := val.Interface()
			json.Unmarshal(value, valP)

			input = append(input, val.Elem())
		}

		rawOutput := solverValue.Call(input)
		var output []interface{}
		for _, value := range rawOutput {
			output = append(output, value.Interface())
		}

		if !reflect.DeepEqual(output, expectedOutput) {
			t.Logf("\nExpected: %+v\n Found  : %+v", expectedOutput, output)
			t.FailNow()
		}
	}
}
