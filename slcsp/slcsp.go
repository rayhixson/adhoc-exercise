package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

const zipsFile = "data/zips.csv"
const plansFile = "data/plans.csv"
const inputFile = "data/slcsp.csv"
const Silver = "Silver"

func main() {
	zips := LoadZips()
	plans := LoadPlans()
	finder := SlcspFinder{zips, plans}

	rows := [][]string{}
	rowHandler := func(record []string) error {
		if len(record) != 2 {
			return errors.New(fmt.Sprintf("Expect 1 field with a zip: %v", record))
		}

		zip := record[0]
		rate := finder.Find(zip)
		strrate := ""
		if rate > 0 {
			strrate = fmt.Sprintf("%.2f", rate)
		}
		rows = append(rows, []string{zip, strrate})
		return nil
	}
	loadFile(inputFile, rowHandler)

	// write the file to stdout
	fmt.Println("zipcode,rate")
	for _, r := range rows {
		fmt.Printf("%s,%s\n", r[0], r[1])
	}
}

type SlcspFinder struct {
	AllZips  Zips
	AllPlans Plans
}

func (s SlcspFinder) Find(givenZip string) float64 {
	matchingZips := s.AllZips.FindInOneRateArea(givenZip)

	plansForZip := Plans{}
	if len(matchingZips) > 0 {
		// all zips are in the same rate area, so just pick one
		plansForZip = s.AllPlans.SilverPlans(matchingZips[0].State, matchingZips[0].RateArea)
	}

	return plansForZip.SecondLowestRate()
}

type Zip struct {
	Code       string
	State      string
	CountyCode string
	Name       string
	RateArea   string
}
type Zips []Zip

// FindInOneRateArea returns all the zip records that match the supplied zip but only
// if they are in the same RateArea
func (zips Zips) FindInOneRateArea(zip string) (matches Zips) {
	area := ""
	for _, z := range zips {
		if z.Code == zip {
			if area == "" {
				area = z.RateArea
			}
			if area != z.RateArea {
				// ambiguous situation, return none
				return Zips{}
			}

			matches = append(matches, z)
		}
	}
	return matches
}

type Plan struct {
	ID         string
	State      string
	MetalLevel string
	Rate       float64
	RateArea   string
}
type Plans []Plan

// SilverPlans returns a list of Silver plans for the state and rateArea
func (plans Plans) SilverPlans(state, rateArea string) (matches Plans) {
	for _, p := range plans {
		if p.State == state &&
			p.RateArea == rateArea &&
			p.MetalLevel == Silver {

			matches = append(matches, p)
		}
	}

	return matches
}

// SecondLowestRate returns the Second Lowest Cost Plan Rate for the set
// It filters out duplicate rates
// It returns 0 if there is only one plan
func (plans Plans) SecondLowestRate() float64 {
	// first remove duplicates by rate
	uniqRates := make(map[float64]bool)
	sorted := Plans{}
	for _, p := range plans {
		if !uniqRates[p.Rate] {
			uniqRates[p.Rate] = true
			sorted = append(sorted, p)
		}
	}

	// if there is more than one then sort and return second lowest
	if len(sorted) > 1 {
		// sort for second last
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Rate < sorted[j].Rate
		})

		return sorted[1].Rate
	}

	return 0
}

func LoadZips() (allZips Zips) {
	deserializer := func(r []string) error {
		if len(r) < 5 {
			return errors.New("Zip expects 5 fields")
		}
		z := Zip{
			Code:       r[0],
			State:      r[1],
			CountyCode: r[2],
			Name:       r[3],
			RateArea:   r[4],
		}
		allZips = append(allZips, z)
		return nil
	}
	loadFile(zipsFile, deserializer)
	return allZips
}

func LoadPlans() (allPlans Plans) {
	deserializer := func(r []string) error {
		if len(r) < 5 {
			return errors.New("Plan expects 5 fields")
		}
		rate, err := strconv.ParseFloat(r[3], 64)
		if err != nil {
			fmt.Println(r)
			return err
		}

		p := Plan{
			ID:         r[0],
			State:      r[1],
			MetalLevel: r[2],
			Rate:       rate,
			RateArea:   r[4],
		}
		allPlans = append(allPlans, p)
		return nil
	}
	loadFile(plansFile, deserializer)
	return allPlans
}

func loadFile(fileName string, deserializeFunc func(r []string) error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic("Expected file: " + fileName)
	}
	defer file.Close()

	r := csv.NewReader(file)
	readHeader := false

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("Failed to read: " + zipsFile)
		}

		if readHeader {
			if err = deserializeFunc(record); err != nil {
				panic("Can't deserialize field: " + err.Error())
			}
		} else {
			readHeader = true
		}
	}
}
