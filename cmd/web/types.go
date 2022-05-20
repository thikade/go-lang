package main

import "log"

type Pets []struct {
	Animal string
	Age    int
}

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	outLog   *log.Logger
}
