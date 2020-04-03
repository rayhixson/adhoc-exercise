package main

import (
	"testing"
)

var zips Zips
var plans Plans

// TestGiven is an example that was given in the requirements
func TestGiven(t *testing.T) {
	zips = LoadZips()
	plans = LoadPlans()
	expect(t, "64148", 245.20)
}

// TestMoreThanOneRateArea checks for 0 because this is ambiguous per requirements
func TestMoreThanOneRateArea(t *testing.T) {
	zips = makeZips()
	plans = makePlans()
	expect(t, "36804", 0)
}

// TestTwoZips specifies two zips in the same rate area - should then match plans
func TestTwoZips(t *testing.T) {
	zips = Zips{
		Zip{"36804", "AL", "01017", "Chambers", "13"},
		Zip{"36804", "AL", "01087", "Macon", "13"},
	}
	plans = makePlans()

	expect(t, "36804", 253.37)
}

// TestFromData is from manual confirming some values
func TestFromData(t *testing.T) {
	zips = LoadZips()
	plans = LoadPlans()

	expect(t, "67118", 212.35)
	expect(t, "48435", 0)
	expect(t, "31551", 290.6)
}

func TestNoPlan(t *testing.T) {
	zips = Zips{{"35956", "AL", "01055", "Etowah", "8"}}
	plans = Plans{}

	expect(t, "35956", 0)
}

func TestOnePlan(t *testing.T) {
	zips = Zips{{"35956", "AL", "01055", "Etowah", "8"}}
	plans = Plans{{"52161YL6358432", "AL", "Silver", 245.82, "8"}}

	expect(t, "35956", 0)
}

func TestTwoPlans(t *testing.T) {
	zips = makeZips()
	plans = makePlans()

	expect(t, "36749", 245.83)
}

func expect(t *testing.T, zip string, rate float64) {
	finder := SlcspFinder{zips, plans}
	foundRate := finder.Find(zip)
	if foundRate != rate {
		t.Errorf("Want: %v -> Got: %v", rate, foundRate)
	}
}

func makeZips() Zips {
	return Zips{
		Zip{"36749", "AL", "01001", "Autauga", "11"},
		Zip{"35035", "AL", "01007", "Bibb", "3"},
		Zip{"35956", "AL", "01055", "Etowah", "8"},
		Zip{"36804", "AL", "01017", "Chambers", "13"},
		Zip{"36804", "AL", "01081", "Lee", "2"},
		Zip{"36804", "AL", "01087", "Macon", "13"},
		Zip{"36804", "AL", "01113", "Russell", "4"},
	}
}

func makePlans() Plans {
	return Plans{
		Plan{"52161YL6358432", "AL", "Silver", 245.82, "11"},
		Plan{"52161YL6358432", "AL", "Silver", 245.83, "11"},
		Plan{"31727UX6116202", "AL", "Gold", 312.06, "9"},
		Plan{"01100AO4222848", "AL", "Silver", 271.77, "5"},
		Plan{"24848KC5063721", "AL", "Silver", 264.84, "1"},
		Plan{"89885YK0256118", "AL", "Silver", 269.11, "8"},
		Plan{"74985TS1756968", "AL", "Silver", 271.77, "8"},
		Plan{"72404YS5031234", "AL", "Silver", 256.21, "8"},
		Plan{"93056UJ0123812", "AL", "Catastrophic", 202.44, "13"},
		Plan{"42009XZ4981402", "AL", "Bronze", 203.6, "13"},
		Plan{"30278PO0677161", "AL", "Bronze", 200.12, "13"},
		Plan{"14561QN7177699", "AL", "Silver", 273.9, "13"},
		Plan{"74948OF2563421", "AL", "Gold", 297.21, "13"},
		Plan{"54828LI9664121", "AL", "Silver", 248.03, "13"},
		Plan{"73838XS0335937", "AL", "Silver", 263.1, "13"},
		Plan{"58894BQ9115557", "AL", "Gold", 290.62, "13"},
		Plan{"49608TQ1551616", "AL", "Platinum", 328.11, "13"},
		Plan{"52295JD3261558", "AL", "Silver", 268.82, "13"},
		Plan{"16912HI8598649", "AL", "Silver", 267.3, "13"},
		Plan{"88913TF7204052", "AL", "Silver", 253.37, "13"},
		Plan{"16793FB8831732", "AL", "Gold", 291.89, "13"},
		Plan{"43897IO5130308", "AL", "Gold", 322.98, "13"},
		Plan{"09860NB8166613", "AL", "Bronze", 209.03, "13"},
		Plan{"51620SI2858327", "AL", "Gold", 321.05, "13"},
		Plan{"62963EU5292623", "AL", "Catastrophic", 173.08, "13"},
		Plan{"64052HO8609985", "AL", "Platinum", 386.34, "13"},
		Plan{"07819PS9333132", "AL", "Bronze", 225.75, "13"},
		Plan{"53153DT4217449", "AL", "Gold", 357.38, "2"},
		Plan{"69867KN9409819", "AL", "Gold", 338.69, "2"},
		Plan{"01919JZ8796954", "AL", "Bronze", 271.47, "2"},
		Plan{"90884WN5801293", "AL", "Silver", 323.25, "2"},
		Plan{"64625YH1566980", "AL", "Silver", 304.67, "2"},
		Plan{"62232AZ0247293", "AL", "Catastrophic", 182.59, "2"},
		Plan{"04751XT0241314", "AL", "Silver", 261.66, "2"},
		Plan{"66776YS7340042", "AL", "Gold", 349.45, "2"},
		Plan{"31595LL1746365", "AL", "Silver", 277.55, "2"},
	}
}
