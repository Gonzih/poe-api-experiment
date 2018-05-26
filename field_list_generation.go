package main

import "log"

func generateFields(dbPath string) error {
	fields := make(map[string]bool, 0)

	err := walkResponses("data.bin", func(r *Response) error {
		for _, stash := range r.Stashes {
			if stash != nil {
				for _, item := range stash.Items {
					if item != nil {
						for _, field := range item.GetAdditionalProperties() {
							fields[field.GetName()] = true
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

	log.Println(fields)

	return nil
}
