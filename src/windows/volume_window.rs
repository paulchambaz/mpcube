use ratatui::{
    layout::Rect,
    prelude::{CrosstermBackend, Stylize, Terminal},
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Paragraph},
    Frame,
};

use crate::mpd_client::{MusicData, StateData};

pub struct VolumeWindow {
    area: Rect,
}

impl VolumeWindow {
    pub fn new() -> VolumeWindow {
        VolumeWindow {
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
        frame.render_widget(
            Paragraph::new("Volume")
                .white()
                .on_yellow(),
            self.area);
    }
}
