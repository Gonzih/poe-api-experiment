package main

import (
	fmt "fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type FieldsForExtraction struct {
	Properties map[string]bool
}

func generateFields(dbPath string) error {
	fields := &FieldsForExtraction{
		Properties: make(map[string]bool, 0),
	}

	c := 0

	err := walkResponses("data.bin", func(r *Response) error {
		for _, stash := range r.Stashes {
			if stash != nil {
				for _, item := range stash.Items {
					if item != nil {
						c++
						if c > 10000 {
							return fmt.Errorf("beep")
						}
						for _, field := range item.GetProperties() {
							fields.Properties[field.GetName()] = true
						}
					}
				}
			}
		}

		return nil
	})

	// if err != nil {
	// 	return err
	// }

	data, err := yaml.Marshal(fields)
	log.Printf("\n%s", string(data))

	if err != nil {
		return err
	}

	return nil
}
