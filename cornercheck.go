package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GET /ventes_immobilieres/offres/rhone_alpes/rhone/?f=a&th=1&ret=1 HTTP/1.1
// Host: www.leboncoin.fr
// Connection: keep-alive
// Pragma: no-cache
// Cache-Control: no-cache
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
// User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.94 Safari/537.36
// Referer: http://www.leboncoin./ventes_immobilieres/offres/rhone_alpes/rhone/?f=a&th=1&ret=1
// Accept-Encoding: gzip,deflate,sdch
// Accept-Language: fr-FR,fr;q=0.8,en-US;q=0.6,en;q=0.4
// Cookie:
// layout=0;
// 	xtvrn=$266818$;
// 	s=red1x3aa997d8f544b3f4132f0ea0d22046956f12a198;
// 	OAX=Wh1JgFQJ8IcAA7hF;
// 	weboForOas={"weboQueryDate":"2014-09-06-12-26","weboCalls":2,"clusters":"","social_demo":"","oasCalls":1};
// 	lazyLoadCounterAppear=13;
// 	location_search_22_169_anse_69480=Anse%2069480:8;
// 	sq=ca=22_s&w=169&c=9&f=a&th=1&ret=1;
// 	cookieFrame=2;
// 	RMFD=011XQDCh6B06BZ|6207QX|6307px|61088K|62089t;
// 	RMFS=011XQHtcU2080N;
// 	is_new_search=1

func main() {
	c := &http.Client{}
	s, err := runReq(c, "http://www.leboncoin.fr/ventes_immobilieres/offres/rhone_alpes/rhone/?f=a&th=1&ret=1")
	if err != nil {
		log.Printf("error reading string: %v", err)
		return
	}
	fmt.Printf("%v\n", s)
}

func addCookies(c *http.Client) {
	url := &c.Jar.SetCookies(u, cookies)
}

func runReq(c *http.Client, url string) (string, error) {
	addCookies(c)

	resp, err := c.Get(url)
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
