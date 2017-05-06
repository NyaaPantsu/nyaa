![nyanpasu~](https://a.safe.moe/dNEEh.png)

# Nyaa replacement [![Build Status](https://travis-ci.org/ewhal/nyaa.svg?branch=master)](https://travis-ci.org/ewhal/nyaa)

## Motivation
The aim of this project is to write a fully featured nyaa replacement in golang
that anyone will be able to deploy locally or remotely.

# Requirements
* Golang

# Installation
* Install [Golang](https://golang.org/doc/install)
* `go get github.com/ewhal/nyaa`
* `go build`
* Download DB and place it in your root folder named as "nyaa.db"
* `./nyaa`
* You can now access your local site over on [localhost:9999](http://localhost:9999)

## Usage

Type `./nyaa -h` for the list of options.

After modifying the files in `./templates`, run `go generate ./... && go build`.

## Systemd

* Edit the unit file `os/nyaa.service` to your liking
* Copy the package's content so that your unit file can find them.
* Copy the unit file in `/usr/lib/systemd/system`
* `systemctl daemon-reload`
* `systemctl start nyaa`

The provided unit file uses options directly; if you prefer a config file, do the following:

* `./nyaa -print-defaults > /etc/nyaa.conf`
* Edit `nyaa.conf` to your liking
* Replace in the unit file the options by `-conf /etc/nyaa.conf`


## Docker

We support docker for easy development and deployment. Simply install docker and
docker-compose by following the instructions [here](https://docs.docker.com/engine/installation/linux/ubuntu/#install-using-the-repository).

Once you've successfully installed docker, make sure you have the database file
in the project's directory as nyaa.db. Then, follow these steps to build and run
the application.

```sh
$ cd deploy/
# You may choose another backend by pointing to the
# appropriate docker-compose file.
$ docker-compose -f docker-compose.sqlite.yml build
$ docker-compose -f docker-compose.sqlite.yml up 
```

Access the website by going to [localhost:9999](http://localhost:9999).

## TODO

### Features until stable release
* Uploading of new Magnet links
* Make sukebei db schema compatible with current code
* Fix up comments display
* custom 404 page
## Other features
* Add support for postgres dbs
* Torrent data scraping from definable tracker
* Report Feature and Moderation System
  * blocking upload of torrent hashes
* API improvement
* Scraping
* Accounts and Registration System
* Daily DB dumps
* Site theme
* p2p sync of dbs?

# LICENSE
This project is licensed under the MIT License - see the LICENSE.md file for details
