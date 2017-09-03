package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kanocz/geocode"
)

func printJSON(j interface{}) {
	data, err := json.MarshalIndent(j, "", "  ")
	if nil != err {
		log.Fatalln("JSON encode error:", err)
	}

	fmt.Println(string(data))
}

func main() {

	if len(os.Args) != 3 && len(os.Args) != 2 {
		fmt.Printf("Usage: %s <address> [<components>]\n", os.Args[0])
		fmt.Printf("Example: %s 'Dortmunder Straße 2, Berlin, Germany'", os.Args[0])
		fmt.Printf("Example: %s 'Dortmunder Straße 2, Berlin' 'country:DE|postal_code:10555'", os.Args[0])
		return
	}

	addr := os.Args[1]
	components := ""
	if len(os.Args) == 3 {
		components = os.Args[2]
	}

	grqst := geocode.Request{
		Language:   "EN",
		Components: components,
		Address:    addr,
	}
	gresp, err := grqst.Lookup(nil)
	if nil != err {
		log.Fatalln("Google maps lookup error:", err)
	}

	fmt.Println("Result from google:")
	printJSON(gresp)

	parsed := gresp.Parse()
	fmt.Println("Parsed result from google:")
	printJSON(parsed)
}
