package search

import (
	"fmt"
	"strings"

	"github.com/NyaaPantsu/nyaa/models/torrents"
)

func tagsToID(tags string) []string {
	var ids []string
	query := &Query{}
	for _, tag := range strings.Split(tags, ",") {
		if tag != "" {
			ta := strings.Split(tag, ":")
			if len(ta) == 2 {
				query.Append("type = ? AND tag = ?", ta[0], ta[1])
			}
		}
	}
	_, params := query.ToDBQuery()
	if len(params) > 0 {
		tID, err := torrents.GetIDs(query)
		if err == nil && len(tID) > 0 {
			for _, id := range tID {
				ids = append(ids, fmt.Sprintf("%d", id))
			}
		}
	}
	return ids
}
