use std::time::Duration;
use crate::music::music_data::{MusicData, StateData};

pub struct Client {
    pub client: mpd::Client,
    pub data: MusicData,
    pub state: StateData,
    cache_path: String,
}
impl Client {
    pub fn new(address: &str, port: u16, cache_path: String) -> Client {
        let mut client = mpd::Client::connect(format!("{}:{}", address, port))
            .expect("Could not connect to mpd");

        let data = if std::path::Path::new(&cache_path).exists() {
            MusicData::from_cache(&cache_path)
        } else {
            let music_data = MusicData::new(&mut client);
            music_data.save_cache(&cache_path);
            music_data
        };

        let state = StateData::new(&mut client, &data);

        Self {
            client,
            data,
            state,
            cache_path,
        }
    }

    pub fn full_sync(&mut self) {
        self.data = MusicData::new(&mut self.client);
        self.data.save_cache(&self.cache_path);
        self.state = StateData::new(&mut self.client, &self.data);
    }

    pub fn start_album(&mut self, album_id: usize) {
        // TODO: probably an await for this clear
        self.clear();
        let album = self
            .data
            .albums
            .get(album_id)
            .expect("Could not get album at id");
        for song in &album.songs {
            let real_song = mpd::Song {
                file: song.uri.clone(),
                name: None,
                title: None,
                last_mod: None,
                artist: None,
                duration: None,
                place: None,
                range: None,
                tags: Vec::new(),
            };
            if self.client.push(real_song).is_err() {
                return;
            }
        }
        if self.client.play().is_ok() {
            self.state.update_status(&mut self.client);
            self.state.update_song(&mut self.client, &self.data);
        }
    }

    pub fn start_title(&mut self, album_id: usize, title_id: usize) {
        self.clear();
        let album = self
            .data
            .albums
            .get(album_id)
            .expect("Could not get album at id");
        for song in &album.songs {
            let real_song = mpd::Song {
                file: song.uri.clone(),
                name: None,
                title: None,
                last_mod: None,
                artist: None,
                duration: None,
                place: None,
                range: None,
                tags: Vec::new(),
            };
            if self.client.push(real_song).is_err() {
                return;
            }
        }
        if self.client.play().is_err() {
            return;
        }

        for _ in 0..title_id {
            if self.client.next().is_err() {
                self.state.update_status(&mut self.client);
                self.state.update_song(&mut self.client, &self.data);
                return;
            }
        }

        self.state.update_status(&mut self.client);
        self.state.update_song(&mut self.client, &self.data);
    }

    pub fn toggle(&mut self) {
        if self.state.playing {
            if self.client.pause(true).is_ok() {
                self.state.update_status(&mut self.client);
            }
        } else if self.client.play().is_ok() {
            self.state.update_status(&mut self.client);
        }
    }

    pub fn next(&mut self) {
        if self.client.next().is_ok() {
            self.state.update_status(&mut self.client);
            self.state.update_song(&mut self.client, &self.data);
        }
    }

    pub fn previous(&mut self) {
        if self.client.prev().is_ok() {
            self.state.update_status(&mut self.client);
            self.state.update_song(&mut self.client, &self.data);
        }
    }

    pub fn volume_up(&mut self) {
        let new_volume = i8::min(100, self.state.volume + 10);
        if self.client.volume(new_volume).is_ok() {
            self.state.update_status(&mut self.client);
        }
    }

    pub fn volume_down(&mut self) {
        let new_volume = i8::max(0, self.state.volume - 10);
        if self.client.volume(new_volume).is_ok() {
            self.state.update_status(&mut self.client);
        }
    }

    pub fn shuffle(&mut self) {
        if self.client.random(!self.state.shuffle).is_ok() {
            self.state.update_status(&mut self.client);
        }
    }

    pub fn repeat(&mut self) {
        if self.client.repeat(!self.state.repeat).is_ok() {
            self.state.update_status(&mut self.client);
        }
    }

    pub fn clear(&mut self) {
        if self.client.clear().is_ok() {
            self.state.update_status(&mut self.client);
            self.state.update_song(&mut self.client, &self.data);
        }
    }

    pub fn seek_forward(&mut self) {
        if let Some(position) = self.state.position {
            if self.client.rewind(position + Duration::new(5, 0)).is_ok() {
                self.state.update_status(&mut self.client);
                self.state.update_song(&mut self.client, &self.data);
            }
        }
    }

    pub fn seek_backward(&mut self) {
        if let Some(position) = self.state.position {
            if self.client.rewind(position - Duration::new(5, 0)).is_ok() {
                self.state.update_status(&mut self.client);
                self.state.update_song(&mut self.client, &self.data);
            }
        }
    }
}
