from selenium import webdriver
from bs4 import BeautifulSoup
import discord
from discord.ext import commands


class E621(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    async def e621(self, ctx):
        embed = discord.Embed(
            title="Search Results",
            description="Here's a list of jobs I found on Indeed, just for you!",
            colour=discord.Colour.orange(),
        )

        browser = webdriver.Firefox()
        browser.get("https://www.indeed.com/jobs?q=cleaner&l=New%20York")
        soup = BeautifulSoup(browser.page_source, "html.parser")
        browser.close()

        for job in soup.find_all("div", {"class": "job_seen_beacon"}):
            job_title = (
                job.find("h2", {"class": "jobTitle"}).find("span").text.strip()
                or "Job Title"
            )
            company_name = (
                job.find("span", {"class": "companyName"}).text.strip()
                or "Company Name"
            )
            company_location = (
                job.find("div", {"class": "companyLocation"}).text.strip() or "Location"
            )

            embed.add_field(
                name=job_title,
                value=f"{company_name} - {company_location}",
                inline=False,
            )

        await ctx.send(embed=embed)
