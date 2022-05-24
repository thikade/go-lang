package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// var regexEmail = regexp.MustCompile(".+@.+\\..+")
var regexNumber = regexp.MustCompile("[0-9]+")

func (obj *SearchObj) Validate() bool {
	obj.Errors = make(map[string]string)

	match := regexNumber.Match([]byte(obj.Days))
	if match == false {
		obj.Errors["Days"] = "Please enter valid number of days!"
		// log.Println("DEBUG: errors.days set to nonzero value")
	}
	if strings.TrimSpace(obj.Token) == "" {
		obj.Errors["Token"] = "Please enter a Search Token!"
		// log.Println("DEBUG: errors.Token set to nonzero value")
	}

	return len(obj.Errors) == 0
}

func (obj *SearchObj) ExecuteSearch() string {
	var jql string = `project = Infrastructure AND resolutiondate > "-500d"`
	queryParams := fmt.Sprintf("maxResults=1&jql=%s", url.QueryEscape(jql))
	jiraUrl := url.URL{
		// User:        &url.Userinfo{},
		Scheme:     "https",
		Host:       "jira.spring.io",
		Path:       "/rest/api/2/search",
		RawQuery:   queryParams,
		ForceQuery: false,
		// RawPath:     "",
	}
	strUrl := jiraUrl.String()
	fmt.Printf("curl: %s\n", strUrl)
	resp, err := http.Get(strUrl)
	if err != nil {
		fmt.Printf("curl error")
		return ""
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("body read error")
	// 	return nil
	// }
	// return body

	retValue := ""

	// parse json result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if val, ok := result["total"]; ok {
		// switch val.(type) {
		// case bool:
		// 	fmt.Printf("bool: %v\n", val)
		// case string:
		// 	fmt.Printf("string: %v\n", val)
		// case int:
		// 	fmt.Printf("int: %v\n", val)
		// case float64:
		// 	fmt.Printf("float64: %v\n", val)
		// default:
		// 	fmt.Printf("unknown: %v\n", val)
		// }
		retValue = fmt.Sprintf("%v", val)
	}

	return retValue

}
