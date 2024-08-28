# Automate /now

📌 [A deep dive into how I built this script.](https://akashgoswami.dev/posts/automating-my-now-page/)

A script to automate my [/now page](https://akashgoswami.com/now).

Inspired by [Robb Knight](https://rknight.me/blog/automating-my-now-page/) and [Sophie Koonin](https://localghost.dev/blog/everything-should-have-an-api-adventures-in-trying-to-automate-stuff/). 

## Usage

1. Rename `.env_example` to `.env`
2. Update `.env` with your username for each service and your Oku collection ID
3. Rename `static_example.md` to `static.md`
4. Update the contents of `static.md` to your liking
5. Run script with `go run ./`

## Data sources
Getting data from various services isn't easy. A lot of services and platforms do not offer a viable way of retrieving your data. 

Instead I've found a number of other services that allow me to obtain and track this data.

### 🍿 [Letterboxd](https://letterboxd.com/)
Letterboxd provides an RSS feed for each of its users. 

**Adding movies to the RSS feed**

To add movies to your Letterboxd RSS feed, you need to click the 'Review or log' button for a movie and add a date for when you watched the movie.

**Viewing the RSS feed**

Via the Letterboxd website, log in and then view your profile. 
There is an RSS icon at the end of the profile navigation menu. Clicking this will take you to the profile's RSS feed.

The RSS url is formatted in the following format:

`https://letterboxd.com/USERNAME_HERE/rss/`

Simply replace the USERNAME_HERE part with any Letterboxd username to view their RSS feed.

### 📚 [Oku](https://oku.club)
Oku provides [an RSS feed for each collection](https://oku.club/blog/oku-has-rss-feeds) for a user. By default users of Oku will have three collections; To Read, Reading and Read.

Finding this feed url isn't super straight forward as there's no icon for it in the web app. Intead to get this url, you'll need to right click the page and select 'Inspect' (or View Page Source) and then under the body tag you should see a line like this:

`<link rel="alternate" type="application/rss+xml" href="https://oku.club/rss/collection/UNIQUE_STRING_HERE">`

The url listed here is the RSS feed for your collection - Most folks will want to use the default 'Reading' collection for this script.

### 📺 [Serializd](https://www.serializd.com/)
Serializd does **not** provide an RSS feed for users but is it possible to retrieve data in a JSON format using their API.

For this script I've used the Diary API endpoint (`https://www.serializd.com/api/user/USERNAME_HERE/diary`) to retrieve recently watched TV shows that have been logged.

**Logging episodes to your diary**

To add an episode to your diary, click the 'Log/Review' button for the episode and add a date for when you watched the episode.

### 🎮 [Backloggd](https://www.backloggd.com)
Backloggd has no RSS feed or an API that be be used to get the games a user is currently playing. Instead this data can be obtained via web scraping a user's public 'playing' page.

**Adding games to your playing page**

Click the 'Playing' button on a games page on Backlogged to add it to your 'playing' page.

### ✈️ [Nomadlist](https://nomadlist.com/)
Nomadlist has a public endpoint for user profiles that displays data in JSON. User profiles contain a history of past countries visited/trips.

This data can be accessed via `https://nomadlist.com/USERNAME_HERE.json`

**Adding a trip to your trip history**

Go to your profile and fill out the trip start, end date, city and then click Add trip. 