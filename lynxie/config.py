import os
from discord import Object
from lynxie.utils import get_env_or_error

DISCORD_TOKEN = get_env_or_error("DISCORD_TOKEN")
DISCORD_GUILD_ID = Object(id=1040757387033849976)
LYNXIE_PREFIX = "?"

DATA_PATH = "data"
ASSETS_PATH = "assets"

DATABASE_URI = f"sqlite:///" + os.path.join(DATA_PATH, "lynxie.db")

# https://tinyfox.dev/docs/
# TODO: Get list from API instead of hardcoding
# https://api.tinyfox.dev/img?animal=animal&json
TINYFOX_ANIMALS = [
    "chi",
    "bear",
    "dog",
    "capy",
    "caracal",
    "chee",
    "yote",
    "bleat",
    "dook",
    "fox",
    "yeen",
    "jaguar",
    "leo",
    "lynx",
    "mane",
    "puma",
    "poss",
    "ott",
    "manul",
    "marten",
    "bun",
    "racc",
    "wah",
    "serval",
    "shiba",
    "skunk",
    "snek",
    "snep",
    "tig",
    "woof",
]

IMAGE_EXTENSIONS = ["png", "jpg", "jpeg", "webp"]

IMAGE_OVERLAYS = {
    "bubble": {
        "path": os.path.join(ASSETS_PATH, "bubble.png"),
        "options": [
            "default",  # Positioned at top
            "bottom",  # Positioned at bottom
            "mask",  # Positioned at top, but transparent
            "mask-bottom",  # Positioned at bottom, but transparent
        ],
    },
    "gang": {
        "path": os.path.join(ASSETS_PATH, "gang.png"),
        "options": ["default"],
    },
    "bandicam": {
        "path": os.path.join(ASSETS_PATH, "bandicam.png"),
        "options": ["default"],
    },
    "jerma": {
        "path": os.path.join(ASSETS_PATH, "jerma.png"),
        "options": ["default"],
    },
    "jerm-a": {
        "path": os.path.join(ASSETS_PATH, "jerm-a.png"),
        "options": ["default"],
    },
}
