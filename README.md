# mpcube

mpcube is a straightforward music player client for [MPD (Music Player Daemon)](https://github.com/MusicPlayerDaemon/MPD), drawing inspiration from [musikube](https://github.com/clangen/musikcube). It offers an intuitive, terminal-based interface, allowing users to efficiently manage and play their music.

![](./demo.gif)

mpcube is a music player client optimized for those who prefer enjoying albums from back to back. Designed to work with [MPD](https://github.com/MusicPlayerDaemon/MPD), it offers a straightforward, text-based interface, emphasizing album playback over shuffled songs or playlists. The tool is ideal for users who have an organized collection of albums and seek a focused, album-centric listening experience. Simple yet functional, it provides an uncluttered environment to dive into your music, one album at a time. The interface is mostly derived from the great music player [musikcube](https://github.com/clangen/musikcube).

## Installation

Currently, the easiest way to install this project is with nix. In the future, it will be added to the Arch Linux user repository. The project can also be installed via Cargo.

### Nix

You can try out the program with `nix-shell` :

```sh
nix-shell -p mpcube
```

It's a good way to ensure it works as intended. Once you're satisfied, you may add it to your `configuration.nix`.

### Cargo

`mpcube` is hosted on [crates.io](https://crates.io/crates/mpcube). To install it, simply :

```sh
cargo install mpcube
```

### Manual

To install the project manually, please consult the [Building](#Building) section.

## Usage

In order to use this client, you will first need to configure the `mpd` server. Make sure to follow their instructions first.

To understand how to use this program, please consult the `man` page. All instructions are detailed there. You may also read them from the `mpcube.1.scd` file.

Here is a brief overview of the program :

```sh
Simple album focused mpd terminal client

Usage: mpcube [OPTIONS]

Options:
      --mpd-host <MPD_HOST>  Ip address of the mpd host [default: 127.0.0.1]
      --mpd-port <MPD_PORT>  Port number of the mpd host [default: 6600]
      --cache <CACHE>        Cache file location [default ~/.cache/mpcube/cache]
  -h, --help                 Print help
  -V, --version              Print version
```

You can also add these to `~/.config/mpcube/config.toml`, so that they are stored.

## Building

### Nix

This project uses [nix](https://github.com/NixOS/nix) for development. If you want to contribute, it is recommended to install nix (not NixOS) to access the development shell.

```sh
git clone https://github.com/paulchambaz/mpcube.git
cd mpcube
nix develop
nix build # to build the project
nix shell # to enter a shell where the built project is installed
just --list # to list the dev commands
```

### Manual

If you want to manually build the project without nix, it should not be too hard.

You will need to following program and libraries :

- `libmpdclient` - for runtime
- `scdoc` - to compile the man page

To develop, you will probably need the following programs :

- `just` - to have access to the dev commands
- `cargo-tarpaulin` - to run the coverage metrics
- `vhs` - to produce the gif at the top of this page

```sh
just run
just build
just fmt
just coverage
just watch-test
```

## Contribution

Contributions to mpcube are welcome. Whether it's feature suggestions, bug reports, or code contributions, feel free to share your input. Please use the project's GitHub repository to open issues or submit pull requests.

## License

mpcube is released under the GPLv3 License. For more details, refer to the LICENSE file included in the repository.
