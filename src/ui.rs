//! This module manages the user interface. It updates each subwindow of the
//! program, resizes it for the terminal area and renders them

use crate::{
    music::music_client::Client,
    windows::{
        album_window::AlbumWindow, bar_window::BarWindow, info_window::InfoWindow,
        status_window::StatusWindow, title_window::TitleWindow, volume_window::VolumeWindow,
    },
};
use ratatui::Frame;
use std::sync::Arc;
use tokio::sync::Mutex;

/// Represents the user interface
pub struct Ui {
    /// Is used for safe connection with the mpd client
    pub client: Arc<Mutex<Client>>,
    /// Tracks on what window the user is
    pub on_album: bool,
    /// The window to display the albums
    pub album_window: AlbumWindow,
    /// The window to display the songs
    pub title_window: TitleWindow,
    /// The window to display information about which song is playing
    pub info_window: InfoWindow,
    /// The window to display volume information
    pub volume_window: VolumeWindow,
    /// The window to display the currnet position in the track playing
    pub bar_window: BarWindow,
    /// The window to display whether shuffle mode and repeat mode are selected
    pub status_window: StatusWindow,
}

impl Ui {
    /// Initialises the user interface
    ///
    /// - `client`: The client used for the mpd connection
    pub fn new(client: Client) -> Self {
        Ui {
            client: Arc::new(Mutex::new(client)),
            on_album: true,
            album_window: AlbumWindow::new(),
            title_window: TitleWindow::new(),
            info_window: InfoWindow::new(),
            volume_window: VolumeWindow::new(),
            bar_window: BarWindow::new(),
            status_window: StatusWindow::new(),
        }
    }

    /// Renders the entire user interface
    /// First we safely update the user interface
    /// Secondly we resize the user interface
    /// Thirdly we render the user interface
    ///
    /// - `frame`: The ratatui frame on which we can render
    pub fn render(&mut self, frame: &mut Frame) {
        let a = frame.size();
        let on_album = self.on_album;

        if let Ok(client) = self.client.try_lock() {
            let data = &client.data;
            let state = &client.state;
            self.album_window.update(on_album, data, state);
            self.title_window.update(on_album, data, state);
            self.info_window.update(on_album, data, state);
            self.volume_window.update(on_album, data, state);
            self.bar_window.update(on_album, data, state);
            self.status_window.update(on_album, data, state);
        }

        // resizing the windows
        let side_width = if a.width > 100 { 40 } else { 2 * a.width / 5 };
        let volume_width = if a.width > 90 { 30 } else { a.width / 3 };

        self.album_window
            .update_area(0, 0, side_width, a.height - 2);
        self.title_window
            .update_area(side_width, 0, a.width - side_width, a.height - 2);
        self.info_window
            .update_area(0, a.height - 2, a.width - 9, 1);
        self.volume_window
            .update_area(0, a.height - 1, volume_width, 1);
        self.bar_window
            .update_area(volume_width, a.height - 1, a.width - volume_width - 9, 1);
        self.status_window
            .update_area(a.width - 9, a.height - 2, 8, 2);

        // rendering of windows
        self.album_window.render(frame);
        self.title_window.render(frame);
        self.info_window.render(frame);
        self.volume_window.render(frame);
        self.bar_window.render(frame);
        self.status_window.render(frame);
    }
}
