# LyricTerm

Lyrics in the terminal

## Installation

Install with go:  
`go install github.com/ItzAfroBoy/lt`  

Make from source:  

```shell
git clone https://github.com/ItzAfroBoy/lt
cd lt
go build
./lt
```

## Usage

`$ lt [--raw] [--album] [--skip] [--artist ARTIST] [--title TITLE]`  

Flags:

- `--raw`: Prints undecorated output to terminal. Not functional with `--album`  
- `--album`: Fetches lyrics for entire album  
- `--skip`: Skips in-program input. Use when `--artist` and `--title` are passed  

Powered by [Genius](https://genius.com)
