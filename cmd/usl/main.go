package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alaturka/gbreve/net/usl"
	"github.com/alaturka/gbreve/text/textutil"
)

var (
	version string //nolint
	build   string //nolint
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s USL [flags...] [attributes...]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")

	flag.PrintDefaults()

	os.Exit(2)
}

func cry(message ...interface{}) {
	fmt.Fprintln(os.Stderr, append([]interface{}{"usl:"}, message...)...)
}

func die(message ...interface{}) {
	cry(message...)

	os.Exit(1)
}

func wanted(defaultAttributes []string, attributes ...string) []string {
	if len(attributes) > 0 {
		return attributes
	}

	return defaultAttributes
}

// Print should be commented
func Print(us *usl.USL, templateMap map[string]string, attributes ...string) {
	m, ks := us.MapCustom(templateMap)

	var pairs []string

	for _, attribute := range wanted(ks, attributes...) {
		if value, ok := m[attribute]; ok {
			s := fmt.Sprintf("%s='%s'", attribute, value)
			pairs = append(pairs, s)
		}
	}

	fmt.Println(strings.Join(pairs, " "))
}

// Bash should be commented
func Bash(variable string, us *usl.USL, templateMap map[string]string, attributes ...string) {
	m, ks := us.MapCustom(templateMap)

	var pairs []string

	pairs = append(pairs, variable+"=(")

	for _, attribute := range wanted(ks, attributes...) {
		if value, ok := m[attribute]; ok {
			s := fmt.Sprintf("[%s]='%s'", attribute, value)
			pairs = append(pairs, s)
		}
	}

	pairs = append(pairs, ")")

	fmt.Println(strings.Join(pairs, " "))
}

type varFlags []string

func (v *varFlags) String() string {
	return fmt.Sprintf("%v", *v)
}

func (v *varFlags) Set(value string) error {
	*v = append(*v, value)

	return nil
}

func main() {
	var variables varFlags

	flag.Usage = usage

	allowLocalPath := flag.Bool("local", false, "Allow local paths while parsing.")
	bashArray := flag.String("bash", "", "Print result as a Bash associated array with the given name.")
	flag.Var(&variables, "var", `Set variable template as 'variable="template"'.`)

	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	args := flag.Args()

	parser := usl.Parse
	if *allowLocalPath {
		parser = usl.ParseMayLocalPath
	}

	us, err := parser(args[0])
	if err != nil {
		die(err)
	}

	templateMap := map[string]string{}

	for _, expr := range variables {
		kv := map[string]string{}
		err := textutil.ParseAssignment(expr, kv)

		if err != nil {
			die(err)
		}

		for k, v := range kv {
			templateMap[k] = v
		}
	}

	if *bashArray != "" {
		Bash(*bashArray, us, templateMap, args[1:]...)
	} else {
		Print(us, templateMap, args[1:]...)
	}
}
