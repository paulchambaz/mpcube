//! This module is used to define all input commands and the action associated
//! with them. It is used by the interface module inside the main loop

use crossterm::event::KeyCode;

use crate::ui::Ui;

/// Used to represent the control flow of the input
pub enum InputControl {
    /// We should continue running the program
    Continue,
    /// The user has specified they want to stop the program
    Break,
}

/// Used to represent the input
pub struct Input;

impl Input {
    /// Handles the input by starting the appropriate task given the keycode the user has pressed
    ///
    /// - `key`: The KeyCode to be handled
    /// - `ui`: A reference to the ui struct
    ///
    /// Returns an `InputControl`, which reprensents the control flow of the program (continue,
    /// break)
    pub async fn handle(key: KeyCode, ui: &mut Ui) -> InputControl {
        match key {
            KeyCode::Char('q') | KeyCode::Char('Q') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.clear();
                });
                return InputControl::Break;
            },
            KeyCode::Char('j') | KeyCode::Down => {
                if ui.on_album {
                    ui.album_window.down();
                    ui.title_window
                        .update_titles(ui.album_window.album_selected);
                    ui.title_window.reset_selected();
                } else {
                    ui.title_window.down();
                }
            }
            KeyCode::Char('k') | KeyCode::Up => {
                if ui.on_album {
                    ui.album_window.up();
                    ui.title_window
                        .update_titles(ui.album_window.album_selected);
                    ui.title_window.reset_selected();
                } else {
                    ui.title_window.up();
                }
            }
            KeyCode::Char('l') | KeyCode::Right => {
                ui.on_album = false;
            }
            KeyCode::Char('h') | KeyCode::Left => {
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
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.toggle();
                });
            }
            KeyCode::Char('n') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.next();
                });
            }
            KeyCode::Char('p') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.previous();
                });
            }
            KeyCode::Char('=') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.volume_up();
                });
            }
            KeyCode::Char('-') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.volume_down();
                });
            }
            KeyCode::Char('x') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.clear();
                });
            }
            KeyCode::Char('.') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.seek_forward();
                });
            }
            KeyCode::Char(',') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.seek_backward();
                });
            }
            KeyCode::Char('s') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.shuffle();
                });
            }
            KeyCode::Char('r') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.repeat();
                });
            }
            KeyCode::Char('U') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.full_sync();
                });
            }
            KeyCode::Char('R') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.random();
                });
            }
            _ => {}
        }

        InputControl::Continue
    }
}
