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
    async def animal(self, ctx, animal):
        animal = animal.lower().strip() or "racc"

        if animal not in TINYFOX_ANIMALS:
            await ctx.reply(
                embed=error_message(
                    f"That animal doesn't exist! Try one of these:\n"
                    f"`{', '.join(TINYFOX_ANIMALS)}`"
                )
            )
            return

        async with ctx.typing():
            request = requests.get(f"https://api.tinyfox.dev/img?animal={animal}&json")
            animal_image = BytesIO(request.content)
            animal_image.seek(0)
            animal_file = discord.File(animal_image, filename="image.png")

            embed = discord.Embed(
                title=animal.capitalize(),
                colour=discord.Colour.orange(),
            ).set_image(
                url="attachment://image.png"
            )

            await ctx.reply(embed=embed, file=animal_file, mention_author=False)
