package main

//import "github.com/Lexty/yandexdnsapi"
import (
	"os"

	"fmt"

	"github.com/voxelbrain/goptions"
)

func main() {
	options := struct {
		Token  string        `goptions:"-t, --token, description='Your token'"`
		Domain string        `goptions:"-d, --domain, description='Domain name'"`
		Help   goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
		GetToken struct{} `goptions:"get-token"`
		List     struct {
			Command string   `goptions:"--command, mutexgroup='input', description='Command to exectute', obligatory"`
			Script  *os.File `goptions:"--script, mutexgroup='input', description='Script to exectute', rdonly"`
		} `goptions:"list"`
		//		Delete struct {
		//			Path  string `goptions:"-n, --name, obligatory, description='Name of the entity to be deleted'"`
		//			Force bool   `goptions:"-f, --force, description='Force removal'"`
		//		} `goptions:"delete"`
	}{ // Default values goes here

	}

	goptions.ParseAndFail(&options)

	fmt.Println(options.GetToken)
}
