//! This module is used to manage the Status window

use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::Paragraph,
    Frame,
};

use crate::music::music_data::{MusicData, StateData};

/// Used to display the Status window
pub struct StatusWindow {
    /// Whether or not shuffle mode is activated
    shuffle: bool,
    /// Whether or not repeat mode is activated
    repeat: bool,
    /// The area reserved for drawing the window
    area: Rect,
}

impl StatusWindow {
    /// Creates a new window
    pub fn new() -> Self {
        StatusWindow {
            shuffle: false,
            repeat: false,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    /// Updates the value of the window from the client
    ///
    /// - `on_album`: Whether or not the user is on the Album window
    /// - `music_data`: The mpd library representation
    /// - `state_data`: The current mpd state reprensentation
    pub fn update(&mut self, _: bool, _: &MusicData, state_data: &StateData) {
        // Update if shuffle mode is activated
        self.shuffle = state_data.shuffle;
        // Update if repeat mode is activated
        self.repeat = state_data.repeat;
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

        // Used to get tyle information
        let get_style = |on: bool| {
            if on {
                Style::default()
                    .fg(Color::Green)
                    .add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::DarkGray)
            }
        };

        // Used to display the text in correct format
        let mut render_widget = |text: &str, style: Style, x: u16, y: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y + y, text.len() as u16, 1),
            );
        };

        // Renders shuffle mode
        render_widget("shuffle", get_style(self.shuffle), 1, 0);

        // Renders repeat mode
        render_widget("repeat", get_style(self.repeat), 2, 1);
    }
}
