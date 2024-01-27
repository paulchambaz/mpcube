use ratatui::{layout::Rect, prelude::Stylize, style::Color, widgets::Paragraph, Frame};

use crate::music::music_data::{MusicData, StateData};

pub struct VolumeWindow {
    volume: i8,
    area: Rect,
}

impl VolumeWindow {
    pub fn new() -> VolumeWindow {
        VolumeWindow {
            volume: 0,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, _: bool, _: &MusicData, state_data: &StateData) {
        self.volume = state_data.volume;
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let a = self.area;
        frame.render_widget(Paragraph::new(" Vol").fg(Color::DarkGray), self.area);

        let start = a.x + 5;
        let end = start + a.width - 11;

        for i in start..=end {
            frame.render_widget(
                Paragraph::new("─").fg(Color::DarkGray),
                Rect::new(i, a.y, 1, 1),
            );
        }

        let ratio = self.volume as f32 / 100.;
        let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;

        frame.render_widget(
            Paragraph::new("█").fg(Color::DarkGray),
            Rect::new(cursor, a.y, 1, 1),
        );

        frame.render_widget(
            Paragraph::new(format!("{}%", self.volume)).fg(Color::DarkGray),
            Rect::new(a.x + a.width - 4, a.y, 4, 1),
        );
    }
}
