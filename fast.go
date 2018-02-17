package main

/* --- OPTIMIZATIONS ---

1. regexp.match -> strings.contains
2. ==TODO easyjson (user struct)
3. file: readAll -> readString
4. ==TODO for browser := range browsers: 2x -> 1x

*/

import (
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"regexp"
	"strings"
	// "log"
	"bufio"
	//easyjson "github.com/mailru/easyjson"
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

	//fileContents, err := ioutil.ReadAll(file) //TODO read when needed?
	//if err != nil {
	//	panic(err)
	//}
	reader := bufio.NewReader(file)

	r := regexp.MustCompile("@") //TODO compile in init? OR use string.find
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	//lines := strings.Split(string(fileContents), "\n") //TODO (split by line) read by line?

	users := make([]map[string]interface{}, 0)
	//for _, line := range lines {
	var line string
	var err0 error
	for {
		line, err0 = reader.ReadString('\n') //TODO remove '\n'?

		user := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user) //TODO easyjson: struct - ???
		if err != nil {
			panic(err)
		}
		users = append(users, user) //TODO allocate earlier

		if err0 != nil {
			break
		}
	}

	//TODO 2x iterations -> 1x
	for i, user := range users { //TODO i is needed to be ouput

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		//TODO range browsers #1 (merge with #2)
		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

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
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func main() {
	FastSearch(os.Stdout)
}
