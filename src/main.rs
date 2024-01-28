//! *mpcube* - Simple album focused mpd terminal client
//!
//! *mpcube* is a lightweight, terminal-based client for the *Music Player Daemon
//! (MPD)*, designed to provide an efficient and focused music listening experience.
//! Emphasizing album-centric playback, *mpcube* allows users to navigate and  play
//! their music collection with a simple and intuitive interface. Inspired by
//! *musikube*, it aims to cater to users who prefer structured album listening
//! sessions over shuffled tracks or playlists.

#![warn(missing_docs)]
#![warn(clippy::missing_docs_in_private_items)]

extern crate mpd;

mod config;
mod input;
mod interface;
mod music;
mod ui;
mod windows;

use config::load_config;
use interface::Interface;
use music::music_client::Client;
use tokio::time::error::Error;

/// Entry point of the program
#[tokio::main]
async fn main() -> Result<(), Error> {
    let config = load_config();

    let client = Client::new(config.mpd_host, config.mpd_port, config.cache);
    Interface::new().render(client).await;

    Ok(())
}
