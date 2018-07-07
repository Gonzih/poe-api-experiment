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
	Fields [][]float64
	FName  string
}

func (in *MLInput) initFName() {
	if len(in.FName) == 0 {
		in.FName = defaultMLInputFname
	}
}

func (in *MLInput) Save() error {
	in.initFName()

	f, err := os.OpenFile(in.FName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("could not open %s: %v", in.FName, err)
	}

	enc := gob.NewEncoder(f)
	return enc.Encode(*in)
}

func (in *MLInput) Load() error {
	in.initFName()

	f, err := os.OpenFile(in.FName, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("could not open %s: %v", in.FName, err)
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

var priceTable = map[string]float64{
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

func parsePriceInChrom(input string) (float64, bool) {
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
		log.Printf(`Unable to parse "%s" in to float64`, match)
		return 0, false
	}

	var rate float64

	for k, p := range priceTable {
		if strings.Contains(input, k) {
			rate = p
			break
		}
	}

	return float64(n) * rate, true
}

func extractProperties(props []*Property, names []string) ([]float64, bool) {
	parsedProps := make(map[string]float64, len(props))
	out := make([]float64, len(names))
	var v float64

	for _, prop := range props {

		// if len(prop.Values) > 1 {
		// 	log.Println(len(prop.Values))
		// 	log.Println(prop.GetName())
		// 	log.Println(prop.Values)
		// }
		if len(prop.Values) == 1 {
			v, _ = parseModString(prop.Values[0].Value)
		}

		parsedProps[prop.GetName()] = v
	}

	for i, name := range names {
		out[i] = parsedProps[name]
	}

	return out, true
}

func extractMods(mods []string, names []string) ([]float64, bool) {
	parsedMods := make(map[string]float64, len(mods))
	out := make([]float64, len(names))

	for _, mod := range mods {
		value, name := parseModString(mod)
		parsedMods[name] = value
	}

	for i, name := range names {
		out[i] = parsedMods[name]
	}

	return out, true
}

func extractFeaturesFromAnItem(item *Item, fieldsConfiguration *fieldsForExtraction) ([]float64, bool) {
	note := item.GetNote()

	price, ok := parsePriceInChrom(note)

	if item.GetFrameType() == 2 && ok {
		// log.Printf(`"%s": %3.3f`, note, price)

		var corrupted float64
		if item.GetCorrupted() {
			corrupted = 1
		}

		cat := item.GetCategory()

		fts := []float64{
			price,
			float64(item.GetIlvl()),
			float64(len(item.GetSockets())),
			float64(numOfLinkedSockets(item.GetSockets())),
			corrupted,
			float64(len(cat.GetAccessories())),
			float64(len(cat.GetArmour())),
			float64(len(cat.GetJewels())),
			float64(len(cat.GetWeapons())),
			float64(len(cat.GetGems())),
			float64(len(cat.GetFlasks())),
			float64(len(cat.GetMaps())),
			float64(len(cat.GetCurrency())),
			float64(len(cat.GetCards())),
		}

		exMods, ok := extractMods(item.GetExplicitMods(), fieldsConfiguration.ExplicitMods)
		if ok {
			fts = append(fts, exMods...)
		}

		imMods, ok := extractMods(item.GetImplicitMods(), fieldsConfiguration.ImplicitMods)
		if ok {
			fts = append(fts, imMods...)
		}

		props, ok := extractProperties(item.GetProperties(), fieldsConfiguration.Properties)
		if ok {
			fts = append(fts, props...)
		}

		return fts, true
	}

	return []float64{}, false
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

	log.Printf("Limiting to %d items", loopLimit)

	fieldsConfiguration, err := loadFieldsConfiguration()

	if err != nil {
		return mlInput, err
	}

	err = walkResponses(dbPath, generateMLInputFromResponse(loopLimit, mlInput, fieldsConfiguration))

	return mlInput, err
}

func walkFeatures(dbPath string, loopLimit int, walkFn func([]float64) error) error {
	log.Printf("Limiting to %d items", loopLimit)

	fieldsConfiguration, err := loadFieldsConfiguration()

	if err != nil {
		return err
	}

	err = walkResponses(dbPath, lazilyWalkFeatures(loopLimit, fieldsConfiguration, walkFn))

	return err

}

func lazilyWalkFeatures(loopLimit int, fieldsConfiguration *fieldsForExtraction, walkFn func([]float64) error) func(*Response) error {
	var loopCounter int
	var err error

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

							err = walkFn(features)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}

		return nil
	}
}
