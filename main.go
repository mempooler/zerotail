package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		filename = flag.String("f", "", "file to tail")
		follow   = flag.Bool("follow", false, "follow the file as it grows")
		lines    = flag.Int("n", 10, "number of lines to tail")
	)
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

	var lineCount int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
		if lineCount > *lines {
			continue
		}

		var fields map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &fields); err != nil {
			log.Error().Err(err).Msgf("failed to parse line: %q", scanner.Text())
			continue
		}

		logger := zerolog.New(os.Stdout).With().Fields(fields).Logger()
		logger.WithLevel(zerolog.Disabled).Msg("")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to read from file")
	}

	if *follow {
		followFile(file)
	}
}

func followFile(file *os.File) {
	file.Seek(0, os.SEEK_END)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var fields map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &fields); err != nil {
			log.Error().Err(err).Msgf("failed to parse line: %q", scanner.Text())
			continue
		}

		logger := zerolog.New(os.Stdout).With().Fields(fields).Logger()
		logger.WithLevel(zerolog.Disabled).Msg("")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to read from file")
	}
}
