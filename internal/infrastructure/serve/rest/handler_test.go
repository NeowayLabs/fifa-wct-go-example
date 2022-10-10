//go:build unit

package rest_test

import (
	"log"
	"os"
)

var testLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
