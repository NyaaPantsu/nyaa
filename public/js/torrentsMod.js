var TorrentsMod = {
    show_hide_button: "show_actions",
    btn_class_action: "cb_action",
    btn_class_submit: "cb_submit",
    selected: [],
    queued: [],
    Create: function() {
        var sh_btn = document.getElementById(TorrentsMod.show_hide_button);
        var btn_actions = document.getElementsByClassName(this.btn_class_action)
        var btn_submit = document.getElementsByClassName(this.btn_class_submit)
        btn_submit[0].disabled = true;
        for (var i=0; i < btn_actions.length; i++) {
            btn_actions[i].disabled = true;
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
        this.queued.push(QueueAction);
    },
    RemoveFromQueue: function(i) {
        this.queued.splice(i, 1);
    },
    formatSelectionToQuery: function() {
        return (this.selected.length > 0 ) ? "torrent_id="+this.selected.join("&torrent_id=") : ""
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
    Delete: function() {
        var withReport = prompt("Do you want delete the reports along the selected torrents?")
        this.AddToQueue({ action: "delete", withReport: withReport, selection: this.selected, queryPost: this.formatSelectionToQuery() });
    },
};

// Load torrentMods when DOM is ready
document.addEventListener("DOMContentLoaded", function() {
    TorrentsMod.checkboxes =  document.querySelectorAll("input[name='torrent_id']");
    TorrentsMod.Create();
});