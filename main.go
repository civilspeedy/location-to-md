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
	driver     *playwright.Playwright
	browser    playwright.Browser
	page       playwright.Page
	mapLinks   []string
	urlArr     []string
	searchList []string
	fileLoc    string
)

const URL string = "https://www.google.com/maps/search/"

// Reads through in.txt and returns each line as element of array.
func getSearchList() error {
	file, err := os.Open(fileLoc)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		searchList = append(searchList, scanner.Text())
	}
	return nil
}

// Loops through provides location list, replaces spaces with "+"
// and then adds onto the end of the search url.
func createSearchStrings() {
	for _, location := range searchList {
		withPluses := strings.ReplaceAll(location, " ", "+")
		urlArr = append(urlArr, URL+withPluses)
	}
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

// Loops through array of provided links. Formats and writes to txt.
func outputLinks() error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}

	var copy string
	for index, link := range mapLinks {
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

func main() {
	mapLinks = nil
	urlArr = nil
	fileLoc = "in.txt"

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

	err = getSearchList()
	if err != nil {
		log.Fatalln("Unable to get searchList: ", err)
	}

	createSearchStrings()

	for _, url := range urlArr {
		err := getMapsLink(url)
		if err != nil {
			log.Fatalln("Failed getting link for ", url, ", err: ", err)
		}
	}

	err = outputLinks()
	if err != nil {
		log.Fatalln("Failed writing out links")
	}

	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err := driver.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

	fmt.Println("\x1b[1;35mDone! Copied to clipboard, or view out.txt\x1b[0m")
}
