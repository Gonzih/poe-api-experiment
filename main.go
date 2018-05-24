package main

import (
	"log"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
)

// type responseData struct {
// 	NextChangeID string `json:"next_change_id"`
// 	Stashes      []struct {
// 		ID                string  `json:"id"`
// 		Public            bool    `json:"public"`
// 		AccountName       *string `json:"accountName"`
// 		LastCharacterName *string `json:"lastCharacterName"`
// 		Stash             *string `json:"stash"`
// 		StashType         string  `json:"stashType"`
// 		Items             []struct {
// 			AbyssJewel bool `json:"abyssJewel"`
// 		} `json:"items"`
// 	} `json:"stashes"`
// }

func main() {
	resp, err := http.Get("http://www.pathofexile.com/api/public-stash-tabs?")
	if err != nil {
		panic(err)
	}
	var data Response

	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(body))
	// err = json.Unmarshal(body, &data)

	unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: false}
	err = unmarshaller.Unmarshal(resp.Body, &data)

	if err != nil {
		panic(err)
	}

	// log.Println(data)

	log.Println("Done!")

	for _, stash := range data.Stashes {
		if stash != nil {
			for _, item := range stash.Items {
				if item != nil {
					log.Println()
					for _, property := range item.GetAdditionalProperties() {
						if property != nil {
							for _, value := range property.Values {
								log.Printf("%s -> %#v", property.GetName(), value)
							}
						}
					}
				}
			}
		}
	}

}
