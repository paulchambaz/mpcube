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

*--mpd-host* <HOST>
	MPD host address [default: 127.0.0.1]

*--mpd-port* <PORT>
	MPD port number [default: 6600]

*--volume-step* <STEP>
	Volume adjustment step [default: 10]

*--seek-duration* <MS>
	Seek duration in milliseconds [default: 5000]

*--tick-interval* <MS>
	UI refresh interval in milliseconds [default: 100]

*--max-retry-delay* <SECONDS>
	Max MPD reconnect delay in seconds [default: 30]

*--scroll-padding* <LINES>
	Scroll padding lines [default: 5]

*--wide-threshold* <COLS>
	Terminal width for wide layout [default: 100]

*--album-width* <COLS>
	Album panel width in wide layout [default: 40]

*--volume-bar-threshold* <COLS>
	Terminal width for wide volume bar [default: 90]

*--volume-bar-width* <COLS>
	Volume bar width in wide layout [default: 30]

*-h*, *--help*
	Print help

*--version*
	Print version

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

Configuration of *mpcube* is handled through the *~/.config/mpcube/config.toml*
file. All fields are optional — defaults are used for any missing value.

Settings are loaded in this order, with later sources overriding earlier ones:
defaults, config file, environment variables, CLI arguments.

```
mpd_host = "127.0.0.1"
mpd_port = 6600
volume_step = 10
seek_duration = 5000
tick_interval = 100
max_retry_delay = 30
scroll_padding = 5
wide_threshold = 100
album_width = 40
volume_bar_threshold = 90
volume_bar_width = 30
```

*mpd_host* is the IP address of your MPD instance. If you have installed your
MPD server on the same computer as *mpcube*, you do not need to change this
setting. [default: 127.0.0.1]

*mpd_port* is the port of your MPD instance. If you have followed this guide or
did not change the MPD default port, you do not need to change this setting.
[default: 6600]

*volume_step* is the percentage by which the volume increases or decreases each
time you press the volume up or volume down key. For example, with the default
value of 10, each keypress adjusts the volume by 10%. [default: 10]

*seek_duration* is the amount of time in milliseconds that the playback position
moves when you press the seek forward or seek backward key. The default value
of 5000 means each keypress seeks by 5 seconds. [default: 5000]

*tick_interval* is the interval in milliseconds at which *mpcube* polls the MPD
server for state updates such as track position, volume changes, and playback
status. Lower values make the UI more responsive but use slightly more
resources. [default: 100]

*max_retry_delay* is the maximum number of seconds *mpcube* will wait between
reconnection attempts when the MPD server is unreachable. Retries use
exponential backoff up to this limit. [default: 30]

*scroll_padding* controls how many lines of context remain visible above or
below the cursor when scrolling through the album or track list. A higher value
keeps the cursor further from the edge of the panel. On small terminals, the
padding is reduced proportionally. [default: 5]

*wide_threshold* is the terminal width in columns above which *mpcube* switches
to a wide layout with a fixed-width album panel. Below this threshold, the
album panel takes up a proportional fraction of the terminal width instead.
[default: 100]

*album_width* is the width in columns of the album panel when the terminal is
wider than *wide_threshold*. This has no effect when the terminal is narrower
than the threshold. [default: 40]

*volume_bar_threshold* is the terminal width in columns above which the volume
bar uses a fixed width. Below this threshold, the volume bar takes up a
proportional fraction of the terminal width instead. [default: 90]

*volume_bar_width* is the width in columns of the volume bar when the terminal
is wider than *volume_bar_threshold*. This has no effect when the terminal is
narrower than the threshold. [default: 30]

Environment variables can also be used to override settings. *MPD_HOST* and
*MPD_PORT* set the MPD connection. All other settings use the *MPCUBE_* prefix
with the setting name in uppercase (e.g. *MPCUBE_VOLUME_STEP*,
*MPCUBE_SEEK_DURATION*).

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

# MODES

*mpcube* is a modal application. Different modes change how keyboard input is
interpreted. Some keys are available in all modes, while others are specific to
a given mode.

The default mode is *normal* mode, which is used for browsing and playback
control. Other modes are entered and exited through specific key bindings.

# SEARCH

Pressing */* in normal mode enters *search* mode. The information bar is
replaced by a search prompt. As you type, *mpcube* performs a fuzzy match
against your library and jumps to the best match. The search matches against
both *artist album* and *album artist* combinations, so you can type in either
order. Results are ranked by match quality — consecutive character matches and
matches at the start of words score higher. You do not need to type special
characters like hyphens or periods, as the fuzzy matching will skip over them.
Press *Enter* to confirm your search and enter *searching* mode, or press
*Escape* to cancel and return to where you were.

In *searching* mode, you can cycle through all matching albums with *n* and *p*.
The information bar shows your position in the match list. Press *Enter* to play
the selected album, *Escape* to cancel and return to your previous position, or
*/* to edit your query.

# KEYBOARD CONTROLS

*mpcube* is designed for efficient keyboard navigation.

## Normal mode

- *j*, *Down* - Move down the list in the *Album* or *Title* window
- *k*, *Up* - Move up the list in the *Album* or *Title* window
- *l*, *Right* - Move to the *Title* window from the *Album* window
- *h*, *Left* - Move to the *Album* window from the *Title* window
- *g*, *Home* - Go to top of list
- *G*, *End* - Go to bottom of list
- *Enter* - Start an album in the *Album* window, a song in the *Title* window
- *Space* - Toggle play / pause
- *n* - Play the next track
- *p* - Play the previous track
- *x* - Clears the *mpd* queue
- *s* - Toggle *shuffle* mode
- *r* - Toggle *repeat* mode
- *U* - Full scan of the *mpd* database
- *R* - Play one album at random
- */* - Enter search mode
- *q*, *Q* - Quits *mpcube* and clears the *mpd* queue

## Search mode

- Typing - Appends to the search query, jumps to best match
- *Backspace* - Deletes the last character from the query
- *Enter* - Confirms the search and enters searching mode
- *Escape* - Cancels the search and restores your previous position

## Searching mode

- *n* - Jump to the next match
- *p* - Jump to the previous match
- *Enter* - Play the selected album and return to normal mode
- *Escape* - Cancel and restore your previous position
- */* , *i* - Edit the search query
- *Space* - Toggle play / pause
- *q*, *Q* - Quits *mpcube* and clears the *mpd* queue

## All modes

- *.* - Seek forward by 5 seconds
- *,* - Seek backward by 5 seconds
- *=* - Increase volume by 10 percent
- *-* - Decrease volume by 10 percent
- *Ctrl+c* - Quits *mpcube* and clears the *mpd* queue

# AUTHORS

This program was written by Paul Chambaz in 2024. It is licensed under the
GPLv3, for more information, please read the *LICENSE* file.
