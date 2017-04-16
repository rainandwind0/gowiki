package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/fatih/color"
)

/*

|=============================================================|
| ______     ______     __     __     __     __  __     __    |
|/\  ___\   /\  __ \   /\ \  _ \ \   /\ \   /\ \/ /    /\ \   |
|\ \ \__ \  \ \ \/\ \  \ \ \/ ".\ \  \ \ \  \ \  _"-.  \ \ \  |
| \ \_____\  \ \_____\  \ \__/".~\_\  \ \_\  \ \_\ \_\  \ \_\ |
|  \/_____/   \/_____/   \/_/   \/_/   \/_/   \/_/\/_/   \/_/ |
|                                                             |
|------ A command line utility for Wikipedia in GO Lang ------|
|=============================================================|

*/

// === Application Setup ====================================================

const AppVersion = "0.5"

// create a new http client for making requests
var client = &http.Client{Timeout: 10 * time.Second}

// temporary structure for holding json in an intermediate format
type tmp [][]string

// colored ouput functions
var infoColor = color.New(color.FgHiBlue).PrintlnFunc()
var inputColor = color.New(color.FgRed).PrintlnFunc()
var titleColor = color.New(color.FgGreen, color.Bold).PrintlnFunc()

// The final structure for holding the JSON results
type result struct {
	search string
	titles []string
	descs  []string
	links  []string
}

/*type query struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages struct {
			Num struct {
				Pageid  int
				Ns      int
				Title   string
				Extract string
			}
		}
	}
}*/

type queryIndex map[string]map[string][]string
type query map[string]map[string]map[string]map[string]interface{}

var fin = new(result)
var limitFlag = flag.Int("l", 5, "the number of results to be displayed")

// use for writing to file
var writeToFile = flag.Bool("f", false, "write aritcle to file instead of console")
var DEFAULT_FILE_NAME = "gowiki_search_"

func main() {

	// === Flags and arguments ===================================================

	// The flag package provides a default help printer via -h switch

	versionFlag := flag.Bool("v", false, "Print the version number.")
	searchFlag := flag.String("s", "", "a search value")
	flag.Parse() // Scan the arguments list

	// === Header ================================================================

	printHeader()

	// === Search ================================================================

	if *versionFlag {
		infoColor("Version:\t", AppVersion)
	}

	if *searchFlag != "" {
		searchWiki(*searchFlag, *limitFlag)
	} else {
		searchPrompt()
	}

}

// === Helper Functions ======================================================

/*
	|	searchWiki 	- a function to run a search on Wikipedia based on a search string
	|		search 	- string 		: the string to be searched for
	|		limit 	- int 			: the number of results to be displayed (the limit flag defaults to 5)
*/
func searchWiki(search string, limit int) {

	fin.search = search
	inputColor("Searching for:\t", search)

	esc := url.QueryEscape(search)

	body := getResults(esc, strconv.Itoa(limit))

	s, err := loadSearch(body)
	if err != nil {}

	printResults(s)

	if len(fin.titles) == 0 {

		text := ""

		reader := bufio.NewReader(os.Stdin)
		inputColor("No results found. Try new search? (y/n) ")
		text, _ = reader.ReadString('\n')

		if text[0] == 'y' {
			searchPrompt()
		} else {
			os.Exit(0)
		}
	}

	readArticle()
}

/*
	|	getResults 	- get the results from the wikipedia server
	|		search 	- string 		: the string to be searched for
	|		limit 	- int 			: the number of results to be displayed (the limit flag defaults to 5)
	|		return 	- byte[] 		: a json result in a byte array format
*/
func getResults(search, limit string) []byte {
	res, err := http.Get("https://en.wikipedia.org/w/api.php?action=opensearch&search=" + search + "&limit=" + limit + "&namespace=0&format=json")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	return []byte(body)
}

/*
	|	loadSearch 	- unmarshals a json byte array into a temporary 2d array of strings
	|		body 	- []byte 		: the byte array to unmarshal
	|		return 	- [][]string 	: a 2d array of json data
*/
func loadSearch(body []byte) (*tmp, error) {

	var s = new(tmp)

	err := json.Unmarshal(body, &s)

	if err != nil {
		// un-comment bellow to see errors for json parsing
		//fmt.Println("whoops:", err)
	}

	return s, err
}

/*
	|	printResults - print the results that are in the final struct for the user to read
*/
func printResults(s *tmp) {
	for i, entry := range *s {

		if i == 1 {
			fin.titles = entry
		}
		if i == 2 {
			fin.descs = entry
		}
		if i == 3 {
			fin.links = entry
		}

	}

	inputColor("")

	for i, entry := range fin.titles {
		var title = entry
		var desc = fin.descs[i]
		var link = fin.links[i]

		titleColor(strconv.Itoa(i) + ": " + title)
		titleColor("(" + link + ")")
		titleColor("================================================================================")
		inputColor("")
		infoColor(desc)
		inputColor("")
	}
}

/*
	|	getArticle 	- request the article from the wikipedia api using the index provided
	|		index	- The index of the article to read from the list of search results
	|		return	- The byte array of the articles data that is ready to be unmarshaled
*/
func getArticle(index int) []byte {
	var esc = url.QueryEscape(fin.titles[index])
	res, err := client.Get("https://en.wikipedia.org/w/api.php?format=json&action=query&indexpageids=&prop=extracts&explaintext=&titles=" + esc + "&format=json")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}
	return body
}

