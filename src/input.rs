use crossterm::event::KeyCode;

use crate::ui::Ui;

pub enum InputControl {
    Continue,
    Break,
}

pub struct Input;

impl Input {
    pub async fn handle(key: KeyCode, ui: &mut Ui) -> InputControl {
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
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.toggle();
                    }
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
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.volume_up();
                    }
                });
            }
            KeyCode::Char('-') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.volume_down();
                    }
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
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.seek_forward();
                    }
                });
            }
            KeyCode::Char(',') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.seek_backward();
                    }
                });
            }
            KeyCode::Char('s') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.shuffle();
                    }
                });
            }
            KeyCode::Char('r') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.repeat();
                    }
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
