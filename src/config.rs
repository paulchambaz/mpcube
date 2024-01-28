//! This module is used to define the configuration that the user can bring to
//! mpcube. Users can specify their config in either a config file or with the
//! argument in the command line

use clap::Parser;
use serde::Deserialize;
use std::{net::IpAddr, path::PathBuf, str::FromStr};

use directories::ProjectDirs;
use std::fs;

/// Used to represent the actual configuration of the user
#[derive(Debug)]
pub struct Config {
    /// The Ip address of the mpd server
    pub mpd_host: IpAddr,
    /// The port of the mpd server
    pub mpd_port: u16,
    /// The path to the cache file
    pub cache: PathBuf,
}

/// Used to represent the toml configuration the user may have
#[derive(Deserialize, Debug)]
struct TomlConfig {
    /// The Ip address of the mpd server
    mpd_host: Option<IpAddr>,
    /// The port of the mpd server
    mpd_port: Option<u16>,
    /// The path to the cache file
    cache: Option<PathBuf>,
}

/// Used to represent the argument which can be parsed through cli
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

    /// Cache file location [default: ~/.cache/mpcube/cache]
    #[clap(long)]
    cache: Option<PathBuf>,
}

/// Used to load the configuration of the user
/// Specifies the default
/// Overwrites them with a toml if it is present
/// Overwrites them with the arguments passed through cli if it is present
pub fn load_config() -> Config {
    let project_dirs =
        ProjectDirs::from("", "", "mpcube").expect("Could not get standard directories");

    let mut mpd_host: IpAddr =
        IpAddr::from_str("127.0.0.1").expect("Could not parse default ip address");
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
