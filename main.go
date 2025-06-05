package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
	"golang.design/x/clipboard"
)

var (
	driver   *playwright.Playwright
	browser  playwright.Browser
	page     playwright.Page
	mapLinks []string
)

const URL string = "https://www.google.com/maps/search/"

// Reads through in.txt and returns each line as element of array.
func getSearchList() []string {
	file, err := os.Open("in.txt")
	if err != nil {
		log.Fatalln("Unable to open file in.txt: ", err)
	}

	scanner := bufio.NewScanner(file)

	var strArr []string
	for scanner.Scan() {
		strArr = append(strArr, scanner.Text())
	}

	return strArr
}

// Loops through array of provided links. Formats and writes to txt.
func outputLinks(links []string, searchList []string) error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}

	var copy string
	for index, link := range links {
		str := fmt.Sprintf("- [%v](%v)", searchList[index], link)
		copy += "\n" + str
		_, err := fmt.Fprintln(file, str)
		if err != nil {
			return err
		}
	}

	err = clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(copy))

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

// Searches for a location on google maps and fetches the share link which is appended to shareLinks.
func getMapsLink(url string) error {
	if _, err := page.Goto(url,
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		}); err != nil {
		return err
	}

	currentUrl := page.URL()
	if strings.Contains(currentUrl, "consent.google.com") {
		rejectBtn := page.GetByLabel("Reject all").First()
		err := rejectBtn.Click()
		if err != nil {
			return err
		}
	}

	shareButton := page.GetByLabel("share").First()
	err := shareButton.Click()
	if err != nil {
		return err
	}

	link, err := page.Locator(".vrsrZe:not(.vrsrZe--disabled)").First().InputValue()
	if err != nil {
		return err
	}

	mapLinks = append(mapLinks, link)

	return nil
}

func main() {
	mapLinks = nil

	err := playwright.Install()
	if err != nil {
		log.Fatalln("Err install playwright: ", err)
	}

	driver, err = playwright.Run()

	if err != nil {
		log.Fatalln("Error starting playwright: ", err)
	}

	browser, err = driver.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false)})
	if err != nil {
		log.Fatalln("Error launching browser: ", err)
	}

	page, err = browser.NewPage()
	if err != nil {
		log.Fatalln("Error create new page: ", err)
	}

	searchList := getSearchList()
	urls := createSearchStrings(searchList)

	for _, url := range urls {
		err := getMapsLink(url)
		if err != nil {
			log.Fatalln("Failed getting link for ", url, ", err: ", err)
		}
	}

	err = outputLinks(mapLinks, searchList)
	if err != nil {
		log.Fatalln("Failed writing out links")
	}

	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := driver.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	log.Println("\x1b[1;35mDone! Copied to clipboard, or view out.txt\x1b[0m")
}
