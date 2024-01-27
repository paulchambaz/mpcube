use std::time::Duration;

use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::music::{music_client::Client, music_data::{MusicData, StateData}};

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
        let a = self.area;
        let block = Block::default().borders(Borders::ALL).border_style(
            Style::default()
                .fg(if self.selected {
                    Color::LightRed
                } else {
                    Color::DarkGray
                })
                .add_modifier(if self.selected {
                    Modifier::BOLD
                } else {
                    Modifier::empty()
                }),
        );

        for (i, (title, duration)) in self
            .title_names
            .iter()
            .zip(self.title_durations.iter())
            .enumerate() // Enumerate after zipping
            .skip(self.offset) // Skip after enumerating
            .take(a.height as usize - 2)
        // Take after skipping
        {
            let style = if self.selected {
                if let (Some(album_id), Some(title_id)) = (self.album_playing, self.title_playing) {
                    if self.album_selected == album_id && i == title_id && i == self.title_selected
                    {
                        Style::default().fg(Color::Black).bg(Color::Cyan)
                    } else if self.album_selected == album_id && i == title_id {
                        Style::default().fg(Color::Black).bg(Color::Green)
                    } else if i == self.title_selected {
                        Style::default().fg(Color::Black).bg(Color::LightBlue)
                    } else {
                        Style::default().fg(Color::DarkGray)
                    }
                } else if i == self.title_selected {
                    Style::default().fg(Color::Black).bg(Color::LightBlue)
                } else {
                    Style::default().fg(Color::DarkGray)
                }
            } else if let (Some(album_id), Some(title_id)) =
                (self.album_playing, self.title_playing)
            {
                if self.album_selected == album_id && i == title_id {
                    Style::default().fg(Color::Black).bg(Color::Green)
                } else {
                    Style::default().fg(Color::DarkGray)
                }
            } else {
                Style::default().fg(Color::DarkGray)
            };

            let secs = duration.as_secs();
            let left = format!("{:2} - {}", i + 1, title);
            let right = format!("  {:02}:{:02} {}", secs / 60, secs % 60, self.title_author);
            let left_len = left.len() as u16;
            let right_len = right.len() as u16;

            frame.render_widget(
                Paragraph::new("").style(style),
                Rect::new(a.x + 1, a.y + 1 + (i - self.offset) as u16, a.width - 2, 1),
            );

            frame.render_widget(
                Paragraph::new(right).style(style),
                Rect::new(
                    a.x + a.width - u16::min(right_len, a.width - 4) - 2,
                    a.y + 1 + (i - self.offset) as u16,
                    u16::min(right_len, a.width - 4),
                    1,
                ),
            );
            frame.render_widget(
                Paragraph::new(left).style(style),
                Rect::new(
                    a.x + 2,
                    a.y + 1 + (i - self.offset) as u16,
                    u16::min(left_len, a.width - 4),
                    1,
                ),
            );
        }

        frame.render_widget(Paragraph::new("").block(block), a);
        let style = Style::default().fg(if self.selected {
            Color::LightRed
        } else {
            Color::DarkGray
        });
        frame.render_widget(
            Paragraph::new(" Titles ").style(style),
            Rect::new(a.x + 4, 0, a.width - 4, 1),
        );
    }

    const BORDER: usize = 5;

    pub fn down(&mut self) {
        if self.title_names.is_empty() {
            return;
        }

        if self.title_selected < self.title_names.len() - 1 {
            self.title_selected += 1;
        }

        if self.title_selected > self.offset + self.area.height as usize - 3 - Self::BORDER
            && self.offset < self.title_names.len() - self.area.height as usize + 2 {
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

        if self.album_selected < self.offset + Self::BORDER 
            && self.offset > 0 {
            self.offset -= 1;
        }
    }

    pub fn update_titles(&mut self, album_selected: usize) {
        self.album_selected = album_selected;
    }

    pub fn reset_selected(&mut self) {
        self.title_selected = 0;
    }

    pub fn play(&mut self, client: &mut Client) {
        client.start_title(self.album_selected, self.title_selected);
    }
}
