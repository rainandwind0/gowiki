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

var fin = new(result)

func main() {

	// === Flags and arguments ===================================================

	// The flag package provides a default help printer via -h switch

	versionFlag := flag.Bool("v", false, "Print the version number.")
	searchFlag := flag.String("s", "", "a search value")
	limitFlag := flag.Int("l", 5, "the number of results to be displayed")
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

		reader := bufio.NewReader(os.Stdin)
		inputColor("Enter a search string: ")
		search, _ := reader.ReadString('\n')

		searchWiki(search, *limitFlag)

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
	if err != nil {

	}

	printResults(s)

	reader := bufio.NewReader(os.Stdin)
	inputColor("Enter an index to read more: ")
	text, _ := reader.ReadString('\n')
	inputColor("reading entry:", text)
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
