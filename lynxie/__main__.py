import asyncio

import discord
from discord.ext import commands
from discord.gateway import DiscordWebSocket

from lynxie.config import DISCORD_TOKEN, LYNXIE_PREFIX, E621_BLACKLIST
from lynxie.database import CommandHistory, Database
from lynxie.utils import mobile_status
from lynxie.commands import Help, Ping, Hello, Music, Animals, Img, E621


db = Database()
DiscordWebSocket.identify = mobile_status
lynxie = commands.Bot(
    intents=discord.Intents.all(),
    command_prefix=LYNXIE_PREFIX,
    help_command=None,
)


@lynxie.event
async def on_ready():
    print(f"Logged in as {lynxie.user} (ID: {lynxie.user.id})")

    in_guilds = "In Guilds: " + str(len(lynxie.guilds))
    commands_used = "Commands called: " + str(
        db.session.query(CommandHistory.user).count()
    )
    unique_users = "Unique Users: " + str(
        db.session.query(CommandHistory.user).distinct().count()
    )
    blacklisted_words = "Blacklisted Words: " + str(len(E621_BLACKLIST))

    bar_len = (
        max(
            len(in_guilds),
            len(commands_used),
            len(unique_users),
            len(blacklisted_words),
        )
        + 4
    )

    print("---- Stats " + "-" * (bar_len - 11))
    print(f"| {in_guilds}{' ' * (bar_len - len(in_guilds) - 3)}|")
    print(f"| {commands_used}{' ' * (bar_len - len(commands_used) - 3)}|")
    print(f"| {unique_users}{' ' * (bar_len - len(unique_users) - 3)}|")
    print(f"| {blacklisted_words}{' ' * (bar_len - len(blacklisted_words) - 3)}|")
    print("-" * bar_len)


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


async def run():
    async with lynxie:
        await lynxie.add_cog(Help(lynxie))
        await lynxie.add_cog(Ping(lynxie))
        await lynxie.add_cog(Hello(lynxie))
        await lynxie.add_cog(Music(lynxie))
        await lynxie.add_cog(Animals(lynxie))
        await lynxie.add_cog(Img(lynxie))
        await lynxie.add_cog(E621(lynxie))
        await lynxie.start(DISCORD_TOKEN)


if __name__ == "__main__":
    asyncio.run(run())
