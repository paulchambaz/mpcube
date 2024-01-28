//! This module is used to manage the Volume window

use ratatui::{
    layout::Rect,
    style::{Color, Style},
    widgets::Paragraph,
    Frame,
};

use crate::music::music_data::{MusicData, StateData};

/// Used to display the Album window
pub struct VolumeWindow {
    /// The value of the volume used to display
    volume: i8,
    /// The area reserved for drawing the window
    area: Rect,
}

impl VolumeWindow {
    /// Creates a new window
    pub fn new() -> Self {
        VolumeWindow {
            volume: 0,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, _: bool, _: &MusicData, state_data: &StateData) {
        // Update the volume
        self.volume = state_data.volume;
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

        // Used to display the text correctly
        let mut render_widget = |text: &str, x: u16, w: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y, w, 1),
            );
        };

        // Render window title
        render_widget("Vol", 1, 3);

        // Render volume value
        render_widget(&format!("{}%", self.volume), area.width - 4, 4);

        // Render the bar
        let start = 5;
        let end = area.width - 6;
        for i in start..=end {
            render_widget("─", i, 1);
        }

        // Render the cursor of the volume level
        let ratio = self.volume as f32 / 100.;
        let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;
        render_widget("█", cursor, 1);
    }
}
