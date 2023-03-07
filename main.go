package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mempooler/zerolog"
	"github.com/mempooler/zerolog/log"
)

func main() {
	filename = flag.String("f", "", "file to tail")
	flag.Parse()

	if *filename == "" {
		fmt.Println("usage: tail -f <filename>")
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open file")
	}
	defer file.Close()

	file.Seek(0, os.SEEK_END)
	scanner := bufio.NewScanner(file)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	for scanner.Scan() {
		var fields map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &fields); err != nil {
			log.Error().Err(err).Msgf("failed to parse line: %q", scanner.Text())
			continue
		}

		log.Log().Fields(fields).Msg("")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to read from file")
	}
}
