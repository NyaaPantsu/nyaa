# coding: utf-8
from elasticsearch import Elasticsearch, helpers
import psycopg2, pprint, sys, time, os

CHUNK_SIZE = 10000

dbparams = ''
pantsu_index = ''

try:
    dbparams = os.environ['PANTSU_DBPARAMS']
except:
    print('[Error]: Environment variable PANTSU_DBPARAMS not defined.')
    sys.exit(1)

try:
    pantsu_index = os.environ['PANTSU_ELASTICSEARCH_INDEX']
except:
    print('[Error]: Environment variable PANTSU_ELASTICSEARCH_INDEX not defined.')
    sys.exit(1)

es = Elasticsearch()
pgconn = psycopg2.connect(dbparams)

cur = pgconn.cursor()
cur.execute("""SELECT torrent_id, torrent_name, category, sub_category, status, 
                      torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed
               FROM torrents
               WHERE deleted_at IS NULL""")

fetches = cur.fetchmany(CHUNK_SIZE)
while fetches:
    actions = list()
    for torrent_id, torrent_name, category, sub_category, status, torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed in fetches:
        doc = {
          'id': torrent_id,
          'name': torrent_name.decode('utf-8'),
          'category': str(category),
          'sub_category': str(sub_category),
          'status': status,
          'hash': torrent_hash,
          'date': date,
          'uploader_id': uploader,
          'downloads': downloads,
          'filesize': filesize,
          'seeders': seeders,
          'leechers': leechers,
          'completed': completed
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
