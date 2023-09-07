import discord
from discord.ext import commands
from discord import app_commands
from config import DISCORD_TOKEN, DISCORD_GUILD_ID
from commands import hello


bot = commands.Bot(intents=discord.Intents.all(), command_prefix=":")
bot.add_cog(hello.Hello(bot))


@bot.event
async def on_ready():
    user = bot.user
    guild = bot.guilds[0]

    print(f"Logged in as {user} (ID: {user.id})")
    print(f"Guild: {guild} (ID: {guild.id})")


if __name__ == "__main__":
    bot.run(DISCORD_TOKEN)
