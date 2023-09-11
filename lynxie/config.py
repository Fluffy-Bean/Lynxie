import os
import requests
import dotenv
from discord import Object


LYNXIE_PREFIX = "?"

DISCORD_TOKEN = (
    dotenv.dotenv_values(".env").get("DISCORD_TOKEN")
    or os.environ.get("DISCORD_TOKEN")
    or None
)
DISCORD_GUILD_ID = Object(id=1040757387033849976)
LYNXIE_PREFIX = "?"

DATA_PATH = os.path.join("lynxie", "data")
ASSETS_PATH = os.path.join("lynxie", "assets")

DATABASE_URI = f"sqlite:///" + os.path.join(DATA_PATH, "lynxie.db")


# https://tinyfox.dev/docs/
tinyfox_animals = requests.get("https://api.tinyfox.dev/img?animal=animal&json").json()
TINYFOX_ANIMALS = tinyfox_animals["available"]


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
with open(os.path.join(ASSETS_PATH, "e621_blacklist.txt"), "r") as file:
    for line in file.readlines():
        if word := line.strip():
            E621_BLACKLIST.add(word)
