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

#[tokio::main]
async fn main() -> Result<(), Error> {
    let config = load_config();

    let client = Client::new(config.mpd_host, config.mpd_port, config.cache);
    Interface::new().render(client).await;

    Ok(())
}
