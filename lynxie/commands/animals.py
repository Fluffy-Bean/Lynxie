import requests
from io import BytesIO

import discord
from discord.ext import commands

from lynxie.config import TINYFOX_ANIMALS
from lynxie.utils import error_message


class Animals(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def animal(self, ctx, animal_choice: str = ""):
        animal_choice = animal_choice.lower().strip() or None

        if not animal_choice:
            error = f"You need to specify an animal! Try one of these: {', '.join(TINYFOX_ANIMALS)}"
            await ctx.reply(embed=error_message(error))
            return

        if animal_choice not in TINYFOX_ANIMALS:
            error = f"That animal doesn't exist! Try one of these: {', '.join(TINYFOX_ANIMALS)}"
            await ctx.reply(embed=error_message(error))
            return

        async with ctx.typing():
            request = requests.get("https://api.tinyfox.dev/img?animal=" + animal_choice)

            with BytesIO(request.content) as response:
                response.seek(0)
                animal_file = discord.File(response, filename="image.png")

                embed = discord.Embed(
                    title=animal_choice.capitalize(),
                    colour=discord.Colour.orange(),
                ).set_image(url="attachment://image.png")

                await ctx.reply(embed=embed, file=animal_file, mention_author=False)
