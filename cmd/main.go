package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Qiryl/gurrency"
	"github.com/Qiryl/gurrency/service/fixer"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// cli
	base := pflag.StringP("base", "b", "EUR", "")
	refs := pflag.StringP("reference", "r", "CAD,USD", "")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	fixerSource := fixer.NewFixerSource(os.Getenv("FIXER_KEY"), os.Getenv("FIXER_URL"), *base, *refs)
	srcs := []*gurrency.Source{gurrency.NewSource(fixerSource)}

	// concurrent requests to the sources
	wg := sync.WaitGroup{}
	for _, s := range srcs {
		wg.Add(1)
		go s.GetRate(&wg)
	}
	wg.Wait()
}
