use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::mpd_client::{MusicData, StateData};

pub struct AlbumWindow {
    selected: bool,
    offset: usize,
    album_playing: Option<usize>,
    album_selected: usize,
    album_names: Vec<String>,
    area: Rect,
}

// TODO: with time, im thinking a more general window class would be great to group a lot of the
// logic toghether (mainly the logic of the size and the logic of the update

// TODO: add mechanism to check if there has been an update and the window needs to be redrawn
impl AlbumWindow {
    pub fn new() -> AlbumWindow {
        AlbumWindow {
            selected: true,
            offset: 0,
            album_playing: None,
            album_selected: 0,
            album_names: vec![
            ],
            area: Rect {
                x: 0,
                y: 0,
                width: 0,
                height: 0,
            },
        }
    }

    pub fn update(&mut self, music_data: &MusicData, state_data: &StateData) {
        self.album_names.clear();
        for album in &music_data.albums {
            self.album_names.push(album.album.clone());
        }
        self.album_playing = state_data.id;
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        if self.area.x != x {
            self.area.x = x;
        }
        if self.area.y != y {
            self.area.y = y;
        }
        if self.area.width != width {
            self.area.width = width;
        }
        if self.area.height != height {
            self.area.height = height;
        }
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let area = self.area;
        let block = Block::default()
            .borders(Borders::ALL)
            .title(" Albums ")
            .border_style(
                Style::default()
                    .fg(if self.selected {
                        Color::LightRed
                    } else {
                        Color::White
                    })
                    .add_modifier(if self.selected {
                        Modifier::BOLD
                    } else {
                        Modifier::empty()
                    }),
            );

        // TODO: this is probably inneficient and we can maybe draw everything at once
        // right now we do n renders for each album and 1 render for the box itself
        // maybe we have to seperate the box from the album names but its quite likely we can
        // draw every text at once
        // a key point is keeping the Rect since they provide excellent protection for
        // overflowing
        // maybe another way to solve this issue is rather to tell the renderer to not render
        // for a while, give it a description of what exists, then render everything at once at
        // the end to reduce draw calls
        // that being said, its not so bad if it is slow since its only a very low number of
        // items being drawn (n ~= 20)
        for (i, album) in self
            .album_names
            .iter()
            .enumerate()
            .skip(self.offset)
            .take(area.height as usize - 2)
        {
            let style = if let Some(playing) = self.album_playing {
                if i == playing && i == self.album_selected {
                    Style::default().fg(Color::Black).bg(Color::Cyan)
                } else if i == playing {
                    Style::default().fg(Color::Black).bg(Color::Green)
                } else if i == self.album_selected {
                    Style::default().fg(Color::Black).bg(Color::LightBlue)
                } else {
                    Style::default()
                }
            } else if i == self.album_selected {
                Style::default().fg(Color::Black).bg(Color::LightBlue)
            } else {
                Style::default()
            };

            let rect = Rect {
                x: area.x + 1,
                y: area.y + 1 + (i - self.offset) as u16,
                width: area.width - 2,
                height: 1,
            };

            frame.render_widget(Paragraph::new(album.clone()).style(style), rect);
        }

        frame.render_widget(Paragraph::new("").block(block), self.area);
    }

    const BORDER: usize = 5;

    pub fn down(&mut self) {
        if self.album_names.is_empty() {
            return;
        }

        if self.album_selected < self.album_names.len() - 1 {
            self.album_selected += 1;
        }

        // when do we want it to update
        if self.album_selected > self.area.height as usize - 3 - Self::BORDER
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

        if self.album_selected < self.offset + Self::BORDER
            && self.offset > 0
        {
            self.offset -= 1;
        }
    }
}
