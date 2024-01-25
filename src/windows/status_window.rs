use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::Paragraph,
    Frame,
};

use crate::mpd_client::{Client, MusicData, StateData};

pub struct StatusWindow {
    shuffle: bool,
    repeat: bool,
    area: Rect,
}

impl StatusWindow {
    pub fn new() -> StatusWindow {
        StatusWindow {
            shuffle: false,
            repeat: false,
            area: Rect::new(0, 0, 0, 0),
        }
    }

    pub fn update(&mut self, _: bool, _: &MusicData, state_data: &StateData) {
        self.shuffle = state_data.shuffle;
        self.repeat = state_data.repeat;
    }

    pub fn update_area(&mut self, x: u16, y: u16, width: u16, height: u16) {
        self.area.x = x;
        self.area.y = y;
        self.area.width = width;
        self.area.height = height;
    }

    pub fn render(&mut self, frame: &mut Frame) {
        let area = self.area;
        let style_shuffle = if self.shuffle {
            Style::default()
                .fg(Color::Green)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };
        let rect_shuffle = Rect {
            x: area.x,
            y: area.y,
            width: 9,
            height: 1,
        };
        frame.render_widget(
            Paragraph::new(" shuffle").style(style_shuffle),
            rect_shuffle,
        );

        let style_repeat = if self.repeat {
            Style::default()
                .fg(Color::Green)
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(Color::DarkGray)
        };
        let rect_repeat = Rect {
            x: area.x,
            y: area.y + 1,
            width: 9,
            height: 1,
        };
        frame.render_widget(Paragraph::new("  repeat").style(style_repeat), rect_repeat);
    }

    pub fn shuffle(&mut self, client: &mut Client) {
        client.shuffle();
    }

    pub fn repeat(&mut self, client: &mut Client) {
        client.repeat();
    }
}
