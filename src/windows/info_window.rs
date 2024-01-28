//! This module is used to manage the Information window

use ratatui::{
    layout::Rect,
    style::{Color, Style},
    widgets::Paragraph,
    Frame,
};

use crate::music::music_data::{MusicData, StateData};

/// Used to display the Information window
pub struct InfoWindow {
    /// Whether or not a song is currently playing
    playing: bool,
    /// The title of the song playing (if any)
    title: Option<String>,
    /// The artist of the song playing (if any)
    artist: Option<String>,
    /// The album of the song playing (if any)
    album: Option<String>,
    /// The area reserved for drawing the window
    area: Rect,
}

impl InfoWindow {
    /// Creates a new window
    pub fn new() -> Self {
        InfoWindow {
            playing: false,
            title: None,
            artist: None,
            album: None,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, _: bool, music_data: &MusicData, state_data: &StateData) {
        // TODO: In the event of a clear, we should default back to None for the album, artist and title

        // Update if we are playing
        self.playing = state_data.playing;
        // If a song is playing
        if let (Some(album_id), Some(title_id)) = (state_data.album_id, state_data.title_id) {
            // We get the album playing
            let album = music_data
                .albums
                .get(album_id)
                .expect("Could not get album value from its id");

            // And use it to update the album, artist and title names
            self.album = Some(album.album.clone());
            self.artist = Some(album.artist.clone());
            self.title = Some(
                album
                    .songs
                    .get(title_id)
                    .expect("Could not get title value from its id")
                    .title
                    .clone(),
            );
        }
    }

    /// Updates the area reserved for the window, used to resize correctly
    ///
    /// - `x`: The x coordinate of the window
    /// - `y`: The y coordinate of the window
    /// - `width`: The width of the window
    /// - `height`: The height of the window
    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    /// Render the window to the ratatui frame
    ///
    /// - `frame`: The ratatui frame to render
    pub fn render(&mut self, frame: &mut Frame) {
        // Get the area reserved for drawing
        let area = self.area;

        // Used to display the text in correct format
        let mut render_widget = |text: &str, highlight: bool, x: u16, w: u16| {
            frame.render_widget(
                Paragraph::new(text).style(Style::default().fg(if highlight {
                    Color::Green
                } else {
                    Color::DarkGray
                })),
                Rect::new(area.x + x, area.y, w, 1),
            );
        };

        // If a song is playing
        if let (Some(title), Some(artist), Some(album)) = (&self.title, &self.artist, &self.album) {
            // Create the list of elements to be displayed
            let str: [&str; 6] = [
                if self.playing { "Playing" } else { "Paused" },
                title,
                "by",
                artist,
                "from",
                album,
            ];
            let len: Vec<u16> = str.iter().map(|s| s.len() as u16).collect();

            // Render the list of elements to be displayed in correct format and position
            let mut sum = 1;
            for (i, (&text, &length)) in str.iter().zip(len.iter()).enumerate() {
                if sum > area.width {
                    return;
                }
                render_widget(text, i % 2 != 0, sum, u16::min(length, area.width - sum));
                sum += length + 1;
            }
        } else {
            // If no song is playing, render default text
            render_widget("Not playing", false, 1, 11);
        }
    }
}
