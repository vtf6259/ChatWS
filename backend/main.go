package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

const version = "Indev 0.0.2"
const sensOn = false
const inProd = false

var motd string = "Welcome to the ChatWS development server! Expect bugs\n"
var customMotd bool = false
var genTls bool = false
var useTls bool = false
var pgppubkeycrypto *crypto.Key

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
var pgpprivkey string = ""
var pgppubkey string = ""

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
	fmt.Println("	--port: The port to run the server on. [REQUIRED]")
	fmt.Println("	--gentls: Generate a tls cert")
	fmt.Println("	--usetls: Use tls")
	fmt.Println("	--tlspubkey: The tls public key")
	fmt.Println("	--tlsprivkey: The tls private key")
	fmt.Println("	--pgppubkey: The pgp public key. [REQUIRED]")
	fmt.Println("	--pgpprivkey: The pgp private key. [REQUIRED]")
	fmt.Println("	--motd: The message of the day default in prod is \"Welcome to a ChatWS instance!\\n\" default in dev is \"Welcome to the ChatWS development server! Expect bugs\\n\"")
	fmt.Println("	--nomotd: Have no motd")
	fmt.Printf("Example for tls: %s --tlspubkey <PUBLIC KEY> --tlsprivkey <PRIVATE KEY> --usetls --port 8098\n", progname(os.Args[0]))
	fmt.Printf("Example for no tls: %s --port 8098 --pgppubkey dev.pub.pgp --pgpprivkey dev.priv.pgp\n", progname(os.Args[0]))
}

func main() {
	if len(os.Args) == 1 {
		displayHelp()
		os.Exit(0)
	}
	var indexToIgnore []int = []int{0}
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
			indexToIgnore = append(indexToIgnore, index+1)
			if err != nil {
				fmt.Printf("Invalid port %s\n", os.Args[index+1])
				os.Exit(0)
			}
			port = portToSet
		case "--motd":
			indexToIgnore = append(indexToIgnore, index+1)
			motd = os.Args[index+1]
			customMotd = true
		case "--nomotd":
			motd = "*NOMOTD*"
			customMotd = true
		case "--pgppubkey":
			prodLogf("pgppubkey: %s\n", sens(os.Args[index+1]))
			indexToIgnore = append(indexToIgnore, index+1)
			pgppubkey = os.Args[index+1]
		case "--pgpprivkey":
			prodLogf("pgpprivkey: %s\n", sens(os.Args[index+1]))
			indexToIgnore = append(indexToIgnore, index+1)
			pgpprivkey = os.Args[index+1]
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
			//prodLogf("indexToIgnore: ")
			//prodLogln(indexToIgnore)
			//prodLogf("index: %d\n", index)
			//prodLogf("")
			var aIndex int
			for _, i := range indexToIgnore {
				//prodLogln(index == i)
				if index == i {
					aIndex = i
					break
				}
				if index == 0 {
					continue
				}
			}
			if index == aIndex { // there is certently a better way to do this
				continue
			}
			fmt.Printf("Invalid argument %s\n", arg)
			os.Exit(1)
			continue
		}
	}
	if pgppubkey == "" {
		fmt.Println("Please provide the public key")
		os.Exit(1)
	}
	if pgpprivkey == "" {
		fmt.Println("Please provide the private key")
		os.Exit(1)
	}
	pgpprivkeydata, err := os.ReadFile(pgpprivkey)
	if err != nil {
		fmt.Println("Failed to read pgpprivkey")
		os.Exit(1)
	}
	pgpprivkey = string(pgpprivkeydata)
	pgppubkeydata, err := os.ReadFile(pgppubkey)
	if err != nil {
		fmt.Println("Failed to read pgppubkey")
		os.Exit(1)
	}
	pgppubkeycrypto, err = crypto.NewKeyFromArmored(string(pgppubkeydata))
	fmt.Printf("%s", logoansi)
	fmt.Printf("Running on port: %d\n", port) // dntger
	runServer(port)
}
