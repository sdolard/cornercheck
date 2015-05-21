package annonce

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	dateYesterday = "Hier"
	dateToday     = "Aujourd'hui"
)

func getLbcShortMonths() map[string]time.Month {
	return map[string]time.Month{
		"jan":  time.January,
		"fev":  time.February,
		"mars": time.March,
		"avr":  time.April,
		"mai":  time.May,
		"juin": time.June,
		"juil": time.July,
		"août": time.August,
		"sept": time.September,
		"oct":  time.October,
		"nov":  time.November,
		"déc":  time.December,
	}
}

var (
	reLbcID        *regexp.Regexp
	reLbcPlacement *regexp.Regexp
	reLbcPrice     *regexp.Regexp
)

// Annonce struct
type Annonce struct {
	HRef            string
	Title           string
	Time            time.Time
	TimeString      string
	Category        string
	MaxPrice        int
	MinPrice        int
	PriceString     string
	Town            string
	Area            string
	PlacementString string
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
	if dayS == dateYesterday {
		// Hier 13:52
		d := now.AddDate(0, 0, -1)
		year = d.Year()
		month = d.Month()
		day = d.Day()
	} else if dayS == dateToday {
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
			panic(fmt.Sprintf("Invalid month: %v. dayS: %v; hourS: %v", decomposedDay[1], dayS, hourS))
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

func lbcPriceToInt(price string) (int, int) {
	if reLbcPrice == nil {
		reLbcPrice = regexp.MustCompile("[ €\u00a0]")
	}
	p := reLbcPrice.ReplaceAllLiteralString(price, "")
	if p == "" {
		return 0, 0
	}
	if strings.Contains(p, "-") {
		prices := strings.Split(p, "-")
		if len(prices) != 2 {
			panic(fmt.Sprintf("Invalid price format: '%v'", p))
		}
		MinPrice64, err := strconv.ParseInt(prices[0], 10, 0)
		if err != nil {
			log.Fatalf("Min price: '%v'; %v", prices[0], err)
		}
		MaxPrice64, err := strconv.ParseInt(prices[1], 10, 0)
		if err != nil {
			panic(fmt.Sprintf("Max price: '%v'; %v", prices[1], err))
		}
		return int(MinPrice64), int(MaxPrice64)
	}
	price64, err := strconv.ParseInt(p, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("price: '%v'; %v", price, err))
	}
	return int(price64), int(price64)
}

func getAnnoncePlacement(annNode *html.Node) (string, string, string) {
	placement := ""
	collect := false
	var f func(*html.Node)
	f = func(n *html.Node) {
		if collect {
			if n.Type == html.TextNode {
				placement = strings.TrimSpace(n.Data)
				return
			}
		} else {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "placement" {
						collect = true
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if placement != "" {
				break
			}
		}
	}
	f(annNode)

	if reLbcPlacement == nil {
		reLbcPlacement = regexp.MustCompile("[\r\n\t\\s]")
	}
	placement = reLbcPlacement.ReplaceAllLiteralString(placement, "")
	if strings.Contains(placement, "/") {
		places := strings.Split(placement, "/")
		return strings.TrimSpace(places[0]), strings.TrimSpace(places[1]), placement
	}
	return "", placement, placement
}

func getAnnoncePrice(annNode *html.Node) (int, int, string) {
	price := ""
	collect := false
	var f func(*html.Node)
	f = func(n *html.Node) {
		if collect {
			if n.Type == html.TextNode {
				price = strings.TrimSpace(n.Data)
				return
			}
		} else {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "price" {
						collect = true
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if price != "" {
				break
			}
		}
	}
	f(annNode)
	min, max := lbcPriceToInt(price)
	return min, max, price
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

// LbcID return annonce ID
func (a Annonce) LbcID() string {
	// http://www.leboncoin.fr/voitures/719527156.htm?ca=22_s
	u, err := url.Parse(a.HRef)
	if err != nil {
		log.Fatal(err)
	}
	if reLbcID == nil {
		reLbcID = regexp.MustCompile(".*/(\\d+)\\.htm")
	}
	subs := reLbcID.FindStringSubmatch(u.Path)
	if len(subs) < 2 || subs[1] == "" {
		panic("Format error in LbcID")
	}
	return subs[1]
}

// ExtractAnnoncesData returns array of annonce
func ExtractAnnoncesData(annNodes []*html.Node, category string) []Annonce {
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
			annonces[i].MinPrice, annonces[i].MaxPrice, annonces[i].PriceString = getAnnoncePrice(annNodes[i])
			annonces[i].Town, annonces[i].Area, annonces[i].PlacementString = getAnnoncePlacement(annNodes[i])

		} else {
			panic("format change")
		}
	}

	return annonces
}
