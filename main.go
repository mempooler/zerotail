package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/mempooler/zerolog"
)

func main() {
	f := flag.String("f", "", "file to tail")
	n := flag.Int("n", 10, "number of lines to tail")
	debug := flag.Bool("debug", true, "")
	trace := flag.Bool("trace", true, "")
	flag.Parse()

	if *f == "" {
		fmt.Println("usage: tail -f <filename>")
		os.Exit(1)
	}

	t, err := tail.TailFile(*f, tail.Config{
		Location: &tail.SeekInfo{
			Offset: -(*n),
			Whence: os.SEEK_END,
		},
		Follow: true,
	})
	if err != nil {
		panic(err)
	}

	w := zerolog.NewFilteredWriter(
		zerolog.MultiLevelWriter(ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}),
		getLevel(*trace, *debug),
	)
	for line := range t.Lines {
		w.Write([]byte(line.Text))
	}
}

func getLevel(trace, debug bool) zerolog.Level {
	if trace {
		return zerolog.TraceLevel
	} else if debug {
		return zerolog.DebugLevel
	} else {
		return zerolog.InfoLevel
	}
}
