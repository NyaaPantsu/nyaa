// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
// Get the modal
var Modal = {
  active: 0,
  Init: function (params) {
    var elements = params.elements
    var button = (params.button != undefined) ? params.button : false
    if (elements.innerHTML != undefined) {

    } else {
      var nbEl = elements.length
      for (var i=0; i < nbEl; i++) {
        var modal = elements[i];
        this.addModal(modal, button, i, params.before, params.after, params.close)
      }
    }
  },
  addModal: function(modal, btn, i, before_callback, after_callback, close_callback) {
    var isBtnArray = false;
    // Get the button that opens the modal
    if (!btn) {
      btn = document.getElementById("modal_btn_"+modal.id)
    } else if (btn.match(/^#/)) {
      btn = document.getElementById(btn.substr(1));
    } else if (btn.match(/^\./)) {
      btn = document.getElementsByClassName(btn.substr(1));
      isBtnArray = true;
    } else {
      console.error("Couldn't find the button")
      return
    }
    if ((isBtnArray) && (i > 0) && (btn.length > 0) && (btn.length > i)) {
      btn[i].addEventListener("click", function(e) {
        if (before_callback != undefined) before_callback()
        modal.style.display = "block";
        Modal.active = modal;
        if (after_callback != undefined) after_callback()
        e.preventDefault();
      });
    } else {
      btn = (isBtnArray) ? btn[0] : btn;
      // When the user clicks on the button, open the modal
      btn.addEventListener("click", function(e) {
        if (before_callback != undefined) before_callback()
        modal.style.display = "block";
        Modal.active = modal;
        if (after_callback != undefined) after_callback()
        e.preventDefault();
      });
    }
    // Get the <span> element that closes the modal
    var span = document.querySelectorAll("#"+modal.id+" .close")[0]
    // When the user clicks on <span> (x), close the modal
    span.addEventListener("click", function(e) {
      modal.style.display = "none";
      Modal.active = 0;
      if (close_callback != undefined) close_callback()
      e.preventDefault();
    });
    // When the user clicks anywhere outside of the modal, close it
    window.addEventListener("click", function(event) {
      if (event.target == modal) {
        modal.style.display = "none";
        Modal.active = 0;
        if (close_callback != undefined) close_callback()
      }
    });
  },
  CloseActive: function() {
    if (this.active != 0) {
      this.active.style.display= "none";
      this.active = 0;
    }
  },
  GetActive: function() {
    return this.active;
  },
  Open: function(q) {
    var modal = document.querySelector(q);
    if (modal != undefined) {
      modal.style.display= "none";
      this.active = modal;
    }
  }
};
// @license-end
