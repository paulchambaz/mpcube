extern crate mpd;
mod mpd_client;
mod ui;
mod windows;

use mpd_client::Client;
use ui::Interface;

fn main() {
    // connect to the mpd instance
    let mut client = Client::new("127.0.0.1", 6600);
    client.full_sync();

    let mut interface = Interface::new(client);

    interface.render();
}
