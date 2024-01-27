use std::time::Duration;
use std::collections::HashMap;
use std::fs::File;
use serde::{Deserialize, Serialize};
use mpd::State;
use std::io::{Read, Write};

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
struct RawMusicEntry {
    uri: String,
    title: String,
    artist: String,
    album: String,
    date: i32,
    track: u32,
    duration: Duration,
}


impl StateData {
    pub fn new(client: &mut mpd::Client, music_data: &MusicData) -> Self {
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
}

impl MusicData {
    pub fn new(client: &mut mpd::Client) -> Self {
        let songs = client.listallinfo().expect("Could not list all songs");

        let mut albums = Vec::new();
        let mut current_album: Option<Album> = None;

        for song in songs {
            let entry = RawMusicEntry::new(song);
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

        let mut music_data = MusicData { albums };
        music_data.sort();

        music_data
    }

    fn sort(&mut self) {
        self.albums.sort_by(Album::cmp);
        for album in self.albums.iter_mut() {
            album.sort();
        }
    }

    pub fn from_cache(cache_path: &str) -> MusicData {
        let mut file = File::open(cache_path).expect("Could not load cache");
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer).expect("Could not read file");
        bincode::deserialize(&buffer[..]).expect("Deserialization failed")
    }

    pub fn save_cache(&self, path: &str) {
        let serialized_data = bincode::serialize(self).expect("Serialization failed");
        let mut file = File::create(path).expect("Could not save create cache file");
        file.write_all(&serialized_data)
            .expect("Could not save cache");
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

impl RawMusicEntry {
    fn new(song: mpd::Song) -> Self {
        println!("{:?}", song);
        let uri = song.file.clone();

        let tags: HashMap<_, _> = song
            .tags
            .into_iter()
            .map(|(k, v)| (k.to_string(), v.to_string()))
            .collect();

        let title = song.title.unwrap_or(uri.clone());

        let artist = song.artist.unwrap_or("Unknown Artist".to_string());

        let album = tags
            .get("ALBUM")
            .or(tags.get("Album"))
            .map(|s| s.to_string())
            .unwrap_or("Unknown Album".to_string());

        let date = tags
            .get("DATA")
            .or(tags.get("date"))
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
            .or(tags.get("track"))
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(0);

        let duration = song.duration.unwrap_or(Duration::new(0, 0));

        RawMusicEntry {
            uri,
            title,
            artist,
            album,
            date,
            track,
            duration,
        }
    }
}
