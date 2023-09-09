import discord
from discord.ext import commands


class Hello(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def hello(self, ctx):
        embed = discord.Embed(
            title="Hello!",
            description="I'm Lynxie, a multipurpose Discord bot written in Python!",
            color=discord.Color.orange(),
        )

        await ctx.reply(embed=embed)
