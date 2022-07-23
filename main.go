package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
)

func GetPageCount(link string) (Report, error) {
	res, err := http.Get(link)
	if err != nil {
		return Report{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Report{}, err
	}

	rs := bytes.NewReader(body)
	count, err := api.PageCount(rs, nil)
	if err != nil {
		return Report{}, err
	}

	return Report{link, count}, nil
}

type Report struct {
	Link  string
	Count int
}

func main() {

	webPage := "https://cure53.de/"
	doc, err := goquery.NewDocument(webPage)
	if err != nil {
		log.Fatal(err)
	}

	var links []string
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if strings.HasPrefix(link, "https://") && strings.Contains(link, "report") && strings.HasSuffix(link, ".pdf") {
			links = append(links, link)
		}
	})
	links = lo.Uniq(links)

	reports := parallel.Map(links, func(link string, _ int) Report {
		report, _ := GetPageCount(link)
		return report
	})

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Count > reports[j].Count
	})

	fmt.Println("link,page_count")
	for _, report := range reports {
		fmt.Printf("%s,%d\n", report.Link, report.Count)
	}
}
