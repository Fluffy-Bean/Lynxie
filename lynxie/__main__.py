import asyncio
import discord
from discord.ext import commands
from config import DISCORD_TOKEN
from commands import Hello, Music


lynxie = commands.Bot(
    intents=discord.Intents.all(),
    command_prefix="AAAA ",
)


@lynxie.event
async def on_ready():
    print(f"Logged in as {lynxie.user} (ID: {lynxie.user.id})")


async def run():
    async with lynxie:
        await lynxie.add_cog(Hello(lynxie))
        await lynxie.add_cog(Music(lynxie))
        await lynxie.start(DISCORD_TOKEN)


if __name__ == "__main__":
    asyncio.run(run())
