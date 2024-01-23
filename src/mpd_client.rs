use std::collections::HashMap;

pub struct Client {
    pub client: mpd::Client,
    pub data: Option<MusicData>,
    pub state: Option<StateData>,
}

#[derive(Debug)]
pub struct StateData {
    pub playing: bool,
    pub id: Option<usize>,
    pub duration: Option<u32>,
    pub volume: u32,
    pub shuffle: bool,
    pub repeat: bool,
}

#[derive(Debug)]
pub struct MusicData {
    pub albums: Vec<Album>,
}

#[derive(Debug)]
pub struct Album {
    pub artist: String,
    pub album: String,
    pub date: i32,
    pub songs: Vec<Song>,
}

#[derive(Debug)]
pub struct Song {
    pub id: String,
    pub title: String,
    pub track: u32,
    pub duration: u32,
}

#[derive(Debug)]
struct RawMusicData {
    entries: Vec<RawMusicEntry>,
}

#[derive(Debug)]
struct RawMusicEntry {
    id: String,
    title: String,
    artist: String,
    album: String,
    date: i32,
    track: u32,
    duration: u32,
}

impl Album {
    pub fn sort(&mut self) {
        self.songs.sort_by(|a, b| {
            // Use track number to order
            match a.track.cmp(&b.track) {
                std::cmp::Ordering::Less => return std::cmp::Ordering::Less,
                std::cmp::Ordering::Greater => return std::cmp::Ordering::Greater,
                std::cmp::Ordering::Equal => {}
            };

            // Handle equal track id by using id
            a.id.cmp(&b.id)
        });
    }
}

impl MusicData {
    fn from_raw(raw_music_data: RawMusicData) -> MusicData {
        let mut albums = Vec::new();
        let mut current_album: Option<Album> = None;

        for entry in raw_music_data.entries {
            match current_album.as_mut() {
                Some(album) if album.album == entry.album => {
                    album.songs.push(Song {
                        id: entry.id,
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
                            id: entry.id,
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
        self.albums.sort_by(|a, b| {
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
        });
        for album in self.albums.iter_mut() {
            album.sort();
        }
    }
}

impl RawMusicData {
    pub fn from_client(client: &mut mpd::Client) -> RawMusicData {
        let songs = client.listall().expect("Could not list all songs");

        let mut out = Vec::new();

        for song in songs {
            let id = song.file.clone();

            let tags: HashMap<_, _> = client
                .readcomments(song)
                .expect("Could not read comments from song")
                .flatten()
                .collect();

            // println!("{:?}", tags);

            let title = tags
                .get("TITLE")
                .or_else(|| tags.get("title"))
                .map(|s| s.to_string())
                .unwrap_or_else(|| id.clone());

            let artist = tags
                .get("ALBUMARTIST")
                .or_else(|| tags.get("ARTIST"))
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

            // TODO: implement getting the actual duration of the song
            let duration = 215535;

            out.push(RawMusicEntry {
                id,
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

    pub fn full_sync(&mut self) {
        let raw_music_data = RawMusicData::from_client(&mut self.client);
        let mut music_data = MusicData::from_raw(raw_music_data);
        music_data.sort();
        self.data = Some(music_data);
        self.sync_state();
    }

    pub fn sync_state(&mut self) {
        // TODO: implement proper status update
        self.state = Some(StateData {
            playing: false,
            id: Some(0),
            duration: None,
            volume: 0,
            shuffle: false,
            repeat: false,
        });
    }

    pub fn start_album(&mut self, album_id: u32) {
        // clear music playing
        // add all the songs from the album to the queue
        // start playing
    }

    pub fn start_title(&mut self, album_id: u32, title_id: u32) {
        // clear music playing
        // add all the songs from the album to the queue
        // start playing
        // skip to title
    }

    pub fn toggle(&mut self) {
        // if playing { pause } else { play }
    }

    pub fn next(&mut self) {
        // if playing { next }
    }

    pub fn previous(&mut self) {
        // if playing { previous }
    }

    pub fn volume_up(&mut self) {
        // if playing { volume_up }
    }

    pub fn volume_down(&mut self) {
        // if playing { volume_down }
    }

    pub fn shuffle(&mut self) {
        // run_random shuffle
    }

    pub fn repeat(&mut self) {
        // run_repeat repeat
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
