# Automate /now

My attempts at automating my [/now page](https://akashgoswami.com/now).

Inspired by [Robb Knight](https://rknight.me/blog/automating-my-now-page/) and [Sophie Koonin](https://localghost.dev/blog/everything-should-have-an-api-adventures-in-trying-to-automate-stuff/). 

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