package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"go-neo4j-query-builder/querybuilder"
)

func main() {
	// Define command-line flags
	jsonFilePath := flag.String("json", "", "Path to the JSON file")
	jsonString := flag.String("data", "", "JSON string")

	// Parse the flags
	flag.Parse()

	// Variable to hold the query data
	var queryData querybuilder.MiniQuery

	// Check if JSON file or JSON string is provided
	if *jsonFilePath != "" {
		// If file path is provided, read the file
		fileData, err := ioutil.ReadFile(*jsonFilePath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			return
		}
		// Parse the JSON from the file
		if err := json.Unmarshal(fileData, &queryData); err != nil {
			fmt.Println("Error parsing JSON from file:", err)
			return
		}
	} else if *jsonString != "" {
		// If JSON string is provided, parse it
		if err := json.Unmarshal([]byte(*jsonString), &queryData); err != nil {
			fmt.Println("Error parsing JSON string:", err)
			return
		}
	} else {
		// If no JSON input is provided, print usage and exit
		fmt.Println("Please provide either a JSON file path or a JSON string.")
		flag.Usage()
		return
	}

	// Generate the Cypher query
	cypherQuery, err := querybuilder.BuildCypher(queryData)
	if err != nil {
		fmt.Println("Error building Cypher query:", err)
		return
	}

	// Print the generated Cypher query
	fmt.Println(cypherQuery)
}
