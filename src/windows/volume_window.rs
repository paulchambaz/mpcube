use ratatui::{layout::Rect, prelude::Stylize, style::Color, widgets::Paragraph, Frame};

use crate::mpd_client::{Client, MusicData, StateData};

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

        let volume_bar_width = a.width - 10;
        for i in 0..volume_bar_width {
            frame.render_widget(
                Paragraph::new("─").fg(Color::DarkGray),
                Rect::new(a.x + 5 + i, a.y, 1, 1),
            );
        }
        let start = 5;
        let end = 5 + volume_bar_width;
        let mut volume_indicator_position = start + ((end - start) * self.volume as u16 / 100);
        if volume_indicator_position > end - 1 {
            volume_indicator_position = end - 1;
        }
        frame.render_widget(
            Paragraph::new("█").fg(Color::DarkGray),
            Rect::new(a.x + volume_indicator_position, a.y, 1, 1),
        );

        frame.render_widget(
            Paragraph::new(format!("{}%", self.volume)).fg(Color::DarkGray),
            Rect::new(a.x + a.width - 4, a.y, 4, 1),
        );
    }

    pub fn volume_up(&mut self, client: &mut Client) {
        client.volume_up();
    }

    pub fn volume_down(&mut self, client: &mut Client) {
        client.volume_down();
    }
}
