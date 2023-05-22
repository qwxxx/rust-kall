use std::time::Duration;
use anyhow::Result;
use sqlx::{postgres::PgPoolOptions};
use tokio::sync::Mutex;
use std::sync::Arc;

mod config;
mod sharkscope;

async fn init_program() -> Result<config::Config> {
    dotenv::dotenv().ok();
    let filter = tracing_subscriber::EnvFilter::from_env("RUST_LOG");
    tracing_subscriber::fmt().with_env_filter(filter).init();

    let file = std::fs::File::open("config.yaml")?;
    let config: config::Config = serde_yaml::from_reader(file)?;

    Ok(config)
}

async fn notify_thread(
        config: config::Config,
        shark_scope: Arc<Mutex<sharkscope::SharkScope>>,
    ) {
        let db=PgPoolOptions::new().max_connections(5).connect(&config.db_conn).await.unwrap();
        sqlx::migrate!("./src/migrations").run(&db).await.unwrap();
        loop {
            let (tournaments, all_tournaments_count) = shark_scope.lock().await.get_tournaments().await.unwrap_or((vec![],0));
            let date = chrono::offset::Utc::now();
            
            tracing::info!("{:?}/{}", tournaments, all_tournaments_count);

            for tournament in tournaments {
                let _result = sqlx::query("insert into five_tournaments (tournament_id, date) values ($1, $2)")
                    .bind(&tournament)
                    .bind(&date)
                    .execute(&db)
                    .await;
            }

            if all_tournaments_count == 0 {
                tokio::time::sleep(Duration::from_secs(60*3)).await;
            } else {
                tokio::time::sleep(Duration::from_secs(20)).await;
            }
        }
}

#[tokio::main]
async fn main() -> Result<()> {
    let config_main = init_program().await?;
    let config_notify = config_main.clone();
    let mut shark_scope = Arc::new(Mutex::new(sharkscope::SharkScope::new(config_main.sharkscope.clone())));
    let mut shark_scope_main=Arc::clone(&shark_scope);

    tokio::spawn(async move {
        notify_thread(config_notify, shark_scope).await
    }).await?;

    tokio::spawn(async move {
        notify_thread(config_main, shark_scope_main).await
    }).await?;

    Ok(())
}