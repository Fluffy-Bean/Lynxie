import os
import dotenv
import discord


def error_message(error: str) -> discord.Embed:
    print("Error: " + error)

    embed = discord.Embed(
        title="Error",
        description=f"`{error}`",
        colour=discord.Colour.red(),
    )
    embed.set_footer(text="For more information, use the help command.")

    return embed


def get_env_or_error(env: str) -> str:
    from_file = dotenv.dotenv_values(".env").get(env)
    from_env = os.environ.get(env)

    if from_file is None and from_env is None:
        raise KeyError(f"Environment variable {env} not found")
    elif from_file is None:
        return from_env
    else:
        return from_file
