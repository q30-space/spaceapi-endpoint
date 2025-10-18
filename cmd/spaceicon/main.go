package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var debug = os.Getenv("SPACEICON_DEBUG") == "1"

const (
	iconOpen   = "󰯉"
	iconError  = ""
	colorRed   = "\x1b[31m"
	colorGreen = "\x1b[32m"
	colorReset = "\x1b[0m"
)

func printRed(s string)   { fmt.Printf("%s%s%s\n", colorRed, s, colorReset) }
func printGreen(s string) { fmt.Printf("%s%s%s\n", colorGreen, s, colorReset) }

func printGreenI3Blocks(icon string) {
	fmt.Println(icon)
	fmt.Println("SpaceAPI")
	fmt.Println("#228800")
}

func printRedI3Blocks(icon string) {
	fmt.Println(icon)
	fmt.Println("SpaceAPI")
	fmt.Println("#FF0F0F")
}

func fetchOpen(url string) (bool, bool) { // (value, ok)
	client := &http.Client{Timeout: 6 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, false
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, false
	}
	if debug {
		fmt.Fprintf(os.Stderr, "BODY: %s\n", string(body))
	}

	var root map[string]any
	decErr := json.Unmarshal(body, &root)
	if decErr != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "JSON ERROR: %v\n", decErr)
		}
		return false, false
	}

	stateRaw, ok := root["state"]
	if !ok || stateRaw == nil {
		if debug {
			fmt.Fprintf(os.Stderr, "NO state\n")
		}
		return false, false
	}

	stateMap, ok := stateRaw.(map[string]any)
	if !ok {
		if debug {
			fmt.Fprintf(os.Stderr, "state not object\n")
		}
		return false, false
	}

	openRaw, ok := stateMap["open"]
	if !ok {
		if debug {
			fmt.Fprintf(os.Stderr, "NO open\n")
		}
		return false, false
	}
	if debug {
		fmt.Fprintf(os.Stderr, "openRaw=%T %v\n", openRaw, openRaw)
	}

	switch v := openRaw.(type) {
	case bool:
		return v, true
	case string:
		if v == "true" || v == "1" {
			return true, true
		}
		if v == "false" || v == "0" {
			return false, true
		}
		return false, false
	case float64:
		return v != 0, true
	default:
		return false, false
	}
}

func main() {
	i3block := flag.Bool("i3block", false, "output for i3blocks: 1st line icon, 2nd line hex color")
	flag.Parse()

	// Expect exactly one remaining arg: URL
	if flag.NArg() != 1 {
		if *i3block {
			printRedI3Blocks(iconError)
		} else {
			printRed(iconError)
		}
		return
	}
	url := flag.Arg(0)

	val, ok := fetchOpen(url)
	if *i3block {
		if ok {
			if val {
				printGreenI3Blocks(iconOpen)
			} else {
				printRedI3Blocks(iconOpen)
			}
		} else {
			printRedI3Blocks(iconError)
		}
		return
	}

	if ok {
		if val {
			printGreen(iconOpen)
		} else {
			printRed(iconOpen)
		}
	} else {
		printRed(iconError)
	}
}
