import json
from base64 import b64encode
import requests

import discord
from discord.ext import commands

from lynxie.config import E621_API_KEY, E621_USERNAME, E621_BLACKLIST
from lynxie.utils import error_message


_E621_API_URL = "https://e621.net/"
_E621_AUTH = f"{E621_USERNAME}:{E621_API_KEY}".encode("utf-8")
_E621_API_HEADERS = {
    "Accept": "application/json",
    "Content-Type": "application/json",
    "User-Agent": f"Lynxie/1.0 (by {E621_USERNAME} on e621)",
    "Authorization": str(b"Basic " + b64encode(_E621_AUTH), "utf-8"),
}


class E621(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def porb(self, ctx, *tags):
        # Base url for the request
        url = _E621_API_URL + "posts.json/?limit=1&tags=order:random+rating:e+"
        caught_tags = []

        for tag in tags:
            tag = tag.lower()
            url += tag + "+"
            if tag in E621_BLACKLIST:
                caught_tags.append(tag)

        for tag in E621_BLACKLIST:
            url += f"-{tag}+"

        if caught_tags:
            error = (
                "An error occurred while fetching the image! "
                f"{', '.join(caught_tags)} is a blacklisted tag!"
            )
            await ctx.reply(embed=error_message(error))
            return

        request = requests.get(url, headers=_E621_API_HEADERS)
        response = json.loads(request.text)

        if request.status_code != 200:
            error = (
                "An error occurred while fetching the image! "
                f"(Error code: {str(request.status_code)})"
            )
            await ctx.reply(embed=error_message(error))
            return

        if not response["posts"]:
            error = "No results found for the given tags! " f"(Tags: {', '.join(tags)})"
            await ctx.reply(embed=error_message(error))
            return

        embed = discord.Embed(
            title="E621",
            description=response["posts"][0]["description"]
            or "No description provided.",
            colour=discord.Colour.orange(),
        )

        embed.add_field(
            name="Score",
            value=f"^ {response['posts'][0]['score']['up']} | "
                  f"v {response['posts'][0]['score']['down']}",
        )
        embed.add_field(
            name="Favorites",
            value=response["posts"][0]["fav_count"],
        )

        embed.add_field(
            name="Source",
            value=", ".join(response["posts"][0]["sources"]) or "No source provided.",
            inline=False,
        )
        embed.add_field(
            name="Tags",
            value=", ".join(response["posts"][0]["tags"]["general"]) or "No tags provided.",
            inline=False,
        )

        embed.set_footer(
            text=f"ID: {response['posts'][0]['id']} | "
                 f"Created: {response['posts'][0]['created_at']}"
        )

        embed.set_image(url=response["posts"][0]["file"]["url"])

        await ctx.reply(embed=embed)
