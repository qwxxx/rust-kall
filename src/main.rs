use std::time::Duration;
use anyhow::Result;

use sqlx::{postgres::PgPoolOptions};
mod config;
mod sharkscope;

async fn init_program()->Result<config::Config>{
    dotenv::dotenv().ok();
    let filter=tracing_subscriber::EnvFilter::from_env("RUST_LOG_DISABLE_CRATES");
    tracing_subscriber::fmt().with_env_filter(filter).init();

    let f = std::fs::File::open("config.yaml")?;
    let config:config::Config= serde_yaml::from_reader(f)?;

    Ok(config)
}

#[tokio::main]
async fn main() -> Result<()> {
    
    let config=init_program().await?;
    let mut shark_scope = sharkscope::SharkScope::new(config.sharkscope);

    let db=PgPoolOptions::new()
    .max_connections(5)
    .connect(&config.db_conn).await?;
    sqlx::migrate!("./src/migrations").run(&db).await?;

    loop {
    
        let (tournaments, all_tournaments_count) = shark_scope.get_tournaments().await?;
        let date = chrono::offset::Utc::now();
        
        tracing::info!("{:?}/{}", tournaments, all_tournaments_count);
        for tournament in tournaments {
                let result = sqlx::query("insert into five_tournaments (tournament_id, date) values ($1, $2)")
                .bind(&tournament)
                .bind(&date)
                .execute(&db)
                .await;
                tracing::info!("{:?}",result);
        }

        if all_tournaments_count == 0 {
            tokio::time::sleep(Duration::from_secs(60*3)).await;
        }else{
            tokio::time::sleep(Duration::from_secs(20)).await;
        }
    }
}