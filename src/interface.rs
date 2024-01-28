//! This module defines the main loop of the program, it creates the ui, displays
//! the ui, polls the input and starts background tasks

use crossterm::{
    event::{self, KeyEventKind},
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
    ExecutableCommand,
};
use ratatui::prelude::{CrosstermBackend, Terminal};
use std::{
    io::{stdout, Stdout},
    time::{Duration, Instant},
};

use crate::input::{Input, InputControl};
use crate::music::music_client::Client;
use crate::ui::Ui;

/// Stores key information about the interface
pub struct Interface {
    /// Which terminal is used for the user interface
    terminal: Terminal<CrosstermBackend<Stdout>>,
}

impl Interface {
    /// Creates a new interface and initializes it
    pub fn new() -> Interface {
        stdout()
            .execute(EnterAlternateScreen)
            .expect("Could not enter alternate screen");

        enable_raw_mode().expect("Could not enable raw mode");

        let mut terminal =
            Terminal::new(CrosstermBackend::new(stdout())).expect("Could not create terminal");

        terminal.clear().expect("Could not clear terminal");

        Interface { terminal }
    }

    /// Renders the interface
    /// First we poll the events
    /// Then we render the user interface
    /// Finally we start background recurring tasks
    ///
    /// - `client`: The client used for the mpd connection
    pub async fn render(&mut self, client: Client) {
        let mut ui = Ui::new(client);
        let mut i = 0;
        loop {
            let start = Instant::now();

            if event::poll(Duration::new(0, 0)).expect("Could not poll events") {
                if let event::Event::Key(key) = event::read().expect("Could not read event") {
                    if key.kind == KeyEventKind::Press {
                        match Input::handle(key.code, &mut ui).await {
                            InputControl::Continue => continue,
                            InputControl::Break => break,
                        }
                    }
                }
            }

            self.terminal
                .draw(|frame| {
                    ui.render(frame);
                })
                .expect("Could not draw frame");

            if i % 48 == 0 {
                i = 0;
                let client_lock = ui.client.clone();
                tokio::spawn(async move {
                    if let Ok(client) = client_lock.try_lock().as_mut() {
                        client.sync();
                    }
                });
            }
            i += 1;

            let end = start.elapsed();
            let mut wait = Duration::from_micros(16667);
            if wait > end {
                wait -= end;
            }
            tokio::time::sleep(wait).await;
        }
    }
}

impl Drop for Interface {
    /// Destroys the user interface and clears the screen
    fn drop(&mut self) {
        stdout()
            .execute(LeaveAlternateScreen)
            .expect("Could not leave alternate screen");
        disable_raw_mode().expect("Could not disable raw mode");
    }
}
