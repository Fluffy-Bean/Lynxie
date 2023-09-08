import asyncio

import discord
from discord.ext import commands

from lynxie.config import DISCORD_TOKEN, LYNXIE_PREFIX, LYNXIE_DB
from lynxie.commands import Help, Ping, Hello, Music, Animals
from lynxie.database import CommandHistory, Database


lynxie = commands.Bot(
    intents=discord.Intents.all(),
    command_prefix=LYNXIE_PREFIX,
    help_command=None,
)

db = Database()


@lynxie.event
async def on_ready():
    print(f"Logged in as {lynxie.user} (ID: {lynxie.user.id})")
    print("------ Stats ------")
    print(f"Guilds: {len(lynxie.guilds)}")
    print(f"Users: {db.session.query(CommandHistory.user).distinct().count()}")


@lynxie.event
async def on_command(ctx):
    if ctx.author == lynxie.user or ctx.author.bot:
        return

    query = CommandHistory(
        command=ctx.command.name,
        user=ctx.author.id,
        channel=ctx.channel.id,
        guild=ctx.guild.id,
        timestamp=ctx.message.created_at,
    )

    db.session.add(query)
    db.session.commit()


@lynxie.event
async def on_message_edit(before, after):
    if (
        before.author == lynxie.user
        or before.author.bot
        or before.content == after.content
    ):
        return

    await before.channel.send(
        f"@{before.author} edited their message!!!\n"
        f'"{before.content}" --> "{after.content}"'
    )

    # await lynxie.process_commands(after)


async def run():
    async with lynxie:
        await lynxie.add_cog(Help(lynxie))
        await lynxie.add_cog(Ping(lynxie))
        await lynxie.add_cog(Hello(lynxie))
        await lynxie.add_cog(Music(lynxie))
        await lynxie.add_cog(Animals(lynxie))
        await lynxie.start(DISCORD_TOKEN)


if __name__ == "__main__":
    asyncio.run(run())
