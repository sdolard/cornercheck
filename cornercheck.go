// Package cornercheck collect le bon coin data
package main

import (
	"cornercheck/annonce"
	"cornercheck/regions"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"runtime"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

const (
	baseURL              = "http://www.leboncoin.fr"
	defaultCategoryIndex = 0 // _vehicules_
	lbcHTMLCharset       = "ISO 8859-15"
	timeLayout           = "02 Jan 06 15:04"
)

type appParams struct {
	Category string
	Region   string
	Area     string
	NumCPU   int
}

func getCategories() []string {
	return []string{
		"_vehicules_",
		"voitures",
		"motos",
		"_immobilier_",
		"_multimedia_",
		"_maison_",
		"_loisirs_",
		"_materiel_professionnel_",
		"_emploi_services_",
		"_",
		"autres",
	}
}

func request(c *http.Client, u string) (string, error) {
	resp, err := c.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	reader, err := charset.NewReader(resp.Body, lbcHTMLCharset)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusProxyAuthRequired {
		return "", errors.New(resp.Status)
	}
	return string(body), nil
}

func initHTTPClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar: cookieJar,
	}
}

func categoriesIndexOf(v string) int {
	for i, category := range getCategories() {
		if v == category {
			return i
		}
	}
	return -1
}

func initFlags() (appParams, error) {
	params := appParams{
		Category: getCategories()[defaultCategoryIndex],
		Region:   regions.DefaultRegion,
		NumCPU:   runtime.NumCPU(), // logical CPUs on the local machine
	}

	flag.StringVar(&params.Category, "category", params.Category, "\r\n\tValues: "+strings.Join(getCategories(), ", "))
	flag.StringVar(&params.Region, "region", params.Region, regions.ToHelpString())
	flag.IntVar(&params.NumCPU, "numcpu", params.NumCPU, "Used cpu")

	flag.Parse()

	// category
	if categoriesIndexOf(params.Category) == -1 {
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

func buildURL(params appParams, page int) string {
	url := fmt.Sprintf("%v/%v/offres/", baseURL, params.Category)

	if params.Area == "" {
		url += params.Region
	} else {
		url += params.Region + "/" + params.Area
	}
	url += "/?"
	// th=1 : enable thumb display
	url += "th=0"
	// it=1 : search only in title
	url += "&it=0"
	// o=   : page number
	if page > 1 {
		url += fmt.Sprintf("&o=%v", page)
	}
	// q=   : query. TODO ?
	//log.Printf("url: %v", url)
	return url
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func getListRootNode(root *html.Node) *html.Node {
	var list *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "list-lbc" {
					list = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
	return list
}

func getAnnonceNodes(listRoot *html.Node, url string) []*html.Node {
	var nodes []*html.Node
	for c := listRoot.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "a" {
				nodes = append(nodes, c)
			} else if c.Data == "div" {
				for _, a := range c.Attr {
					if !strings.Contains("classid", a.Key) && !strings.Contains("clearoas-x", a.Val) {
						panic(fmt.Sprintf("Unexpected annonces root node (format change ? Key='%v', Val='%v', Url='%v')", a.Key, a.Val, url))
					}
				}
			} else {
				panic(fmt.Sprintf("Unexpected annonces root node (format change ? Data='%v', Url='%v')", c.Data, url))
			}
		}
	}
	return nodes
}

func parseRequestedHTMLPage(page string, category string, url string) int {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	listRootNode := getListRootNode(doc)
	if listRootNode == nil {
		fmt.Printf("No annonce found\n")
		return 0
	}

	nodes := getAnnonceNodes(listRootNode, url)
	fmt.Printf("Annonces: %v\n", len(nodes))
	annnonces := annonce.ExtractAnnoncesData(nodes, category)
	for _, ann := range annnonces {
		fmt.Printf("%v# %v: %v, %v-%v (%v), %v, %v, %v\n",
			ann.Time.Format(timeLayout),
			ann.Category,
			ann.Title,
			ann.MinPrice,
			ann.MaxPrice,
			ann.PriceString,
			ann.PlacementString,
			ann.LbcID(),
			ann.HRef)
	}
	return len(nodes)
}

func main() {
	appParams, err := initFlags()
	if err != nil {
		log.Printf("%v", err)
		printUsage()
		return
	}

	httpClient := initHTTPClient()

	page := 0
	cAnnoncesCount := make(chan int)
	quit := false
	for {
		for i := 0; i < appParams.NumCPU; i++ {
			go func(page int, done chan int) {
				url := buildURL(appParams, page)
				s, err := request(httpClient, url)
				if err != nil {
					log.Printf("Error running request: %v", err)
					return
				}

				done <- parseRequestedHTMLPage(s, appParams.Category, url)
			}(page, cAnnoncesCount)
			page++
		}

		for i := 0; i < appParams.NumCPU; i++ {
			if <-cAnnoncesCount == 0 && !quit {
				quit = true
			}
		}
		if quit {
			break
		}
	}
}
