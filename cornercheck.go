package main

import (
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

const (
	BASE_URL               = "http://www.leboncoin.fr"
	DEFAULT_CATEGORY_INDEX = 0 // voitures
	DEFAULT_REGION_INDEX   = 0 // rhone_alpes
)

var (
	category string
	region   string
	area     string
)

func getCategories() []string {
	return []string{
		"voitures",
		"motos",
	}
}

type Region struct {
	Name  string
	Areas []string
}

func getRegions() []Region {
	var RhoneAlpesAreas = []string{
		"ain",
		"ardeche",
		"drome",
		"isere",
		"loire",
		"rhone",
		"savoie",
		"haute_savoie",
	}

	return []Region{
		{
			"rhone_alpes",
			RhoneAlpesAreas,
		},
	}
}

// GET /voitures/offres/rhone_alpes/rhone/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000 HTTP/1.1
// Host: www.leboncoin.fr
// Connection: keep-alive
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.101 Safari/537.36
// Referer: http://www.leboncoin.fr/voitures/offres/rhone_alpes/ain/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000
// Accept-Encoding: gzip,deflate,sdch
// Accept-Language: fr-FR,fr;q=0.8,en-US;q=0.6,en;q=0.4
// Cookie: xtvrn=$266818$;
//  OAX=Wh1ajFQmtwAADbeR;
//  location_search_22_1_toutes_les_communes_01600=Toutes%20les%20communes%2001600:3;
//  location_search_22_1_toutes_les_communes_69480=Toutes%20les%20communes%2069480:1;
//  sli=1;
//  lazyLoadCounterAppear=39;
//  weboForOas={"weboQueryDate":"2014-10-09-09-38","weboCalls":2,"clusters":"","audiences":"","social_demo":"","oasCalls":1};
//  RMFD=011Xc8J2O205fc!O107aY!O307aZ!O108WN!S208ZY!B508fi;
//  layout=0;
//  s=red1x490e57f2ad31c85f5fc275aa0354d81abdde2062;
//  sq=ca=22_s&w=101&c=2&f=p&th=1&ps=8&pe=9&ms=50000&me=125000;
//  cookieFrame=2;
//  is_new_search=1

// GET /voitures/offres/rhone_alpes/ain/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000 HTTP/1.1
// Host: www.leboncoin.fr
// Connection: keep-alive
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.101 Safari/537.36
// Referer: http://www.leboncoin.fr/voitures/offres/rhone_alpes/rhone/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000
// Accept-Encoding: gzip,deflate,sdch
// Accept-Language: fr-FR,fr;q=0.8,en-US;q=0.6,en;q=0.4
// Cookie: xtvrn=$266818$;
//  OAX=Wh1ajFQmtwAADbeR;
//  location_search_22_1_toutes_les_communes_01600=Toutes%20les%20communes%2001600:3;
//  location_search_22_1_toutes_les_communes_69480=Toutes%20les%20communes%2069480:1;
//  sli=1;
//  lazyLoadCounterAppear=39;
//  weboForOas={"weboQueryDate":"2014-10-09-09-38","weboCalls":2,"clusters":"","audiences":"","social_demo":"","oasCalls":1};
//  RMFD=011Xc8J2O205fc!O107aY!O307aZ!O108WN!S208ZY!B508fi;
//  layout=0;
//  sq=ca=22_s&w=169&c=2&f=p&th=1&ps=8&pe=9&ms=50000&me=125000;
//  cookieFrame=2;
//  s=red1x490e57f2ad31c85f5fc275aa0354d81abdde2062;
//  is_new_search=1

func addCookies(c *http.Client, u string) error {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		log.Printf("error parsing string: %v", u)
		return err
	}
	cookies := []*http.Cookie{
		{Name: "sq", Value: "ca=22_s&w=101&c=2&f=p&th=1&ps=8&pe=9&ms=50000&me=125000", Path: "/", Domain: ".leboncoin.fr"},
	}

	c.Jar.SetCookies(parsedUrl, cookies)
	return nil
}

func request(c *http.Client, u string) (string, error) {
	//err := addCookies(c, u)
	// if err != nil {
	// 	return "", err
	// }

	resp, err := c.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func initHttpClient() *http.Client {
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

func getRegionAndArea(v string) (string, string, error) {
	for _, r := range getRegions() {
		if v == r.Name {
			return r.Name, "", nil
		}

		for _, area := range r.Areas {
			if v == area {
				return r.Name, area, nil
			}
		}
	}
	return "", "", fmt.Errorf("Invalid region: '%v'", v)
}

func initFlags() error {
	flag.StringVar(&category, "category", getCategories()[DEFAULT_CATEGORY_INDEX], "Categories")
	flag.StringVar(&region, "region", getRegions()[DEFAULT_REGION_INDEX].Name, "Regions")

	flag.Parse()

	// category
	if categoriesIndexOf(category) == -1 {
		return fmt.Errorf("Invalid category: '%v'", category)
	}
	log.Printf("category: %v", category)

	// region
	r, a, err := getRegionAndArea(region)
	if err != nil {
		return err
	}
	region = r
	area = a
	log.Printf("region: %v; area: %v", region, area)

	return nil
}

func buildUrl() string {
	url := fmt.Sprintf("%v/%v/offres/", BASE_URL, category)

	if area == "" {
		url += region
	} else {
		url += region + "/" + area
	}
	url += "/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000"
	log.Printf("url: %v", url)
	return url
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	err := initFlags()
	if err != nil {
		log.Printf("%v", err)
		printUsage()
		return
	}

	s, err := request(initHttpClient(), buildUrl())
	if err != nil {
		log.Printf("Error running request: %v", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(s)) // TODO: Parse node
	if err != nil {
		log.Printf("Error parsing response: %v", err)
		return
	}

	fmt.Printf("%v\n", s)
}
