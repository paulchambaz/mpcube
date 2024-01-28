use std::time::Duration;

use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};
use std::sync::Arc;
use tokio::sync::Mutex;

use crate::music::{
    music_client::Client,
    music_data::{MusicData, StateData},
};

pub struct TitleWindow {
    selected: bool,
    offset: usize,
    album_playing: Option<usize>,
    album_selected: usize,
    title_playing: Option<usize>,
    title_selected: usize,
    title_names: Vec<String>,
    title_durations: Vec<Duration>,
    title_author: String,
    area: Rect,
}

impl TitleWindow {
    pub fn new() -> TitleWindow {
        TitleWindow {
            selected: false,
            offset: 0,
            album_playing: None,
            album_selected: 0,
            title_playing: None,
            title_selected: 0,
            title_names: vec![],
            title_durations: vec![],
            title_author: String::new(),
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, on_album: bool, music_data: &MusicData, state_data: &StateData) {
        self.selected = !on_album;
        self.title_names.clear();
        self.title_durations.clear();
        for song in &music_data.albums[self.album_selected].songs {
            self.title_names.push(song.title.clone());
            self.title_durations.push(song.duration);
        }
        self.title_author = music_data.albums[self.album_selected].artist.clone();
        self.album_playing = state_data.album_id;
        self.title_playing = state_data.title_id;
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let area = self.area;

        let mut render_widget = |left: &str, right: &str, style: Style, y: u16| {
            frame.render_widget(
                Paragraph::new(left),
                Rect::new(
                    area.x + 2,
                    area.y + 1 + y - self.offset as u16,
                    u16::min(left.len() as u16, area.width - 4),
                    1,
                ),
            );
            frame.render_widget(
                Paragraph::new(right),
                Rect::new(
                    area.x + area.width - u16::min(right.len() as u16, area.width - 4) - 2,
                    area.y + 1 + y - self.offset as u16,
                    u16::min(right.len() as u16, area.width - 4),
                    1,
                ),
            );
            frame.render_widget(
                Paragraph::new("").style(style),
                Rect::new(
                    area.x + 1,
                    area.y + 1 + y - self.offset as u16,
                    area.width - 2,
                    1,
                ),
            );
        };

        for (i, (title, duration)) in self
            .title_names
            .iter()
            .zip(self.title_durations.iter())
            .enumerate()
            .skip(self.offset)
            .take(area.height as usize - 2)
        {
            let selected = self.selected;
            let playing_album = self
                .album_playing
                .map_or(false, |playing| playing == self.album_selected);
            let playing_title = self.title_playing.map_or(false, |playing| playing == i);
            let selected_title = self.title_selected == i;

            let style = match (selected, playing_album, playing_title, selected_title) {
                (true, true, true, true) => Style::default().fg(Color::Black).bg(Color::Cyan),
                (_, true, true, _) => Style::default().fg(Color::Black).bg(Color::Green),
                (true, _, _, true) => Style::default().fg(Color::Black).bg(Color::LightBlue),
                _ => Style::default().fg(Color::DarkGray),
            };

            let secs = duration.as_secs();
            render_widget(
                &format!("{:2} - {}", i + 1, title),
                &format!(" {:02}:{:02} {}", secs / 60, secs % 60, self.title_author),
                style,
                i as u16,
            );
        }

        let text = "Title";

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
            Rect::new(area.x + 5, area.y, len + 2, 1),
        );
    }

    pub fn down(&mut self) {
        if self.title_names.is_empty() {
            return;
        }

        if self.title_selected < self.title_names.len() - 1 {
            self.title_selected += 1;
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

        if self.title_selected > self.offset + self.area.height as usize - 3 - border
            && self.offset < self.title_names.len() - self.area.height as usize + 2
        {
            self.offset += 1;
        }
    }

    pub fn up(&mut self) {
        if self.title_names.is_empty() {
            return;
        }

        if self.title_selected > 0 {
            self.title_selected -= 1;
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

    pub fn update_titles(&mut self, album_selected: usize) {
        self.album_selected = album_selected;
    }

    pub fn reset_selected(&mut self) {
        self.title_selected = 0;
    }

    pub fn play(&mut self, client: &mut Arc<Mutex<Client>>) {
        let client_lock = client.clone();
        let album_selected = self.album_selected;
        let title_selected = self.title_selected;
        tokio::spawn(async move {
            let mut client = client_lock.lock().await;
            client.start_title(album_selected, title_selected);
        });
    }
}
