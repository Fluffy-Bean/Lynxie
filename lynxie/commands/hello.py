import discord
from discord.ext import commands


class Hello(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def hello(self, interaction: discord.Interaction):
        """Says hello!"""
        await interaction.response.send_message(f"Balls, {interaction.user.mention}")
