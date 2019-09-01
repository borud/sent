package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/logrusorgru/aurora"
	sentiment "github.com/vmarkovtsev/BiDiSentiment"
)

type progOpts struct {
	Verbose           bool    `short:"v" long:"verbose" description:"Verbose mode, show individual score for lines of file"`
	NegativeThreshold float32 `short:"n" long:"negative-threshold" description:"Threshold for when something is deemed negative" default:"0.600"`
	PositiveThreshold float32 `short:"p" long:"positive-threshold" description:"Threshold for when something is deemed positive" default:"0.400"`
}

var opts progOpts

func main() {
	args, err := flags.NewParser(&opts, flags.Default).Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Error parsing flags: %v", err)
	}

	if opts.NegativeThreshold > 1.00 || opts.NegativeThreshold < 0.0 {
		log.Fatal("Negative threshold must be between 0.0 and 1.0")
	}

	if opts.PositiveThreshold > 1.00 || opts.PositiveThreshold < 0.0 {
		log.Fatal("PositiveThreshold must be between 0.0 and 1.0")
	}

	if (opts.NegativeThreshold + opts.PositiveThreshold) > 1.0 {
		log.Fatal("PositiveThreshold + NegativeThreshold cannot be greater than 1.0")
	}

	batchSize := sentiment.GetBatchSize()

	// Print out a bit of info if in verbose mode
	if opts.Verbose {
		log.Printf("       Batch size: %d", batchSize)
		log.Printf("NegativeThreshold: %4.3f", opts.NegativeThreshold)
		log.Printf("PositiveThreshold: %4.3f", opts.PositiveThreshold)
	}

	session, _ := sentiment.OpenSession()
	defer session.Close()

	var lines []string

	eval := func(fn string) float32 {
		score, err := sentiment.Evaluate(lines, session)
		if err != nil {
			log.Printf("Error running sentiment analysis on %s: %v ", fn, err)
			return 0.5
		}

		var sum float32
		for n, line := range lines {
			// Skip empty lines
			if line == "" {
				continue
			}

			// Skip NaN
			if math.IsNaN(float64(score[n])) {
				continue
			}

			sum += score[n]

			if opts.Verbose {
				switch {
				case score[n] >= opts.NegativeThreshold:
					fmt.Printf("%v %4.3f | %s\n", score[n], aurora.BrightRed(score[n]), line)

				case score[n] <= opts.PositiveThreshold:
					fmt.Printf("%v %4.3f | %s\n", score[n], aurora.BrightBlue(score[n]), line)

				default:
					fmt.Printf("%v %4.3f | %s\n", score[n], aurora.BrightYellow(score[n]), line)
				}
			}
		}

		lines = lines[:0]

		return sum
	}

	// Process files
	for _, filename := range args {
		var scoreSum float32
		numLines := 0

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Unable to open %s: %v", filename, err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			numLines++

			if len(lines) >= batchSize {
				scoreSum += eval(filename)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error scanning %s: %v ", filename, err)
		}
		file.Close()

		// deal with partial batch
		if len(lines) > 0 {
			scoreSum += eval(filename)
		}

		average := scoreSum / float32(numLines)
		switch {
		case average >= opts.NegativeThreshold:
			log.Printf("%s = %4.3f", filename, aurora.BrightRed(scoreSum/float32(numLines)))

		case average <= opts.PositiveThreshold:
			log.Printf("%s = %4.3f", filename, aurora.BrightBlue(scoreSum/float32(numLines)))

		default:
			log.Printf("%s = %4.3f", filename, aurora.BrightYellow(scoreSum/float32(numLines)))
		}

	}

}
