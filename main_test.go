package main

import (
	"net/http"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestRealHttpRequest(t *testing.T) {

	resp, err := http.Get("http://www.pathofexile.com/api/public-stash-tabs?")
	if err != nil {
		panic(err)
	}
	var data Response

	defer resp.Body.Close()
	unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: false}
	err = unmarshaller.Unmarshal(resp.Body, &data)

	assert.Nil(t, err)

	if err == nil {
		assert.True(t, true)
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
