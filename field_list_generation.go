package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type fieldsForExtraction struct {
	Properties   []string
	ImplicitMods []string
	ExplicitMods []string
}

func loadFieldsConfiguration() (*fieldsForExtraction, error) {
	fields := &fieldsForExtraction{}

	data, err := ioutil.ReadFile("fields.yaml")
	if err != nil {
		return fields, err
	}

	err = yaml.Unmarshal(data, &fields)
	if err != nil {
		return fields, err
	}

	return fields, nil
}

func saveFieldsOnDisk(fields *fieldsForExtraction) error {
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

var numRegexp = regexp.MustCompile(`\+?\d+`)

func parseModString(input string) (float32, string) {
	var n float32
	nums := numRegexp.FindAllString(input, -1)

	if len(nums) > 0 {
		for _, num := range nums {
			i, err := strconv.ParseInt(num, 10, 64)

			if err != nil {
				log.Fatalf(`Unable to parse "%s" in to float32`, num)
			}

			n += float32(i)
		}

		n /= float32(len(nums))
	} else {
		n = 1
	}

	return n, numRegexp.ReplaceAllString(input, `\d+`)
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
						if item.GetFrameType() == 2 {

							// c++
							// if c > 10000 {
							// 	return fmt.Errorf("beep")
							// }

							if mods := item.GetExplicitMods(); len(mods) > 0 {
								for _, mod := range mods {
									_, name := parseModString(mod)
									exMods[name] = true
								}
							}

							if mods := item.GetImplicitMods(); len(mods) > 0 {
								for _, mod := range mods {
									_, name := parseModString(mod)
									imMods[name] = true
								}
							}

							for _, field := range item.GetProperties() {
								props[field.GetName()] = true
							}
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

	fields := &fieldsForExtraction{
		Properties:   make([]string, 0),
		ExplicitMods: make([]string, 0),
		ImplicitMods: make([]string, 0),
	}

	for k := range props {
		fields.Properties = append(fields.Properties, k)
	}

	for k := range exMods {
		fields.ExplicitMods = append(fields.ExplicitMods, k)
	}

	for k := range imMods {
		fields.ImplicitMods = append(fields.ImplicitMods, k)
	}

	err = saveFieldsOnDisk(fields)

	if err != nil {
		return err
	}

	return nil
}
