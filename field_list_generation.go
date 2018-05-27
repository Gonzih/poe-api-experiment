package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type FieldsForExtraction struct {
	Properties   []string
	ImplicitMods []string
	ExplicitMods []string
}

func saveFieldsOnDisk(fields *FieldsForExtraction) error {
	f, err := os.OpenFile("fields.yaml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(fields)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

var numRegexp = regexp.MustCompile(`\d+`)

func parseModString(input string) (float32, string) {
	var n float32
	nums := numRegexp.FindAllString(input, -1)

	if len(nums) == 1 {
		i, err := strconv.ParseInt(nums[0], 10, 64)

		if err != nil {
			log.Fatalf(`Unable to parse "%s" in to float32`, nums[0])
		}

		n = float32(i)
	}

	return n, numRegexp.ReplaceAllString(input, "")
}

func generateFields(dbPath string) error {
	props := make(map[string]bool, 0)
	imMods := make(map[string]bool, 0)
	exMods := make(map[string]bool, 0)

	_, _ = imMods, exMods

	// c := 0

	err := walkResponses("data.bin", func(r *Response) error {
		for _, stash := range r.Stashes {
			if stash != nil {
				for _, item := range stash.Items {
					if item != nil {
						// c++
						// if c > 10000 {
						// 	return fmt.Errorf("beep")
						// }

						if d := item.GetExplicitMods(); len(d) > 0 {
							log.Println(d)
						}
						for _, field := range item.GetProperties() {
							props[field.GetName()] = true
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	fields := &FieldsForExtraction{
		Properties: make([]string, 0),
	}

	for k := range props {
		fields.Properties = append(fields.Properties, k)
	}

	err = saveFieldsOnDisk(fields)

	if err != nil {
		return err
	}

	return nil
}
