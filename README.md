# VPLAY

Vplay is designed to load videos into a video grid


## Folder Player

`vplay_folder <FOLDER>`

will load all the videos into vplay wall

## Videos from Browser

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

6. now browser streaming videos site they'll load up in vplay wall

