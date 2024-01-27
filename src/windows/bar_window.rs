use std::time::Duration;

use ratatui::{
    layout::Rect,
    style::{Color, Style},
    widgets::Paragraph,
    Frame,
};

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
                .expect("Could not find album playing");

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
        let area = self.area;
        let style = Style::default().fg(Color::DarkGray);
        let format_time =
            |time: Duration| format!("{:02}:{:02}", time.as_secs() / 60, time.as_secs() % 60);

        let mut render_widget = |text: &str, x: u16, w: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y, w, 1),
            );
        };

        render_widget(
            &format_time(self.position.unwrap_or(Duration::new(0, 0))),
            1,
            5,
        );
        render_widget(
            &format_time(self.duration.unwrap_or(Duration::new(0, 0))),
            area.width - 5,
            5,
        );

        let start = 7;
        let end = area.width - 7;
        for i in start..=end {
            render_widget("─", i, 1);
        }

        if let (Some(position), Some(duration)) = (self.position, self.duration) {
            let ratio = position.as_millis() as f32 / duration.as_millis() as f32;
            let cursor = ((1. - ratio) * start as f32 + ratio * end as f32) as u16;
            render_widget("█", cursor, 1);
        }
    }
}
