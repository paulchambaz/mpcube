# mpcube

mpcube is a lightweight, terminal-based client for [MPD (Music Player Daemon)](https://github.com/MusicPlayerDaemon/MPD), focused on album-centric playback. Navigate and play your music collection with a simple, vim-inspired interface. Inspired by [musikcube](https://github.com/clangen/musikcube).

![](./demo.gif)

## Installation

### Nix

```sh
nix shell github:paulchambaz/mpcube
```

Once you're satisfied, you may add it to your `configuration.nix`.

### Go

```sh
go install github.com/paulchambaz/mpcube@latest
```

## Usage

You will need a running MPD server. Make sure to follow [MPD's documentation](https://www.musicpd.org/) first.

```
Usage: mpcube [OPTIONS]

Options:
  --mpd-host          MPD host address (default: 127.0.0.1)
  --mpd-port          MPD port number (default: 6600)
  --volume-step       volume adjustment step (default: 10)
  --seek-duration     seek duration in milliseconds (default: 5000)
  --tick-interval     UI refresh interval in milliseconds (default: 100)
  --max-retry-delay   max MPD reconnect delay in seconds (default: 30)
  --scroll-padding    scroll padding lines (default: 5)
  --wide-threshold    terminal width for wide layout (default: 100)
  --album-width       album panel width in wide layout (default: 40)
  --volume-bar-threshold  terminal width for wide volume bar (default: 90)
  --volume-bar-width  volume bar width in wide layout (default: 30)
  --version           print version
  -h, --help          print help
```

All options can also be set in `~/.config/mpcube/config.toml`. See the man page for full details.

## Building

### Nix

This project uses [nix](https://github.com/NixOS/nix) for development.

```sh
git clone https://github.com/paulchambaz/mpcube.git
cd mpcube
nix develop
nix build       # build the project
nix shell       # enter a shell with mpcube installed
just --list     # list dev commands
```

### Manual

You will need Go and optionally:

- `scdoc` — to compile the man page
- `just` — to use the dev commands
- `vhs` — to produce the demo GIF

```sh
just run
just build
just fmt
just test
```

## Contribution

Contributions to mpcube are welcome. Whether it's feature suggestions, bug reports, or code contributions, feel free to open issues or submit pull requests.

## License

mpcube is released under the GPLv3 License. For more details, refer to the LICENSE file included in the repository.
