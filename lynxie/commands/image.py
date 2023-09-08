import os
import requests
from io import BytesIO
import datetime

from PIL import Image

import discord
from discord.ext import commands

from lynxie.config import IMAGE_EXTENSIONS, IMAGE_OVERLAYS, ASSETS_PATH
from lynxie.utils import error_message


class Img(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self._overlays = {
            "bubble": Image.open(os.path.join(ASSETS_PATH, "bubble.png")).convert(
                "RGBA"
            ),
            "gang": Image.open(os.path.join(ASSETS_PATH, "gang.png")).convert("RGBA"),
        }

    @commands.command()
    async def overlay(self, ctx, style: str = None):
        start_time = datetime.datetime.now()

        style = style.lower().strip() if style else None

        # Check if image should be processed
        async with ctx.typing():
            if not style or style not in IMAGE_OVERLAYS:
                error = (
                    "That is not a valid option! Valid options are:\n"
                    f"`{', '.join(IMAGE_OVERLAYS)}`"
                )
                await ctx.reply(embed=error_message(error))
                return

            if not ctx.message.attachments:
                error = "You need to attach an image to use this command!"
                await ctx.reply(embed=error_message(error))
                return

            # Extracts file extension from filename
            if (
                not ctx.message.attachments[0].filename.split(".")[-1].lower()
                in IMAGE_EXTENSIONS
            ):
                error = (
                    "Unsupported file type! Supported file types are:\n"
                    f"`{', '.join(IMAGE_EXTENSIONS)}`"
                )
                await ctx.reply(embed=error_message(error))
                return

            if ctx.message.attachments[0].size > 8_000_000:
                error = (
                    "That image is too big! Please use an image that is less than 8MB."
                )
                await ctx.reply(embed=error_message(error))
                return

            if (
                not 0 < ctx.message.attachments[0].width <= 3500
                or not 0 < ctx.message.attachments[0].height <= 3500
            ):
                error = "Image must be at least 1x1 and under 3500x3500!"
                await ctx.reply(embed=error_message(error))
                return

        # Process image
        async with ctx.typing():
            response = requests.get(ctx.message.attachments[0].url)
            message_attachment = Image.open(BytesIO(response.content))

            if message_attachment.width < message_attachment.height:
                message_attachment.thumbnail((200, message_attachment.height))
            else:
                message_attachment.thumbnail((message_attachment.width, 200))

            if style == "bubble":
                # The bubble is resized twice as for some reason .copy() doesn't work
                message_attachment.paste(
                    self._overlays["bubble"].resize(
                        (message_attachment.width, self._overlays["bubble"].height)
                    ),
                    (0, 0),
                    self._overlays["bubble"].resize(
                        (message_attachment.width, self._overlays["bubble"].height)
                    ),
                )
            elif style == "gang":
                message_attachment.paste(
                    self._overlays["gang"],
                    (
                        (
                            (message_attachment.width - self._overlays["gang"].width)
                            // 2
                        ),
                        (message_attachment.height - self._overlays["gang"].height),
                    ),
                    self._overlays["gang"],
                )

            message_file = BytesIO()
            message_attachment.save(message_file, format="PNG")
            message_file.seek(0)
            message_file = discord.File(message_file, filename="image.png")

            time_taken = datetime.datetime.now() - start_time
            embed = (
                discord.Embed(title=style.capitalize(), colour=discord.Colour.orange())
                .set_image(url="attachment://image.png")
                .set_footer(
                    text=f"{message_attachment.width}x{message_attachment.height}, "
                    f"{time_taken.microseconds / 1000}ms"
                )
            )

            await ctx.reply(embed=embed, file=message_file, mention_author=False)
