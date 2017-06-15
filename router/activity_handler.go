package router

import (
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/service/activity"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"

	"github.com/gorilla/mux"
)

// ActivityListHandler : Show a list of activity
func ActivityListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	pagenum := 1
	offset := 100
	userid := r.URL.Query().Get("userid")
	filter := r.URL.Query().Get("filter")
	defer r.Body.Close()
	var err error
	messages := msg.GetMessages(r)
	currentUser := getUser(r)
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var conditions []string
	var values []interface{}
	if userid != "" && userPermission.HasAdmin(currentUser) {
		conditions = append(conditions, "user_id = ?")
		values = append(values, userid)
	}
	if filter != "" {
		conditions = append(conditions, "filter = ?")
		values = append(values, filter)
	}

	activities, nbActivities := activity.GetAllActivities(offset, (pagenum-1)*offset, strings.Join(conditions, " AND "), values...)
	common := newCommonVariables(r)
	common.Navigation = navigation{nbActivities, offset, pagenum, "activity_list"}
	htv := modelListVbs{common, activities, messages.GetAllErrors(), messages.GetAllInfos()}
	err = activityList.ExecuteTemplate(w, "index.html", htv)
	log.CheckError(err)
}
