# Fansly datascraper

## Config
you need to setup a config so you can use this tool.

windows: `%APPDATA%\fanslyio\config.json`

macOS: `$HOME/Library/Application Support/fanslyio/config.json`

*nix: `$HOME/.config/fanslyio/config.json`


### To get `user-agent` text
login to the page. open `Developer tools` go to the console:
type this: `alert(window.navigator.userAgent)` copy the text from that box
### To get `authorization` text
type this: `alert(JSON.parse(localStorage.getItem("session_active_session")).token)` copy the text from that box

```json
{
    "user-agent": "Add user-agent from browser",
    "authorization": "Add authorization value from browser"
}
```

If you only want to download from specific users, add `allowlist` to the file like this:
```json
{
    "user-agent": "Add user-agent from browser",
    "authorization": "Add authorization value from browser",
    "allowlist": ["username", "username2"]
}
```

If you don't want to download from specific users, add `blocklist` to the file like this:
```json
{
    "user-agent": "Add user-agent from browser",
    "authorization": "Add authorization value from browser",
    "blocklist": ["username3"]
}
```
