use serenity::{
    builder,
    model::prelude::command::CommandOptionType,
};

pub fn register(
    command: &mut builder::CreateApplicationCommand,
) -> &mut builder::CreateApplicationCommand {
    command
        .name("numberinput")
        .description("Test command for number input")
        .create_option(|option| {
            option
                .name("int")
                .description("An integer from 1 to 69")
                .kind(CommandOptionType::Integer)
                .min_int_value(1)
                .max_int_value(69)
                .required(true)
        })
        .create_option(|option| {
            option
                .name("float")
                .description("A float from -2.1 to 621.69")
                .kind(CommandOptionType::Number)
                .min_number_value(-2.1)
                .max_number_value(621.69)
                .required(true)
        })
}
