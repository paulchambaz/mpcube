use ratatui::{
    layout::Rect,
    style::{Color, Modifier, Style},
    widgets::Paragraph,
    Frame,
};

use crate::music::music_data::{MusicData, StateData};

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

        let get_style = |on: bool| {
            if on {
                Style::default()
                    .fg(Color::Green)
                    .add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::DarkGray)
            }
        };

        let mut render_widget = |text: &str, style: Style, x: u16, y: u16| {
            frame.render_widget(
                Paragraph::new(text).style(style),
                Rect::new(area.x + x, area.y + y, text.len() as u16, 1),
            );
        };

        render_widget("shuffle", get_style(self.shuffle), 1, 0);
        render_widget("repeat", get_style(self.repeat), 2, 1);
    }
}
