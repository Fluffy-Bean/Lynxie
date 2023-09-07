import discord
from discord.ext import commands


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
            description="Lynxie's prefix is `AAAA `",
            colour=discord.Colour.orange(),
        )

        for command, description in self.help_commands.items():
            embed.add_field(
                name=command,
                value=description,
                inline=False,
            )

        await ctx.send(embed=embed)
