# GoWiki
### A command line utility in GO Lang for reading Wikipedia documents.

## Setting-up GoWiki
1. Install Go for your preferred [platform](https://golang.org/doc/install)  
2. Install the [Color](https://github.com/fatih/color) package as follow:  
    `go get install https://github.com/fatih/color`
3. Clone the GoWiki repository:  
    `git clone https://github.com/rainandwind0/gowiki.git`  
4. From the root of the project directory issue:  
    `go run src/main.go [OPTIONS]`
## Usage
To see usage statements: `go run src/main.go -h`  
    `-f write the article to file instead of console (default: false)`  
    `-l the of search results to be displayed (default: 5)`  
    `-s a search term`  
    `-v print version number`

### Chain any of the following flags together to get your results
Run with prompts: `go run src/main.go`  
Preset search term: `go run src/main.go -s="search string"`  
Limit the number of results to be displayed: `go run src/main.go -l=10`  
See the version number: `go run src/main.go -v`  
Search with a preset term with a set limit and the selected article is written to file:  
`go run src/main.go -f -l 10 -s "Oprah"`

*NOTE:* flag values can be set with the equal (`-s="Oprah"`) or without (`-s "Oprah"`)  
## Todo:
- [ ] Add option to select only intro of returned article  
- [x] Write to file  
- [x] Pretty print (especially with nested heading)  
- [ ] Binary installer
