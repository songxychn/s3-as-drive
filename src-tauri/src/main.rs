// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod minio;

use std::env::args;
use std::fs;
use std::fs::File;
use std::io::{Read, Write};
use std::path::PathBuf;
use ::minio::s3::args::{ListObjectsArgs, ListObjectsV2Args};
use ::minio::s3::client::Client;
use ::minio::s3::creds::StaticProvider;
use ::minio::s3::http::BaseUrl;
use ::minio::s3::response::ListObjectsV2Response;
use dirs::home_dir;
use serde::{Deserialize, Serialize};
use serde_json::from_str;
use crate::minio::{GetFileListResp, S3Config};

#[tauri::command]
fn get_s3_config() -> Result<S3Config, String> {
    let mut path = PathBuf::from(get_base_dir().unwrap());
    path.push("s3-config.json");

    // 打开文件
    let mut file = File::open(path).expect("Failed to open file");
    let mut contents = String::new();
    file.read_to_string(&mut contents).expect("Failed to read file");

    // 解析 JSON 数据
    let config: S3Config = from_str(&contents).expect("Failed to parse JSON");
    return Ok(config);
}

#[tauri::command]
fn update_s3_config(s3_config: S3Config) -> Result<(), String> {
    let mut path = PathBuf::from(get_base_dir().unwrap());
    path.push("s3-config.json");
    let mut file = File::create(path).unwrap();
    file.write_all(serde_json::to_string(&s3_config).unwrap().as_bytes()).unwrap();
    return Ok(());
}

fn get_base_dir() -> Result<PathBuf, String> {
    let home_dir = home_dir().expect("Failed to get home directory");
    let mut base_dir = PathBuf::from(home_dir);
    base_dir.push("s3-as-drive/");
    // 如果不存在就创建
    fs::create_dir_all(&base_dir).expect("Failed to create base directory");
    return Ok(base_dir);
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![get_s3_config, update_s3_config])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}