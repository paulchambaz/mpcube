use crate::{
    music::music_client::Client,
    windows::{
        album_window::AlbumWindow, bar_window::BarWindow, info_window::InfoWindow,
        status_window::StatusWindow, title_window::TitleWindow, volume_window::VolumeWindow,
    },
};
use ratatui::Frame;

pub struct Ui {
    pub client: Client,
    pub on_album: bool,
    pub album_window: AlbumWindow,
    pub title_window: TitleWindow,
    pub info_window: InfoWindow,
    pub volume_window: VolumeWindow,
    pub bar_window: BarWindow,
    pub status_window: StatusWindow,
}

impl Ui {
    pub fn new(client: Client) -> Self {
        Ui {
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

    pub fn render(&mut self, frame: &mut Frame) {
        let a = frame.size();
        let on_album = self.on_album;
        let data = &self.client.data;
        let state = &self.client.state;

        // TODO: updating this each and every time might be slow..
        // furthermore if we do async operations we want to add some kind of
        // protection so that it does not get read while this is getting updated
        // and the ui is being changed...

        // updating windows
        self.album_window.update(on_album, data, state);
        self.title_window.update(on_album, data, state);
        self.info_window.update(on_album, data, state);
        self.volume_window.update(on_album, data, state);
        self.bar_window.update(on_album, data, state);
        self.status_window.update(on_album, data, state);

        // resizing the windows
        let side_width = if a.width > 80 { 40 } else { a.width / 2 };
        let volume_width = if a.width > 90 { 30 } else { a.width / 3 };

        self.album_window
            .update_area(0, 0, side_width, a.height - 2);
        self.title_window
            .update_area(side_width, 0, a.width - side_width, a.height - 2);
        self.info_window
            .update_area(0, a.height - 2, a.width - 9, 1);
        self.volume_window
            .update_area(0, a.height - 1, volume_width, 1);
        self.bar_window
            .update_area(volume_width, a.height - 1, a.width - volume_width - 9, 1);
        self.status_window
            .update_area(a.width - 9, a.height - 2, 8, 2);

        // rendering of windows
        self.album_window.render(frame);
        self.title_window.render(frame);
        self.info_window.render(frame);
        self.volume_window.render(frame);
        self.bar_window.render(frame);
        self.status_window.render(frame);
    }
}
