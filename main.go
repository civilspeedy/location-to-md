package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const URL string = "https://www.google.com/maps/search/"

// Reads through search.txt and returns each line as element of array.
func getSearchList() []string {
	file, err := os.Open("search.txt")
	if err != nil {
		log.Fatalln("Unable to open file search.txt: ", err)
	}

	scanner := bufio.NewScanner(file)

	var strArr []string
	for scanner.Scan() {
		strArr = append(strArr, scanner.Text())
	}

	return strArr
}

// Loops through provides location list, replaces spaces with "+"
// and then adds onto the end of the search url.
// List of URLS is returned.
func createSearchStrings(locations []string) []string {
	var urlArr []string
	for _, location := range locations {
		withPluses := strings.ReplaceAll(location, " ", "+")
		urlArr = append(urlArr, URL+withPluses)
	}
	return urlArr
}

func takeScreenshot(tab playwright.Page) {
	if _, err := tab.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshot.png"),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
}

func runScraper() {
	pw, err := playwright.Run()

	if err != nil {
		log.Fatalln("Error starting playwrite: ", err)
	}

	ff, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false)})
	if err != nil {
		log.Fatalln("Error launching firefox: ", err)
	}

	tab, err := ff.NewPage()
	if err != nil {
		log.Fatalln("Error create new page: ", err)
	}

	if _, err = tab.Goto("https://www.google.com/maps/search/Eiffel+Tower+Paris+France",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); err != nil {
		log.Fatalln("Could not go to url: ", err)
	}

	url := tab.URL()
	if strings.Contains(url, "consent.google.com") {
		rejectBtn := tab.GetByLabel("Reject all").First()
		err := rejectBtn.Click()
		if err != nil {
			log.Fatalln("Unable to click: ", err)
		}

		takeScreenshot(tab)
		time.Sleep(100 * time.Second)
	}

	if err = ff.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func main() {
	searchList := getSearchList()
	urls := createSearchStrings(searchList)
	fmt.Println(urls)
}
