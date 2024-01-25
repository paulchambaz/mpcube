use ratatui::{layout::Rect, prelude::Stylize, style::Color, widgets::Paragraph, Frame};

use crate::mpd_client::{MusicData, StateData};

pub struct BarWindow {
    area: Rect,
}

impl BarWindow {
    pub fn new() -> BarWindow {
        BarWindow {
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, _: bool, _: &MusicData, _: &StateData) {}

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let a = self.area;
        frame.render_widget(
            Paragraph::new("00:00").fg(Color::DarkGray),
            Rect::new(a.x + 1, a.y, 5, 1),
        );
        frame.render_widget(
            Paragraph::new("00:00").fg(Color::DarkGray),
            Rect::new(a.x + a.width - 5, a.y, 5, 1),
        );
        let volume_bar_width = a.width - 13;
        for i in 0..volume_bar_width {
            frame.render_widget(
                Paragraph::new("─").fg(Color::DarkGray),
                Rect::new(a.x + 7 + i, a.y, 1, 1),
            );
        }
    }
}
