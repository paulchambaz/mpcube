use crossterm::{
    event::{self, KeyCode, KeyEventKind},
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
    ExecutableCommand,
};
use ratatui::{
    layout::Rect,
    prelude::{CrosstermBackend, Stylize, Terminal},
    style::{Color, Modifier, Style},
    text::Span,
    widgets::{Block, Borders, Paragraph, Wrap},
    Frame,
};
use std::io::{stdout, Stdout};

use crate::mpd_client::Client;
use crate::windows::{
    album_window::AlbumWindow,
    title_window::TitleWindow,
    info_window::InfoWindow,
    volume_window::VolumeWindow,
    bar_window::BarWindow,
    status_window::StatusWindow,
};

pub struct Interface {
    terminal: Terminal<CrosstermBackend<Stdout>>,
    client: Client,
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
                        self.album_window.update(music_data, state_data);
                        self.title_window.update(music_data, state_data);
                        self.info_window.update(music_data, state_data);
                        self.volume_window.update(music_data, state_data);
                        self.bar_window.update(music_data, state_data);
                        self.status_window.update(music_data, state_data);
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
                        if area.width > 80 { area.width - 40 } else { area.width - area.width / 2 },
                        area.height - 2,
                    );
                    self.info_window.update_area(
                        0,
                        area.height - 2,
                        area.width - 9,
                        1,
                    );
                    self.volume_window.update_area(
                        0,
                        area.height - 1,
                        30,
                        1,
                    );
                    // TODO: the lack of resizing can cause negative sizes which results in panic
                    self.bar_window.update_area(
                        30,
                        area.height - 1,
                        area.width - 30 - 9,
                        1,
                    );
                    self.status_window.update_area(
                        area.width - 9,
                        area.height - 2,
                        8,
                        2,
                    );

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
                            KeyCode::Char('j') => self.album_window.down(),
                            KeyCode::Char('k') => self.album_window.up(),
                            KeyCode::Char('s') => self.status_window.shuffle(),
                            KeyCode::Char('r') => self.status_window.repeat(),
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
