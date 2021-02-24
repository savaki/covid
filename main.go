package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Handler struct {
}

type Location struct {
	Zip                  string
	Street               string
	StoreNumber          int
	State                string
	OpenTimeslots        int
	OpenAppointmentSlots int
	Name                 string
	Latitude             float64
	Longitude            float64
	City                 string
}

func newHandler() *Handler {
	return &Handler{}
}

func findOpenLocations() ([]Location, error) {
	resp, err := http.Get("https://heb-ecom-covid-vaccine.hebdigital-prd.com/vaccine_locations.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vaccine locations: %w", err)
	}
	defer resp.Body.Close()

	var content struct {
		Locations []Location
	}
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode vaccine locations: %w", err)
	}

	var got []Location
	for _, location := range content.Locations {
		if location.OpenAppointmentSlots > 0 || location.OpenTimeslots > 0 {
			got = append(got, location)
		}
	}

	return got, nil
}

func main() {
	locations, err := findOpenLocations()
	if err != nil {
		log.Fatalln(err)
	}

	n := len(locations)
	fmt.Printf("found %v locations with openings\n", n)
	if n > 0 {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(locations)
	}
}
