import datetime
import requests
from io import BytesIO

from PIL import Image, ImageEnhance

import discord
from discord.ext import commands

from lynxie.config import IMAGE_EXTENSIONS, IMAGE_OVERLAYS
from lynxie.utils import error_message


class Img(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def overlay(self, ctx, overlay_choice: str = None, overlay_style: str = "default"):
        start_time = datetime.datetime.now()

        overlay_choice = overlay_choice.lower().strip() if overlay_choice else None
        overlay_style = overlay_style.lower().strip() if overlay_style else "default"

        image_attachments = None
        if ctx.message.attachments:
            image_attachments = ctx.message.attachments[0]
        elif ctx.message.reference and ctx.message.reference.resolved.attachments:
            image_attachments = ctx.message.reference.resolved.attachments[0]
        elif ctx.message.embeds and ctx.message.embeds[0].image:
            image_attachments = ctx.message.embeds[0].image
        else:
            channel = ctx.guild.get_channel(ctx.channel.id)
            async for message in channel.history(limit=10):
                if message.attachments:
                    image_attachments = message.attachments[0]
                    break
                if message.embeds and message.embeds[0].image:
                    image_attachments = message.embeds[0].image
                    break

        async with ctx.typing():
            if not image_attachments:
                error = "No image was found!"
                await ctx.reply(embed=error_message(error))
                return

            if not overlay_choice or overlay_choice not in IMAGE_OVERLAYS:
                error = (
                    "That is not a valid option! Valid options are:\n"
                    f"`{', '.join(IMAGE_OVERLAYS)}`"
                )
                await ctx.reply(embed=error_message(error))
                return

            if overlay_style not in IMAGE_OVERLAYS[overlay_choice]["options"]:
                error = (
                    "That is not a valid option! Valid options are:\n"
                    f"`{', '.join(IMAGE_OVERLAYS[overlay_choice]['options'])}`"
                )
                await ctx.reply(embed=error_message(error))
                return

            # Defaults to gwa as I cant be asked to make a better error handler
            filename = image_attachments.filename or image_attachments.url or "image.gwa"
            if not filename.split(".")[-1].lower() in IMAGE_EXTENSIONS:
                error = (
                    "Unsupported file type! Supported file types are "
                    ", ".join(IMAGE_EXTENSIONS)
                )
                await ctx.reply(embed=error_message(error))
                return

            if image_attachments.size and image_attachments.size > 8 * 1024 * 1024:
                error = (
                    "That image is too big! "
                    "Please use an image that is less than 8MB."
                )
                await ctx.reply(embed=error_message(error))
                return

            if (
                not 10 < image_attachments.width <= 3500
                or not 10 < image_attachments.height <= 3500
            ):
                error = "Image must be at least 10x10 and under 3500x3500!"
                await ctx.reply(embed=error_message(error))
                return

            request = requests.get(image_attachments.url)
            attachment = Image.open(BytesIO(request.content))
            width, height = attachment.width, attachment.height

            if width < height:
                attachment.thumbnail((200, height))
            else:
                attachment.thumbnail((width, 200))

            width, height = attachment.width, attachment.height

            if overlay_choice == "bubble":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay = overlay.resize((width, overlay.height))

                if overlay_style in ["default", "top"]:
                    attachment.paste(overlay, (0, 0), overlay)
                elif overlay_style in ["bottom"]:
                    overlay = overlay.rotate(180)
                    attachment.paste(overlay, (0, height - overlay.height), overlay)
                elif overlay_style in ["mask", "mask-bottom"]:
                    # This is a lazy method of creating a mask
                    # 1. Reduce brightness of overlay to 0 (black)
                    # 2. Create a white square the size of the image
                    # 3. Paste the overlay onto the white square

                    overlay = ImageEnhance.Brightness(overlay).enhance(0)

                    mask = Image.new("RGB", (width, height), (255, 255, 255))
                    mask.paste(overlay, (0, 0), overlay)

                    if overlay_style == "mask-bottom":
                        mask = mask.rotate(180)

                    mask = mask.convert("L")

                    attachment.putalpha(mask)
            elif overlay_choice == "gang":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                position = ((width - overlay.width) // 2, (height - overlay.height))
                attachment.paste(overlay, position, overlay)
            elif overlay_choice == "bandicam":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(overlay, ((width-overlay.width)//2, 0), overlay)
            elif overlay_choice == "jerma":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(overlay, (width-overlay.width, height-overlay.height), overlay)
            elif overlay_choice == "jerm-a":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(overlay, ((width-overlay.width)//2, height-overlay.height), overlay)
            with BytesIO() as response:
                attachment.save(response, format="PNG")

                response.seek(0)
                response = discord.File(response, filename="image.png")

                time_taken = (datetime.datetime.now() - start_time).microseconds / 1000

                embed = (
                    discord.Embed(title=overlay_choice.capitalize(), colour=discord.Colour.orange())
                    .set_image(url="attachment://image.png")
                    .set_footer(text=f"{width}x{height}, {time_taken}ms")
                )

                await ctx.reply(embed=embed, file=response, mention_author=False)
