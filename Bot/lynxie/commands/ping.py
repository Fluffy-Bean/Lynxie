import discord
from discord.ext import commands


class Ping(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def ping(self, ctx):
        embed = discord.Embed(
            title="Pong!",
            description=f"{round(self.bot.latency * 1000)}ms",
            color=discord.Color.orange(),
        )
        await ctx.reply(embed=embed, mention_author=False)
