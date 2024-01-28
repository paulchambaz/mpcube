//! This module is used to manage the Title window

use std::time::Duration;

use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::music::{
    music_client::Client,
    music_data::{MusicData, StateData},
};

/// Used to display the Title window
pub struct TitleWindow {
    /// Whether or not the Title window is selected
    selected: bool,
    /// The offset in the scroll list
    offset: usize,
    /// The id of the album playing (if any)
    album_playing: Option<usize>,
    /// The id of the album selected
    album_selected: usize,
    /// The id of the song playing (if any)
    title_playing: Option<usize>,
    /// The id of the song selected
    title_selected: usize,
    /// The list of song names used to display
    title_names: Vec<String>,
    /// The list of song duration used to display
    title_durations: Vec<Duration>,
    /// The name of the song author used to display
    title_author: String,
    /// The area reserved for drawing the window
    area: Rect,
}

impl TitleWindow {
    /// Creates a new window
    pub fn new() -> Self {
        TitleWindow {
            selected: false,
            offset: 0,
            album_playing: None,
            album_selected: 0,
            title_playing: None,
            title_selected: 0,
            title_names: vec![],
            title_durations: vec![],
            title_author: String::new(),
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, on_album: bool, music_data: &MusicData, state_data: &StateData) {
        // Update if we are not on the Album window
        self.selected = !on_album;
        // Reconstructs the list of title names and durations of the selected album
        self.title_names.clear();
        self.title_durations.clear();
        for song in &music_data.albums[self.album_selected].songs {
            self.title_names.push(song.title.clone());
            self.title_durations.push(song.duration);
        }
        // Update the author name of the selected album
        self.title_author = music_data.albums.get(self.album_selected).expect("Could not find album selected").artist.clone();

        // Finally update if an album and song is playing
        self.album_playing = state_data.album_id;
        self.title_playing = state_data.title_id;
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
        let mut render_widget = |left: &str, right: &str, style: Style, y: u16| {
            // Render the left part
            frame.render_widget(
                Paragraph::new(left),
                Rect::new(
                    area.x + 2,
                    area.y + 1 + y - self.offset as u16,
                    u16::min(left.len() as u16, area.width - 4),
                    1,
                ),
            );

            // Render the right part
            frame.render_widget(
                Paragraph::new(right),
                Rect::new(
                    area.x + area.width - u16::min(right.len() as u16, area.width - 4) - 2,
                    area.y + 1 + y - self.offset as u16,
                    u16::min(right.len() as u16, area.width - 4),
                    1,
                ),
            );

            // Render the style in the scroll list
            frame.render_widget(
                Paragraph::new("").style(style),
                Rect::new(
                    area.x + 1,
                    area.y + 1 + y - self.offset as u16,
                    area.width - 2,
                    1,
                ),
            );
        };

        // For all songs
        for (i, (title, duration)) in self
            .title_names
            .iter()
            .zip(self.title_durations.iter())
            .enumerate()
            .skip(self.offset)
            .take(area.height as usize - 2)
        {
            // Do we select the Title window, is it the playing album, is it the playing title, is it the selected title
            let selected = self.selected;
            let playing_album = self
                .album_playing
                .map_or(false, |playing| playing == self.album_selected);
            let playing_title = self.title_playing.map_or(false, |playing| playing == i);
            let selected_title = self.title_selected == i;

            // Match style
            let style = match (selected, playing_album, playing_title, selected_title) {
                (true, true, true, true) => Style::default().fg(Color::Black).bg(Color::Cyan),
                (_, true, true, _) => Style::default().fg(Color::Black).bg(Color::Green),
                (true, _, _, true) => Style::default().fg(Color::Black).bg(Color::LightBlue),
                _ => Style::default().fg(Color::DarkGray),
            };

            // Render list item
            let secs = duration.as_secs();
            render_widget(
                &format!("{:2} - {}", i + 1, title),
                &format!(" {:02}:{:02} {}", secs / 60, secs % 60, self.title_author),
                style,
                i as u16,
            );
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
        let text = "Title";
        let len = text.len() as u16;
        frame.render_widget(
            Paragraph::new(format!(" {} ", text)).style(border_style),
            Rect::new(area.x + 5, area.y, len + 2, 1),
        );
    }

    /// Go down the list of songs
    pub fn down(&mut self) {
        if self.title_names.is_empty() {
            return;
        }

        // Update the song selected
        if self.title_selected < self.title_names.len() - 1 {
            self.title_selected += 1;
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

        if self.title_selected > self.offset + self.area.height as usize - 3 - padding
            && self.offset < self.title_names.len() - self.area.height as usize + 2
        {
            self.offset += 1;
        }
    }

    /// Go up the list of songs
    pub fn up(&mut self) {
        if self.title_names.is_empty() {
            return;
        }

        // Update the song selected
        if self.title_selected > 0 {
            self.title_selected -= 1;
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

    /// Update the index of the album selected
    pub fn update_titles(&mut self, album_selected: usize) {
        self.album_selected = album_selected;
    }

    /// Reset the title selected
    pub fn reset_selected(&mut self) {
        self.title_selected = 0;
    }

    /// Play a given song
    ///
    /// - `client`: The client protected for thread safety
    pub fn play(&mut self, client: &mut Arc<Mutex<Client>>) {
        let album_selected = self.album_selected;
        let title_selected = self.title_selected;
        // Start a background thread for the operation
        let client_lock = client.clone();
        tokio::spawn(async move {
            let mut client = client_lock.lock().await;
            client.start_title(album_selected, title_selected);
        });
    }
}
