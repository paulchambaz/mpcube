//! This module is used to parse the mpd library data into coherent structures to be used by the
//! rest of the program
use mpd::State;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::File;
use std::io::{Read, Write};
use std::time::Duration;

/// Represents the current state of the mpd server
#[derive(Debug)]
pub struct StateData {
    /// If it is playing
    pub playing: bool,
    /// What album (if any) is playing
    pub album_id: Option<usize>,
    /// What song (if any) is playing
    pub title_id: Option<usize>,
    /// What's the position in the song (if any) playing
    pub position: Option<Duration>,
    /// What is the volume for mpd
    pub volume: i8,
    /// Whether we are in shuffle mode
    pub shuffle: bool,
    /// Whether we are in repeat mode
    pub repeat: bool,
}

/// Represents the library data from mpd once organised
#[derive(Debug, Serialize, Deserialize)]
pub struct MusicData {
    /// The list of albums in the library
    pub albums: Vec<Album>,
}

/// Represents a given album in the mpd library
#[derive(Debug, Serialize, Deserialize)]
pub struct Album {
    /// The artist of the album
    pub artist: String,
    /// The name of the album
    pub album: String,
    /// The date for the release of the album
    pub date: i32,
    /// The list of songs in the album
    pub songs: Vec<Song>,
}

/// Represents a given song in the mpd library
#[derive(Debug, Serialize, Deserialize)]
pub struct Song {
    /// The uri in the mpd library for access to the song
    pub uri: String,
    /// The title of the song
    pub title: String,
    /// The track number for the song
    pub track: u32,
    /// The duration of the song
    pub duration: Duration,
}

/// Represents the raw output from the mpd library
#[derive(Debug)]
struct RawMusicEntry {
    /// The uri in the mpd library for access to the song
    uri: String,
    /// The title of the song
    title: String,
    /// The artist of the song
    artist: String,
    /// The album of the song
    album: String,
    /// The date of the song
    date: i32,
    /// The track number of the song
    track: u32,
    /// The duration of the song
    duration: Duration,
}

impl StateData {
    /// Creates a new representation of the state of the mpd server
    /// 
    /// - `client` - The mpd client used to get state info from
    /// - `music_data` - The representation of the mpd library
    ///
    /// Returns a new up-to-date mpd state representation
    pub fn new(client: &mut mpd::Client, music_data: &MusicData) -> Self {
        // Get status and current song by connecting to mpd
        let status = client.status().expect("Could not connect to mpd");
        let song = client.currentsong();

        // Get info from the song by finding it in the database
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

        // Return album_id and title_id to the correct format
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

    /// Updates the current state of the mpd state representation
    /// This function is very similar to StateData::new(), but does not create a new object
    ///
    /// - `client` - The mpd client used to get state info from
    /// - `music_data` - The representation of the mpd library
    pub fn update(&mut self, client: &mut mpd::Client, music_data: &MusicData) {
        let status = client.status().expect("Could not connect to mpd");
        self.playing = status.state == State::Play;
        self.position = status.elapsed;
        self.volume = status.volume;
        self.shuffle = status.random;
        self.repeat = status.repeat;

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
    /// Creates a new mpd library representation
    ///
    /// - `client` - The mpd client used to get library info from
    pub fn new(client: &mut mpd::Client) -> Self {
        // First we do a heavy call to mpd to get all the songs in the database and their metadata
        let songs = client.listallinfo().expect("Could not list all songs");

        let mut albums = Vec::new();
        let mut current_album: Option<Album> = None;

        // Then for each song
        for song in songs {
            // First we parse them to a valid entry
            let entry = RawMusicEntry::new(song);
            match current_album.as_mut() {
                // If we have a given album for them already, we add it to the list of songs
                Some(album) if album.album == entry.album => {
                    // TODO: it would be good to update the author to the shortest author on the
                    // track list, given that it most often is the author (without featuring) and
                    // is therefore most often the real author name
                    album.songs.push(Song {
                        uri: entry.uri,
                        title: entry.title,
                        track: entry.track,
                        duration: entry.duration,
                    });
                }
                _ => {
                    // If we already have an album ready, add it on the list of albums
                    if let Some(album) = current_album.take() {
                        albums.push(album);
                    }
                    // If not, we create it
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

        // Making sure the last album is correctly added
        if let Some(album) = current_album {
            albums.push(album);
        }

        // Finally we create the actual object and sort it
        let mut music_data = MusicData { albums };
        music_data.sort();

        music_data
    }

    /// Sorts the list of albums
    fn sort(&mut self) {
        // First we sort the list of albums
        self.albums.sort_by(Album::cmp);
        // Then for each album, we sort the songs themselves
        for album in self.albums.iter_mut() {
            album.sort();
        }
    }

    /// Used to recreate the full data from cache directly instead of re-doing the heavy operation
    /// of MusicData::new()
    ///
    /// - `cache_path`: Path to the cache
    pub fn from_cache(cache_path: &str) -> MusicData {
        let mut file = File::open(cache_path).expect("Could not load cache");
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer).expect("Could not read file");
        bincode::deserialize(&buffer[..]).expect("Deserialization failed")
    }

    /// Saves the MusicData to a cache file for later usage
    ///
    /// - `cache_path`: Path to the cache
    pub fn save_cache(&self, cache_path: &str) {
        let serialized_data = bincode::serialize(self).expect("Serialization failed");
        let mut file = File::create(cache_path).expect("Could not save create cache file");
        file.write_all(&serialized_data)
            .expect("Could not save cache");
    }

    /// Given an album title and a uri, find the album id and title id in the album and title lists
    /// of a given song
    ///
    /// - `album_title`: The name of the album being searched
    /// - `uri`: The unique identifier of the song in the mpd database
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
    /// Used to sort the album list, compares two albums
    ///
    /// - `a`: The first album
    /// - `b`: The second album
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

    /// Sorts the list of songs from a given album
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
    /// Creates a new RawMusicEntry from the output of the mpd listallinfo command
    ///
    /// - `song`: The mpd song
    fn new(song: mpd::Song) -> Self {
        // First we get the unique uri identifier
        let uri = song.file.clone();

        // Then we get the tags as they contain a lot of useful information
        let tags: HashMap<_, _> = song
            .tags
            .into_iter()
            .map(|(k, v)| (k.to_string(), v.to_string()))
            .collect();

        // We get the title, if no title is found, we reuse the full uri as a fallback
        let title = song.title.unwrap_or(uri.clone());

        // We get the artist, if no artist is found, we use 'Unknown Artist'
        let artist = song.artist.unwrap_or("Unknown Artist".to_string());

        // We get the album, if no album is found, we use 'Unknown Album'
        let album = tags
            .get("ALBUM")
            .or(tags.get("Album"))
            .or(tags.get("album"))
            .map(|s| s.to_string())
            .unwrap_or("Unknown Album".to_string());

        // We get the date, if no date is found, we assume a date of 0
        let date = tags
            .get("DATA")
            .or(tags.get("Date"))
            .or(tags.get("date"))
            .and_then(|s| {
                if let Some(year) = s.split('-').next() {
                    year.parse::<i32>().ok()
                } else {
                    s.parse::<i32>().ok()
                }
            })
            .unwrap_or(0);

        // We get the track number, if no track is found, we assume a track number of 0
        let track = tags
            .get("TRACKNUMBER")
            .or(tags.get("Track"))
            .or(tags.get("track"))
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(0);

        // Finally we get the duration of the song
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
