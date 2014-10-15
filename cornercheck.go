package main

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/charset"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_URL               = "http://www.leboncoin.fr"
	DEFAULT_CATEGORY_INDEX = 0 // voitures
	DEFAULT_REGION_INDEX   = 0 // rhone_alpes
	DATE_YESTERDAY         = "Hier"
	DATE_TODAY             = "Aujourd'hui"
	LBC_HTML_CHARSET       = "ISO 8859-15"
)

var (
	category string
	region   string
	area     string
	parse    bool
	page     int
)

func getCategories() []string {
	return []string{
		"voitures",
		"motos",
	}
}

type Annonce struct {
	HRef       string
	Title      string
	Time       time.Time
	TimeString string
	Category   string
}

type Region struct {
	Name  string
	Areas []string
}

type LbcDate struct {
	Day  string
	Hour string
}

func getLbcShortMonths() map[string]time.Month {
	return map[string]time.Month{
		"jan":  time.January,
		"fev":  time.February,
		"mar":  time.March,
		"avr":  time.April,
		"mai":  time.May,
		"juin": time.June,
		"juil": time.July,
		"août": time.August,
		"sept": time.September,
		"oct":  time.October,
		"nov":  time.November,
		"dec":  time.December,
	}
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

func initFlags() error {
	flag.StringVar(&category, "category", getCategories()[DEFAULT_CATEGORY_INDEX], "Categories")
	flag.StringVar(&region, "region", getRegions()[DEFAULT_REGION_INDEX].Name, "Regions")
	flag.IntVar(&page, "page", 0, "Page")
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
	//url += "/?f=p&th=1&ps=8&pe=9&ms=50000&me=125000"
	//url += "/?f=p&th=1&ps=8&pe=9"
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

func parseRequestedHTMLPage(page string) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	listRootNode := getListRootNode(doc)
	if listRootNode == nil {
		fmt.Printf("No annonce found\n")
		return
	}

	nodes := getAnnonceNodes(listRootNode)
	fmt.Printf("Annonces: %v\n", len(nodes))
	annnonces := extractAnnoncesData(nodes)
	for _, ann := range annnonces {
		fmt.Printf("%v# %v: %v, %v\n", ann.Time, ann.Category, ann.Title, ann.HRef)
	}
}

func lbcDateToTime(dayS, hourS string) (time.Time, string) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()

	// Hours > 13:52
	decomposedHour := strings.Split(hourS, ":")
	hour64, err := strconv.ParseInt(decomposedHour[0], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	hour := int(hour64)
	min64, err := strconv.ParseInt(decomposedHour[1], 10, 0)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	min := int(min64)

	// Day
	if dayS == DATE_YESTERDAY {
		// Hier 13:52
		d := now.AddDate(0, 0, -1)
		year = d.Year()
		month = d.Month()
		day = d.Day()
	} else if dayS == DATE_TODAY {
		// Aujourd'hui 13:52
		// Initialized data are valid for this case
	} else {
		// 28 sept
		decomposedDay := strings.Split(dayS, " ")
		day64, err := strconv.ParseInt(decomposedDay[0], 10, 0)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		day = int(day64)
		month = getLbcShortMonths()[decomposedDay[1]]
		if month == 0 {
			panic(fmt.Sprintf("Invalid month: %v", decomposedDay[1]))
		}
	}

	return time.Date(
		year,
		month,
		day,
		hour,
		min,
		0, // sec
		0, // nsec,
		now.Location(),
	), fmt.Sprintf("%v %v", dayS, hourS)
}

func getAnnonceDate(annNode *html.Node) (time.Time, string) {
	var date []string
	collect := false
	level := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if collect {
			if n.Type == html.TextNode {
				data := strings.TrimSpace(n.Data)
				if data != "" {

					date = append(date, data)
				}
			}
		} else {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "date" {
						collect = true
						level = 0
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			level++
			f(c)
			level--
			if level < 0 {
				collect = false
			}
		}
	}
	f(annNode)
	return lbcDateToTime(date[0], date[1])
}

func extractAnnoncesData(annNodes []*html.Node) []Annonce {
	annonces := make([]Annonce, len(annNodes))

	for i, annNode := range annNodes {
		if annNode.Data == "a" {
			annonces[i].Category = category
			for _, att := range annNode.Attr {
				switch att.Key {
				case "href":
					annonces[i].HRef = att.Val
				case "title":
					annonces[i].Title = att.Val
				}
			}

			annonces[i].Time, annonces[i].TimeString = getAnnonceDate(annNodes[i])
		} else {
			panic("format change")
		}
	}

	return annonces
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
