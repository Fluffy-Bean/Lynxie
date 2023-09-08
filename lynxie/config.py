import os
from discord import Object
from lynxie.utils import get_env_or_error

DATA_PATH = "data"

LYNXIE_PREFIX = "~"
LYNXIE_DB = f"sqlite:///" + os.path.join(DATA_PATH, "lynxie.db")

DISCORD_GUILD_ID = Object(id=1040757387033849976)
DISCORD_TOKEN = get_env_or_error("DISCORD_TOKEN")

TINYFOX_ANIMALS = [
    "chi", "bear", "dog", "capy", "caracal", "chee",
    "yote", "bleat", "dook", "fox", "yeen", "jaguar",
    "leo", "lynx", "mane", "puma", "poss", "ott",
    "manul", "marten", "bun", "racc", "wah", "serval",
    "shiba", "skunk", "snek", "snep", "tig", "woof",
]
