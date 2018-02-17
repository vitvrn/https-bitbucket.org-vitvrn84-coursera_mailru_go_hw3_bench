package main

/* --- OPTIMIZATIONS ---

1. regexp.match -> strings.contains
2. easyjson (user struct)
3. file: readAll -> readString
4. for browser := range browsers: 2x -> 1x

*/

import (
	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	//	"regexp"
	"strings"
	// "log"
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

	fileContents, err := ioutil.ReadAll(file) //TODO read when needed?
	if err != nil {
		panic(err)
	}

	//	r := regexp.MustCompile("@") //TODO compile in init? OR use string.find
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	//	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n") //TODO (split by line) read by line?

	//var line string //???
	//var buffer []byte //TODO set capacity?
	user := &User{}
	fmt.Fprintln(out, "found users:")
	for i := range lines {
		user.UnmarshalJSON([]byte(lines[i])) //???

		isAndroid := false
		isMSIE := false

		//TODO range browsers #1 (merge with #2)
		//		for _, browserRaw := range browsers {
		for _, browser := range user.Browsers {
			//TODO ??? if Android -> else MSIE //can't be both in one browser line
			if strings.Contains(browser, "Android") { //+
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}

			if strings.Contains(browser, "MSIE") { //+
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
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
		//==email := r.ReplaceAllString(user.Email, " [at] ") //TODO use strings module?
		//==foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)

		//				foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, r.ReplaceAllString(user.Email, " [at] ")) //TODO use strings module?
		//		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, r.ReplaceAllString(user.Email, " [at] "))

		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, strings.Replace(user.Email, "@", " [at] ", -1))
	}

	//fmt.Fprintln(out, "found users:\n"+foundUsers)
	//	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func main() {
	FastSearch(os.Stdout)
}
