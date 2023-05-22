#[derive(serde::Deserialize,Clone)]
pub struct SharkScope {
    pub username: String,
    pub password: String,
    pub appname: String,
    pub appkey: String,
}

#[derive(serde::Deserialize,Clone)]
pub struct Config {
    pub sharkscope: SharkScope,
    pub db_conn: String,
}