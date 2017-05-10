# Deployment

## Docker

The docker-compose files can be used to quickly deploy the website without
needing other dependencies than docker and docker-compose.

We offer two database back-end (but that might change to only postgresql later).

> NOTE: Use the *-prod* version to deploy in production. See the section
> [production](#production).

### Usage

The first step depends on the back-end chosen.

For the **mysql** back-end, you need the database file inside the project's
top-level directory and named *nyaa.db*.

For the **postgresql** back-end, you need the database dump inside the project's
top-level directory and named nyaa\_psql.backup.

You may now start the container as such.

```
$ export GOPATH=$HOME/.go
$ mkdir -p $HOME/.go/src/github.com/ewhal
$ cd $HOME/.go/src/github.com/ewhal
$ git clone https://github.com/ewhal/nyaa
$ cd nyaa/deploy
$ docker-compose -f <docker_compose_file> up
```

The website will be available at [http://localhost:9999](http://localhost:9999).

> NOTE: The website might crash if the database takes longer than the amount of
> time sleeping in the [init.sh](init.sh) script.

> NOTE: The docker-compose file uses version 3, but doesn't yet use any feature
> from the version 3. If you're getting an error because your version of docker
> is too low, you can try changing the version to '2' in the compose file.

### Production

This is specific to the
[docker-compose.postgres-prod.yml](docker-compose.postgres-prod.yml) compose
file. This should be used in production.

This setup uses an external postgresql database configured on the host machine
instead of using a container. You must therefore install and configure
postgresql in order to use this compose file.

Set the correct database parameters in [postgres-prod.env](postgres-prod.env).
You can then follow the steps above.

### Cleaning docker containers

Docker can end up taking a lot of space very quickly. The script
[prune\_docker.sh](prune_docker.sh) will get rid of unused docker images and
volumes.

## Ansible

> IMPORTANT: Make sure the website connects to pgpool's port. Otherwise, no
> caching will be done. Ansible assume you have a user on the remote that has
> sudo (no password).

You'll have to change a few variables in [hosts](host)

### Setup server

This script is installs docker (might be dropped?). It also install postgresql
with pgpool-ii and configure it. There is also a backup script installed that
creates dump of the db and seeds it.

> NOTE: The backup script needs to have access to a GPG key to sign the dumps.
> It also needs a file with the passphrase, see
> [group_vars/all](group_vars/all).

```
$ cd ansible/
$ ansible-playbook -i hosts setup_server.yml
```

## TODOs
- Delete .torrents after X days
- Add public keys to db (?)
- Show public keys and link to .torrents on the website 
