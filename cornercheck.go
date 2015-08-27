// Package cornercheck collect le bon coin data
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/sdolard/cornercheck/annonce"
	"github.com/sdolard/cornercheck/categories"
	"github.com/sdolard/cornercheck/db"
	"github.com/sdolard/cornercheck/regions"
	"github.com/sdolard/cornercheck/request"
)

func initFlags() (request.AppParams, error) {
	params := request.AppParams{
		Category: categories.Get()[categories.DefaultCategoryIndex],
		Region:   regions.DefaultRegion,
		NumCPU:   runtime.NumCPU(), // logical CPUs on the local machine
	}

	flag.StringVar(&params.Category, "category", params.Category, "\r\n\tValues: "+strings.Join(categories.Get(), ", "))
	flag.StringVar(&params.Region, "region", params.Region, regions.ToHelpString())
	flag.IntVar(&params.NumCPU, "numcpu", params.NumCPU, "Used cpu")

	flag.Parse()

	// category
	if categories.IndexOf(params.Category) == -1 {
		return params, fmt.Errorf("Invalid category: '%v'", params.Category)
	}
	log.Printf("category: %v", params.Category)

	// region
	r, a, err := regions.GetRegionAndArea(params.Region)
	if err != nil {
		return params, err
	}
	params.Region = r
	params.Area = a
	log.Printf("region: %v; area: %v", params.Region, params.Area)

	// NumCPU
	if params.NumCPU < 1 {
		params.NumCPU = 1
	}

	return params, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func createDb() *sql.DB {
	dbInstance := db.Open()
	annonce.CreateAnnonceTable(dbInstance)
	return dbInstance
}

func main() {
	appParams, err := initFlags()
	if err != nil {
		log.Printf("%v", err)
		printUsage()
		return
	}

	dbInstance := createDb()
	defer db.Close()

	page := 1
	cAnnonces := make(chan []annonce.Annonce)
	for {
		for i := 0; i < appParams.NumCPU; i++ {
			go request.GetPage(page, appParams, cAnnonces)
			page++
		}
		if annonce.Insert(cAnnonces, dbInstance, appParams.NumCPU) == 0 {
			break
		}
	}

}
