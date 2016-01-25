package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pilu/robin"
)

func main() {
	var (
		configPath    string
		scraperName   string
		logLevelError bool
		logLevelInfo  bool
		logLevelDebug bool
	)

	runOpts := robin.DefaultRunnerOptions

	flag.StringVar(&configPath, "f", "", "path to the configuration file")
	flag.StringVar(&scraperName, "s", "", "scraper to run")
	flag.BoolVar(&logLevelError, "v", false, "log level error")
	flag.BoolVar(&logLevelInfo, "vv", false, "log level info")
	flag.BoolVar(&logLevelDebug, "vvv", false, "log level debug")
	flag.Parse()

	if configPath == "" {
		fmt.Printf("-c flag is mandatory\n")
		os.Exit(1)
	}

	switch {
	case logLevelDebug:
		runOpts.LogLevel = robin.LogLevelDebug
	case logLevelInfo:
		runOpts.LogLevel = robin.LogLevelInfo
	case logLevelError:
		runOpts.LogLevel = robin.LogLevelError
	}

	c, err := robin.NewConfigFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	scrapers, err := c.Scrapers()
	if err != nil {
		log.Fatal(err)
	}

	if len(scrapers) == 0 {
		log.Fatal("no scrapers found in config file")
	}

	if len(scrapers) > 1 && scraperName == "" {
		log.Fatal("config file contains multiple scrapers. you must specify the one you want to run with the -s flag")
	} else if len(scrapers) == 1 {
		for name, _ := range scrapers {
			scraperName = name
		}
	}

	r := robin.NewRunner(scrapers, runOpts)
	exp := robin.NewStdoutExporter()
	err = r.Run(scraperName, exp)
	if err != nil {
		log.Fatal(err)
	}
}
