# Olive Helps Loop: Network Monitor

This is a sample Loop for Olive Helps. Copy any URL to the clipboard and this will send a whisper if the URL becomes unavailable.

It is useful for having it monitor an internal company site, and it will notify you if your VPN connection is lost.

## Requirements

You will need to install [Olive Helps](https://oliveai.com/olive-helps/).

Install [golang 1.16](https://golang.org/), clone this repo, and build the project with:
```shell
make build
```

After that is complete, you can start Olive Helps and install the contents of the `./build` directory as a Local Loop.

Sample output (read bottom to top):

![Whispers from Netmon Olive Helps Loop](netmon_whispers.png)
