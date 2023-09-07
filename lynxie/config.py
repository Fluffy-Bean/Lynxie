import dotenv
from discord import Object

try:
    DISCORD_TOKEN = dotenv.dotenv_values(".env")["DISCORD_TOKEN"]
except KeyError:
    print("DISCORD_TOKEN not found in .env file")
    exit(1)

DISCORD_GUILD_ID = Object(id=1040757387033849976)
