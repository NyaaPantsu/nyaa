var TorrentsMod = {
    show_hide_button: "show_actions",
    btn_class_action: "cb_action",
    btn_class_submit: "cb_submit",
    selected: [],
    queued: [],
    unique_id:1,
    Create: function() {
        var sh_btn = document.getElementById(TorrentsMod.show_hide_button);
        var btn_actions = document.getElementsByClassName(this.btn_class_action)
        var btn_submit = document.getElementsByClassName(this.btn_class_submit)
        btn_submit[0].disabled = true;
        for (var i=0; i < btn_actions.length; i++) {
            btn_actions[i].disabled = true;
            switch (btn_actions[i].id) {
                case "delete":
                    btn_actions[i].addEventListener("click", this.Delete)
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
            var display = "inline"
            var divActions = this.nextElementSibling;
            console.log(divActions)
            if (divActions.style.display == "inline") {
            display = "none";
            }
            divActions.style.display = display;
            var td_cbs = document.getElementsByClassName("tr-cb")
            for (var i=0; i < td_cbs.length; i++) {
                td_cb = td_cbs[i];
                td_cb.style.display = (display == "inline") ? "table-cell" : "none";
            }
            var toggleText = this.dataset.toggleText;
            this.dataset.toggleText = this.innerText;
            this.innerText = toggleText;
        });
    },
    getId: function(){
        return this.unique_id++;
    },

    // UI Methods
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
    disableBtnSubmit: function() {
        var btn_submit = document.getElementsByClassName(this.btn_class_submit)
        btn_submit[0].disabled = true;
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
                selection.key = i
                listHTML += Templates.Render("torrents."+this.queued[i].action+".item", selection);
            }
            this.queued[i].list = listHTML;
            this.queued[i].key = i;
            div[this.queued[i].action] += Templates.Render("torrents."+this.queued[i].action+".block", this.queued[i])
        }
        document.querySelector(".modal .edit_changes").innerHTML = div["edit"]
        document.querySelector(".modal .delete_changes").innerHTML = div["delete"]
    },
    toggleList: function(el) {
        el.parentNode.nextSibling.style.display = (el.parentNode.nextSibling.style.display != "block") ? "block" : "none"
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
    RemoveFromQueue: function(i) {
        for (t in this.queued[i].selection) {
            this.RemoveItemFromQueue(i, t);
        }
        return false;
    },
    RemoveFromQueueAction: function(i) {
        this.removeFromParent(document.getElementById("list_"+this.queued[i].unique_id));
        this.queued.splice(i, 1);
        if (this.queued.length == 0) {
             this.disableBtnSubmit();
             Modal.CloseActive();
        }
    },
    formatSelectionToQuery: function() {
        var format = "";
        for (s in this.selected) {
            format += "&torrent_id="+this.selected[s].id
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
        if (test == 0) this.RemoveFromQueueAction(i);
        return false;
    },

    // Event Handlers
    checkboxEventHandler: function(e) {
        var el = e.target;
        var name = el.dataset.name;
        var id = el.value;
        if (el.checked) TorrentsMod.addToSelection({name:name, id:id});
        else TorrentsMod.removeFromSelection({name:name, id:id});
        if (TorrentsMod.selected.length > 0) TorrentsMod.enableBtnActions();
        else TorrentsMod.disableBtnActions();
    },

    // Action Methods
    Delete: function(e) {
        var withReport = confirm("Do you want to delete the reports along the selected torrents?")
        var selection = TorrentsMod.selected;
        TorrentsMod.AddToQueue({ action: "delete", withReport: withReport, selection: selection, queryPost: TorrentsMod.formatSelectionToQuery() });
        for (i in selection) document.getElementById("torrent_"+i).style.display="none";
        TorrentsMod.selected = []
        TorrentsMod.disableBtnActions();
        e.preventDefault();
    },
};

// Load torrentMods when DOM is ready
document.addEventListener("DOMContentLoaded", function() {
    TorrentsMod.checkboxes =  document.querySelectorAll("input[name='torrent_id']");
    TorrentsMod.Create();
});