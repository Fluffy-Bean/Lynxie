import discord
from discord.ext import commands
from lynxie.config import LYNXIE_PREFIX


class Help(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.help_commands = {
            "help": "Show this message",
            "ping": "Pong!",
            "hello": "Say hello to Lynxie!",
            "join": "Join the voice channel you're in",
            "play <url>": "Play a song from YouTube, SoundCloud, etc.",
            "stop": "Stop the current song and leave the voice channel",
            "e621": "Search e621.net",
        }

    @commands.command()
    async def help(self, ctx):
        embed = discord.Embed(
            title="Help",
            description=f"Lynxie's prefix is `{LYNXIE_PREFIX}`",
            colour=discord.Colour.orange(),
        )

        for command, description in self.help_commands.items():
            embed.add_field(
                name=command,
                value=description,
                inline=False,
            )

        await ctx.reply(embed=embed)
