import requests

from discord.ext import commands


class E621(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def e621(self, ctx):
        await ctx.reply(f":3")
