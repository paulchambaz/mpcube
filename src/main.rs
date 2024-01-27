extern crate mpd;

mod input;
mod interface;
mod music;
mod ui;
mod windows;

use interface::Interface;
use music::music_client::Client;
use tokio::time::error::Error;

#[tokio::main]
async fn main() -> Result<(), Error> {
    // TODO: this should come from the settings or default to the user ~/.cache/mpcube/cache.bin
    let client = Client::new("127.0.0.1", 6600, "mpcube.bin".to_string());
    Interface::new().render(client).await;

    Ok(())
}
