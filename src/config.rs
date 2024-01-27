use serde::Deserialize;
use std::{path::PathBuf, net::IpAddr, str::FromStr};
use clap::Parser;

use directories::ProjectDirs;
use std::fs;

#[derive(Debug)]
pub struct Config {
    pub mpd_host: IpAddr,
    pub mpd_port: u16,
    pub cache: PathBuf,
}

#[derive(Deserialize, Debug)]
struct TomlConfig {
    mpd_host: Option<IpAddr>,
    mpd_port: Option<u16>,
    cache: Option<PathBuf>,
}

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
#[clap(name = "mpcube")]
struct Args {
    /// Ip address of the mpd host [default: 127.0.0.1]
    #[clap(long)]
    mpd_host: Option<IpAddr>,

    /// Port number of the mpd host [default: 6600]
    #[clap(long)]
    mpd_port: Option<u16>,

    /// Cache file location [default ~/.cache/mpcube/cache]
    #[clap(long)]
    cache: Option<PathBuf>,
}

pub fn load_config() -> Config {
    let project_dirs = ProjectDirs::from("", "", "mpcube").expect("Could not get standard directories");

    let mut mpd_host: IpAddr = IpAddr::from_str("127.0.0.1").expect("Could not parse default ip address");
    let mut mpd_port: u16 = 6600;
    let mut cache: PathBuf = project_dirs.cache_dir().join("cache");

    let config_file = project_dirs.config_dir().join("config.toml");

    if config_file.exists() {
        let content = fs::read_to_string(config_file).expect("Could not read config file");
        let config: TomlConfig = toml::from_str(&content).expect("Could not parse config file");

        mpd_host = config.mpd_host.unwrap_or(mpd_host);
        mpd_port = config.mpd_port.unwrap_or(mpd_port);
        cache = config.cache.unwrap_or(cache);
    }

    let args = Args::parse();

    mpd_host = args.mpd_host.unwrap_or(mpd_host);
    mpd_port = args.mpd_port.unwrap_or(mpd_port);
    cache = args.cache.unwrap_or(cache);

    if let Some(parent) = cache.parent() {
        fs::create_dir_all(parent).expect("Could not create cache directory");
    }

    Config {
        mpd_host,
        mpd_port,
        cache,
    }
}
