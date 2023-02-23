# VPLAY

Vplay is designed to load videos into a video grid

Built for viewing porn in mind

## Porn Folder Player

`vplay_folder <FOLDER>`

will load all the videos into vplay wall

## Porn Videos from Browser

Works best with pornhub

1. launch chrome with experimental apis

`/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --enable-experimental-extension-apis --enable-logging`

2. load up extension in chrome got to chrome://extensions/ and hit the Load Unpacked point it to the chrome-extension/app folder
 take note of the extension id `ghbmnnjooekpmoecnnnilnnbdlolhkhi`
3. edit chrome-extension/config/com.sample.native_msg_golang.json
    - fix `path` to your vplay_browser executable
    - fix `allowed_origins` with the extension id from before
4. cp /Users/stan/dev/vplay/chrome-extension/config/com.sample.native_msg_golang.json \
   /Users/stan/Library/Application\ Support/Google/Chrome/NativeMessagingHosts/



5. go to chrome-extension://ghbmnnjooekpmoecnnnilnnbdlolhkhi/main.html
hit **Connect To Native Host**

6. now browser pornhub videos they'll load up in vplay wall

## Pornhub url load (WIP)

Load all the videos from a pornhub category or playlist or model page

Instead use these commands to download to a folder and then use `vplay_folder` for now

```
P='https://www.pornhub.org/pornstar/goldie-baby'

yt-dlp --get-id $P | xargs -I {} -P 4 yt-dlp --downloader aria2c -N 10 --retry-sleep 1 -i --fixup never --force-ipv4 --no-check-certificate \
'https://www.pornhub.org/view_video.php?viewkey={}'
````
