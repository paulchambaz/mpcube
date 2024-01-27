use crossterm::event::KeyCode;

use crate::ui::Ui;

pub enum InputControl {
    Continue,
    Break,
}

pub struct Input;

impl Input {
    pub fn handle(key: KeyCode, ui: &mut Ui) -> InputControl {
        match key {
            KeyCode::Char('q') => return InputControl::Break,
            KeyCode::Char('Q') => return InputControl::Break,
            KeyCode::Char('j') => {
                if ui.on_album {
                    ui.album_window.down();
                    // TODO: this should only occur if we actually change position
                    ui.title_window
                        .update_titles(ui.album_window.album_selected);
                    ui.title_window.reset_selected();
                } else {
                    ui.title_window.down();
                }
            }
            KeyCode::Char('k') => {
                if ui.on_album {
                    ui.album_window.up();
                    // TODO: this should only occur if we actually change position
                    ui.title_window
                        .update_titles(ui.album_window.album_selected);
                    ui.title_window.reset_selected();
                } else {
                    ui.title_window.up();
                }
            }
            KeyCode::Char('l') => {
                ui.on_album = false;
            }
            KeyCode::Char('h') => {
                ui.on_album = true;
            }
            KeyCode::Enter => {
                if ui.on_album {
                    ui.album_window.play(&mut ui.client);
                } else {
                    ui.title_window.play(&mut ui.client);
                }
            }
            KeyCode::Char(' ') => {
                ui.client.toggle();
            }
            KeyCode::Char('n') => {
                ui.client.next();
            }
            KeyCode::Char('p') => {
                ui.client.previous();
            }
            KeyCode::Char('=') => {
                ui.client.volume_up();
            }
            KeyCode::Char('-') => {
                ui.client.volume_down();
            }
            KeyCode::Char('x') => {
                ui.client.clear();
            }
            KeyCode::Char('.') => {
                ui.client.seek_forward();
            }
            KeyCode::Char(',') => {
                ui.client.seek_backward();
            }
            KeyCode::Char('s') => {
                ui.client.shuffle();
            }
            KeyCode::Char('r') => {
                ui.client.repeat();
            }
            KeyCode::Char('U') => {
                ui.client.full_sync();
            }
            _ => {}
        }

        InputControl::Continue
    }
}
