package main

/* --- OPTIMIZATIONS ---

1. regexp.match -> strings.contains
2. easyjson (user struct)
3. file: readAll -> readString|readSlice
4. for browser := range browsers: 2x -> 1x
5. seenBrowsers: slice -> map

*/

import (
	"fmt"
	"io"
	"os"
	"strings"

	"bufio"
)

//easyjson:json
type User struct {
	Browsers []string `json:"browsers"`
	//	Company  string   //--- `json:"company"`
	//	Country  string   //--- `json:"country"`
	Email string `json:"email"`
	//	Job      string   //--- `json:"job"`
	Name string `json:"name"`
	//	Phone    string   //--- `json:"phone"`
}

func init() {
	//TODO prepare data structures etc..?
	println("Preparing...")
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath) //TODO (open input file) move to init?
	if err != nil {
		panic(err)
	}
	defer file.Close() //+

	reader := bufio.NewReader(file)

	seenBrowsers := make(map[string]struct{}, 128)

	uniqueBrowsers := 0
	buffer := make([]byte, 0, 1024) //TODO 256? 512? 1024? ???
	user := &User{}

	fmt.Fprintln(out, "found users:")
	var readErr error
	for i := 0; readErr == nil; i++ {
		buffer, readErr = reader.ReadSlice('\n')

		user.UnmarshalJSON(buffer)

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			//TODO ??? if Android -> else MSIE //can't be both in one browser line
			if strings.Contains(browser, "Android") { //+
				isAndroid = true
				if _, seenBefore := seenBrowsers[browser]; !seenBefore {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
				continue //TODO ??? if Android -> else MSIE
			}

			if strings.Contains(browser, "MSIE") { //+
				isMSIE = true
				if _, seenBefore := seenBrowsers[browser]; !seenBefore {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, strings.Replace(string(user.Email), "@", " [at] ", -1))
	}

	fmt.Fprintln(out, "\nTotal unique browsers", uniqueBrowsers)
}

func main() {
	FastSearch(os.Stdout)
}
