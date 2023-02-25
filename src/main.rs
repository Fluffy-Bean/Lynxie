mod commands;

use dotenv::dotenv;
use std::env;

use serenity::{
    async_trait,
//    model::application::command::Command,
    model::application::interaction::{
        Interaction,
        InteractionResponseType,
    },
    model::gateway::Ready,
    model::id::GuildId,
    prelude::*,
};

// Will handle Pisscord events
struct Handler;

#[async_trait]
impl EventHandler for Handler {
    // Create a handler for slash commands
    async fn interaction_create(&self, ctx: Context, interaction: Interaction) {
        if let Interaction::ApplicationCommand(command) = interaction {
            println!("Heard command: {:?}", command);

            // Create a string to hold the response
            let content = match command.data.name.as_str() {
                "ping" => commands::ping::run(&command.data.options),
                "id" => commands::id::run(&command.data.options),
                "attachmentinput" => commands::attachmentinput::run(&command.data.options),
                _ => "Not implemented :c".to_string(),
            };

            // Respond to the command
            if let Err(why) = command
                .create_interaction_response(&ctx.http, |response| {
                    response
                        .kind(InteractionResponseType::ChannelMessageWithSource)
                        .interaction_response_data(|message| message.content(content))
                })
                .await
            {
                // Oups, something went wrong D:
                println!("Error responding to command: {:?}", why);
            }
        }
    }

    // Create a handler for when the bot is ready
    async fn ready(&self, ctx: Context, ready: Ready) {
        // Display that the bot is ready along with the username
        println!("{} is ready Bitch!", ready.user.name);

        // Set guild to register commands in
        // Much faster than waiting for commands to be avaiable globally
        let guild_id = GuildId(
            env::var("GUILD_ID")
                .expect("No GUILD_ID variable found!!!!!")
                .parse()
                .expect("GUILD_ID must be an intiger :<")
        );

        let _commands = GuildId::set_application_commands(&guild_id, &ctx.http, |commands| {
            commands
                .create_application_command(|command| commands::ping::register(command))
                .create_application_command(|command| commands::id::register(command))
                .create_application_command(|command| commands::numberinput::register(command))
                .create_application_command(|command| commands::attachmentinput::register(command))
        })
        .await;

        println!("Guild {} has been registered!", guild_id); // commands

        // Set global commands
        // This is slow and should be commented out when testing
        /* 
        let guild_command = Command::create_global_application_command(&ctx.http, |command| {
            commands::wonderful_command::register(command)
        })
        .await;

        println!("Following comamnds available globally!: {:#?}", guild_command);
        */
    }
}

#[tokio::main]
async fn main() {
    // Load dotenv file and yoink the token
    dotenv().ok();
    let token = env::var("DISCORD_TOKEN").expect("No DISCORD_TOKEN variable found!!!!!");

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