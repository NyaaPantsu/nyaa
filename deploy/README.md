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
$ mkdir -p $HOME/.go/src/github.com/NyaaPantsu
$ cd $HOME/.go/src/github.com/NyaaPantsu
$ git clone https://github.com/NyaaPantsu/nyaa
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

You'll have to change a few variables in [hosts](host). Replace the host:ip
address to the host:ip of the target server. You can also change the user
ansible uses to connect to the server. The user needs to have sudo ALL.

You'll also maybe have to tweak a few variables in
[group_vars/all](group_vars/all) such as the database password, etc (but should
probably be left like this).


### Setup server playbook

This playbook installs and configure:

- postgresql (It also includes pgpool for caching)
- firewalld
- golang
- elasticsearch
- backup system (uses cronjob to do daily backup of the database)

> NOTE: The backup script needs to have access to a GPG key to sign the dumps.
> It also needs a file with the passphrase, see
> [group_vars/all](group_vars/all).

```
$ ansible-playbook -i hosts setup_server.yml
```


### Restore Database Playbook

This playbook restores a database from dump. The dump has to be named
nyaa_psql.backup and needs to be placed in the toplevel project directory *on
your local host*. The database will be copied to the remote host and then will
be restored.

```
$ ansible-playbook -i hosts restore_database.yml
```


### Create Elasticsearch Index Playbook

This playbook creates the elasticsearch index for our database from
[ansible/roles/elasticsearch/files/elasticsearch_settings.yml](ansible/roles/elasticsearch/files/elasticsearch_settings.yml)

```
$ ansible-playbook -i hosts create_elasticsearch_index.yml
```


### Populate Elasticsearch Index Playbook

This playbook uses a python script to populate the elasticsearch index from all
the data inside the database.

> WARNING: Make sure the python script is in sync with the mapping defined in
> the elasticsearch index configuration.

```
$ ansible-playbook -i hosts populate_elasticsearch_index.yml
```

## Playbook Testing

You can easily test these playbooks by using vagrant. Once you have vagrant
installed:

```
# Download centos/7 image
$ vagrant  init centos/7

# Create and boot the vm
$ vagrant up
$ vagrant ssh
```

Now you have to setup your host to be able to connect to the vm using ssh. One
way is to copy your public ssh key to the `~/.ssh/authorized_keys` file. Once
that is done, your local host should be able to connect to the vm using ssh.

You can now tests the playbooks.

## TODOs
- Delete .torrents after X days
- Add public keys to db (?)
- Show public keys and link to .torrents on the website
- Tuning elasticsearch indexing / analyzer
