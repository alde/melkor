# Melkor

[![Build Status][travis-image]](https://travis-ci.org/alde/melkor)
[![Coverage Status][coveralls-image]](https://coveralls.io/r/alde/melkor?branch=master)
[![Go Report][goreport-image]](https://goreportcard.com/report/github.com/alde/melkor)

## Purpose
Melkor is a caching layer for AWS, inspired by [Edda](https://github.com/Netflix/Edda) but intended to be simpler.

**Note: This is early prototype code, with very limited functionality. It might even go nowhere. **

## Crawlers
Crawlers are meant to periodically scrape the AWS api and put it into a cache.

## API
Get all items:

    /v1/aws/{collection}

Get a limited number of items:

    /v1/aws/{collection}?_limit=1

Get a list of expanded items:

    /v1/aws/{collection}?_expand=true

Get a single item:

    /v1/aws/{collection}/{id}

# Contributors
- Rickard Dybeck ([alde](https://github.com/alde))

## License
[Licence](./LICENSE)

[travis-image]: https://img.shields.io/travis/alde/melkor.svg?style=flat
[coveralls-image]: https://img.shields.io/coveralls/alde/melkor.svg?style=flat
[goreport-image]: https://goreportcard.com/badge/github.com/alde/melkor
