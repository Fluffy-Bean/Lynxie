import sys
import discord
from discord.gateway import _log
from lynxie.config import LYNXIE_PREFIX


async def mobile_status(self):
    payload = {
        "op": self.IDENTIFY,
        "d": {
            "token": self.token,
            "properties": {
                "$os": sys.platform,
                "$browser": "Discord Android",
                "$device": "Discord Android",
                "$referrer": "",
                "$referring_domain": "",
            },
            "compress": True,
            "large_threshold": 250,
            "v": 3,
        },
    }

    if self.shard_id is not None and self.shard_count is not None:
        payload["d"]["shard"] = [self.shard_id, self.shard_count]

    state = self._connection
    if state._activity is not None or state._status is not None:
        payload["d"]["presence"] = {
            "status": state._status,
            "game": state._activity,
            "since": 0,
            "afk": False,
        }

    if state._intents is not None:
        payload["d"]["intents"] = state._intents.value

    await self.call_hooks(
        "before_identify", self.shard_id, initial=self._initial_identify
    )
    await self.send_as_json(payload)
    _log.info("Shard ID %s has sent the IDENTIFY payload.", self.shard_id)


def error_message(error: str) -> discord.Embed:
    print("Error: " + error)

    embed = discord.Embed(
        title="Error :(",
        description=error,
        colour=discord.Colour.red(),
    )
    embed.set_footer(
        text=f"For more information, use the " f"{LYNXIE_PREFIX}help command."
    )

    return embed