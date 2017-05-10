![nyanpasu~](https://my.mixtape.moe/aglaxe.png)

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

Windows Users If you get `"standard_init_linux.go:178: exec user process caused "no such file or directory"`
download [dos2unix](https://sourceforge.net/projects/dos2unix/files/latest/download) and run "dos2unix.exe"
on the /deploy/init.sh to convert CR+LF to LF.

```sh
# Make sure the project is in here $GOPATH/src/github.com/ewhal/nyaa
$ cd deploy/
# You may choose another backend by pointing to the
# appropriate docker-compose file.
$ docker-compose -f docker-compose.sqlite.yml build
$ docker-compose -f docker-compose.sqlite.yml up 
```

Access the website by going to [localhost:9999](http://localhost:9999).

> For postgres, place the dump in the toplevel directory and name it to
> nyaa_psql.backup.

## TODO
## High priority
* Torrent data scraping from definable tracker (We have a tracker that the owner is ok for us to scrape from)
  * seeds/leeachers
  * file lists
  * Downloads
* Accounts and Registration System(WIP)
  * blocking upload of torrent hashes

## Lower priority
 * Get code up to standard of go lint recommendations
 * Write tests
 * fix sukebei categories
 * Daily DB dumps
 * Site theme
   * original nyaa theme
 * API improvement
 * Scraping of fan subbing RSS feeds

# LICENSE
This project is licensed under the MIT License - see the LICENSE.md file for details
