# kavideos

This project aims to be an easy to use bot and api for downloading videos from various sources, like static web pages, dynamic web pages that loads their media via JavaScript (as long as they use an reconizable extension), kwai and a lot more (see Cobalt page for full support).

For use with Cobalt, is advivised that you self host your own instance, then just set the envs COBALT_URL and COBALT_TOKEN to their respective values.

## Telegram Bot

You can send any link and the bot will try to match to kwai or one of cobalt supported services. You can also use /download to search in a static webpage for videos, or /browser if the web page is dynamically generated.

Also, you need to set the env var BOT_TOKEN to your bot token.

## API

The api support is very basic, you just need to send a request to /download. The body should be a JSON with the following schema:

- url (required): The url of the website/video that you want to download
- isWebpage (optional): A boolean to tell the service to search for videos on the html of the url
- useBrowser (optional): When combined with `isWebpage`, will search for the page contents loading it on a web browser.

# TODO

- [ ] Use real browser headers
- [ ] Send the initial webpage cookies when downloading the media on webpages
