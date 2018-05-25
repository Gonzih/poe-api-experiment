package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestRealHttpRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test due to short run")
		return
	}

	nextChangeID := "112870934-116856871-113462743-125603141-117476043"
	counter := 0

	for {
		url := fmt.Sprintf("http://www.pathofexile.com/api/public-stash-tabs?id=%s", nextChangeID)

		var data Response

		body := httpGetBody(url)
		defer body.Close()
		unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: false}
		err := unmarshaller.Unmarshal(body, &data)

		assert.Nil(t, err)

		if err == nil {
			assert.True(t, true)
		}

		log.Println("Done!")
		log.Println(counter)

		counter++
		nextChangeID = data.GetNextChangeId()

		if counter >= 3 {
			break
		}
	}

	// for _, stash := range data.Stashes {
	// 	if stash != nil {
	// 		for _, item := range stash.Items {
	// 			if item != nil {
	// 				log.Println()
	// 				for _, property := range item.GetAdditionalProperties() {
	// 					if property != nil {
	// 						for _, value := range property.Values {
	// 							log.Printf("%s -> %#v", property.GetName(), value)
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

}
