# coding: utf-8
from elasticsearch import Elasticsearch, helpers
import psycopg2, pprint, sys, time, os

CHUNK_SIZE = 10000

def getEnvOrExit(var):
    environment = ''
    try:
        environment = os.environ[var]
    except:
        print('[Error]: Environment variable ' + var + ' not defined.')
        sys.exit(1)
    return environment

dbparams = getEnvOrExit('PANTSU_DBPARAMS')
pantsu_index = getEnvOrExit('PANTSU_ELASTICSEARCH_INDEX')
torrent_tablename = getEnvOrExit('PANTSU_TORRENT_TABLENAME')

es = Elasticsearch()
pgconn = psycopg2.connect(dbparams)

cur = pgconn.cursor()
cur.execute("""SELECT torrent_id, torrent_name, description, hidden, category, sub_category, status, 
                      torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed, language
               FROM {torrent_tablename}
               WHERE deleted_at IS NULL""".format(torrent_tablename=torrent_tablename))

fetches = cur.fetchmany(CHUNK_SIZE)
while fetches:
    actions = list()
    for torrent_id, torrent_name, description, hidden, category, sub_category, status, torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed, language in fetches:
        doc = {
          'id': torrent_id,
          'name': torrent_name.decode('utf-8'),
          'category': str(category),
          'sub_category': str(sub_category),
          'status': status,
          'hash': torrent_hash,
          'hidden': hidden,
          'description': description,
          'date': date,
          'uploader_id': uploader,
          'downloads': downloads,
          'filesize': filesize,
          'seeders': seeders,
          'leechers': leechers,
          'completed': completed,
          'language': language
        }
        action = {
            '_index': pantsu_index,
            '_type': 'torrents',
            '_id': torrent_id,
            '_source': doc
        }
        actions.append(action)
    helpers.bulk(es, actions, chunk_size=CHUNK_SIZE, request_timeout=120)
    del(fetches)
    fetches = cur.fetchmany(CHUNK_SIZE)
cur.close()
pgconn.close()
