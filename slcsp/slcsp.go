package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

const Silver = "Silver"

// Given a file with zip codes, determine the second lowest cost silver plan
// for that zip code. Print out the zip along with the rate; blank for the rate
// if there is none or there is no definitive answer.
// See data/README.md for requirement details
func main() {
	// input that we will later print out along with rates
	const inputFile = "data/slcsp.csv"

	// reference data
	const zipsFile = "data/zips.csv"
	const plansFile = "data/plans.csv"

	zips := LoadZips(zipsFile)
	plans := LoadPlans(plansFile)
	slcsp := SlcspFinder{zips, plans}

	// read input file
	rows := [][]string{}
	rowHandler := func(record []string) error {
		if len(record) != 2 {
			return fmt.Errorf("Expecting 2 fields: %v", record)
		}

		zip := record[0]
		rate := slcsp.FindRate(zip)
		stringRate := ""
		if rate > 0 {
			stringRate = fmt.Sprintf("%.2f", rate)
		}
		rows = append(rows, []string{zip, stringRate})
		return nil
	}
	loadFile(inputFile, rowHandler)

	// write the file to stdout
	fmt.Println("zipcode,rate")
	for _, r := range rows {
		fmt.Printf("%s,%s\n", r[0], r[1])
	}
}

// SlcspFinder has a list of zips and a list of plans and finds rates given zips
type SlcspFinder struct {
	AllZips  Zips
	AllPlans Plans
}

// FindRate returns the second lowest cost silver plan for the given zip if there is a definitive one
func (s SlcspFinder) FindRate(givenZip string) float64 {
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
// if they are in the same State and RateArea, otherwise none
func (zips Zips) FindInOneRateArea(zip string) (matches Zips) {
	area := ""
	for _, z := range zips {
		if z.Code == zip {
			if area == "" {
				area = z.State + z.RateArea
			}
			if area != (z.State + z.RateArea) {
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

// SilverPlans returns a list of Silver plans for the state and rateArea.
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

// SecondLowestRate returns the Second Lowest Cost Plan Rate for the set of Plans.
// It filters out duplicate rates.
// It returns 0 if there is only one plan.
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
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Rate < sorted[j].Rate
		})

		// second lowest
		return sorted[1].Rate
	}

	return 0
}

func LoadZips(file string) (allZips Zips) {
	zipHandler := func(record []string) error {
		if len(record) < 5 {
			return fmt.Errorf("Zip expects 5 fields: %v", record)
		}
		allZips = append(allZips, Zip{
			Code:       record[0],
			State:      record[1],
			CountyCode: record[2],
			Name:       record[3],
			RateArea:   record[4],
		})
		return nil
	}
	loadFile(file, zipHandler)
	return allZips
}

func LoadPlans(file string) (allPlans Plans) {
	planHandler := func(record []string) error {
		if len(record) < 5 {
			return fmt.Errorf("Plan expects 5 fields: %v", record)
		}
		rate, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return fmt.Errorf("Processing record %v -> %w", record, err)
		}

		allPlans = append(allPlans, Plan{
			ID:         record[0],
			State:      record[1],
			MetalLevel: record[2],
			Rate:       rate,
			RateArea:   record[4],
		})
		return nil
	}
	loadFile(file, planHandler)
	return allPlans
}

// loadFile loads a csv file passing each row/record to the handler.
// It assumes the first row of the file is a header row
// It panics if file is missing since these are considered critical to the app.
func loadFile(fileName string, recordHandler func(record []string) error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panicln("Expected file:", fileName)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	readHeader := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicln("Failed to read:", fileName, err)
		}

		if readHeader {
			if err = recordHandler(record); err != nil {
				log.Panicf("Can't deserialize field [%v] : %v", record, err)
			}
		} else {
			readHeader = true
		}
	}
}
