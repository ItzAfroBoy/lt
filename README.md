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
</pre>
</div>

## Installation

### Install with go

```shell
go install github.com/ItzAfroBoy/lt
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

`Usage: lt [--raw] [--save] [--load] [--album] [--skip] [--artist ARTIST] [--title TITLE]`  

- `--raw`: Prints undecorated output to terminal. Not functional with `--album`  
- `--album`: Fetches lyrics for entire album  
- `--skip`: Skips in-program input. Use when `--artist` and `--title` are passed  
- `--save`: Save the lyrics for offline use.
- `--load`: Load saved lyrics

Powered by [Genius](https://genius.com)
