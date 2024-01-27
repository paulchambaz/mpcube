use rand::Rng;
use crate::music::music_data::{MusicData, StateData};
use std::{time::Duration, net::IpAddr, path::PathBuf};

pub struct Client {
    pub client: mpd::Client,
    pub data: MusicData,
    pub state: StateData,
    cache_path: String,
}

impl Client {
    pub fn new(address: IpAddr, port: u16, cache_path: PathBuf) -> Client {
        let mut client = mpd::Client::connect(format!("{}:{}", address, port))
            .expect("Could not connect to mpd");

        let cache_path = cache_path.to_str().expect("Could not get cache path").to_string();

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

    pub fn sync(&mut self) {
        self.state.update(&mut self.client, &self.data);
    }

    pub fn start_album(&mut self, album_id: usize) {
        if self.client.clear().is_ok() {
            self.sync();
        }

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
            self.sync();
        }
    }

    pub fn start_title(&mut self, album_id: usize, title_id: usize) {
        if self.client.clear().is_ok() {
            self.sync();
        }

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
                self.sync();
                return;
            }
        }

        self.sync();
    }

    pub fn toggle(&mut self) {
        if self.state.playing {
            if self.client.pause(true).is_ok() {
                self.sync();
            }
        } else if self.client.play().is_ok() {
            self.sync();
        }
    }

    pub fn next(&mut self) {
        if self.client.next().is_ok() {
            self.sync();
        }
    }

    pub fn previous(&mut self) {
        if self.client.prev().is_ok() {
            self.sync();
        }
    }

    pub fn volume_up(&mut self) {
        let new_volume = i8::min(100, self.state.volume + 10);
        if self.client.volume(new_volume).is_ok() {
            self.sync();
        }
    }

    pub fn volume_down(&mut self) {
        let new_volume = i8::max(0, self.state.volume - 10);
        if self.client.volume(new_volume).is_ok() {
            self.sync();
        }
    }

    pub fn shuffle(&mut self) {
        if self.client.random(!self.state.shuffle).is_ok() {
            self.sync();
        }
    }

    pub fn repeat(&mut self) {
        if self.client.repeat(!self.state.repeat).is_ok() {
            self.sync();
        }
    }

    pub fn clear(&mut self) {
        if self.client.clear().is_ok() {
            self.sync();
        }
    }

    pub fn seek_forward(&mut self) {
        if let Some(position) = self.state.position {
            if self.client.rewind(position + Duration::new(5, 0)).is_ok() {
                self.sync();
            }
        }
    }

    pub fn seek_backward(&mut self) {
        if let Some(position) = self.state.position {
            if self.client.rewind(position - Duration::new(5, 0)).is_ok() {
                self.sync();
            }
        }
    }

    pub fn random(&mut self) {
        let mut rng = rand::thread_rng();
        let mut album_id;
        loop {
            album_id = rng.gen_range(0..self.data.albums.len());  
            if let Some(album_playing) = self.state.album_id {
                if album_id != album_playing {
                    break;
                }
            } else {
                break;
            }
        }
        self.start_album(album_id);
    }
}
