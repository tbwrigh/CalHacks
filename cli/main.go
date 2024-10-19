package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Check if there are arguments provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'install' or 'uninstall' command")
		os.Exit(1)
	}

	// Switch between different commands
	switch os.Args[1] {
	case "install":
		installCmd()
	case "uninstall":
		uninstallCmd()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func installCmd() {
	// Define the --codeql flag
	codeqlFlag := flag.NewFlagSet("install", flag.ExitOnError)
	withCodeQL := codeqlFlag.Bool("codeql", false, "Install only CodeQL")

	// Parse flags for the install command
	codeqlFlag.Parse(os.Args[2:])

	// Check the value of --codeql and display the appropriate message
	if *withCodeQL {
		fmt.Println("Installing CodeQL only")
	} else {
		fmt.Println("Installing autolock 🔒")
	}
}

func uninstallCmd() {
	// No flags for uninstall, just a simple function
	fmt.Println("Uninstalling autolock 🔓")
}
