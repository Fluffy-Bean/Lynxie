import os
import dotenv
from discord import Object

DISCORD_TOKEN = (
    dotenv.dotenv_values(".env").get("DISCORD_TOKEN")
    or os.environ.get("DISCORD_TOKEN")
    or None
)
DISCORD_GUILD_ID = Object(id=1040757387033849976)
LYNXIE_PREFIX = "?"

DATA_PATH = os.path.join("lynxie", "data")
ASSETS_PATH = os.path.join("lynxie", "assets")

DATABASE_URI = "sqlite:///" + os.path.join(DATA_PATH, "lynxie.db")

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
        "path": os.path.join(ASSETS_PATH, "overlays", "bubble.png"),
        "options": [
            "default",  # Positioned at top
            "bottom",  # Positioned at bottom
            "mask",  # Positioned at top, but transparent
            "mask-bottom",  # Positioned at bottom, but transparent
        ],
    },
    "gang": {
        "path": os.path.join(ASSETS_PATH, "overlays", "gang.png"),
        "options": ["default"],
    },
    "bandicam": {
        "path": os.path.join(ASSETS_PATH, "overlays", "bandicam.png"),
        "options": ["default"],
    },
    "jerma": {
        "path": os.path.join(ASSETS_PATH, "overlays", "jerma.png"),
        "options": ["default"],
    },
    "jerm-a": {
        "path": os.path.join(ASSETS_PATH, "overlays", "jerm-a.png"),
        "options": ["default"],
    },
    "liveleak": {
        "path": os.path.join(ASSETS_PATH, "overlays", "liveleak.png"),
        "options": ["default"],
    },
}

E621_API_KEY = (
    dotenv.dotenv_values(".env").get("E621_API_KEY")
    or os.environ.get("E621_API_KEY")
    or None
)
E621_USERNAME = (
    dotenv.dotenv_values(".env").get("E621_USERNAME")
    or os.environ.get("E621_USERNAME")
    or None
)
E621_BLACKLIST = set()
with open(os.path.join(ASSETS_PATH, "e621_blacklist.txt"), "r") as f:
    [E621_BLACKLIST.add(line.strip()) for line in f.readlines() if line.strip()]
