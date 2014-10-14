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
	parse    bool
)

func getCategories() []string {
	return []string{
		"voitures",
		"motos",
	}
}

type Annonce struct {
	HRef  string
	Title string
	Date  string
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
	flag.BoolVar(&parse, "parse", true, "Parse")

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

func parseRequestedHTMLPage(page string) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	list := getListRootNode(doc)
	getNodes(list)
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

func getNodes(list *html.Node) {
	var ann *Annonce
	var state string

	var ca = func(a *Annonce) *Annonce {
		if a == nil {
			a = new(Annonce)
		}
		return a
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {

			if n.Data == "div" {
				for _, a := range n.Attr {
					if n.Parent == list && a.Key == "class" && a.Val == "clear" {
						fmt.Printf("%v, %v, %v<br>\n", ann.Title, ann.HRef, ann.Date)
						ann = nil
						break
					}

					if a.Key == "class" && a.Val == "date" {
						state = "date"
						break
					}
				}
			}
			if n.Data == "a" {
				for _, a := range n.Attr {
					switch a.Key {
					case "href":
						ann = ca(ann)
						ann.HRef = a.Val
					case "title":
						ann = ca(ann)
						ann.Title = a.Val
					}
				}
			}
		}
		if n.Type == html.TextNode {
			if state == "date" {
				ann.Date = n.Data
				state = ""
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(list)
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

	if parse {
		parseRequestedHTMLPage(s)
	} else {
		fmt.Printf("%v\n", s)
	}
}
