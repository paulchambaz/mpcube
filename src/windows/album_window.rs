//! This module is used to manage the Album window

use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::music::{
    music_client::Client,
    music_data::{MusicData, StateData},
};
use std::sync::Arc;
use tokio::sync::Mutex;

/// Used to display the Album window
pub struct AlbumWindow {
    /// Whether or not the Album window is selected
    selected: bool,
    /// The offset in the scroll list
    offset: usize,
    /// The id of the album playing (if any)
    album_playing: Option<usize>,
    /// The id of the album selected
    pub album_selected: usize,
    /// The list of album names used to display
    album_names: Vec<String>,
    /// The area reserved for drawing the window
    area: Rect,
}

impl AlbumWindow {
    /// Creates a new window
    pub fn new() -> Self {
        AlbumWindow {
            selected: true,
            offset: 0,
            album_playing: None,
            album_selected: 0,
            album_names: vec![],
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, on_album: bool, music_data: &MusicData, state_data: &StateData) {
        // Update if we are on the Album window
        self.selected = on_album;
        // Reconstructs the list of albums from scratch
        self.album_names.clear();
        for album in &music_data.albums {
            self.album_names.push(album.album.clone());
        }
        // Update the current album playing
        self.album_playing = state_data.album_id;
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

        // Simplifies drawing the individual list members
        let mut render_widget = |text: &str, style: Style, y: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(
                    area.x + 1,
                    area.y + 1 + y - self.offset as u16,
                    area.width - 2,
                    1,
                ),
            );
        };

        // For all the albums
        for (i, album) in self
            .album_names
            .iter()
            .enumerate()
            .skip(self.offset)
            .take(area.height as usize - 2)
        {
            // Is is the playing album or the selected album
            let playing_album = self.album_playing.map_or(false, |playing| playing == i);
            let selected_album = self.album_selected == i;

            // Match style
            let style = match (playing_album, selected_album) {
                (true, true) => Style::default().fg(Color::Black).bg(Color::Cyan),
                (true, false) => Style::default().fg(Color::Black).bg(Color::Green),
                (false, true) => Style::default().fg(Color::Black).bg(Color::LightBlue),
                (false, false) => Style::default().fg(Color::DarkGray),
            };

            // Render list item
            render_widget(album, style, i as u16);
        }

        // Renders the border
        let border_style = if self.selected {
            Style::default()
                .fg(Color::LightRed)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };

        frame.render_widget(
            Paragraph::new("").block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_style(border_style),
            ),
            area,
        );

        // Render border title
        let text = "Album";
        let len = text.len() as u16;
        frame.render_widget(
            Paragraph::new(format!(" {} ", text)).style(border_style),
            Rect::new(area.x + (area.width - len - 2) / 2, area.y, len + 2, 1),
        );
    }

    /// Go down the list of albums
    pub fn down(&mut self) {
        if self.album_names.is_empty() {
            return;
        }

        // Update the album selected
        if self.album_selected < self.album_names.len() - 1 {
            self.album_selected += 1;
        }

        // Update the offset given a position and a padding
        let padding = match self.area.height as usize - 2 {
            0..=3 => 0,
            4..=7 => 1,
            8..=11 => 2,
            12..=15 => 3,
            16..=19 => 4,
            20..=usize::MAX => 5,
            _ => 0,
        };

        if self.album_selected > self.offset + self.area.height as usize - 3 - padding
            && self.offset < self.album_names.len() - self.area.height as usize + 2
        {
            self.offset += 1;
        }
    }

    /// Go up the list of albums
    pub fn up(&mut self) {
        if self.album_names.is_empty() {
            return;
        }

        // Update the album selected
        if self.album_selected > 0 {
            self.album_selected -= 1;
        }

        // Update the offset given a position and a padding
        let padding = match self.area.height as usize - 2 {
            0..=3 => 0,
            4..=7 => 1,
            8..=11 => 2,
            12..=15 => 3,
            16..=19 => 4,
            20..=usize::MAX => 5,
            _ => 0,
        };

        if self.album_selected < self.offset + padding && self.offset > 0 {
            self.offset -= 1;
        }
    }

    /// Plays a given album
    ///
    /// - `client`: The client protected for thread safety
    pub fn play(&mut self, client: &mut Arc<Mutex<Client>>) {
        let album_selected = self.album_selected;
        // Start a background thread for the operation
        let client_lock = client.clone();
        tokio::spawn(async move {
            let mut client = client_lock.lock().await;
            client.start_album(album_selected);
        });
    }
}
