use crossterm::{
    event::{self, KeyEventKind},
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
    ExecutableCommand,
};
use ratatui::prelude::{CrosstermBackend, Terminal};
use std::io::{stdout, Stdout};

use crate::input::{Input, InputControl};
use crate::music::music_client::Client;
use crate::ui::Ui;

pub struct Interface {
    terminal: Terminal<CrosstermBackend<Stdout>>,
}

impl Interface {
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

    pub fn render(&mut self, client: Client) {
        let mut ui = Ui::new(client);
        loop {
            // main canvas
            self.terminal
                .draw(|frame| {
                    ui.render(frame);
                })
                .expect("Could not draw frame");

            if event::poll(std::time::Duration::from_millis(16)).expect("Could not poll events") {
                // this is where we will start the background update
                // every 60 frames, we should run an update status and update song update

                if let event::Event::Key(key) = event::read().expect("Could not read event") {
                    if key.kind == KeyEventKind::Press {
                        match Input::handle(key.code, &mut ui) {
                            InputControl::Continue => {}
                            InputControl::Break => break,
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
