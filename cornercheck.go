package main

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/charset"
	"flag"
	"fmt"
	"github.com/sdolard/cornercheck/annonce"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"runtime"
	"strings"
)

const (
	BASE_URL               = "http://www.leboncoin.fr"
	DEFAULT_CATEGORY_INDEX = 0 // voitures
	DEFAULT_REGION_INDEX   = 0 // rhone_alpes
	LBC_HTML_CHARSET       = "ISO 8859-15"
	TIME_LAYOUT            = "02 Jan 06 15:04"
)

type AppParams struct {
	Category string
	Region   string
	Area     string
	Parse    bool
	MaxPage  int
	NumCpu   int
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

type Region struct {
	Name  string
	Areas []string
}

type LbcDate struct {
	Day  string
	Hour string
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

	var IleDeFranceAreas = []string{
		"paris",
		"seine_et_marne",
		"yvelines",
		"essonne",
		"hauts_de_seine",
		"seine_saint_denis",
		"val_de_marne",
		"val_d_oise",
	}

	return []Region{
		{
			"rhone_alpes",
			RhoneAlpesAreas,
		}, {
			"ile_de_france",
			IleDeFranceAreas,
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
	reader, err := charset.NewReader(resp.Body, LBC_HTML_CHARSET)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(reader)
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

func initFlags() (AppParams, error) {
	appParams := AppParams{
		Category: getCategories()[DEFAULT_CATEGORY_INDEX],
		Region:   getRegions()[DEFAULT_REGION_INDEX].Name,
		Parse:    true,
		MaxPage:  -1,               // no limits,
		NumCpu:   runtime.NumCPU(), // logical CPUs on the local machine
	}

	flag.StringVar(&appParams.Category, "category", appParams.Category, "Categories: todo")
	flag.StringVar(&appParams.Region, "region", appParams.Region, "Regions: todo")
	flag.BoolVar(&appParams.Parse, "parse", appParams.Parse, "Parse: todo")
	flag.IntVar(&appParams.MaxPage, "maxpage", appParams.MaxPage, "MaxPage: todo")
	flag.IntVar(&appParams.NumCpu, "numcpu", appParams.NumCpu, "NumCpu: todo")

	flag.Parse()

	// category
	if categoriesIndexOf(appParams.Category) == -1 {
		return appParams, fmt.Errorf("Invalid category: '%v'", appParams.Category)
	}
	log.Printf("category: %v", appParams.Category)

	// region
	r, a, err := getRegionAndArea(appParams.Region)
	if err != nil {
		return appParams, err
	}
	appParams.Region = r
	appParams.Area = a
	log.Printf("region: %v; area: %v", appParams.Region, appParams.Area)

	// NumCpu
	if appParams.NumCpu < 1 {
		appParams.NumCpu = 1
	}

	return appParams, nil
}

func buildUrl(appParams AppParams, page int) string {
	url := fmt.Sprintf("%v/%v/offres/", BASE_URL, appParams.Category)

	if appParams.Area == "" {
		url += appParams.Region
	} else {
		url += appParams.Region + "/" + appParams.Area
	}
	//url += "/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000"
	//url += "/?f=p&th=1&ps=8&pe=9"
	//th=1 : enable thumb display
	if page <= 1 {
		url += "/?f=p&th=1&ps=8&pe=9"
	} else {
		url += fmt.Sprintf("/?o=%v&th=1&ps=8&pe=9", page)
	}

	log.Printf("url: %v", url)
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

func getAnnonceNodes(listRoot *html.Node) []*html.Node {
	var nodes []*html.Node
	for c := listRoot.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {

			if c.Data == "a" {
				nodes = append(nodes, c)
			} else if c.Data == "div" {
				for _, a := range c.Attr {
					if a.Key != "class" && a.Val != "clear" {
						panic(fmt.Sprintf("Unexpected annonces root node (format change ? Key = %v, Val=)", a.Key, a.Val))
					}
				}
			} else {
				panic(fmt.Sprintf("Unexpected annonces root node (format change ?, Data: %v)", c.Data))
			}
		}
	}
	return nodes
}

func parseRequestedHTMLPage(page string, category string) int {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	listRootNode := getListRootNode(doc)
	if listRootNode == nil {
		fmt.Printf("No annonce found\n")
		return 0
	}

	nodes := getAnnonceNodes(listRootNode)
	fmt.Printf("Annonces: %v\n", len(nodes))
	annnonces := annonce.ExtractAnnoncesData(nodes, category)
	for _, ann := range annnonces {
		fmt.Printf("%v# %v: %v, %v-%v (%v), %v, %v\n", ann.Time.Format(TIME_LAYOUT), ann.Category, ann.Title, ann.MinPrice, ann.MaxPrice, ann.PriceString, ann.LbcId(), ann.HRef)
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

	httpClient := initHttpClient()

	page := 0
	cAnnoncesCount := make(chan int)
	quit := false
	for {
		for i := 0; i < appParams.NumCpu; i++ {
			go func(page int, done chan int) {
				s, err := request(httpClient, buildUrl(appParams, page))
				if err != nil {
					log.Printf("Error running request: %v", err)
					return
				}

				done <- parseRequestedHTMLPage(s, appParams.Category)
			}(page, cAnnoncesCount)
			page++
		}

		for i := 0; i < appParams.NumCpu; i++ {
			if <-cAnnoncesCount == 0 && !quit {
				quit = true
			}
		}
		if quit {
			break
		}
	}
}
