# Location to MarkDown

## Description

Executable written in GoLang that uses browser automation to fetch the "Share" link for locations on google maps. This "Share" link is then formatted into a markdown link with the location name and output into a txt file, as well as copied to clipboard.

## Reason

While I was planning my travel break, it became rather tedious having to manually copy and format the locations I wanted to visit into my notes app ([Obsidian](https://obsidian.md/)). I felt this would be a good opportunity to explore a new programming language, and Golang caught my eye, as well as to learn about browser automation.

## External Libraries:

- [Playwright](https://github.com/playwright-community/playwright-go)

## Usage

- Ensure in.txt is in the same directory as the executable.
- List locations in in.txt with new lines between each one.
- Putting the city and country in the location name will yield faster and more accurate results.
- If Playwright is not installed on your device, the executable may take a moment to install the browser drivers.
- An out.txt file will be created with the formatted locations inside.
