use ratatui::{
    layout::Rect,
    style::{Color, Style},
    widgets::Paragraph,
    Frame,
};

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
        let area = self.area;
        let style = Style::default().fg(Color::DarkGray);

        let mut render_widget = |text: &str, x: u16, w: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y, w, 1),
            );
        };

        render_widget("Vol", 1, 3);
        render_widget(&format!("{}%", self.volume), area.width - 4, 4);

        let start = 5;
        let end = area.width - 6;
        for i in start..=end {
            render_widget("─", i, 1);
        }

        let ratio = self.volume as f32 / 100.;
        let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;
        render_widget("█", cursor, 1);
    }
}
