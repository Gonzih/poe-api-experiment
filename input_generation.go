package main

import (
	"log"
)

func mapFrameType(t int64) string {
	switch t {
	case 0:
		return "normal"
	case 1:
		return "magic"
	case 2:
		return "rare"
	case 3:
		return "unique"
	case 4:
		return "gem"
	case 5:
		return "currency"
	case 6:
		return "divination card"
	case 7:
		return "quest item"
	case 8:
		return "prophecy"
	case 9:
		return "relic"
	default:
		return ""
	}
}

func numOfLinkedSockets(sockets []*Socket) int64 {
	if len(sockets) == 0 {
		return 0
	}

	var maxGroup, g int64

	for _, socket := range sockets {
		if socket != nil {
			g = socket.GetGroup()

			if g > maxGroup {
				maxGroup = g
			}
		}
	}

	return maxGroup + 1
}

var propertyKeys = make(map[string]bool, 0)

func generateMLInputFromResponse(r *Response) error {
	for _, stash := range r.Stashes {
		if stash != nil {
			for _, item := range stash.Items {
				if item != nil {
					if len(item.GetNote()) > 0 && item.GetFrameType() == 2 {
						// log.Printf("%v -> %s", mapFrameType(item.GetFrameType()), item.GetNote())

						// log.Printf("Sockets: %d, groups: %d", len(item.GetSockets()), numOfLinkedSockets(item.GetSockets()))

						for _, property := range item.GetProperties() {
							if property != nil {
								propertyKeys[property.GetName()] = true
								// for _, value := range property.Values {
								// 	log.Printf("%s -> %#v", property.GetName(), value)
								// }
							}
						}

						// log.Println("=================>")
					}
				}
			}
		}
	}

	return nil
}

func generateMLInput(dbPath string) {
	walkResponses(dbPath, generateMLInputFromResponse)
	for k := range propertyKeys {
		log.Println(k)
	}
}
