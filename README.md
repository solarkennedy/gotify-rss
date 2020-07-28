# gotify-rss [![Build Status](https://travis-ci.org/solarkennedy/gotify-rss.svg?branch=master)](https://travis-ci.org/solarkennedy/gotify-rss)

A plugin for [gotify/server](https://github.com/gotify/server) which polls RSS feed. Forked from [gotify-archsec](https://github.com/buckket/gotify-archsec).

## Building

For building the plugin gotify/build docker images are used to ensure compatibility with 
[gotify/server](https://github.com/gotify/server).

`GOTIFY_VERSION` can be a tag, commit or branch from the gotify/server repository.

This command builds the plugin for amd64, arm-7 and arm64. 
The resulting shared object will be compatible with gotify/server version 2.0.5.
```bash
$ make GOTIFY_VERSION="v2.0.5" FILE_SUFFIX="for-gotify-v2.0.5" build
```

## Installation

Copy built shared object to the gotify plugin directory and restart gotify.

## Configuration

- `refresh_interval`: Polling interval in seconds
- `feed_url`: URL of the RSS feed

## License

GNU GPLv3+
