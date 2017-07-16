// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
var TorrentsMod = {
  // Variables that can be modified to change the dom interactions
  show_hide_button: "show_actions",
  btn_class_action: "cb_action",
  btn_class_submit: "cb_submit",
  progress_bar_id: "progress_modtool",
  status_input_name: "status_id",
  owner_input_name: "owner_id",
  category_input_name: "category_id",
  delete_btn: "delete",
  lock_delete_btn: "lock_delete",
  edit_btn: "edit",
  refreshTimeout: 3000,
  // Internal variables used for processing the request
  selected: [],
  queued: [],
  unique_id:1,
  error_count:0,
  progress_count: 0,
  progress_max: 0,
  pause: false,
  enabled: false,

  // Init method
  Create: function() {
    var sh_btn = document.getElementById(TorrentsMod.show_hide_button);
    var btn_actions = document.getElementsByClassName(this.btn_class_action)
    var btn_submit = document.getElementsByClassName(this.btn_class_submit)
    btn_submit[0].disabled = true;
    for (var i=0; i < btn_actions.length; i++) {
      btn_actions[i].disabled = true;
      switch (btn_actions[i].id) {
        case this.delete_btn:
        btn_actions[i].addEventListener("click", this.Delete)
        break;
        case this.lock_delete_btn:
        btn_actions[i].addEventListener("click", this.LockDelete)
        break;
        case this.edit_btn:
        btn_actions[i].addEventListener("click", this.Edit)
        break;
        default:
        break;
      }
    }
    for (var i=0; i < this.checkboxes.length; i++) {
      checkbox = this.checkboxes[i];
      checkbox.addEventListener("change", this.checkboxEventHandler)
    }
    sh_btn.addEventListener("click", function(e) {
      var divActions = this.nextElementSibling;
      if (divActions.style.display == "inline") {
        TorrentsMod.enabled = false;
      } else {
        TorrentsMod.enabled = true;
      }
      divActions.style.display = (TorrentsMod.enabled) ? "inline" : "none";
      var td_cbs = document.getElementsByClassName("tr-cb")
      for (var i=0; i < td_cbs.length; i++) {
        td_cb = td_cbs[i];
        td_cb.style.display = (TorrentsMod.enabled) ? "table-cell" : "none";
      }
      var toggleText = this.dataset.toggleText;
      this.dataset.toggleText = this.innerText;
      this.innerText = toggleText;
    });
  },
  // generate a unique id for a query
  getId: function(){
    return this.unique_id++;
  },

  // UI Methods
  selectAll: function(bool) {
    var l = TorrentsMod.checkboxes.length;
    for (var i = 0; i < l; i++) {
      TorrentsMod.checkboxes[i].checked = bool;
      TorrentsMod.checkboxEventHandlerFunc(TorrentsMod.checkboxes[i]);
    }
  },
  disableBtnActions: function() {
    var btn_actions = document.getElementsByClassName(this.btn_class_action)
    for (var i=0; i < btn_actions.length; i++) {
      btn_actions[i].disabled = true;
    }
  },
  enableBtnActions: function() {
    var btn_actions = document.getElementsByClassName(this.btn_class_action)
    for (var i=0; i < btn_actions.length; i++) {
      btn_actions[i].disabled = false;
    }
  },
  enableBtnSubmit: function() {
    var btn_submit = document.getElementsByClassName(this.btn_class_submit)
    btn_submit[0].disabled = false;
  },
  enableApplyChangesBtn: function() {
    var btn_apply_changes = document.getElementById("confirm_changes");
    btn_apply_changes.disabled=false;
  },
  disableBtnSubmit: function() {
    var btn_submit = document.getElementsByClassName(this.btn_class_submit)
    btn_submit[0].disabled = true;
  },
  disableApplyChangesBtn: function() {
    var btn_apply_changes = document.getElementById("confirm_changes");
    btn_apply_changes.disabled=true;
  },
  removeDivFromList: function(i) {
    var queueAction = this.queued[i];
    var parentDiv = document.getElementById(queueAction.unique_id).parentNode;
    parentDiv.removeChild(document.getElementById(queueAction.unique_id));
  },
  removeFromParent: function(el) {
    var parentDiv = el.parentNode;
    parentDiv.removeChild(el);
  },
  generatingModal: function() {
    listLength = this.queued.length;
    var div = {"edit": "", "delete": ""};
    for (var i=0; i < listLength; i++) {
      var listHTML = "";
      for(key in this.queued[i].selection) {
        var selection = this.queued[i].selection[key];
        selection.key = i;
        listHTML += Templates.Render("torrents."+this.queued[i].action+".item", selection);
      }
      this.queued[i].list = listHTML;
      this.queued[i].key = i;
      div[this.queued[i].action] += Templates.Render("torrents."+this.queued[i].action+".block", this.queued[i]);
    }
    this.progress_count = 0;
    this.progress_max = listLength;
    document.querySelector(".modal .edit_changes").innerHTML = div["edit"];
    document.querySelector(".modal .delete_changes").innerHTML = div["delete"];
  },
  toggleList: function(el) {
    el.parentNode.nextSibling.style.display = (el.parentNode.nextSibling.style.display != "block") ? "block" : "none"
  },
  addToLog: function(type, msg) {
    var logDiv = document.querySelector(".modal .logs_mess");
    if (logDiv.style.display == "none") logDiv.style.display = "block"
    logDiv.innerHTML += Templates.Render("torrents.logs."+type, msg);
  },
  updateProgressBar: function() {
    document.querySelector("#"+this.progress_bar_id).style.display = "block";
    var progress_green = document.querySelector("#"+this.progress_bar_id+" .progress-green");
    var perc = this.progress_count/this.progress_max*100;
    progress_green.style.width=perc+"%";
    progress_green.innerText = this.progress_count+"/"+this.progress_max;
  },
  resetModal: function() {
    var logDiv = document.querySelector(".modal .logs_mess");
    logDiv.style.display = "none";
    document.querySelector("#"+this.progress_bar_id).style.display = "none";
    logDiv.innerHTML = "";
    document.querySelector(".modal .edit_changes").innerHTML = "";
    document.querySelector(".modal .delete_changes").innerHTML = "";
    this.enableApplyChangesBtn();
  },
  statusToClassName: function(status) {
    var className = ["", "normal", "remake", "trusted", "aplus", ""]
    return className[status];
  },

  // Selection Management Methods
  addToSelection: function (torrent) {
    this.selected[torrent.id] = torrent;
  },
  removeFromSelection: function(torrent) {
    delete this.selected[torrent.id];
    for (t in this.selected) {
      return
    }
    this.selected = [];
  },

  // Query Queue Management Methods
  AddToQueue: function(QueueAction) {
    QueueAction.unique_id = this.getId(); // used for DOM interaction
    this.queued.push(QueueAction);
    this.enableBtnSubmit()
  },
  RemoveFromQueueAfterItems: function(i) {
    for (t in this.queued[i].selection) {
      this.RemoveItemFromQueue(i, t);
    }
  },
  RemoveFromQueue: function(i) {
    this.progress_max = (this.progress_max - 1 >= 0) ? this.progress_max-1 : 0;
    this.RemoveFromQueueAction(i)
    return false;
  },
  RemoveFromQueueAction: function(i) {
    this.removeFromParent(document.getElementById("list_"+this.queued[i].unique_id));
    this.queued.splice(i, 1);
    if (this.queued.length == 0) {
      this.disableBtnSubmit();
      if (this.progress_max>0) {
        this.disableApplyChangesBtn();
      } else {
        Modal.CloseActive();
      }
    }
  },
  formatSelectionToQuery: function(selection) {
    var format = "";
    for (s in selection) {
      format += "&torrent_id="+selection[s].id
    }
    return (format != "") ? format.substr(1) : ""
  },
  RemoveItemFromQueue: function(i, id) {
    this.removeFromParent(document.getElementById("list_item_"+id));
    delete this.queued[i].selection[id];
    document.getElementById("torrent_cb_"+id).checked=false;
    document.getElementById("torrent_"+id).style.display="";
    var test = 0;
    for (t in this.queued[i].selection) {
      test++
      break;
    }
    if (test == 0) this.RemoveFromQueue(i);
    return false;
  },
  newQueryAttempt: function(queryUrl, queryPost, callback) {
    Query.Post(queryUrl, queryPost, function (response) {
      if ((response.length == 0)||(!response.ok)) { // Query has failed
        var errorMsg = response.errors.join("<br>");
        TorrentsMod.addToLog("error", errorMsg);
        TorrentsMod.error_count++;
        if (TorrentsMod.error_count < 2) {
          TorrentsMod.addToLog("success", T.r("try_new_attempt"));
          TorrentsMod.newQueryAttempt(queryUrl, queryPost, callback)
        } else {
          TorrentsMod.addToLog("error", T.r("query_is_broken", queryUrl, queryPost));
          if (callback != undefined) {
            TorrentsMod.addToLog("error", "Passing to the next query...");
            callback(response) // So we can query only one item
          }
        }
      } else {
        var succesMsg = (response.infos != null) ? response.infos.join("<br>") : T.r("query_executed_success");
        TorrentsMod.addToLog("success", succesMsg)
        if (callback != undefined) callback(response) // So we can query only one item
      }
    });
  },
  QueryQueue: function(i, callback) {
    if (this.queued.length > 0) {
      var QueueAction = this.queued[i]; // we clone it so we can delete it safely
      this.RemoveFromQueueAction(i);
      var queryPost = "";
      var queryUrl = "/mod/api/torrents";
      QueueAction.queryPost = TorrentsMod.formatSelectionToQuery(QueueAction.selection)
      if (QueueAction.action == "delete") {
        queryPost="action="+QueueAction.action;
        queryPost+="&withreport="+QueueAction.withReport;
        queryPost += "&status="+((QueueAction.status != undefined) ? QueueAction.status : "");
        queryPost += "&"+QueueAction.queryPost; // we add torrent id
      } else if (QueueAction.action == "edit") {
        queryPost="action=multiple";
        queryPost += "&status="+QueueAction.status;
        queryPost += "&owner="+QueueAction.owner;
        queryPost += "&category="+QueueAction.category;
        queryPost += "&"+QueueAction.queryPost; // we add torrent id
      }
      TorrentsMod.newQueryAttempt(queryUrl, queryPost, callback)
    } else {
      TorrentsMod.addToLog("success", T.r("all_operations_done"))
      if (TorrentsMod.refreshTimeout > 0) {
        TorrentsMod.addToLog("success", T.r("refreshing_in", TorrentsMod.refreshTimeout/1000))
        setTimeout(function(){
          window.location.reload()
        }, TorrentsMod.refreshTimeout);
      }
    }
  },
  QueryLoop: function() {
    if (TorrentsMod.progress_count <= TorrentsMod.progress_max) {
      TorrentsMod.updateProgressBar()
      TorrentsMod.progress_count++;
      if (TorrentsMod.progress_count > TorrentsMod.progress_max) TorrentsMod.progress_count = TorrentsMod.progress_max;
      TorrentsMod.QueryQueue(0, TorrentsMod.QueryLoop);
    }
  },

  // Event Handlers
  checkboxEventHandler: function(e) {
    var el = e.target;
    TorrentsMod.checkboxEventHandlerFunc(el);
  },
  checkboxEventHandlerFunc: function(el) {
    var name = el.dataset.name;
    var id = el.value;
    if (el.checked) TorrentsMod.addToSelection({name:name, id:id});
    else TorrentsMod.removeFromSelection({name:name, id:id});
    if (TorrentsMod.selected.length > 0) TorrentsMod.enableBtnActions();
    else TorrentsMod.disableBtnActions();
  },

  // Action Methods
  DeleteHandler: function(locked) {
    var withReport = confirm(T.r("delete_reports_with_torrents"))
    var selection = TorrentsMod.selected;
    if (locked)
    TorrentsMod.AddToQueue({ action: "delete",
    withReport: withReport,
    selection: selection,
    queryPost: "",
    infos: T.r("with_lock")+ ((withReport) ? T.r("and_reports") : ""),
    status: "5" });
    else TorrentsMod.AddToQueue({
      action: "delete",
      withReport: withReport,
      selection: selection,
      infos: (withReport) ? T.r("with_reports") : "",
      queryPost: ""});
      for (i in selection) document.getElementById("torrent_"+i).style.display="none";
      TorrentsMod.selected = []
      TorrentsMod.disableBtnActions();
    },
    Delete: function(e) {
      TorrentsMod.DeleteHandler(false);
      e.preventDefault();
    },
    LockDelete: function(e) {
      TorrentsMod.DeleteHandler(true);
      e.preventDefault();
    },
    Edit: function(e) {
      var selection = TorrentsMod.selected;
      var status = document.querySelector(".modtools *[name='"+TorrentsMod.status_input_name+"']").value;
      var owner_id = document.querySelector(".modtools *[name='"+TorrentsMod.owner_input_name+"']").value;
      var category = document.querySelector(".modtools *[name='"+TorrentsMod.category_input_name+"']").value;
      var infos = "";
      infos += (status != "") ? T.r("status_js", status) : "";
      infos += (owner_id != "") ? T.r("owner_id_js", owner_id) : "";
      infos += (category != "") ? T.r("category_js", category) : "";
      TorrentsMod.AddToQueue({
        action: "edit",
        selection: selection,
        queryPost: "", // We don't format now, we wait until the query is sent
        infos: (infos != "" ) ? T.r("with_st", infos) : T.r("no_changes"),
        status: status,
        category: category,
        owner: owner_id });
        if (status != "") {
          for (i in selection) document.getElementById("torrent_"+i).className="torrent-info "+TorrentsMod.statusToClassName(status);
        }
        TorrentsMod.selected = []
        TorrentsMod.disableBtnActions();
        e.preventDefault();
      },
      ApplyChanges: function() {
        this.pause = false;
        this.QueryLoop();
      }
    };

    // Load torrentMods when DOM is ready
    document.addEventListener("DOMContentLoaded", function() {
      TorrentsMod.checkboxes =  document.querySelectorAll("input[name='torrent_id']");
      TorrentsMod.Create();
    });
    // @license-end
