use ratatui::{layout::Rect, prelude::Stylize, style::Color, widgets::Paragraph, Frame};

use crate::music::music_data::{MusicData, StateData};

pub struct InfoWindow {
    playing: bool,
    title: Option<String>,
    artist: Option<String>,
    album: Option<String>,
    area: Rect,
}

impl InfoWindow {
    pub fn new() -> InfoWindow {
        InfoWindow {
            playing: false,
            title: None,
            artist: None,
            album: None,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, _: bool, music_data: &MusicData, state_data: &StateData) {
        self.playing = state_data.playing;
        if let (Some(album_id), Some(title_id)) = (state_data.album_id, state_data.title_id) {
            let album = music_data
                .albums
                .get(album_id)
                .expect("Could not get album value from its id");
            self.album = Some(album.album.clone());
            self.artist = Some(album.artist.clone());
            self.title = Some(
                album
                    .songs
                    .get(title_id)
                    .expect("Could not get title value from its id")
                    .title
                    .clone(),
            );
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
        if let (Some(title), Some(artist), Some(album)) = (&self.title, &self.artist, &self.album) {
            let str: [&str; 6] = [
                if self.playing { "Playing" } else { "Paused" },
                title,
                "by",
                artist,
                "from",
                album,
            ];
            let len: Vec<u16> = str.iter().map(|s| s.len() as u16).collect();

            let mut sum = 1;

            str.iter()
                .zip(len.iter())
                .enumerate()
                .for_each(|(i, (&text, &length))| {
                    if sum > a.width {
                        return;
                    }
                    if i % 2 == 0 {
                        frame.render_widget(
                            Paragraph::new(text).fg(Color::DarkGray),
                            Rect::new(a.x + sum, a.y, u16::min(length, a.width - sum), a.height),
                        );
                    } else {
                        frame.render_widget(
                            Paragraph::new(text).fg(Color::Green),
                            Rect::new(a.x + sum, a.y, u16::min(length, a.width - sum), a.height),
                        );
                    }
                    sum += length + 1;
                });
        } else {
            frame.render_widget(
                Paragraph::new(" Not playing").fg(Color::DarkGray),
                self.area,
            );
        }
    }
}
