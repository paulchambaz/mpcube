extern crate mpd;
mod mpd_client;
mod ui;
mod windows;

use mpd_client::Client;
use ui::Interface;

const CACHE_FILE_PATH: &str = "mpcube.bin";

fn main() {
    // connect to the mpd instance
    let mut client = Client::new("127.0.0.1", 6600);
    // only do in the absence of the cache
    client.init_sync(CACHE_FILE_PATH);

    let mut interface = Interface::new(client);

    interface.render();
}
