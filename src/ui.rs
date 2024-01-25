use crossterm::{
    event::{self, KeyCode, KeyEventKind},
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
    ExecutableCommand,
};
use ratatui::prelude::{CrosstermBackend, Terminal};
use std::io::{stdout, Stdout};

use crate::mpd_client::Client;
use crate::windows::{
    album_window::AlbumWindow, bar_window::BarWindow, info_window::InfoWindow,
    status_window::StatusWindow, title_window::TitleWindow, volume_window::VolumeWindow,
};

pub struct Interface {
    terminal: Terminal<CrosstermBackend<Stdout>>,
    client: Client,
    on_album: bool,
    album_window: AlbumWindow,
    title_window: TitleWindow,
    info_window: InfoWindow,
    volume_window: VolumeWindow,
    bar_window: BarWindow,
    status_window: StatusWindow,
}

impl Interface {
    pub fn new(client: Client) -> Interface {
        stdout()
            .execute(EnterAlternateScreen)
            .expect("Could not enter alternate screen");

        enable_raw_mode().expect("Could not enable raw mode");

        let mut terminal =
            Terminal::new(CrosstermBackend::new(stdout())).expect("Could not create terminal");

        terminal.clear().expect("Could not clear terminal");

        Interface {
            terminal,
            client,
            on_album: true,
            album_window: AlbumWindow::new(),
            title_window: TitleWindow::new(),
            info_window: InfoWindow::new(),
            volume_window: VolumeWindow::new(),
            bar_window: BarWindow::new(),
            status_window: StatusWindow::new(),
        }
    }

    pub fn render(&mut self) {
        loop {
            // main canvas
            self.terminal
                .draw(|frame| {
                    let area = frame.size();

                    if let (Some(music_data), Some(state_data)) =
                        (self.client.data.as_ref(), self.client.state.as_ref())
                    {
                        // TODO: updating this each and every time might be slow..
                        // furthermore if we do async operations we want to add some kind of
                        // protection so that it does not get read while this is getting updated
                        // and the ui is being changed...
                        self.album_window
                            .update(self.on_album, music_data, state_data);
                        self.title_window
                            .update(self.on_album, music_data, state_data);
                        self.info_window
                            .update(self.on_album, music_data, state_data);
                        self.volume_window
                            .update(self.on_album, music_data, state_data);
                        self.bar_window
                            .update(self.on_album, music_data, state_data);
                        self.status_window
                            .update(self.on_album, music_data, state_data);
                    }

                    self.album_window.update_area(
                        0,
                        0,
                        if area.width > 80 { 40 } else { area.width / 2 },
                        area.height - 2,
                    );

                    self.title_window.update_area(
                        if area.width > 80 { 40 } else { area.width / 2 },
                        0,
                        if area.width > 80 {
                            area.width - 40
                        } else {
                            area.width - area.width / 2
                        },
                        area.height - 2,
                    );
                    self.info_window
                        .update_area(0, area.height - 2, area.width - 9, 1);

                    self.volume_window.update_area(
                        0,
                        area.height - 1,
                        if area.width > 90 { 30 } else { area.width / 3 },
                        1,
                    );

                    self.bar_window.update_area(
                        if area.width > 90 { 30 } else { area.width / 3 },
                        area.height - 1,
                        area.width - if area.width > 90 { 30 } else { area.width / 3 } - 9,
                        1,
                    );
                    self.status_window
                        .update_area(area.width - 9, area.height - 2, 8, 2);

                    // rendering of windows
                    self.album_window.render(frame);
                    self.title_window.render(frame);
                    self.info_window.render(frame);
                    self.volume_window.render(frame);
                    self.bar_window.render(frame);
                    self.status_window.render(frame);
                })
                .expect("Could not draw frame");

            // main event poll
            if event::poll(std::time::Duration::from_millis(16)).expect("Could not poll events") {
                if let event::Event::Key(key) = event::read().expect("Could not read event") {
                    // TODO: manage user input for all windows
                    if key.kind == KeyEventKind::Press {
                        match key.code {
                            KeyCode::Char('q') => break,
                            KeyCode::Char('Q') => break,
                            KeyCode::Char('j') => {
                                if self.on_album {
                                    self.album_window.down();
                                    self.title_window
                                        .update_titles(self.album_window.album_selected);
                                    self.title_window.reset_selected();
                                } else {
                                    self.title_window.down();
                                }
                            }
                            KeyCode::Char('k') => {
                                if self.on_album {
                                    self.album_window.up();
                                    self.title_window
                                        .update_titles(self.album_window.album_selected);
                                    self.title_window.reset_selected();
                                } else {
                                    self.title_window.up();
                                }
                            }
                            KeyCode::Char('l') => {
                                self.on_album = false;
                                // if main { go to title }
                            }
                            KeyCode::Char('h') => {
                                self.on_album = true;
                                // if !main { go to album }
                            }
                            KeyCode::Enter => {
                                // if main { start_album } else { start_song }
                            }
                            KeyCode::Char(' ') => {
                                self.info_window.toggle(&mut self.client);
                            }
                            KeyCode::Char('n') => {
                                // next
                            }
                            KeyCode::Char('p') => {
                                // previous
                            }
                            KeyCode::Char('=') => {
                                self.volume_window.volume_up(&mut self.client);
                            }
                            KeyCode::Char('-') => {
                                self.volume_window.volume_down(&mut self.client);
                            }
                            KeyCode::Char('x') => {
                                // clear
                            }
                            KeyCode::Char('.') => {
                                // seek forward
                            }
                            KeyCode::Char(',') => {
                                // seek backward
                            }
                            KeyCode::Char('s') => {
                                self.status_window.shuffle(&mut self.client);
                            }
                            KeyCode::Char('r') => {
                                self.status_window.repeat(&mut self.client);
                            }
                            _ => {}
                        }
                    }
                }
            }
        }
    }
}

impl Drop for Interface {
    fn drop(&mut self) {
        stdout()
            .execute(LeaveAlternateScreen)
            .expect("Could not leave alternate screen");
        disable_raw_mode().expect("Could not disable raw mode");
    }
}
