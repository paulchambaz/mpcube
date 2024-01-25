use mpd::State;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::File;
use std::io::{Read, Write};
use std::time::Duration;

pub struct Client {
    pub client: mpd::Client,
    pub data: Option<MusicData>,
    pub state: Option<StateData>,
}

#[derive(Debug)]
pub struct StateData {
    pub playing: bool,
    pub album_id: Option<usize>,
    pub title_id: Option<usize>,
    pub position: Option<Duration>,
    pub volume: i8,
    pub shuffle: bool,
    pub repeat: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MusicData {
    pub albums: Vec<Album>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Album {
    pub artist: String,
    pub album: String,
    pub date: i32,
    pub songs: Vec<Song>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Song {
    pub uri: String,
    pub title: String,
    pub track: u32,
    pub duration: Duration,
}

#[derive(Debug)]
struct RawMusicData {
    entries: Vec<RawMusicEntry>,
}

#[derive(Debug)]
struct RawMusicEntry {
    uri: String,
    title: String,
    artist: String,
    album: String,
    date: i32,
    track: u32,
    duration: Duration,
}

impl Album {
    pub fn cmp(a: &Album, b: &Album) -> std::cmp::Ordering {
        // Handle empty or missing artist names
        match (a.artist.is_empty(), b.artist.is_empty()) {
            (true, false) => return std::cmp::Ordering::Less,
            (false, true) => return std::cmp::Ordering::Greater,
            _ => {}
        };

        // Skip "The " at the beginning of artist names
        let a_artist = a.artist.strip_prefix("The ").unwrap_or(&a.artist);
        let b_artist = b.artist.strip_prefix("The ").unwrap_or(&b.artist);

        // Compare by artist
        match a_artist.cmp(b_artist) {
            std::cmp::Ordering::Equal => {}
            non_eq => return non_eq,
        }

        // Compare by date
        match a.date.cmp(&b.date) {
            std::cmp::Ordering::Equal => {}
            non_eq => return non_eq,
        }

        // Handle empty or missing album names
        match (a.album.is_empty(), b.album.is_empty()) {
            (true, false) => std::cmp::Ordering::Less,
            (false, true) => std::cmp::Ordering::Greater,
            _ => a.album.cmp(&b.album), // Compare by album
        }
    }

    pub fn sort(&mut self) {
        self.songs.sort_by(|a, b| {
            // Use track number to order
            match a.track.cmp(&b.track) {
                std::cmp::Ordering::Less => return std::cmp::Ordering::Less,
                std::cmp::Ordering::Greater => return std::cmp::Ordering::Greater,
                std::cmp::Ordering::Equal => {}
            };

            // Handle equal track uri by using uri
            a.uri.cmp(&b.uri)
        });
    }
}

impl StateData {
    pub fn update_status(&mut self, client: &mut mpd::Client) {
        let status = client.status().expect("Could not connect to mpd");
        self.playing = status.state == State::Play;
        self.position = status.elapsed;
        self.volume = status.volume;
        self.shuffle = status.random;
        self.repeat = status.repeat;
    }

    pub fn update_song(&mut self, client: &mut mpd::Client, music_data: &MusicData) {
        let song = client.currentsong();
        let id: Option<(usize, usize)> = if let Ok(Some(song)) = song {
            let album_opt = song
                .tags
                .iter()
                .find(|(key, _value)| key == "ALBUM" || key == "Album");

            if let Some((_key, album)) = album_opt {
                let album = album.to_string(); // Convert &String to String
                let uri = song.file.to_string(); // Convert &String to String
                music_data.find(album, uri)
            } else {
                None
            }
        } else {
            None
        };

        let (album_id, title_id) = if let Some((album_id, title_id)) = id {
            (Some(album_id), Some(title_id))
        } else {
            (None, None)
        };

        self.album_id = album_id;
        self.title_id = title_id;
    }

    pub fn from_client(client: &mut mpd::Client, music_data: &MusicData) -> StateData {
        let status = client.status().expect("Could not connect to mpd");
        let song = client.currentsong();

        let id: Option<(usize, usize)> = if let Ok(Some(song)) = song {
            let album_opt = song
                .tags
                .iter()
                .find(|(key, _value)| key == "ALBUM" || key == "Album");

            if let Some((_key, album)) = album_opt {
                let album = album.to_string(); // Convert &String to String
                let uri = song.file.to_string(); // Convert &String to String
                music_data.find(album, uri)
            } else {
                None
            }
        } else {
            None
        };

        let (album_id, title_id) = if let Some((album_id, title_id)) = id {
            (Some(album_id), Some(title_id))
        } else {
            (None, None)
        };

        StateData {
            playing: status.state == State::Play,
            album_id,
            title_id,
            position: status.elapsed,
            volume: status.volume,
            shuffle: status.random,
            repeat: status.repeat,
        }
    }
}

impl MusicData {
    // i need to add the path for this function
    fn from_cache(path: &str) -> MusicData {
        let mut file = File::open(path).expect("Could not load cache");
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer).expect("Could not read file");
        bincode::deserialize(&buffer[..]).expect("Deserialization failed")
    }

    // i need to add the path for this function
    fn save_cache(&mut self, path: &str) {
        let serialized_data = bincode::serialize(self).expect("Serialization failed");
        let mut file = File::create(path).expect("Could not save create cache file");
        file.write_all(&serialized_data)
            .expect("Could not save cache");
    }

    fn from_raw(raw_music_data: RawMusicData) -> MusicData {
        let mut albums = Vec::new();
        let mut current_album: Option<Album> = None;

        for entry in raw_music_data.entries {
            match current_album.as_mut() {
                Some(album) if album.album == entry.album => {
                    album.songs.push(Song {
                        uri: entry.uri,
                        title: entry.title,
                        track: entry.track,
                        duration: entry.duration,
                    });
                }
                _ => {
                    if let Some(album) = current_album.take() {
                        albums.push(album);
                    }
                    current_album = Some(Album {
                        artist: entry.artist,
                        album: entry.album,
                        date: entry.date,
                        songs: vec![Song {
                            uri: entry.uri,
                            title: entry.title,
                            track: entry.track,
                            duration: entry.duration,
                        }],
                    });
                }
            }
        }

        if let Some(album) = current_album {
            albums.push(album);
        }

        Self { albums }
    }

    pub fn sort(&mut self) {
        self.albums.sort_by(Album::cmp);
        for album in self.albums.iter_mut() {
            album.sort();
        }
    }

    pub fn find(&self, album_title: String, uri: String) -> Option<(usize, usize)> {
        for (album_idx, album) in self.albums.iter().enumerate() {
            if album.album == album_title {
                if let Some(title_idx) = album.songs.iter().position(|song| song.uri == uri) {
                    return Some((album_idx, title_idx));
                } else {
                    return None;
                }
            }
        }
        None
    }
}

impl RawMusicData {
    pub fn from_client(client: &mut mpd::Client) -> RawMusicData {
        let songs = client.listall().expect("Could not list all songs");

        let mut out = Vec::new();

        for song in songs {
            let uri = song.file.clone();
            println!("Song: {:?}", song);

            // let full_song = client.find("file", uri).expect("Could not fetch song details");

            let tags: HashMap<_, _> = client
                .readcomments(song.clone())
                .expect("Could not read comments from song")
                .flatten()
                .collect();
            println!("Tags: {:?}", tags);

            // TODO: i dont think i have to use unwrap or else and or_else just or and unwrap_or
            let title = tags
                .get("TITLE")
                .or_else(|| tags.get("title"))
                .map(|s| s.to_string())
                .unwrap_or_else(|| uri.clone());

            let artist = tags
                .get("ARTIST")
                .or_else(|| tags.get("artist"))
                .map(|s| s.to_string())
                .unwrap_or_else(|| "Unknown Artist".to_string());

            let album = tags
                .get("ALBUM")
                .or_else(|| tags.get("album"))
                .map(|s| s.to_string())
                .unwrap_or_else(|| "Unknown Album".to_string());

            let date = tags
                .get("DATE")
                .or_else(|| tags.get("date"))
                .and_then(|s| {
                    if let Some(year) = s.split('-').next() {
                        year.parse::<i32>().ok()
                    } else {
                        s.parse::<i32>().ok()
                    }
                })
                .unwrap_or(0);

            let track = tags
                .get("TRACKNUMBER")
                .or_else(|| tags.get("track"))
                .and_then(|s| s.parse::<u32>().ok())
                .unwrap_or(0);

            let duration = Duration::new(0, 0);
            // TODO: implement getting the actual duration of the song
            // let duration = Duration::new(0, 0);
            // let duration = song.duration;
            // println!("Title: {:?}", title);
            // println!("Duration: {:?}", duration);
            // println!();

            out.push(RawMusicEntry {
                uri,
                title,
                artist,
                album,
                date,
                track,
                duration,
            });
        }

        RawMusicData { entries: out }
    }
}

impl Client {
    pub fn new(address: &str, port: u16) -> Client {
        let client = mpd::Client::connect(format!("{}:{}", address, port))
            .expect("Could not connect to mpd");
        Self {
            client,
            data: None,
            state: None,
        }
    }

    pub fn init_sync(&mut self, cache_path: &str) {
        if std::path::Path::new(cache_path).exists() {
            let music_data = MusicData::from_cache(cache_path);
            self.data = Some(music_data);
            self.sync_state();
        } else {
            self.full_sync(cache_path);
        }
    }

    pub fn full_sync(&mut self, cache_path: &str) {
        let raw_music_data = RawMusicData::from_client(&mut self.client);
        let mut music_data = MusicData::from_raw(raw_music_data);

        music_data.sort();

        music_data.save_cache(cache_path);
        self.data = Some(music_data);
        self.sync_state();
    }

    pub fn sync_state(&mut self) {
        if let Some(data) = &self.data {
            let state_data = StateData::from_client(&mut self.client, data);
            self.state = Some(state_data);
        }
    }

    pub fn start_album(&mut self, album_id: usize) {
        // clear music playing
        // add all the songs from the album to the queue
        // start playing
    }

    pub fn start_title(&mut self, album_id: usize, title_id: usize) {
        // clear music playing
        // add all the songs from the album to the queue
        // start playing
        // skip to title
    }

    pub fn toggle(&mut self) {
        if let Some(state) = &mut self.state {
            if self.client.toggle_pause().is_ok() {
                state.update_status(&mut self.client);
            }
        }
    }

    pub fn next(&mut self) {
        // if playing { next }
    }

    pub fn previous(&mut self) {
        // if playing { previous }
    }

    pub fn volume_up(&mut self) {
        if let Some(state) = &mut self.state {
            let new_volume = i8::min(100, state.volume + 10);
            if self.client.volume(new_volume).is_ok() {
                state.update_status(&mut self.client);
            }
        }
    }

    pub fn volume_down(&mut self) {
        if let Some(state) = &mut self.state {
            let new_volume = i8::max(0, state.volume - 10);
            if self.client.volume(new_volume).is_ok() {
                state.update_status(&mut self.client);
            }
        }
    }

    pub fn shuffle(&mut self) {
        if let Some(state) = &mut self.state {
            if self.client.random(!state.shuffle).is_ok() {
                state.update_status(&mut self.client);
            }
        }
    }

    pub fn repeat(&mut self) {
        if let Some(state) = &mut self.state {
            if self.client.repeat(!state.repeat).is_ok() {
                state.update_status(&mut self.client);
            }
        }
    }

    pub fn clear(&mut self) {
        // if playing { clear }
    }

    pub fn seek_forward(&mut self) {
        // if playing { run_seek_current +seek_incr }
    }

    pub fn seek_backward(&mut self) {
        // if playing { run_seek_current -seek_incr }
    }
}
