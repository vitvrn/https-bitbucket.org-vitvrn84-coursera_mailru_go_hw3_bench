package main

/* --- OPTIMIZATIONS ---

1. regexp.match -> strings.contains
2. easyjson (user struct)
3. file: readAll -> readString
4. for browser := range browsers: 2x -> 1x
5. ==TODO read file by line (don't use ReadAll)

*/

import (
	//"encoding/json"
	"fmt"
	"io"
	//	"io/ioutil"
	"os"
	//	"regexp"
	"strings"
	// "log"

	"bufio"
)

//easyjson:json
type User struct { //TODO use easyjson
	Browsers []string `json:"browsers"`
	//	Company  string   //--- `json:"company"`
	//	Country  string   //--- `json:"country"`
	Email string `json:"email"`
	//	Job      string   //--- `json:"job"`
	Name string `json:"name"`
	//	Phone    string   //--- `json:"phone"`
}

func init() {
	//TODO prepare data structures etc...
	println("Hello")
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath) //TODO (open input file) move to init?
	if err != nil {
		panic(err)
	}
	defer file.Close() //+

	reader := bufio.NewReader(file)

	seenBrowsers := make([]string, 0, 128) //TODO 32? 64? 128? ???
	uniqueBrowsers := 0
	buffer := make([]byte, 0, 1024) //TODO 256? 512? 1024? ???
	user := &User{}
	var readErr error

	fmt.Fprintln(out, "found users:")
	for i := 0; readErr == nil; i++ {
		buffer, readErr = reader.ReadSlice('\n')

		user.UnmarshalJSON(buffer)

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			//TODO ??? if Android -> else MSIE //can't be both in one browser line
			if strings.Contains(browser, "Android") { //+
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers { //TODO use map[string]struct{} ??
					if item == browser {
						notSeenBefore = false
						break //+
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
				continue //??? if Android -> else MSIE
			}

			if strings.Contains(browser, "MSIE") { //+
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
						break //+
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, strings.Replace(user.Email, "@", " [at] ", -1))
	}

	//	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func main() {
	FastSearch(os.Stdout)
}
