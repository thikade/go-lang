package search

import (
	"fmt"
	"math/rand"
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

func (obj *SearchObj) ExecuteSearch() bool {
	var results int = rand.Intn(30) + 1
	obj.TotalResults = results
	obj.Results = make([]string, results, results)
	for index := range obj.Results {
		obj.Results[index] = fmt.Sprintf("%d", (index+1)*3)
	}
	return true
}
