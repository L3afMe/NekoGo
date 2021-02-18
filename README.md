<div align="center">
    <h1>NekoGo</h1>
    <p>A simple to use Discord selfbot with a focus on speed (that isn't written in Python or Javascript!)</p>
    <p>
        <a><img alt="Latest Release" src="https://img.shields.io/github/v/release/L3afMe/NekoGo?color=%23bd93f9&include_prereleases&style=for-the-badge"></a>
        <a href="https://github.com/L3afMe/NekoGo/actions?query=workflow%3A%22Go+Build%22"><img alt="Build Status Badge" src="https://img.shields.io/github/workflow/status/L3afMe/NekoGo/Go%20Build/master?color=%23bd93f9&style=for-the-badge"></a>
        <img alt="Lines of code" src="https://img.shields.io/tokei/lines/github/L3afMe/NekoGo?color=%23bd93f9&style=for-the-badge">
        <a href="https://github.com/L3afMe/NekoGo/blob/master/LICENSE"><img alt="License Badge" src="https://img.shields.io/github/license/L3afMe/NekoGo?color=%23bd93f9&style=for-the-badge"></a>
    </div>
</div>

## Contents
- [Warning](#warning)
- [Features](#features)
- [Commands](#commands)
- [How To Use](#how-to-use)
  - [Prebuilt Binaries](#prebuilt-binaries)
  - [Building From Source](#building-from-source)

## Warning
### Using a selfbot is explicitly against Discord's TOS, by using NekoGo you acknowledge that this may lead to your account being permanently terminated. While this has never happened to anyone during testing, there will always be a risk that it could happen at any time.

### NekoGo is currently in pre-release. This means there may be bugs, unexpected behaviour, and features that don't work or that get removed.

## Features
#### Features with a ~~strikethrough~~ aren't currently implemented and will be added in the coming future
- Written in Golang so it is much faster than anything in Python/JavaScript (Which most selfbots are in)
- Thoroughly documented commands. Every command, subcommand, subsubcommand, etc has it's own help menu with examples
- Aliases galore. Too lazy to type the whole command name? Check help and there will likely be a much shorter alias
- Everything is easily configurable through chat so no need to mess with setting up config files. Setup your token the first run and you never have to worry about it again
- Interaction gifs like `kiss`, `slap`, `hug`, `poke`, etc
- ~~Iterpeters for several languages including JavaScript, Python, Brainfuck and more~~
- ~~Image generation using a mentioned user's profile picture~~
- ~~Mention and keyword logging as well as message sniper~~
- ~~Information commands like `serverinfo`, `userinfo`, `roleinfo`, `channelinfo`~~
- ~~Search words with `urbandictionary`, `wikipedia`, `dictionary`~~
- ~~Steal emotes in chat and add them to your own server with `emotestealer`~~
- ~~Automatically switch your avatar, name, and tag (if you have Nitro) at a certain interval~~

## Commands
A list of commands can be found [in the Wiki](https://github.com/L3afMe/NekoGo/wiki/Commands). Please note this may not be updated, to get a list of all current commands, check the help menu in Discord.

## How To Use
### Prebuilt Binaries
#### Windows/Linux
- Prerequisites
  - None!
- Steps
  1) Download the latest binaries from the [Releases](https://github.com/L3afMe/NekoGo/releases) page
  2) Run the downladed file and input your token when prompted
  3) Profit?
- Notes
  - On Linux you will need to run `chmod +x NekoGo-linux` to make the file executable

#### Linux Server
If you don't want to use `pm2` and prefer something like `screen`, skip step `i` and `vii`

- Prerequisites
  - NodeJS (Used for pm2)
- Steps
  1) Install `pm2`, this will ensure that NekoGo keeps running after leaving SSH and (although it shouldn't happen) restart NekoGo if it crashes
    - `npm i -g pm2`
  2) Make a new directory and move into it
    - `mkdir NekoGo && cd NekoGo`
  3) Download the latest release (Replace the URL with the latest from [Releases](https://github.com/L3afMe/NekoGo/releases) if I forget to update it)
    - `wget URL -o NekoGo`
  4) Allow the file to be run
    - `chmod +x NekoGo`
  5) Run it once to setup the config
    - `./NekoGo`
  6) After config has been set up and it's running. Press Control-C to stop it so we can run it in pm2 now
  7) Start a new `pm2` process
    - `pm2 start NekoGo`
  8) Profit?
- Notes
  - To check logs use `pm2 logs NekoGo`

### Building From Source
- Prerequisites
  - [Golang](https://golang.org/doc/install)
  - [Git](https://git-scm.com/downloads)
  - Make (Optional - Used for cross-platform building on Linux)
- Steps
  1) Clone the repository
    - `git clone https://github.com/L3afMe/NekoGo.git`
  2) Move into the repository
    - `cd NekoGo`
  3) Build the binarie(s)
    - Linux (Cross-platform): `make build`
    - Linux (Host platform): `go build -o bin/NekoGo-linux *.go`
    - Windows (Host platform): `go build -o bin/NekoGo-windows.exe *.go`
  4) The built file will be in `bin/`
  5) Continue from [running Prebuilt Binaries](#prebuilt-binaries)
