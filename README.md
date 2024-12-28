<div align="center">
<pre>
 __                                 ______                             
/\ \                     __        /\__  _\                            
\ \ \      __  __  _ __ /\_\    ___\/_/\ \/    __   _ __    ___ ___    
 \ \ \  __/\ \/\ \/\`'__\/\ \  /'___\ \ \ \  /'__`\/\`'__\/' __` __`\  
  \ \ \L\ \ \ \_\ \ \ \/ \ \ \/\ \__/  \ \ \/\  __/\ \ \/ /\ \/\ \/\ \ 
   \ \____/\/`____ \ \_\  \ \_\ \____\  \ \_\ \____\\ \_\ \ \_\ \_\ \_\
    \/___/  `/___/> \/_/   \/_/\/____/   \/_/\/____/ \/_/  \/_/\/_/\/_/
               /\___/                                                  
               \/__/                                                   
<br>
Lyrics in the terminal
<br>
<img alt="GitHub License" src="https://img.shields.io/github/license/ItzAfroBoy/lt"> <img alt="GitHub tag (with filter)" src="https://img.shields.io/github/v/tag/ItzAfroBoy/lt?label=version"> <a href="https://www.codefactor.io/repository/github/itzafroboy/lt"><img src="https://www.codefactor.io/repository/github/itzafroboy/lt/badge" alt="CodeFactor" /></a> <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/ItzAfroBoy/lt">
</pre>
</div>

## Installation

### Install with go

```shell
go install github.com/ItzAfroBoy/lt@latest
lt ...
```

### Build from source

```shell
git clone https://github.com/ItzAfroBoy/lt
cd lt
go install
lt ...
```

## Usage

`Usage: lt [--raw] [--export | import] [--album] [--artist ARTIST] [--title TITLE]`  

- `--raw`: Prints undecorated output to terminal. Not functional with `--album`  
- `--album`: Fetches lyrics for entire album  
- `--export`: Save the lyrics for offline use. Not functional with `--album`  
- `--import`: Load saved lyrics

Powered by:

- [Genius](https://genius.com)
- [Bubbletea](https://github.com/charmbracelet/bubbletea/)
