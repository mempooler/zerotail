package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hpcloud/tail"
	"github.com/mempooler/zerolog"
)

func main() {
	filename := flag.String("f", "", "file to tail")
	flag.Parse()

	if *filename == "" {
		fmt.Println("usage: tail -f <filename>")
		os.Exit(1)
	}

	t, err := tail.TailFile(*filename, tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}

	w := zerolog.ConsoleWriter{Out: os.Stderr}
	for line := range t.Lines {
		w.Write([]byte(line.Text))

		/*
			var fields map[string]interface{}
			if err := json.Unmarshal([]byte(line.Text), &fields); err != nil {
				log.Error().Err(err).Msgf("failed to parse line: %q", line.Text)
				continue
			}

			log.Log().Fields(fields).Msg("")
		*/
	}
}
