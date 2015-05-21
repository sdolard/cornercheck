package request

import (
	"cornercheck/annonce"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

var httpClient = initHTTPClient()

const (
	lbcHTMLCharset = "ISO 8859-15"
	baseURL        = "http://www.leboncoin.fr"
)

// AppParams ...
type AppParams struct {
	Category string
	Region   string
	Area     string
	NumCPU   int
}

// GetPage ...
func GetPage(page int, params AppParams, done chan []annonce.Annonce) {
	url := buildURL(params, page)
	log.Printf("Getting page %v (%v)", page, url)
	s, err := request(httpClient, url)
	if err != nil {
		log.Printf("Error running request: %v", err)
		return
	}

	annnonces := parseRequestedHTMLPage(s, params.Category, url)
	log.Printf("Getting page %v (%v) DONE", page, url)
	done <- annnonces
}

func buildURL(params AppParams, page int) string {
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
	if page >= 1 {
		url += fmt.Sprintf("&o=%v", page)
	}
	// q=   : query. TODO ?
	//log.Printf("url: %v", url)
	return url
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

func parseRequestedHTMLPage(page string, category string, url string /*, db *sql.DB*/) []annonce.Annonce {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	listRootNode := getListRootNode(doc)
	if listRootNode == nil {
		fmt.Printf("No annonce found\n")
		return nil
	}

	nodes := getAnnonceNodes(listRootNode, url)
	annonces := annonce.ExtractAnnoncesData(nodes, category)
	return annonces
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