/*
	|	parseArticle 	- unmarshals the json byte array of the article to be read to get the page data
	|		body		- The body to be unmarshaled
	|		return		- The unmarshaled pages
*/
func parseArticle(body []byte) (*query, error) {
	var q = new(query)
	err := json.Unmarshal(body, &q)
	if err != nil {
		// un-comment bellow to see errors for json parsing
		// fmt.Println("whoops:", err)
	}

	return q, err
}

/*
	|	parseArticleIndex 	- unmarshals the json byte array of the article to be read to get the page indexes
	|		body			- The body to be unmarshaled
	|		return			- The unmarshaled indexes
*/
func parseArticleIndex(body []byte) (*queryIndex, error) {
	var q = new(queryIndex)
	err := json.Unmarshal(body, &q)
	if err != nil {
		// fmt.Println("whoops:", err)
	}

	return q, err
}

/*
	|	readArticle 	- collects the index from the user to begin loading the article to read
*/
func readArticle() {
	var text = ""
	var check = 0
	check, err := strconv.Atoi(text)

	for err != nil {
		reader := bufio.NewReader(os.Stdin)
		inputColor("Enter an index to read more: ")
		text, _ = reader.ReadString('\n')
		check, err = strconv.Atoi(text[0 : len(text) - 1])
		if err != nil {
			inputColor(err)
		}
		if check < 0 || check > len(fin.titles) {
			inputColor("Invalid index.")
			os.Exit(1)
		}
	}

	body := getArticle(check)
	q, err := parseArticle(body)
	qi, err := parseArticleIndex(body)

	if *writeToFile {
		writeFile(qi, q)
	} else {
		inputColor("reading entry:", text)
		for i, entry := range (*qi)["query"]["pageids"] {
			if i == 0 {
				titleColor("================================================================================")
				titleColor("=", (*q)["query"]["pages"][entry]["title"])
				titleColor("================================================================================")
				prettyPrintPage((*q)["query"]["pages"][entry]["extract"].(string))
				// infoColor((*q)["query"]["pages"][entry]["extract"])
			}
		}
	}
	reader := bufio.NewReader(os.Stdin)
	inputColor("Return to results? (y/n) ")
	text, _ = reader.ReadString('\n')

	if text[0] == 'y' {
		searchWiki(fin.search, *limitFlag)
	} else {
		os.Exit(0)
	}

	inputColor("page: ", (*q)["query"])
}

/*
	|	writeFile 	- write selected article to file instead of standard out
	|		qi			- The unmarshaled indexes of articles to be read
	|		q			- The unmarshaled pages to write to file
*/
func writeFile(qi *queryIndex, q *query) {
	fileWriter, fileError := os.Create(DEFAULT_FILE_NAME)
	var articleTitle = ""

	// catch exception
	if fileError != nil {
		panic(fileError)
	}

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer fileWriter.Close()

	for i, entry := range (*qi)["query"]["pageids"] {
		if i == 0 {
			fileWriter.WriteString("================================================================================\n")
			fileWriter.WriteString("= " + (*q)["query"]["pages"][entry]["title"].(string) + "\n")
			fileWriter.WriteString("================================================================================\n")
			//prettyPrintPage((*q)["query"]["pages"][entry]["extract"].(string))
			fileWriter.WriteString((*q)["query"]["pages"][entry]["extract"].(string) + "\n")
		}
		articleTitle = (*q)["query"]["pages"][entry]["title"].(string)
	}
	// rename default name to one that contains the article title
	os.Rename(DEFAULT_FILE_NAME, (DEFAULT_FILE_NAME + articleTitle))
}

/*
	|	prettyPrintPage - Print the page information in a beautified format (highlight the titles)
*/
func prettyPrintPage(page string) {

	scanner := bufio.NewScanner(strings.NewReader(page))

	for scanner.Scan() == true {
		if strings.Contains(scanner.Text(), "==") { // is some sort of heading
			if !strings.Contains(scanner.Text(), "===") { // is a main heading
				titleColor(strings.Replace(scanner.Text(), "=", "", -1))
				titleColor("--------------------------------------------------------------------------------")
			} else if !strings.Contains(scanner.Text(), "====") { // is a subheading
				titleColor(">" + strings.Replace(scanner.Text(), "=", "", -1))
			} else { // is a sub-subheading or lower
				titleColor("   >>" + strings.Replace(scanner.Text(), "=", "", -1))
			}

			scanner.Scan()
		}

		infoColor(scanner.Text())
	}
}

/*
	|	searchPrompt - prompts user for an input string and calls searchWiki on that string
*/
func searchPrompt() {

	reader := bufio.NewReader(os.Stdin)
	inputColor("Enter a search string: ")
	search, _ := reader.ReadString('\n')

	searchWiki(search, *limitFlag)
}

/*
	|	printHeader - print the fancy ASCII art header
*/
func printHeader() {
	infoColor(`|=============================================================|`)
	infoColor(`| ______     ______     __     __     __     __  __     __    |`)
	infoColor(`|/\  ___\   /\  __ \   /\ \  _ \ \   /\ \   /\ \/ /    /\ \   |`)
	infoColor(`|\ \ \__ \  \ \ \/\ \  \ \ \/ ".\ \  \ \ \  \ \  _"-.  \ \ \  |`)
	infoColor(`| \ \_____\  \ \_____\  \ \__/".~\_\  \ \_\  \ \_\ \_\  \ \_\ |`)
	infoColor(`|  \/_____/   \/_____/   \/_/   \/_/   \/_/   \/_/\/_/   \/_/ |`)
	infoColor(`|                                                             |`)
	infoColor(`|------ A command line utility for Wikipedia in GO Lang ------|`)
	infoColor(`|=============================================================|`)
	println("")
}
