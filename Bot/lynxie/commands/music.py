import yt_dlp
import discord
from discord.ext import commands
from lynxie.utils import error_message


ffmpeg_options = {"options": "-vn"}
ydl_opts = {"format": "bestaudio"}
ytdl = yt_dlp.YoutubeDL(ydl_opts)


class Music(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def join(self, ctx, *, channel: discord.VoiceChannel):
        if ctx.voice_client is not None:
            return await ctx.voice_client.move_to(channel)
        await channel.connect()

    @commands.command()
    async def play(self, ctx, *, url):
        async with ctx.typing():
            song_info = ytdl.extract_info(url, download=False)
            ctx.voice_client.play(
                discord.FFmpegPCMAudio(song_info["url"], **ffmpeg_options)
            )

        embed = discord.Embed(
            title="Now playing",
            description=f"[{song_info['title']}]({song_info['webpage_url']})",
            color=discord.Color.orange(),
        )

        await ctx.reply(embed=embed, mention_author=False)

    @commands.command()
    async def stop(self, ctx):
        await ctx.voice_client.disconnect()

    @play.before_invoke
    async def ensure_voice(self, ctx):
        if ctx.voice_client is None:
            if ctx.author.voice:
                await ctx.author.voice.channel.connect()
            else:
                error = "You are not connected to a voice channel."
                await ctx.reply(
                    embed=error_message(error),
                    mention_author=False,
                )
                return
        elif ctx.voice_client.is_playing():
            ctx.voice_client.stop()
