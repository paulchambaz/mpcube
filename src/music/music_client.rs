//! This module is used to define the connection with mpd. It handles the connection and the
//! commands to mpd

use crate::music::music_data::{MusicData, StateData};
use rand::Rng;
use std::{net::IpAddr, path::PathBuf, time::Duration};

/// The client used to connect to mpd
pub struct Client {
    /// The actual connection with mpd
    pub client: mpd::Client,
    /// The representation of the music library
    pub data: MusicData,
    /// The representation of the current mpd state
    pub state: StateData,
    /// The path to the cache - used to save and load
    cache_path: String,
}

impl Client {
    /// Creates a new client give the address of the mpd connection, the port and a cache path
    ///
    /// - `address`: The address of the mpd server
    /// - `port`: The port of the mpd server
    /// - `cache_path`: The path to the cache
    ///
    /// Returns a valid client
    pub fn new(address: IpAddr, port: u16, cache_path: PathBuf) -> Client {
        // First we connect to the mpd client
        let mut client = mpd::Client::connect(format!("{}:{}", address, port))
            .expect("Could not connect to mpd");

        // Then we get the full path
        let cache_path = cache_path
            .to_str()
            .expect("Could not get cache path")
            .to_string();

        // If the path exists (a cache is already present
        let data = if std::path::Path::new(&cache_path).exists() {
            // We load the music data from cache
            MusicData::from_cache(&cache_path)
        } else {
            // We load the music data from the mpd connection
            let music_data = MusicData::new(&mut client);
            // Then we save it to cache for future usage
            music_data.save_cache(&cache_path);
            music_data
        };

        // Finally we load the state data from the mpd connection and the music data
        let state = StateData::new(&mut client, &data);

        Self {
            client,
            data,
            state,
            cache_path,
        }
    }

    /// Does a full sync by telling mpd to update its database, then reconstructing the library
    /// representation, saving it to cache for future usage and reloading the state
    pub fn full_sync(&mut self) {
        if self.client.update().is_ok() {
            self.data = MusicData::new(&mut self.client);
            self.data.save_cache(&self.cache_path);
            self.state = StateData::new(&mut self.client, &self.data);
        }
    }

    /// Reloads the current mpd state
    pub fn sync(&mut self) {
        self.state.update(&mut self.client, &self.data);
    }

    /// Starts an album given an id. Clears the current queue, then replace it with the full album
    /// before starting the new album
    ///
    /// - `album_id`: The id of the album to start
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

    /// Starts an title given an album id and song id. Clears the current queue, then replace it
    /// with the full album before starting the new title
    ///
    /// - `album_id`: The id of the album to start
    /// - `title_id`: The id of the title to start
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

    /// Toggles play and pause for mpd
    pub fn toggle(&mut self) {
        if self.state.playing {
            if self.client.pause(true).is_ok() {
                self.sync();
            }
        } else if self.client.play().is_ok() {
            self.sync();
        }
    }

    /// Plays the next track in the mpd queue
    pub fn next(&mut self) {
        if self.client.next().is_ok() {
            self.sync();
        }
    }

    /// Plays the previous track in the mpd queue
    pub fn previous(&mut self) {
        if self.client.prev().is_ok() {
            self.sync();
        }
    }

    /// Increases the volume for mpd
    pub fn volume_up(&mut self) {
        let new_volume = i8::min(100, self.state.volume + 10);
        if self.client.volume(new_volume).is_ok() {
            self.sync();
        }
    }

    /// Decreases the volume for mpd
    pub fn volume_down(&mut self) {
        let new_volume = i8::max(0, self.state.volume - 10);
        if self.client.volume(new_volume).is_ok() {
            self.sync();
        }
    }

    /// Toggles shuffle mode
    pub fn shuffle(&mut self) {
        if self.client.random(!self.state.shuffle).is_ok() {
            self.sync();
        }
    }

    /// Toggles repeat mode
    pub fn repeat(&mut self) {
        if self.client.repeat(!self.state.repeat).is_ok() {
            self.sync();
        }
    }

    /// Clears the mpd queue
    pub fn clear(&mut self) {
        if self.client.clear().is_ok() {
            self.sync();
        }
    }

    /// Seeks forward in the current mpd song
    pub fn seek_forward(&mut self) {
        if let Some(position) = self.state.position {
            if self.client.rewind(position + Duration::new(5, 0)).is_ok() {
                self.sync();
            }
        }
    }

    /// Seeks backward in the current mpd song
    pub fn seek_backward(&mut self) {
        if let Some(position) = self.state.position {
            let target = if position < Duration::new(5, 0) {
                Duration::new(0, 0)
            } else {
                position - Duration::new(5, 0)
            };
            if self.client.rewind(target).is_ok() {
                self.sync();
            }
        }
    }

    /// Selects a random new album and plays it
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
