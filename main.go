package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type CON struct {
	CAR string
	CDR string
}
func (c CON) String() string {
	return fmt.Sprintf("(%s . \"%s\")", c.CAR, c.CDR)
}

type AList []CON
func (al AList) String() string {
	var cons []string
	for _, con := range al {
		cons = append(cons, con.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(cons, " "))
}

type ElispPackageNameFlag string
func (f *ElispPackageNameFlag) String() string {
	return string(*f)
}
func (f *ElispPackageNameFlag) Set(value string) error {
	blank := regexp.MustCompile("[^\\w-_]+")
	*f = ElispPackageNameFlag(blank.ReplaceAll([]byte(value), []byte("-")))

	return nil
}

type ConfigFileNamesFlag []string
func (f *ConfigFileNamesFlag) String() string {
	return strings.Join(*f, ",")
}
func (f *ConfigFileNamesFlag) Set(value string) error {
	for _, fileName := range strings.Split(value, ",") {
		*f = append(*f, fileName)
	}
	return nil
}

func GenElispPackageCode(packageName string, env AList) string {
	return fmt.Sprintf(`;; Code automatically generated with //
(defvar %s/env %v "docstring")

(defun %s/get (name)
)

(provide '%s)
`, packageName, env, packageName, packageName)
}

func LoadOSEnvVariables(list *AList) {
	for _, varStr := range os.Environ() {
		keyValue := strings.Split(varStr, "=")
		*list = append(*list, CON{CAR: keyValue[0], CDR: keyValue[1]})
	}
}

func LoadConfigFileEnvVariables(list *AList, configFileNames []string) {
	env := make(map[string]string)
	for _, fileName := range configFileNames {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	env, err = godotenv.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	}
	for key, value := range env {
		*list = append(*list, CON{CAR: key, CDR: value})
	}
}

func main() {
	var elispPackageName ElispPackageNameFlag
	flag.Var(&elispPackageName, "p", "The name of the generated elisp package.")
	var configFileNames ConfigFileNamesFlag
	flag.Var(&configFileNames, "f", "Comma-separated list of env files to parse.")

	flag.Parse()

	if elispPackageName == "" {
		elispPackageName = "env"
	}
	var list AList

	LoadOSEnvVariables(&list)
	LoadConfigFileEnvVariables(&list, configFileNames)


	code := GenElispPackageCode(string(elispPackageName), list)
	os.WriteFile(fmt.Sprintf("%s.el", elispPackageName), []byte(code), 0644)
}
