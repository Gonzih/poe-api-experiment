package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func httpGetBody(url string) io.ReadCloser {
	resp, err := http.Get("http://www.pathofexile.com/api/public-stash-tabs?")
	if err != nil {
		panic(err)
	}

	return resp.Body
}

func TestRealHttpRequest(t *testing.T) {
	id := ""

	for {
		url := fmt.Sprintf("http://www.pathofexile.com/api/public-stash-tabs?%s", id)

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

		if len(id) > 0 {
			break
		}

		id = data.GetNextChangeId()
	}

	// log.Println("Done!")

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
