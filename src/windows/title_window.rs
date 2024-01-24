use ratatui::{
    layout::Rect,
    prelude::{CrosstermBackend, Stylize, Terminal},
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::mpd_client::{MusicData, StateData};

pub struct TitleWindow {
    selected: bool,
    offset: usize,
    title_playing: Option<usize>,
    title_selected: usize,
    title_names: Vec<String>,
    title_tracks: Vec<u32>,
    title_duration: Vec<String>,
    title_author: String,
    area: Rect,
}

impl TitleWindow {
    pub fn new() -> TitleWindow {
        TitleWindow {
            selected: false,
            offset: 0,
            title_playing: None,
            title_selected: 0,
            title_names: vec![
                "title a".to_string(),
                "title b".to_string(),
                "title c".to_string(),
                "title d".to_string(),
                "title e".to_string(),
            ],
            title_tracks: vec![
                0,
                1,
                2,
                3,
                4,
            ],
            title_duration: vec![
                "05:58".to_string(),
                "06:50".to_string(),
                "03:43".to_string(),
                "02:22".to_string(),
                "01:32".to_string(),
            ],
            title_author: "Author".to_string(),
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, music_data: &MusicData, state_data: &StateData) {
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let area = self.area;
        let block = Block::default()
            .borders(Borders::ALL)
            .title(" Titles ")
            .border_style(
                Style::default()
                    .fg(if self.selected {
                        Color::LightRed
                    } else {
                        Color::Gray
                    })
                    .add_modifier(if self.selected {
                        Modifier::BOLD
                    } else {
                        Modifier::empty()
                    })
            );

        for (i, title) in self.title_names.iter().enumerate().skip(self.offset).take(area.height as usize -2) {
            let style = if self.selected {
                if let Some(playing) = self.title_playing {
                    if i == playing && i == self.title_selected {
                        Style::default().fg(Color::Black).bg(Color::Cyan)
                    } else if i == playing {
                        Style::default().fg(Color::Black).bg(Color::Green)
                    } else if i == self.title_selected {
                        Style::default().fg(Color::Black).bg(Color::LightBlue)
                    } else {
                        Style::default().fg(Color::Gray)
                    }
                } else if i == self.title_selected {
                        Style::default().fg(Color::Black).bg(Color::LightBlue)
                } else {
                        Style::default()
                }
            } else if let Some(playing) = self.title_playing {
                if i == playing {
                    Style::default().fg(Color::Black).bg(Color::Green)
                } else {
                    Style::default()
                }
            } else {
                Style::default()
            };

            let rect = Rect {
                x: area.x + 1,
                y: area.y + 1 + (i - self.offset) as u16,
                width: area.width - 2,
                height: 1,
            };

            frame.render_widget(Paragraph::new(title.clone()).style(style), rect);

        }

        frame.render_widget(Paragraph::new("").block(block), self.area);
    }
}
