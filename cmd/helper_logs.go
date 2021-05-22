package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

// Centralize error handing, simple print message and exit
func logError(err error, message string) {
	fmt.Println()
	fmt.Println()
	fmt.Println(color.RedString("Error: " + message))
	log.Fatalln(err)
}

// Warning but dont exit
func logWarning(err error, message string) {
	fmt.Println()
	fmt.Println()
	fmt.Println(color.YellowString("Error: " + message))
}
