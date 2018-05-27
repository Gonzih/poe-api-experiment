package main

import (
	"encoding/gob"
	fmt "fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const defaultMLInputFname = "data/ml-input.bin"

type MLInput struct {
	Fields [][]float32
	FName  string
}

func (in *MLInput) Save() error {
	var fname string

	if len(in.FName) > 0 {
		fname = in.FName
	} else {
		fname = defaultMLInputFname
	}

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", fname, err)
	}

	enc := gob.NewEncoder(f)
	return enc.Encode(in)
}

func (in *MLInput) Load() error {
	var fname string

	if len(in.FName) > 0 {
		fname = in.FName
	} else {
		fname = defaultMLInputFname
	}

	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", fname, err)
	}

	dec := gob.NewDecoder(f)
	return dec.Decode(in)
}

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

var priceTable = map[string]float32{
	"chrom":   1,
	"alt":     1.1,
	"jew":     2,
	"chance":  2.8,
	"chisel":  3.5,
	"alch":    4,
	"fuse":    7,
	"scour":   7.8,
	"blessed": 18,
	"chaos":   14,
	"regret":  14,
	"vaal":    14,
	"gcp":     25, // gem cutter prism
	"divine":  233,
	"exa":     700,
	// "et":      2800,
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

func parsePriceInChrom(input string) (float32, bool) {
	if len(input) == 0 {
		return 0, false
	}

	if !(strings.HasPrefix(input, "~price") || strings.HasPrefix(input, "~b/o")) {
		// log.Printf(`Price note was not generated automatically, ignoring: "%s"`, input)
		return 0, false
	}

	match := priceRegexp.FindString(input)

	if len(match) == 0 {
		// log.Printf(`Did not find any price string match in "%s"`, input)
		return 0, false
	}

	n, err := strconv.ParseFloat(match, 64)

	if err != nil {
		log.Printf(`Unable to parse "%s" in to float32`, match)
		return 0, false
	}

	var rate float32

	for k, p := range priceTable {
		if strings.Contains(input, k) {
			rate = p
			break
		}
	}

	return float32(n) * rate, true
}

func extractMods(mods []string, names []string) ([]float32, bool) {
	parsedMods := make(map[string]float32, len(mods))
	out := make([]float32, len(names))

	for _, mod := range mods {
		value, name := parseModString(mod)
		parsedMods[name] = value
	}

	for i, name := range names {
		out[i] = parsedMods[name]
	}

	return out, true
}

func extractFeaturesFromAnItem(item *Item, fieldsConfiguration *fieldsForExtraction) ([]float32, bool) {
	note := item.GetNote()

	price, ok := parsePriceInChrom(note)

	if item.GetFrameType() == 2 && ok {
		// log.Printf(`"%s": %3.3f`, note, price)

		var corrupted float32
		if item.GetCorrupted() {
			corrupted = 1
		}

		cat := item.GetCategory()

		fts := []float32{
			price,
			float32(item.GetIlvl()),
			float32(len(item.GetSockets())),
			float32(numOfLinkedSockets(item.GetSockets())),
			corrupted,
			float32(len(cat.GetAccessories())),
			float32(len(cat.GetArmour())),
			float32(len(cat.GetJewels())),
			float32(len(cat.GetWeapons())),
			float32(len(cat.GetGems())),
			float32(len(cat.GetFlasks())),
			float32(len(cat.GetMaps())),
			float32(len(cat.GetCurrency())),
			float32(len(cat.GetCards())),
		}

		exMods, ok := extractMods(item.GetExplicitMods(), fieldsConfiguration.ExplicitMods)
		if ok {
			fts = append(fts, exMods...)
		}

		imMods, ok := extractMods(item.GetImplicitMods(), fieldsConfiguration.ImplicitMods)
		if ok {
			fts = append(fts, imMods...)
		}

		return fts, true
	}

	return []float32{}, false
}

func generateMLInputFromResponse(loopLimit int, mlInput *MLInput, fieldsConfiguration *fieldsForExtraction) func(*Response) error {
	var loopCounter int

	return func(r *Response) error {
		for _, stash := range r.Stashes {
			if stash != nil {
				for _, item := range stash.Items {
					if item != nil {
						features, ok := extractFeaturesFromAnItem(item, fieldsConfiguration)
						if ok {
							loopCounter++
							if loopCounter > loopLimit {
								return exitingTheLoopErr
							}

							mlInput.Fields = append(mlInput.Fields, features)
						}
					}
				}
			}
		}

		return nil
	}
}

func generateMLInput(dbPath string, loopLimit int) (*MLInput, error) {
	mlInput := &MLInput{}
	var fields [][]float32

	log.Printf("Limiting to %d items", loopLimit)

	fieldsConfiguration, err := loadFieldsConfiguration()

	if err != nil {
		return mlInput, err
	}

	err = walkResponses(dbPath, generateMLInputFromResponse(loopLimit, mlInput, fieldsConfiguration))

	mlInput.Fields = fields

	return mlInput, err
}
