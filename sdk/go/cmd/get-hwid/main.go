package main

import (
	"fmt"

	"github.com/cheo/licensemanager/sdk/go/license"
)

func main() {
	// The HWID logic doesn't require connecting to the server,
	// so we can initialize with a dummy config.
	client := license.New(license.Config{})
	fmt.Println("Your Machine ID (HWID):")
	fmt.Println(client.HWID())
}