package main

import (
	fmt "fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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

var priceRegexp = regexp.MustCompile(`\d+`)

func parsePriceInChaos(input string) float32 {
	match := priceRegexp.FindString(input)

	if len(match) == 0 {
		log.Fatalf(`Did not find any price string match in "%s"`, input)
	}

	n, err := strconv.ParseFloat(match, 64)

	if err != nil {
		log.Fatalf(`Unable to parse "%s" in to float32`, match)
	}

	return float32(n)
}

func extractFeaturesFromAnItem(item *Item) (features []float32, ok bool) {
	note := item.GetNote()

	if item.GetFrameType() == 2 && strings.Contains(note, "chaos") && priceRegexp.MatchString(note) {
		features = []float32{
			float32(parsePriceInChaos(item.GetNote())),
			float32(item.GetIlvl()),
			float32(len(item.GetSockets())),
			float32(numOfLinkedSockets(item.GetSockets())),
		}

		ok = true

		return
	}

	return
}

func generateMLInputFromResponse(loopLimit int, mlInput *[][]float32) func(*Response) error {
	var loopCounter int

	return func(r *Response) error {
		for _, stash := range r.Stashes {
			if stash != nil {
				for _, item := range stash.Items {
					if item != nil {
						features, ok := extractFeaturesFromAnItem(item)
						if ok {
							loopCounter++
							if loopCounter > loopLimit {
								return fmt.Errorf("Exiting from the loop, sorry")
							}

							*mlInput = append(*mlInput, features)
						}
					}
				}
			}
		}

		return nil
	}
}

func generateMLInput(dbPath string) [][]float32 {
	var mlInput [][]float32
	loopLimit := 5000

	log.Printf("Limiting to %d items", loopLimit)

	walkResponses(dbPath, generateMLInputFromResponse(loopLimit, &mlInput))

	return mlInput
}
