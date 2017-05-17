![nyanpasu~](https://my.mixtape.moe/aglaxe.png)

# Nyaa replacement [![Build Status](https://travis-ci.org/NyaaPantsu/nyaa.svg?branch=master)](https://travis-ci.org/NyaaPantsu/nyaa)

## Motivation
The aim of this project is to write a fully featured nyaa replacement in golang
that anyone will be able to deploy locally or remotely.

## [Roadmap](https://trello.com/b/gMJBwoRq/nyaa-pantsu-cat-roadmap)
The Roadmap will give you an overview of the features and tasks that the project are currently discussing, working on and have completed.
If you are looking for a feature that is not listed just make a GitHub Issue and it will get added to the trello board.

You can view the public trello board [here](https://trello.com/b/gMJBwoRq/nyaa-pantsu-cat-roadmap) or click on the "Roadmap".

# Requirements
* Golang

# Installation
Ubuntu 17.04 fails to build, use a different OS or docker
* Install [Golang](https://golang.org/doc/install) (version >=1.8)
* `go get github.com/NyaaPantsu/nyaa`
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
# Make sure the project is in here $GOPATH/src/github.com/NyaaPantsu/nyaa
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
* Improve tools for users
  * Editing of uploaded torrents
  * soft deletion of uploaded torrents
* Scraping of fan subbing RSS feeds similar to metainfo_fetcher and scraper
  * nyaa.si
  * anidex.moe
* Site theme
  * original nyaa theme
  * Implement mockup design from /g/anon

* Use elastic search or sphinix search
* Use new db abstraction layer and remove all ORM code
* API improvement
* Get code up to standard of go lint recommendations
* Write tests


# LICENSE
This project is licensed under the MIT License - see the LICENSE.md file for details

# Disclaimer
I take no legal responsibility for anything this code is used for. This is an purely an educational proof of concept.
