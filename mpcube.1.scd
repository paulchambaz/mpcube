MPCUBE(1) "github.com/paulchambaz/mpcube" "General Commands Manual"

# NAME

*mpcube* - Simple album focused mpd terminal client

# SYNOPSIS

*mpcube* [OPTIONS]

# DESCRIPTION

*mpcube* is a lightweight, terminal-based client for the *Music Player Daemon
(MPD)*, designed to provide an efficient and focused music listening experience.
Emphasizing album-centric playback, *mpcube* allows users to navigate and  play
their music collection with a simple and intuitive interface. Inspired by
*musikube*, it aims to cater to users who prefer structured album listening
sessions over shuffled tracks or playlists.

# OPTIONS

*mpd-host <MPD_HOST>*  Ip address of the mpd host [default: 127.0.0.1]

*mpd-port <MPD_PORT>*  Port number of the mpd host [default: 6600]

*--cache <CACHE>*      Cache file location [default: ~/.cache/mpcube/cache]

*-h*, *--help*         Print help

*-V*, *--version*      Print version

# MUSIC PLAYER DAEMON

*mpcube* is just a client for *mpd*, so before we can use *mpcube*, we will have
to configure *mpd*.

For *Arch* based system: install using *sudo pacman -S mpd*.
For *Debian* based system: install using *sudo apt install mpd*.
For *Fedora* based system: install using *sudo dnf install mpd*.
For *NixOS* based system: add the following to your *configuration.nix*:

```
services.mpd = {
	enable = true;
	user = "your-user";
	extraConfig = ''
		# the content of your configuration
	'';
};
```

Now that the program is installed, it is important to configure it. I recommend
having the simplest configuration possible for mpd, here's mine :

```
music_directory "/home/paul/music"
auto_update "yes"
port "6600"
audio_output {
	type "pipewire"
	name "pipewire"
}
```

The only think you will probably change is the *audio_output*, use *also*,
*pulse* or *pipewire* depending on the audio server you have installed on your
system.

Finally, before we can use *mpd*, we must start it : *sudo systemctl enable --now mpd*.

This is just the very start of *mpd* configuration, please read the *Arch
wiki*'s page on mpd, or *NixOS wiki*'s page on mpd for more information.

# CONFIGURATION

Configuration of *mpcube* is handles through the *~/.config/mpcube/config.toml*
file. This file allows you to customize the basic settings of mpcube.

It contains three optional fields, feel free to leave any of them empty if your
mpd configuration does not require any change:

```
mpd_host = "127.0.0.1"
mpd_port = 6600
chahe = "~/.cache/mpcube/cache"
```

*mpd_host* is the ip address of your mpd instance, if you have installed your
*mpd* server on the same computer than *mpcube*, then you do not need to change
this setting. *mpd_port* is the port of your mpd instance, if you have followed
this guide or did not change the *mpd* default port, then you do not need to
change this setting. Finally *cache* is the path to the cache file, a cache file
is kept so that startup times are faster. For an estimation, on my system,
parsing my library of ~1600 files takes ~.2 seconds, with the cache, it is
instant.

# ALBUM AND TITLE WINDOW

The main interface of *mpcube* is divided into two primary windows : the
*Album* and *Title* windows.

The *Album* window displays the list of albums from your music library. You can
use *j* and *k*, the vim bindings to move up and down the window. You can also
use *Up* and *Down*, the arrow keys for the same result. You can also press
*Enter* to start an entire album, your *mpd* queue will be cleared and the album
you have selected will take precedence. If you want to play a specific track of
the album, then you will need to navigate to the *Title* window by pressing *l*
or *Right*.

The *Title* window displays the list of titles from the album you are currently
selecting. Just like with the *Album* window, you can use the same *j*, *k*,
*Up* and *Down* keys to navigate. You can also press *Enter* to start a given
track, when you do, the *mpd* queue will be cleared and the album you are
currently selecting will take precedence and the track you have selected will
start playing. If you want to escape with window and go back to the *Album*
window, you can do so by pressing *h* or *Left*.

# INFORMATION WINDOWS

Aside from the *Album* and *Title* window, *mpcube* provides additional windows
for detailed information and controls. Please understand that you can not
navigate inside these window, they are only a mirror of what is happening with
the mpd server. You may see them as a remote control in a sense. The
*Information* window will display the currenly playing song, and if you are
playing or paused. The *Status* window will display whether or not you are in
*shuffle* mode - in which the songs in your queue go in a random order and whether
or not you are in *repeat* mode - in which your entire queue will start again
after it has ended. The *Volume* window displays the current volume of the *mpd*
playback. Finally, the *Bar* window will display the current position in the
song you are playing.

# KEYBOARD CONTROLS

*mpcube* is designed for efficient keyboard navigation. Here are the key
bindings:

- *j*, *Down* - Move down the list in the *Album* or *Title* window
- *k*, *Up* - Move up the list in the *Album* or *Title* window
- *l*, *Right* - Move to the *Title* window from the *Album* window
- *h*, *Left* - Move to the *Album* window from the *Title* window
- *Enter* - Start an album in the *Album* window, a song in the *Title* window
- *Space* - Toggle play / pause
- *n* - Play the next track
- *p* - Play the previous track
- *=* - Increase volume by 10 percent
- *-* - Decrease volume by 10 percent
- *x* - Clears the *mpd* queue
- *.* - Seek forward by 5 seconds
- *,* - Seek backward by 5 seconds
- *s* - Toggle *shuffle* mode
- *r* - Toggle *repeat* mode
- *U* - Full scan of the *mpd* database
- *R* - Play one album at random
- *q*, *Q* - Quits *mpcube* and clears the *mpd* queue

# AUTHORS

This program was written by Paul Chambaz in 2024. It is licensed under the
GPLv3, for more information, please read the *LICENSE* file.
