#!/usr/bin/env python3
# vim: set fileencoding=utf-8 :

import asyncio

import discord
from discord.ext import commands
from discord.gateway import DiscordWebSocket

from config import DISCORD_TOKEN, LYNXIE_PREFIX
from database import CommandHistory, Database
from utils import mobile_status, error_message
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
async def on_command_error(ctx, error):
    if isinstance(error, commands.CommandNotFound):
        return

    print(error)

    error = "An internal error occurred while processing your command, oopsie..."
    await ctx.reply(embed=error_message(error), delete_after=5)


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
