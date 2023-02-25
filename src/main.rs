use dotenv::dotenv;
use std::env;

use serenity::{
    async_trait,
    model::{channel::Message, gateway::Ready},
    prelude::*,
};

// Will handle Pisscord events
struct Handler;

#[async_trait]
impl EventHandler for Handler {
    // Create a handler for a new message event
    async fn message(&self, ctx: Context, msg: Message) {
        if msg.content == "ping" {
            // Set Reply to "pong"
            let reply: &str = "pong";

            if let Err(why) = msg.channel_id.say(&ctx.http, reply).await {
                // If an error occurs, display it
                println!("Error sending message: {:?}", why);
            }
        }

        if msg.content == "pale" {
            // Set Reply to "We should ban this guy"
            let reply: &str = "We should ban this guy";

            if let Err(why) = msg.channel_id.say(&ctx.http, reply).await {
                // If an error occurs, display it
                println!("Error sending message: {:?}", why);
            }
        }
    }

    // Create a handler for when the bot is ready
    async fn ready(&self, _: Context, ready: Ready) {
        // Display that the bot is ready along with the username
        println!("{} is ready Bitch!", ready.user.name);
    }
}

#[tokio::main]
async fn main() {
    // Load dotenv file and yoink the token
    dotenv().ok();
    let token = env::var("DISCORD_TOKEN").expect("No Env enviroment variable found!");

    // Set the bots intents
    let intents = GatewayIntents::GUILD_MESSAGES
        | GatewayIntents::DIRECT_MESSAGES
        | GatewayIntents::MESSAGE_CONTENT;

    // Login to the Pisscord client as a bot
    let mut client = Client::builder(token, intents).event_handler(Handler).await.expect("Error connecting to client, are you online???");

    // We start a single shart
    if let Err(why) = client.start().await {
        println!("Clinet buggin out, heres why: {:?}", why);
    }
}