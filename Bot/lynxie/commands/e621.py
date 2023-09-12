import json
from base64 import b64encode
import requests

import discord
from discord.ext import commands

from .config import E621_API_KEY, E621_USERNAME, E621_BLACKLIST
from .utils import error_message


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
        url = "https://e621.net/posts.json/?limit=1&tags=order:random+rating:e+"
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
                f"{', '.join(['`'+tag+'`' for tag in caught_tags])} "
                f"is a blacklisted tag!"
            )
            await ctx.reply(embed=error_message(error))
            return

        request = requests.get(url, headers=_E621_API_HEADERS)
        response = json.loads(request.text)

        if request.status_code == 503:
            error = (
                "The bot is currently rate limited! "
                "Wait a while before trying again."
            )
            await ctx.reply(embed=error_message(error))
            return
        if request.status_code != 200:
            error = (
                "An error occurred while fetching the image! "
                f"(Error code: {str(request.status_code)})"
            )
            await ctx.reply(embed=error_message(error))
            return

        if not response["posts"]:
            tags_to_display = range(min(len(tags), 20))
            error = (
                "No results found for the given tags! "
                f"(Tags: {', '.join(['`'+tags[i]+'`' for i in tags_to_display])})"
            )
            await ctx.reply(embed=error_message(error))
            return

        post = response["posts"][0]
        general_tags = post["tags"]["general"]

        embed = discord.Embed(
            title="E621",
            description=post["description"] or "No description provided.",
            colour=discord.Colour.orange(),
        )

        embed.add_field(
            name="Score",
            value=f"⬆️ {post['score']['up']} | ⬇️ {post['score']['down']}",
        )
        embed.add_field(
            name="Favorites",
            value=post["fav_count"],
        )
        embed.add_field(
            name="Comments",
            value=post["comment_count"],
        )

        embed.add_field(
            name="Source(s)",
            value=", ".join(post["sources"]) or "No source provided.",
            inline=False,
        )
        embed.add_field(
            name="Tags",
            value=(
                ", ".join(
                    [
                        "`" + general_tags[i] + "`"
                        for i in range(min(len(general_tags), 20))
                    ]
                )
                or "No tags provided."
            ),
            inline=False,
        )

        embed.set_footer(text=f"ID: {post['id']} | Created: {post['created_at']}")
        embed.set_image(url=post["file"]["url"])

        await ctx.reply(embed=embed)
