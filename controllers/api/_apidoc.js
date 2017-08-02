// ------------------------------------------------------------------------------------------
// General apiDoc documentation blocks and old history blocks.
// ------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------
// Current Success.
// ------------------------------------------------------------------------------------------


// ------------------------------------------------------------------------------------------
// Current Errors.
// ------------------------------------------------------------------------------------------
/**
 * @apiDefine NotFoundError
 * @apiVersion 1.0.0
 * @apiError {String[]} errors List of errors messages with a 404 error message in it.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 404 Not Found
 *     {
 *       "errors": [ "404_not_found", ... ]
 *     }
 */
// ------------------------------------------------------------------------------------------
// History.
// ------------------------------------------------------------------------------------------

/**
 * @api {get} /search/ Search Torrents
 * @apiVersion 1.0.0
 * @apiName FindTorrents
 * @apiGroup Torrents
 *
 * @apiParam {String[]} c In which categories to search.
 * @apiParam {String} q Query to search (torrent name).
 * @apiParam {String} limit Number of results per page.
 * @apiParam {String} userID Uploader ID owning the torrents.
 * @apiParam {String} fromID Show results with torrents ID superior to this.
 * @apiParam {String} s Torrent status.
 * @apiParam {String} maxage Torrents which have been uploaded the last x days.
 * @apiParam {String} toDate Torrents which have been uploaded since x <code>dateType</code>.
 * @apiParam {String} fromDate Torrents which have been uploaded the last x <code>dateType</code>.
 * @apiParam {String} dateType Which type of date (<code>d</code> for days, <code>m</code> for months, <code>y</code> for years).
 * @apiParam {String} minSize Filter by minimal size in <code>sizeType</code>.
 * @apiParam {String} maxSize Filter by maximal size in <code>sizeType</code>.
 * @apiParam {String} sizeType Which type of size (<code>b</code> for bytes, <code>k</code> for kilobytes, <code>m</code> for megabytes, <code>g</code> for gigabytes).
 * @apiParam {String} sort Torrent sorting type (0 = id, 1 = name, 2 = date, 3 = downloads, 4 = size, 5 = seeders, 6 = leechers, 7 = completed).
 * @apiParam {Boolean} order Order ascending or descending (true = ascending).
 * @apiParam {String[]} lang Filter the languages.
 * @apiParam {Number} page Search page.
 *
 * @apiSuccess {Object[]} torrents List of torrent object (see view for the properties).
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *			[...]
 *		}
 *
 * @apiUse NotFoundError
 */
