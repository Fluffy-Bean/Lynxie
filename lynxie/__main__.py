import asyncio
import discord
from discord.ext import commands
from config import DISCORD_TOKEN
from commands import Help, Ping, Hello, Music, E621


lynxie = commands.Bot(
    intents=discord.Intents.all(),
    command_prefix="AAAA ",
    help_command=None,
)


@lynxie.event
async def on_ready():
    print(f"Logged in as {lynxie.user} (ID: {lynxie.user.id})")


@lynxie.event
async def on_message_edit(before, after):
    await before.channel.send(
        f"@{before.author} edited their message!!!\n"
        f'Before: "{before.content}"\n'
        f'After: "{after.content}"'
    )


async def run():
    async with lynxie:
        await lynxie.add_cog(Help(lynxie))
        await lynxie.add_cog(Ping(lynxie))
        await lynxie.add_cog(Hello(lynxie))
        await lynxie.add_cog(Music(lynxie))
        await lynxie.add_cog(E621(lynxie))
        await lynxie.start(DISCORD_TOKEN)


if __name__ == "__main__":
    asyncio.run(run())
