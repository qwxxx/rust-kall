use crate::config;
use reqwest::Response;
use serde_json::Value;
use md5;
use hex;
use anyhow::Result;
pub struct SharkScope {
    appname: String,
    username: String,
    hash: String,
    user_agent: String
}

impl SharkScope {
    pub fn new(config:config::SharkScope) -> Self {
        let password_hash = md5::compute(&config.password);
        let password_hash = hex::encode(password_hash.as_ref());
        let hash = format!("{:?}", md5::compute(format!("{}{}", password_hash, config.appkey)));

        SharkScope {
            appname:config.appname,
            username:config.username,
            hash,
            user_agent: String::from("npm-sharkscope")
        }
    }

    fn url(&self) -> String {
        format!("https://www.sharkscope.com/api/{}/networks/WPN/activeTournaments?Filter=Entrants:1~*;Stake:USD10~250;Type:ST,6MX;Class:SNG", self.appname)
    }

    fn headers(&self) -> reqwest::header::HeaderMap {
        let mut headers = reqwest::header::HeaderMap::new();
        headers.insert(reqwest::header::USER_AGENT, self.user_agent.parse().unwrap());
        headers.insert("Accept", "application/json".parse().unwrap());
        headers.insert("Password", self.hash.parse().unwrap());
        headers.insert("Username", self.username.parse().unwrap());
        headers
    }

    pub async fn get_tournaments(&mut self) -> Result<(Vec<String>, u8)> {
        let client = reqwest::Client::new();
        let response: Response = client.get(self.url()).headers(self.headers()).send().await?;
        let json: Value = response.json().await?;

        let tournaments = &json["Response"]["RegisteringTournamentsResponse"]["RegisteringTournaments"]["RegisteringTournament"];
        let mut tournaments_ids: Vec<String> = Vec::new();
        let mut tournaments_count: u8 = 0;

        if tournaments.is_array() {
            if let Some(tournaments) = tournaments.as_array() {
                for tournament in tournaments {
                    if let Some(id) = tournament["@id"].as_str() {
                        if tournament["@currentEntrants"] == "5" {
                            tournaments_ids.push(String::from(id));
                        }

                        tournaments_count += 1;
                    }
                }
            }
        } else if tournaments.is_object() {
            if let Some(id) = tournaments["@id"].as_str() {
                if tournaments["@currentEntrants"] == "5" {
                    tournaments_ids.push(String::from(id));
                }
                
                tournaments_count = 1;
            }
        }
        
        Ok((tournaments_ids, tournaments_count))
    }
}