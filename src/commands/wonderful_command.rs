use serenity::builder::CreateApplicationCommand;

pub fn _register(command: &mut CreateApplicationCommand) -> &mut CreateApplicationCommand {
    command.name("wonderful_command").description("A wonderful command")
}