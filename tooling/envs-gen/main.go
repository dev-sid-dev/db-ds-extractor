package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var (
	configFile  = os.Args[1]
	outputFile  = os.Args[2]
	packageName = os.Args[3]
)

func init() {

	if configFile == "" || outputFile == "" || packageName == "" {
		panic(fmt.Errorf("you must specify the location of a toml config source, the output file destination, and package name"))
	}

	//update base paths, so we can call via go generate from anywhere.
	config, err := filepath.Abs(configFile)
	if err != nil {
		panic(err)
	}
	configFile = config

	output, err := filepath.Abs(outputFile)
	if err != nil {
		panic(err)
	}

	outputFile = output

}

func main() {

	fmt.Println(fmt.Sprintf("Updating available env variables from %s", configFile))

	// read the whole content of file and pass it to file variable, in case of error pass it to err variable, defer
	// closing is done internally.
	configFile, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("error opening file: %w", err))
	}

	cfg := &AppConfig{PackageName: packageName}

	err = toml.Unmarshal(configFile, cfg)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("table").Parse(tmpl))

	file, err := os.Create(fmt.Sprintf("%s", outputFile))

	if err != nil {
		log.Println(fmt.Sprintf("error creating: %s, %s", outputFile, err))
		return
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(fmt.Sprintf("error closing: %s, %s", outputFile, err))
		}
		for envName, defaultValue := range cfg.Runtime.EnvVariables {
			fmt.Println(fmt.Sprintf("Added environment variable %s with default of \"%s\"", envName, defaultValue))
		}
		fmt.Println("Env variables updated.")
	}()

	err = t.Execute(file, cfg)

	if err != nil {
		log.Println(fmt.Sprintf("error error executing template for: %s, %s", outputFile, err))
	}

}
