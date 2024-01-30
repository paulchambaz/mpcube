# mpcube

mpcube is a straightforward music player client for [MPD (Music Player Daemon)](https://github.com/MusicPlayerDaemon/MPD), drawing inspiration from [musikube](https://github.com/clangen/musikcube). It offers an intuitive, terminal-based interface, allowing users to efficiently manage and play their music.

![](./demo.gif)

mpcube is a lightweight, terminal-based client for the [Music Player Daemon (MPD)](https://github.com/MusicPlayerDaemon/MPD), designed to provide an efficient and focused music listening experience. Emphasizing album-centric playback, mpcube allows users to navigate and  play their music collection with a simple and intuitive interface. Inspired by *musikube*, it aims to cater to users who prefer structured album listening sessions over shuffled tracks or playlists.

## Installation

Currently, the only way to install this project is manually, however, in the close future, I intend to publish it to crates.io, nixpkgs and the arch user repository.

### Manual

To install the project manually, please consult the [**Building**](#Building) section.

### Nix

**Coming soon.** You can try out the program with `nix-shell` :

```sh
nix-shell -p mpcube
```

It's a good way to ensure it works as intended. Once you're satisfied, you may add it to your `configuration.nix`.

### Cargo

**Coming soon.** `mpcube` is hosted on [crates.io](https://crates.io/crates/mpcube). To install it, simply :

```sh
cargo install mpcube
```

### AUR

**Coming soon.** `mpcube` is also hosted on the [Arch Linux User Repository](https://aur.archlinux.org/packages/mpcube). To install it, simply :

```sh
yay -S mpcube
```

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
      --cache <CACHE>        Cache file location [default: ~/.cache/mpcube/cache]
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
