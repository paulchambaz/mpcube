//! This module is used to manage the Bar window

use std::time::Duration;

use ratatui::{
    layout::Rect,
    style::{Color, Style},
    widgets::Paragraph,
    Frame,
};

use crate::music::music_data::{MusicData, StateData};

/// Used to display the Bar window
pub struct BarWindow {
    /// The current position of the song playing (if any)
    position: Option<Duration>,
    /// The current duration of the song playing (if any)
    duration: Option<Duration>,
    /// The area reserved for drawing the window
    area: Rect,
}

impl BarWindow {
    /// Creates a new window
    pub fn new() -> Self {
        BarWindow {
            position: None,
            duration: None,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, _: bool, music_data: &MusicData, state_data: &StateData) {
        // Update the position
        self.position = state_data.position;
        // If an album is playing, get the album
        if let Some(album_id) = state_data.album_id {
            let album = music_data
                .albums
                .get(album_id)
                .expect("Could not find album playing");

            // If a song is playing, get the song
            if let Some(title_id) = state_data.title_id {
                let song = album
                    .songs
                    .get(title_id)
                    .expect("Could not find title playing");

                // Update the duration of the song playing
                self.duration = Some(song.duration);
            }
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

        // Get style
        let style = Style::default().fg(Color::DarkGray);

        // Used to format the time correctly
        let format_time =
            |time: Duration| format!("{:02}:{:02}", time.as_secs() / 60, time.as_secs() % 60);

        // Used to display the text correctly
        let mut render_widget = |text: &str, x: u16, w: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y, w, 1),
            );
        };

        // Render the position time
        render_widget(
            &format_time(self.position.unwrap_or(Duration::new(0, 0))),
            1,
            5,
        );

        // Render the duration time
        render_widget(
            &format_time(self.duration.unwrap_or(Duration::new(0, 0))),
            area.width - 5,
            5,
        );

        // Render the bar
        let start = 7;
        let end = area.width - 7;
        for i in start..=end {
            render_widget("─", i, 1);
        }

        // Render the cursor if a song is playing
        if let (Some(position), Some(duration)) = (self.position, self.duration) {
            let ratio = position.as_millis() as f32 / duration.as_millis() as f32;
            let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;
            render_widget("█", cursor, 1);
        }
    }
}
