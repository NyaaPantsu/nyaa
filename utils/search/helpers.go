package search

import (
	"fmt"
	"strings"

	"github.com/NyaaPantsu/nyaa/models/torrents"
)

func tagsToID(tags string) []string {
	var ids []string
	where := WhereParams{}
	for i, tag := range strings.Split(tags, ",") {
		if tag != "" {
			ta := strings.Split(tag, ":")
			if len(ta) == 2 {
				if i > 0 {
					where.Conditions += " AND "
				}
				where.Conditions += "type = ? AND tag = ?"
				where.Params = append(where.Params, ta[0])
				where.Params = append(where.Params, ta[1])
			}
		}
	}
	if len(where.Params) > 0 {
		tID, err := torrents.GetIDs(&where)
		if err == nil && len(tID) > 0 {
			for _, id := range tID {
				ids = append(ids, fmt.Sprintf("%d", id))
			}
		}
	}
}
