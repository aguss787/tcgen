package tcgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type outputType string

const (
	INPUT  outputType = "in"
	OUTPUT outputType = "out"
)

func save(prefix string, id int, output []byte, mode outputType) error {
	f, err := os.Create(fmt.Sprintf("%s_%d.%s", prefix, id, mode))
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal("Error closing file", err)
		}
	}()

	_, err = f.Write(output)
	return err
}

func read(filename string) ([]json.RawMessage, error) {
	tc, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := tc.Close()
		if err != nil {
			log.Fatal("Error closing file", err)
		}
	}()

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result []json.RawMessage
	err = json.Unmarshal(content, &result)

	return result, err
}
