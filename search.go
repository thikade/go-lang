package main

import (
	"regexp"
	"strings"
)

// var regexEmail = regexp.MustCompile(".+@.+\\..+")
var regexNumber = regexp.MustCompile("[0-9]+")

type SearchObj struct {
	Days   string
	Token  string
	Errors map[string]string
	// search response
	TotalResults int
	Results      []string
}

func (msg *SearchObj) Validate() bool {
	msg.Errors = make(map[string]string)

	match := regexNumber.Match([]byte(msg.Days))
	if match == false {
		msg.Errors["Days"] = "Please enter valid number of days!"
		// log.Println("DEBUG: errors.days set to nonzero value")
	}
	if strings.TrimSpace(msg.Token) == "" {
		msg.Errors["Token"] = "Please enter a Search Token!"
		// log.Println("DEBUG: errors.Token set to nonzero value")
	}

	return len(msg.Errors) == 0
}
