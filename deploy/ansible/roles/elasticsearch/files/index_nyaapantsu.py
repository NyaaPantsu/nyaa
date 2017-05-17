# coding: utf-8
from elasticsearch import Elasticsearch, helpers
import psycopg2, pprint, sys, time, os

PANTSU_INDEX = 'pantsu'
CHUNK_SIZE   = 15000

dbparams = ''
try:
    dbparams = os.environ['PANTSU_DBPARAMS']
except:
    print('[Error]: Environment variable PANTSU_DBPARAMS not defined.')
    sys.exit(1)

es = Elasticsearch()
pgconn = psycopg2.connect(dbparams)

cur = pgconn.cursor()
cur.execute("""select torrent_id, torrent_name, category, sub_category, status 
               from torrents 
               WHERE deleted_at IS NULL""")

fetches = cur.fetchmany(CHUNK_SIZE)
while fetches:
    actions = list()
    for torrent_id, torrent_name, category, sub_category, status in fetches:
        doc = {
          'torrent_id': torrent_id,
          'torrent_name': torrent_name.decode('utf-8'),
          'category': category,
          'sub_category': sub_category,
          'status': status
        }
        action = {
            '_index': PANTSU_INDEX,
            '_type': 'document',
            '_id': torrent_id,
            '_source': doc
        }
        actions.append(action)
    helpers.bulk(es, actions, chunk_size=CHUNK_SIZE, request_timeout=120)
    del(fetches)
    fetches = cur.fetchmany(CHUNK_SIZE)
