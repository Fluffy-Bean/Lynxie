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
            try:
                song_info = ytdl.extract_info(url, download=False)
                print(song_info["url"])
                ctx.voice_client.play(
                    discord.FFmpegPCMAudio(song_info["url"], **ffmpeg_options)
                )
            except Exception:
                await ctx.reply(
                    embed=error_message(
                        "An error occurred while processing this request."
                    ),
                    mention_author=False,
                )
                return

        await ctx.send(f"Now playing: {song_info['title']}")

    @commands.command()
    async def stop(self, ctx):
        await ctx.voice_client.disconnect()

    @play.before_invoke
    async def ensure_voice(self, ctx):
        if ctx.voice_client is None:
            if ctx.author.voice:
                await ctx.author.voice.channel.connect()
            else:
                await ctx.reply(
                    embed=error_message("You are not connected to a voice channel!"),
                    mention_author=False,
                )
        elif ctx.voice_client.is_playing():
            ctx.voice_client.stop()
