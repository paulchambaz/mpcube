use std::time::Duration;

use ratatui::{layout::Rect, prelude::Stylize, style::Color, widgets::Paragraph, Frame};

use crate::music::music_data::{MusicData, StateData};

pub struct BarWindow {
    position: Option<Duration>,
    duration: Option<Duration>,
    area: Rect,
}

impl BarWindow {
    pub fn new() -> BarWindow {
        BarWindow {
            position: None,
            duration: None,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, _: bool, music_data: &MusicData, state_data: &StateData) {
        self.position = state_data.position;
        if let Some(album_id) = state_data.album_id {
            let album = music_data
                .albums
                .get(album_id)
                .expect("Could not find album playing"); // &Album
            if let Some(title_id) = state_data.title_id {
                let song = album
                    .songs
                    .get(title_id)
                    .expect("Could not find title playing");
                self.duration = Some(song.duration);
            }
        }
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let a = self.area;
        if let Some(position) = self.position {
            let seconds = position.as_secs();
            frame.render_widget(
                Paragraph::new(format!("{:02}:{:02}", seconds / 60, seconds % 60))
                    .fg(Color::DarkGray),
                Rect::new(a.x + 1, a.y, 5, 1),
            );
        }
        if let Some(duration) = self.duration {
            let seconds = duration.as_secs();
            frame.render_widget(
                Paragraph::new(format!("{:02}:{:02}", seconds / 60, seconds % 60))
                    .fg(Color::DarkGray),
                Rect::new(a.x + a.width - 5, a.y, 5, 1),
            );
        }
        let start = a.x + 7;
        let end = start + a.width - 14;
        for i in start..=end {
            frame.render_widget(
                Paragraph::new("─").fg(Color::DarkGray),
                Rect::new(i, a.y, 1, 1),
            );
        }

        if let (Some(position), Some(duration)) = (self.position, self.duration) {
            let duration = duration.as_millis() as f32;
            let position = position.as_millis() as f32;

            let ratio = position / duration;

            let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;

            frame.render_widget(
                Paragraph::new("█").fg(Color::DarkGray),
                Rect::new(cursor, a.y, 1, 1),
            );
        }
    }
}
