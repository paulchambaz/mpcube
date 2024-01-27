use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::music::{
    music_client::Client,
    music_data::{MusicData, StateData},
};
use std::sync::Arc;
use tokio::sync::Mutex;

pub struct AlbumWindow {
    selected: bool,
    offset: usize,
    album_playing: Option<usize>,
    pub album_selected: usize,
    album_names: Vec<String>,
    area: Rect,
}

impl AlbumWindow {
    pub fn new() -> AlbumWindow {
        AlbumWindow {
            selected: true,
            offset: 0,
            album_playing: None,
            album_selected: 0,
            album_names: vec![],
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, on_album: bool, music_data: &MusicData, state_data: &StateData) {
        self.selected = on_album;
        self.album_names.clear();
        for album in &music_data.albums {
            self.album_names.push(album.album.clone());
        }
        self.album_playing = state_data.album_id;
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let area = self.area;

        let mut render_widget = |text: &str, style: Style, y: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(
                    area.x + 1,
                    area.y + 1 + y - self.offset as u16,
                    area.width - 2,
                    1,
                ),
            );
        };

        for (i, album) in self
            .album_names
            .iter()
            .enumerate()
            .skip(self.offset)
            .take(area.height as usize - 2)
        {
            let playing_album = self.album_playing.map_or(false, |playing| playing == i);
            let selected_album = self.album_selected == i;

            let style = match (playing_album, selected_album) {
                (true, true) => Style::default().fg(Color::Black).bg(Color::Cyan),
                (true, false) => Style::default().fg(Color::Black).bg(Color::Green),
                (false, true) => Style::default().fg(Color::Black).bg(Color::LightBlue),
                (false, false) => Style::default().fg(Color::DarkGray),
            };

            render_widget(album, style, i as u16);
        }

        let text = "Album";

        let border_style = if self.selected {
            Style::default()
                .fg(Color::LightRed)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };

        frame.render_widget(
            Paragraph::new("").block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_style(border_style),
            ),
            area,
        );

        let len = text.len() as u16;
        frame.render_widget(
            Paragraph::new(format!(" {} ", text)).style(border_style),
            Rect::new(area.x + (area.width - len - 2) / 2, area.y, len + 2, 1),
        );
    }

    pub fn down(&mut self) {
        if self.album_names.is_empty() {
            return;
        }

        if self.album_selected < self.album_names.len() - 1 {
            self.album_selected += 1;
        }

        let border = match self.area.height as usize - 2 {
            0..=3 => 0,
            4..=7 => 1,
            8..=11 => 2,
            12..=15 => 3,
            16..=19 => 4,
            20..=usize::MAX => 5,
            _ => 0,
        };

        if self.album_selected > self.offset + self.area.height as usize - 3 - border
            && self.offset < self.album_names.len() - self.area.height as usize + 2
        {
            self.offset += 1;
        }
    }

    pub fn up(&mut self) {
        if self.album_names.is_empty() {
            return;
        }

        if self.album_selected > 0 {
            self.album_selected -= 1;
        }

        let border = match self.area.height as usize - 2 {
            0..=3 => 0,
            4..=7 => 1,
            8..=11 => 2,
            12..=15 => 3,
            16..=19 => 4,
            20..=usize::MAX => 5,
            _ => 0,
        };

        if self.album_selected < self.offset + border && self.offset > 0 {
            self.offset -= 1;
        }
    }

    pub fn play(&mut self, client: &mut Arc<Mutex<Client>>) {
        let client_lock = client.clone();
        let album_selected = self.album_selected;
        tokio::spawn(async move {
            let mut client = client_lock.lock().await;
            client.start_album(album_selected);
        });
    }
}
