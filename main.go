package main

import (
	"bufio"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"strings"
)

var (
	driver  *playwright.Playwright
	browser playwright.Browser
	page    playwright.Page
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

func outputLinks(links []string, searchList []string) error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}

	for index, link := range links {
		str := fmt.Sprintf("- [%v](%v)", searchList[index], link)
		_, err := fmt.Fprintln(file, str)
		if err != nil {
			return err
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
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

func getMapsLink(url string) (string, error) {
	if _, err := page.Goto(url,
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); err != nil {
		return "", err
	}

	currentUrl := page.URL()
	if strings.Contains(currentUrl, "consent.google.com") {
		rejectBtn := page.GetByLabel("Reject all").First()
		err := rejectBtn.Click()
		if err != nil {
			return "", err
		}
	}
	shareButton := page.GetByLabel("share").First()
	err := shareButton.Click()
	if err != nil {
		return "", err
	}

	link, err := page.Locator(".vrsrZe:not(.vrsrZe--disabled)").First().InputValue()
	if err != nil {
		return "", err
	}

	return link, nil
}

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatalln("Err install playwrite: ", err)
	}

	driver, err = playwright.Run()

	if err != nil {
		log.Fatalln("Error starting playwrite: ", err)
	}

	browser, err = driver.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true)})
	if err != nil {
		log.Fatalln("Error launching firefox: ", err)
	}

	page, err = browser.NewPage()
	if err != nil {
		log.Fatalln("Error create new page: ", err)
	}

	searchList := getSearchList()
	urls := createSearchStrings(searchList)

	var links []string
	for _, url := range urls {
		link, err := getMapsLink(url)
		if err != nil {
			log.Fatalln("Failed geting link for ", url, ", err: ", err)
		}
		links = append(links, link)
	}

	err = outputLinks(links, searchList)
	if err != nil {
		log.Fatalln("Failed writing out links")
	}

	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := driver.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	log.Println("Done! Check out.txt")
}
