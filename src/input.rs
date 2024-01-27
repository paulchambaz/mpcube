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
                    ui.album_window.play(&mut ui.client).await;
                } else {
                    ui.title_window.play(&mut ui.client).await;
                }
            }
            KeyCode::Char(' ') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.toggle().await;
                });
            }
            KeyCode::Char('n') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.next().await;
                });
            }
            KeyCode::Char('p') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.previous().await;
                });
            }
            KeyCode::Char('=') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.volume_up().await;
                });
            }
            KeyCode::Char('-') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.volume_down().await;
                });
            }
            KeyCode::Char('x') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.clear().await;
                });
            }
            KeyCode::Char('.') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.seek_forward().await;
                });
            }
            KeyCode::Char(',') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.seek_backward().await;
                });
            }
            KeyCode::Char('s') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.shuffle().await;
                });
            }
            KeyCode::Char('r') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.repeat().await;
                });
            }
            KeyCode::Char('U') => {
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    let mut client = client_lock.lock().await;
                    client.full_sync().await;
                });
            }
            _ => {}
        }

        InputControl::Continue
    }
}
