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

    @staticmethod
    async def get_image_attachments(ctx):
        if ctx.message.attachments:
            return ctx.message.attachments[0]
        if ctx.message.reference:
            if ctx.message.reference.resolved.attachments:
                return ctx.message.reference.resolved.attachments[0]
            if (
                ctx.message.reference.resolved.embeds
                and ctx.message.reference.resolved.embeds[0].image
            ):
                return ctx.message.reference.resolved.embeds[0].image
        elif ctx.message.embeds and ctx.message.embeds[0].image:
            return ctx.message.embeds[0].image
        else:
            channel = ctx.guild.get_channel(ctx.channel.id)
            async for message in channel.history(limit=10):
                if message.attachments:
                    return message.attachments[0]
                if message.embeds and message.embeds[0].image:
                    return message.embeds[0].image

        return None

    @commands.command()
    async def overlay(
        self, ctx, overlay_choice: str = None, overlay_style: str = "default"
    ):
        start_time = datetime.datetime.now()

        overlay_choice = overlay_choice.lower().strip() if overlay_choice else None
        overlay_style = overlay_style.lower().strip() if overlay_style else "default"

        image_attachments = await self.get_image_attachments(ctx)

        if not image_attachments:
            error = "No image was found!"
            await ctx.reply(embed=error_message(error))
            return

        if not overlay_choice or overlay_choice not in IMAGE_OVERLAYS:
            error = (
                f"Invalid overlay choice! Use one of "
                f"these: {', '.join(IMAGE_OVERLAYS)}"
            )
            await ctx.reply(embed=error_message(error))
            return

        if overlay_style not in IMAGE_OVERLAYS[overlay_choice]["options"]:
            error = (
                f"{overlay_choice} has these "
                f"options: {', '.join(IMAGE_OVERLAYS[overlay_choice]['options'])}"
            )
            await ctx.reply(embed=error_message(error))
            return

        # Defaults to gwa as I cant be asked to make a better error handler
        file_name = (
            image_attachments.filename or image_attachments.url or "balls"
        ).lower()
        file_extension = file_name.split(".")[-1]
        if file_extension not in IMAGE_EXTENSIONS:
            error = (
                f"Unsupported file type! Use one "
                f"of these: {', '.join(IMAGE_EXTENSIONS)}"
            )
            await ctx.reply(embed=error_message(error))
            return

        if image_attachments.size and image_attachments.size > 8 * 1024 * 1024:
            error = "Image must be less than 8MB!"
            await ctx.reply(embed=error_message(error))
            return

        width, height = image_attachments.width, image_attachments.height
        if not 10 < width <= 4500 or not 10 < height <= 4500:
            error = "Image must be at least over 10x10 and under 4500x4500!"
            await ctx.reply(embed=error_message(error))
            return

        async with ctx.typing():
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
                attachment.paste(overlay, ((width - overlay.width) // 2, 0), overlay)
            elif overlay_choice == "jerma":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(
                    overlay, (width - overlay.width, height - overlay.height), overlay
                )
            elif overlay_choice == "jerm-a":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(
                    overlay,
                    ((width - overlay.width) // 2, height - overlay.height),
                    overlay,
                )
            elif overlay_choice == "liveleak":
                overlay = Image.open(IMAGE_OVERLAYS[overlay_choice]["path"])
                overlay.thumbnail((width, overlay.height))
                attachment.paste(overlay, (0, 0), overlay)
            with BytesIO() as response:
                attachment.save(response, format="PNG")

                response.seek(0)
                response = discord.File(response, filename="image.png")

                time_taken = (datetime.datetime.now() - start_time).microseconds / 1000

                embed = (
                    discord.Embed(
                        title=overlay_choice.capitalize(),
                        colour=discord.Colour.orange(),
                    )
                    .set_image(url="attachment://image.png")
                    .set_footer(text=f"{width}x{height}, {time_taken}ms")
                )

                await ctx.reply(embed=embed, file=response, mention_author=False)

    @commands.command()
    async def saveable(self, ctx):
        start_time = datetime.datetime.now()

        image_attachments = await self.get_image_attachments(ctx)

        if not image_attachments:
            error = "No image was found!"
            await ctx.reply(embed=error_message(error))
            return

        # Defaults to gwa as I cant be asked to make a better error handler
        file_name = (
            image_attachments.filename or image_attachments.url or "balls"
        ).lower()
        file_extension = file_name.split(".")[-1]
        if file_extension not in IMAGE_EXTENSIONS:
            error = (
                f"Unsupported file type! Use one "
                f"of these: {', '.join(IMAGE_EXTENSIONS)}"
            )
            await ctx.reply(embed=error_message(error))
            return

        if image_attachments.size and image_attachments.size > 8 * 1024 * 1024:
            error = "Image must be less than 8MB!"
            await ctx.reply(embed=error_message(error))
            return

        width, height = image_attachments.width, image_attachments.height
        if not 10 < width <= 4500 or not 10 < height <= 4500:
            error = "Image must be at least 10x10 and under 4500x4500!"
            await ctx.reply(embed=error_message(error))
            return

        async with ctx.typing():
            request = requests.get(image_attachments.url)
            attachment = Image.open(BytesIO(request.content))
            width, height = attachment.width, attachment.height

            with BytesIO() as response:
                attachment.save(response, format="GIF")

                response.seek(0)
                response = discord.File(response, filename="image.gif")

                time_taken = (datetime.datetime.now() - start_time).microseconds / 1000

                embed = (
                    discord.Embed(
                        title="Saveable",
                        description="Image converted to GIF, click the star to save it :3",
                        colour=discord.Colour.orange(),
                    )
                    .set_image(url="attachment://image.gif")
                    .set_footer(text=f"{width}x{height}, {time_taken}ms")
                )

                await ctx.reply(embed=embed, file=response, mention_author=False)
