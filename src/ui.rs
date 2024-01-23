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
use crate::windows::album_window::AlbumWindow;

pub struct Interface {
    terminal: Terminal<CrosstermBackend<Stdout>>,
    client: Client,
    album_window: AlbumWindow,
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
                        self.album_window.update(music_data, state_data);
                        // TODO: update other windows
                    }

                    self.album_window.update_area(
                        0,
                        0,
                        if area.width > 80 { 40 } else { area.width / 2 },
                        area.height - 2,
                    );
                    // TODO: update other windows area

                    // rendering of windows
                    self.album_window.render(frame);
                    // TODO: render other windows
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
