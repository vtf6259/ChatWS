package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const version = "Indev 0.0.1"
const sensOn = false
const inProd = false

var genTls bool = false
var useTls bool = false

func prodLogln(toLog ...any) (any, error) {
	if inProd {
		return nil, nil
	}
	res, err := fmt.Println(toLog...)
	return res, err
}
func prodLogf(toLog string, formats ...any) (any, error) {
	if inProd {
		return nil, nil
	}
	res, err := fmt.Printf(toLog, formats...)
	return res, err
}

func sens(sensitive any) any {
	if sensOn {
		return "REDACTED"
	} else {
		return sensitive
	}
}

var port int

//go:embed logo.ansi
var logoansi string // generated with https://github.com/cacalabs/libcaca/tree/69a42132350da166a98afe4ab36d89008197b5f2

func progname(input string) string {
	lastSlashIndex := strings.LastIndex(input, "/")
	if lastSlashIndex != -1 {
		return input[lastSlashIndex+1:]
	} else {
		lastSlashIndex := strings.LastIndex(input, "\\")
		if lastSlashIndex != -1 {
			fmt.Println("Windows style paths detectedd (\\) we do not do anything different other than path stuff so do not expect stuff to 100% work on windows")
			return input[lastSlashIndex+1:]
		}
	}
	return input
}

func displayHelp() {
	fmt.Printf("[%s]: The ChatWS Backend\n", progname(os.Args[0]))
	fmt.Println("Options: ")
	fmt.Println("	--help: Display this help message and exit")
	fmt.Println("	--version: Display the version and exit")
	fmt.Println("	--port: The port to run the server on. (The default port for the ChatWS frontend is 8098)")
	fmt.Println("	--gentls: Generate a tls cert")
	fmt.Println("	--usetls: Use tls")
	fmt.Println("	--tlspubkey: The tls public key")
	fmt.Println("	--tlsprivkey: The tls private key")
}

func main() {
	if len(os.Args) == 1 {
		displayHelp()
		os.Exit(0)
	}
	var indexToIgnore int
	for index, arg := range os.Args {
		switch arg {
		case "--help":
			displayHelp()
			os.Exit(0)
		case "--version":
			fmt.Printf("%s Version %s\n", progname(os.Args[0]), version)
			os.Exit(0)
		case "--port":
			portToSet, err := strconv.Atoi(os.Args[index+1])
			indexToIgnore = index + 1
			if err != nil {
				fmt.Printf("Invalid port %s\n", os.Args[index+1])
				os.Exit(0)
			}
			port = portToSet
		case "--gentls":
			genTls = true
			fmt.Println("Unimplemented")
		case "--usetls":
			useTls = true
			fmt.Println("Unimplemented")
		case "--tlspubkey":
			fmt.Println("Unimplemented")
		case "--tlsprivkey":
			fmt.Println("Unimplemented")
		default:
			if index != 0 && index != indexToIgnore {
				fmt.Printf("Unknown option %s\n", arg)
				os.Exit(1)
			}
			continue
		}
	}
	fmt.Printf("%s", logoansi)
	fmt.Printf("Running on port: %d\n", port) // dntger
	runServer(port)
}
