extern crate mpd;

mod input;
mod interface;
mod music;
mod ui;
mod windows;

use interface::Interface;
use music::music_client::Client;

fn main() {
    let client = Client::new("127.0.0.1", 6600, "mpcube.bin".to_string());
    Interface::new().render(client);
}
