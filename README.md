# Sentiment analyzer

Simple sentiment analyzer using
[BiDiSentiment](https://github.com/vmarkovtsev/BiDiSentiment) that
offers a slightly more elaborate interface than the CLI application
that comes with the library.  Mostly made this for my own
entertainment.  I'm not sure it is that useful.

## Build

Make sure you have `libtensorflow` installed since the application
depends on this library. 

Make sure you set `GOOS` and `GOARCH` when building this.  You can set
them and build like this:

    GOOS=linux GOARCH=amd64 make

It is set up to build on macOS per default, so if you are on a Mac
just type:

    make
	

## Command line options

    $ bin/sent -h
    Usage:
      sent [OPTIONS]

    Application Options:
      -v, --verbose             Verbose mode, show individual score for lines of file
      -n, --negative-threshold= Threshold for when something is deemed negative (default: 0.600)
      -p, --positive-threshold= Threshold for when something is deemed positive (default: 0.400)

    Help Options:
      -h, --help                Show this help message
