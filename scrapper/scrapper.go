package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	company  string
	location string
	summary  string
}

func Scrapper(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob
	c := make(chan []extractedJob)

	pages := getPages(baseURL)
	for i := 0; i < pages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < pages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)

	}
	writeJobs(term, jobs)
	fmt.Println("Done! extract ", len(jobs), "jobs")
}

func writeJobs(term string, jobs []extractedJob) {
	c := make(chan []string)
	file, err := os.Create(term + "_jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	defer file.Close()
	headers := []string{
		"Id",
		"Title",
		"Company",
		"Location",
		"Summary",
	}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		go writeJobDetail(job, c)

	}
	for i := 0; i < len(jobs); i++ {
		jobSlice := <-c
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func writeJobDetail(job extractedJob, c chan<- []string) {
	c <- []string{
		"https://kr.indeed.com/viewjob?jk=" + job.id + " ",
		job.title,
		job.company,
		job.location,
		job.summary,
	}
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(50*page)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkStatus(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCard := doc.Find(".jobsearch-SerpJobCard")
	searchCard.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})
	for i := 0; i < searchCard.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)

	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	company := cleanString(card.Find(".company").Text())
	location := cleanString(card.Find(".location").Text())
	summary := cleanString(card.Find(".summary").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {

	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find("#searchCountPages").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		jobs := re.FindAllString(s.Text(), -1)[1]
		jobs = strings.Replace(jobs, ",", "", -1)
		num, _ := strconv.Atoi(jobs)
		pages = num / 50
		if num%50 != 0 {
			pages += 1
		}
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatus(req *http.Response) {
	if req.StatusCode != 200 {
		log.Fatalln("Request failed with Status :", req.StatusCode)
	}
}
