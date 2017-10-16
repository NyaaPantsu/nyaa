define({
  "api": [{
      "type": "get",
      "url": "/search/",
      "title": "Search Torrents",
      "version": "1.1.1",
      "name": "FindTorrents",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "c",
              "description": "<p>In which categories to search.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "q",
              "description": "<p>Query to search (torrent name).</p>"
            },
            {
              "group": "Parameter",
              "type": "Number",
              "optional": false,
              "field": "page",
              "description": "<p>Page of the search results.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "limit",
              "description": "<p>Number of results per page.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "userID",
              "description": "<p>Uploader ID owning the torrents.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "fromID",
              "description": "<p>Show results with torrents ID superior to this.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "s",
              "description": "<p>Torrent status.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "maxage",
              "description": "<p>Torrents which have been uploaded the last x days.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "toDate",
              "description": "<p>Torrents which have been uploaded since x <code>dateType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "fromDate",
              "description": "<p>Torrents which have been uploaded the last x <code>dateType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "dateType",
              "description": "<p>Which type of date (<code>d</code> for days, <code>m</code> for months, <code>y</code> for years).</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "minSize",
              "description": "<p>Filter by minimal size in <code>sizeType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "maxSize",
              "description": "<p>Filter by maximal size in <code>sizeType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "sizeType",
              "description": "<p>Which type of size (<code>b</code> for bytes, <code>k</code> for kilobytes, <code>m</code> for megabytes, <code>g</code> for gigabytes).</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "sort",
              "description": "<p>Torrent sorting type (0 = id, 1 = name, 2 = date, 3 = downloads, 4 = size, 5 = seeders, 6 = leechers, 7 = completed).</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "order",
              "description": "<p>Order ascending or descending (true = ascending).</p>"
            },
            {
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "lang",
              "description": "<p>Filter the languages.</p>"
            }
          ]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
            "group": "Success 200",
            "type": "Object[]",
            "optional": false,
            "field": "torrents",
            "description": "<p>List of torrent object (see view for the properties).</p>"
          }]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n\t\t\t\"torrents\": [...],\n\t\t\t\"queryRecordCount\": 50,\n\t\t\t\"totalRecordCount\": 798414\n\t\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
            "group": "Error 4xx",
            "type": "String[]",
            "optional": false,
            "field": "errors",
            "description": "<p>List of errors messages with a 404 error message in it.</p>"
          }]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "HTTP/1.1 404 Not Found\n{\n  \"errors\": [ \"404_not_found\", ... ]\n}",
          "type": "json"
        }]
      }
    },
    {
      "type": "get",
      "url": "/search/",
      "title": "Search Torrents",
      "version": "1.0.0",
      "name": "FindTorrents",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "c",
              "description": "<p>In which categories to search.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "q",
              "description": "<p>Query to search (torrent name).</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "limit",
              "description": "<p>Number of results per page.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "userID",
              "description": "<p>Uploader ID owning the torrents.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "fromID",
              "description": "<p>Show results with torrents ID superior to this.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "s",
              "description": "<p>Torrent status.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "maxage",
              "description": "<p>Torrents which have been uploaded the last x days.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "toDate",
              "description": "<p>Torrents which have been uploaded since x <code>dateType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "fromDate",
              "description": "<p>Torrents which have been uploaded the last x <code>dateType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "dateType",
              "description": "<p>Which type of date (<code>d</code> for days, <code>m</code> for months, <code>y</code> for years).</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "minSize",
              "description": "<p>Filter by minimal size in <code>sizeType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "maxSize",
              "description": "<p>Filter by maximal size in <code>sizeType</code>.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "sizeType",
              "description": "<p>Which type of size (<code>b</code> for bytes, <code>k</code> for kilobytes, <code>m</code> for megabytes, <code>g</code> for gigabytes).</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "sort",
              "description": "<p>Torrent sorting type (0 = id, 1 = name, 2 = date, 3 = downloads, 4 = size, 5 = seeders, 6 = leechers, 7 = completed).</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "order",
              "description": "<p>Order ascending or descending (true = ascending).</p>"
            },
            {
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "lang",
              "description": "<p>Filter the languages.</p>"
            },
            {
              "group": "Parameter",
              "type": "Number",
              "optional": false,
              "field": "page",
              "description": "<p>Search page.</p>"
            }
          ]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
            "group": "Success 200",
            "type": "Object[]",
            "optional": false,
            "field": "torrents",
            "description": "<p>List of torrent object (see view for the properties).</p>"
          }]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n\t\t\t[...]\n\t\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/_apidoc.js",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
            "group": "Error 4xx",
            "type": "String[]",
            "optional": false,
            "field": "errors",
            "description": "<p>List of errors messages with a 404 error message in it.</p>"
          }]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "HTTP/1.1 404 Not Found\n{\n  \"errors\": [ \"404_not_found\", ... ]\n}",
          "type": "json"
        }]
      }
    },
    {
      "type": "get",
      "url": "/view/:id",
      "title": "Request Torrent information",
      "version": "1.1.1",
      "name": "GetTorrent",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>Torrent unique ID.</p>"
          }]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "id",
              "description": "<p>ID of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "name",
              "description": "<p>Name of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "status",
              "description": "<p>Status of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "hash",
              "description": "<p>Hash of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Date",
              "optional": false,
              "field": "date",
              "description": "<p>Uploaded date of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "filesize",
              "description": "<p>File size in Bytes of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "description",
              "description": "<p>Description of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Object[]",
              "optional": false,
              "field": "comments",
              "description": "<p>Comments of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "sub_category",
              "description": "<p>Sub Category of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "category",
              "description": "<p>Category of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "anidb_id",
              "description": "<p>Anidb ID of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "uploader_id",
              "description": "<p>ID of the torrent uploader.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "uploader_name",
              "description": "<p>Username of the torrent uploader.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "uploader_old",
              "description": "<p>Old username from nyaa of the torrent uploader.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "website_link",
              "description": "<p>External Link of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String[]",
              "optional": false,
              "field": "languages",
              "description": "<p>Languages of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "magnet",
              "description": "<p>Magnet URI of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "String",
              "optional": false,
              "field": "torrent",
              "description": "<p>Download URL of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "seeders",
              "description": "<p>Number of seeders of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "leechers",
              "description": "<p>Number of leechers of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "completed",
              "description": "<p>Downloads completed of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Date",
              "optional": false,
              "field": "last_scrape",
              "description": "<p>Last statistics update of the torrent.</p>"
            },
            {
              "group": "Success 200",
              "type": "Object[]",
              "optional": false,
              "field": "file_list",
              "description": "<p>List of files in the torrent.</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n\t{\n\t\"id\": 952801,\n\t\"name\": \"[HorribleSubs] Uchouten Kazoku S2 [720p]\",\n\t\"status\": 1,\n\t\"hash\": \"6E4D96F7A0B0456672E80B150CCB7C15868CD47D\",\n\t\"date\": \"2017-07-05T11:01:39Z\",\n\t\"filesize\": 4056160259,\n\t\"description\": \"<p>Unofficial batch</p>\\n\",\n\t\"comments\": [],\n\t\"sub_category\": \"5\",\n\t\"category\": \"3\",\n\t\"anidb_id\": \"\",\n\t\"downloads\": 0,\n\t\"uploader_id\": 7177,\n\t\"uploader_name\": \"DarAR92\",\n\t\"uploader_old\": \"\",\n\t\"website_link\": \"http://horriblesubs.info/\",\n\t\"languages\": [\n\t\"en-us\"\n\t],\n\t\"magnet\": \"magnet:?xt=urn:btih:6E4D96F7A0B0456672E80B150CCB7C15868CD47D&dn=%5BHorribleSubs%5D+Uchouten+Kazoku+S2+%5B720p%5D&tr=http://nyaa.tracker.wf:7777/announce&tr=http://nyaa.tracker.wf:7777/announce&tr=udp://tracker.doko.moe:6969&tr=http://tracker.anirena.com:80/announce&tr=http://anidex.moe:6969/announce&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://zer0day.ch:1337&tr=udp://9.rarbg.com:2710/announce&tr=udp://tracker2.christianbro.pw:6969/announce&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://eddie4.nl:6969/announce&tr=udp://tracker.doko.moe:6969/announce\",\n\t\"torrent\": \"https://nyaa.pantsu.cat/download/6E4D96F7A0B0456672E80B150CCB7C15868CD47D\",\n\t\"seeders\": 4,\n\t\"leechers\": 2,\n\t\"completed\": 28,\n\t\"last_scrape\": \"2017-07-07T07:48:32.509635Z\",\n\t\"file_list\": [\n\t{\n\t\"path\": \"[HorribleSubs] Uchouten Kazoku S2 - 01[720p].mkv\",\n\t\"filesize\": 338250895\n\t},\n\t{\n\t\"path\": \"[HorribleSubs] Uchouten Kazoku S2 - 12 [720p].mkv\",\n\t\"filesize\": 338556275\n\t}\n\t]\n\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
            "group": "Error 4xx",
            "type": "String[]",
            "optional": false,
            "field": "errors",
            "description": "<p>List of errors messages with a 404 error message in it.</p>"
          }]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "HTTP/1.1 404 Not Found\n{\n  \"errors\": [ \"404_not_found\", ... ]\n}",
          "type": "json"
        }]
      }
    },
    {
      "type": "get",
      "url": "/head/:id",
      "title": "Request Torrent Head",
      "version": "1.1.1",
      "name": "GetTorrentHead",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>Torrent unique ID.</p>"
          }]
        }
      },
      "success": {
        "examples": [{
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
            "group": "Error 4xx",
            "type": "String[]",
            "optional": false,
            "field": "errors",
            "description": "<p>List of errors messages with a 404 error message in it.</p>"
          }]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "HTTP/1.1 404 Not Found\n{\n  \"errors\": [ \"404_not_found\", ... ]\n}",
          "type": "json"
        }]
      }
    },
    {
      "type": "get",
      "url": "/",
      "title": "Request Torrents index",
      "version": "1.1.1",
      "name": "GetTorrents",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>Torrent unique ID.</p>"
          }]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Object[]",
              "optional": false,
              "field": "torrents",
              "description": "<p>List of torrent object (see view for the properties).</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "queryRecordCount",
              "description": "<p>Number of torrents given.</p>"
            },
            {
              "group": "Success 200",
              "type": "Number",
              "optional": false,
              "field": "totalRecordCount",
              "description": "<p>Total number of torrents.</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n\t\t\t\"torrents\": [...],\n\t\t\t\"queryRecordCount\": 50,\n\t\t\t\"totalRecordCount\": 798414\n\t\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
            "group": "Error 4xx",
            "type": "String[]",
            "optional": false,
            "field": "errors",
            "description": "<p>List of errors messages with a 404 error message in it.</p>"
          }]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "HTTP/1.1 404 Not Found\n{\n  \"errors\": [ \"404_not_found\", ... ]\n}",
          "type": "json"
        }]
      }
    },
    {
      "type": "post",
      "url": "/update/",
      "title": "Update a Torrent",
      "version": "1.1.1",
      "name": "UpdateTorrent",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "username",
              "description": "<p>Torrent uploader name.</p>"
            },
            {
              "group": "Parameter",
              "type": "Number",
              "optional": false,
              "field": "id",
              "description": "<p>Torrent ID.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "name",
              "description": "<p>Torrent name.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "category",
              "description": "<p>Torrent category.</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "remake",
              "description": "<p>Torrent is a remake.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "description",
              "description": "<p>Torrent description.</p>"
            },
            {
              "group": "Parameter",
              "type": "Number",
              "optional": false,
              "field": "status",
              "description": "<p>Torrent status.</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "hidden",
              "description": "<p>Torrent hidden.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "website_link",
              "description": "<p>Torrent website link.</p>"
            },
            {
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "languages",
              "description": "<p>Torrent languages.</p>"
            }
          ]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request is done without failing</p>"
            },
            {
              "group": "Success 200",
              "type": "String[]",
              "optional": false,
              "field": "infos",
              "description": "<p>Messages information relative to the request</p>"
            },
            {
              "group": "Success 200",
              "type": "Object",
              "optional": false,
              "field": "data",
              "description": "<p>The resulting torrent updated (see view for the properties)</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
              "group": "Error 4xx",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request couldn't be done due to some errors.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "String[]",
              "optional": false,
              "field": "errors",
              "description": "<p>List of errors messages.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "Object[]",
              "optional": false,
              "field": "all_errors",
              "description": "<p>List of errors object messages for each wrong field</p>"
            }
          ]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n      \"ok\": false,\n      \"errors\": [ ... ]\n      \"all_errors\": {\n\t\t \t\"username\": [ ... ],\n       }\n    }",
          "type": "json"
        }]
      }
    },
    {
      "type": "post",
      "url": "/upload",
      "title": "Upload a Torrent",
      "version": "1.1.1",
      "name": "UploadTorrent",
      "group": "Torrents",
      "parameter": {
        "fields": {
          "Parameter": [{
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "username",
              "description": "<p>Torrent uploader name.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "name",
              "description": "<p>Torrent name.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "magnet",
              "description": "<p>Torrent magnet URI.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "category",
              "description": "<p>Torrent category.</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "remake",
              "description": "<p>Torrent is a remake.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "description",
              "description": "<p>Torrent description.</p>"
            },
            {
              "group": "Parameter",
              "type": "Number",
              "optional": false,
              "field": "status",
              "description": "<p>Torrent status.</p>"
            },
            {
              "group": "Parameter",
              "type": "Boolean",
              "optional": false,
              "field": "hidden",
              "description": "<p>Torrent hidden.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "website_link",
              "description": "<p>Torrent website link.</p>"
            },
            {
              "group": "Parameter",
              "type": "String[]",
              "optional": false,
              "field": "languages",
              "description": "<p>Torrent languages.</p>"
            },
            {
              "group": "Parameter",
              "type": "File",
              "optional": false,
              "field": "torrent",
              "description": "<p>Torrent file to upload (you have to send a torrent file or a magnet, not both!).</p>"
            }
          ]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request is done without failing</p>"
            },
            {
              "group": "Success 200",
              "type": "String[]",
              "optional": false,
              "field": "infos",
              "description": "<p>Messages information relative to the request</p>"
            },
            {
              "group": "Success 200",
              "type": "Object",
              "optional": false,
              "field": "data",
              "description": "<p>The resulting torrent uploaded (see view for the properties)</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Torrents",
      "error": {
        "fields": {
          "Error 4xx": [{
              "group": "Error 4xx",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request couldn't be done due to some errors.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "String[]",
              "optional": false,
              "field": "errors",
              "description": "<p>List of errors messages.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "Object[]",
              "optional": false,
              "field": "all_errors",
              "description": "<p>List of errors object messages for each wrong field</p>"
            }
          ]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n      \"ok\": false,\n      \"errors\": [ ... ]\n      \"all_errors\": {\n\t\t \t\"username\": [ ... ],\n       }\n    }",
          "type": "json"
        }]
      }
    },
    {
      "type": "post",
      "url": "/login/",
      "title": "Login a user",
      "version": "1.1.1",
      "name": "Login",
      "group": "Users",
      "parameter": {
        "fields": {
          "Parameter": [{
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "username",
              "description": "<p>Username or Email.</p>"
            },
            {
              "group": "Parameter",
              "type": "String",
              "optional": false,
              "field": "password",
              "description": "<p>Password.</p>"
            }
          ]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request is done without failing</p>"
            },
            {
              "group": "Success 200",
              "type": "String[]",
              "optional": false,
              "field": "infos",
              "description": "<p>Messages information relative to the request</p>"
            },
            {
              "group": "Success 200",
              "type": "Object",
              "optional": false,
              "field": "data",
              "description": "<p>The connected user object</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n\t\t{\n\t\t\tdata:\n      \t\t[{\n\t\t\t\t\tuser_id:1,\n\t\t\t\t\tusername:\"username\",\n\t\t\t\t\tstatus:1,\n\t\t\t\t\ttoken:\"token\",\n\t\t\t\t\tmd5:\"\",\n\t\t\t\t\tcreated_at:\"date\",\n\t\t\t\t\tliking_count:0,\n\t\t\t\t\tliked_count:0\n\t\t\t\t}],\n\t\t\tinfos: [\"Logged\", ... ],\n\t\t\tok:true\n\t\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Users",
      "error": {
        "fields": {
          "Error 4xx": [{
              "group": "Error 4xx",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request couldn't be done due to some errors.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "String[]",
              "optional": false,
              "field": "errors",
              "description": "<p>List of errors messages.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "Object[]",
              "optional": false,
              "field": "all_errors",
              "description": "<p>List of errors object messages for each wrong field</p>"
            }
          ]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n      \"ok\": false,\n      \"errors\": [ ... ]\n      \"all_errors\": {\n\t\t \t\"username\": [ ... ],\n       }\n    }",
          "type": "json"
        }]
      }
    },
    {
      "type": "get",
      "url": "/profile/",
      "title": "Get a user profile",
      "version": "1.1.1",
      "name": "Profile",
      "group": "Users",
      "parameter": {
        "fields": {
          "Parameter": [{
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>User ID.</p>"
          }]
        }
      },
      "success": {
        "fields": {
          "Success 200": [{
              "group": "Success 200",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request is done without failing</p>"
            },
            {
              "group": "Success 200",
              "type": "String[]",
              "optional": false,
              "field": "infos",
              "description": "<p>Messages information relative to the request</p>"
            },
            {
              "group": "Success 200",
              "type": "Object",
              "optional": false,
              "field": "data",
              "description": "<p>The user object</p>"
            }
          ]
        },
        "examples": [{
          "title": "Success-Response:",
          "content": "    HTTP/1.1 200 OK\n\t\t{\n\t\t\tdata:\n      \t\t[{\n\t\t\t\t\tuser_id:1,\n\t\t\t\t\tusername:\"username\",\n\t\t\t\t\tstatus:1,\n\t\t\t\t\tmd5:\"\",\n\t\t\t\t\tcreated_at:\"date\",\n\t\t\t\t\tliking_count:0,\n\t\t\t\t\tliked_count:0\n\t\t\t\t}],\n\t\t\tinfos: [\"Logged\", ... ],\n\t\t\tok:true\n\t\t}",
          "type": "json"
        }]
      },
      "filename": "controllers/api/api.go",
      "groupTitle": "Users",
      "error": {
        "fields": {
          "Error 4xx": [{
              "group": "Error 4xx",
              "type": "Boolean",
              "optional": false,
              "field": "ok",
              "description": "<p>The request couldn't be done due to some errors.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "String[]",
              "optional": false,
              "field": "errors",
              "description": "<p>List of errors messages.</p>"
            },
            {
              "group": "Error 4xx",
              "type": "Object[]",
              "optional": false,
              "field": "all_errors",
              "description": "<p>List of errors object messages for each wrong field</p>"
            }
          ]
        },
        "examples": [{
          "title": "Error-Response:",
          "content": "    HTTP/1.1 200 OK\n    {\n      \"ok\": false,\n      \"errors\": [ ... ]\n      \"all_errors\": {\n\t\t \t\"username\": [ ... ],\n       }\n    }",
          "type": "json"
        }]
      }
    }
  ]
});
