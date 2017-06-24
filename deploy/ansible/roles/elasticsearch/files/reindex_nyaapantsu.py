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
# We MUST use NO QUERY CACHE because the values are insert on triggers and
# not through pgppool.
cur.execute('/*NO QUERY CACHE*/ SELECT reindex_torrents_id, torrent_id, action FROM reindex_{torrent_tablename}'.format(torrent_tablename=torrent_tablename))

fetches = cur.fetchmany(CHUNK_SIZE)
while fetches:
    actions = list()
    delete_cur = pgconn.cursor()
    for reindex_id, torrent_id, action in fetches:
        new_action = {
          '_op_type': action,
          '_index': pantsu_index,
          '_type': 'torrents',
          '_id': torrent_id
        }
        if action == 'index':
            select_cur = pgconn.cursor()
            select_cur.execute("""SELECT torrent_id, torrent_name, description, hidden, category, sub_category, status,
                                  torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed, language
                           FROM {torrent_tablename}
                           WHERE torrent_id = {torrent_id}""".format(torrent_id=torrent_id, torrent_tablename=torrent_tablename))
            torrent_id, torrent_name, description, hidden, category, sub_category, status, torrent_hash, date, uploader, downloads, filesize, seeders, leechers, completed, language = select_cur.fetchone()
            doc = {
              'id': torrent_id,
              'name': torrent_name.decode('utf-8'),
              'category': str(category),
              'sub_category': str(sub_category),
              'status': status,
              'hidden': hidden,
              'description': description,
              'hash': torrent_hash,
              'date': date,
              'uploader_id': uploader,
              'downloads': downloads,
              'filesize': filesize,
              'seeders': seeders,
              'leechers': leechers,
              'completed': completed,
              'language': language
            }
            new_action['_source'] = doc
            select_cur.close()
        delete_cur.execute('DELETE FROM reindex_{torrent_tablename} WHERE reindex_torrents_id = {reindex_id}'.format(reindex_id=reindex_id,torrent_tablename=torrent_tablename))
        actions.append(new_action)
    pgconn.commit() # Commit the deletes transaction
    delete_cur.close()
    helpers.bulk(es, actions, chunk_size=CHUNK_SIZE, request_timeout=120)
    del(fetches)
    fetches = cur.fetchmany(CHUNK_SIZE)
cur.close()
pgconn.close()
