package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/hpcloud/tail"
	"github.com/mempooler/zerolog"
	"github.com/mempooler/zerolog/log"
)

func main() {
	filename := flag.String("f", "", "file to tail")
	flag.Parse()

	if *filename == "" {
		fmt.Println("usage: tail -f <filename>")
		os.Exit(1)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	t, err := tail.TailFile(*filename, tail.Config{Follow: true})
	for line := range t.Lines {
		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line.Text), &fields); err != nil {
			log.Error().Err(err).Msgf("failed to parse line: %q", line.Text)
			continue
		}

		log.Log().Fields(fields).Msg("")
	}
}
