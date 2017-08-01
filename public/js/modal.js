// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat
// Get the modal
var Modal = {
  active: 0,
  // Initialise a modal or multiple ones
  // takes as parameter an object params:
  // @param params Object{ elements: NodeList|Node, button: ID(eg. #something)|Class(eg. .something), before: callback, after: callback, close: callback }
  Init: function (params) {
    var elements = params.elements
    var button = (params.button != undefined) ? params.button : false
    if (elements.innerHTML == undefined) {
      var nbEl = elements.length
      for (var i = 0; i < nbEl; i++) {
        var modal = elements[i];
        this.addModal(modal, button, i, params.before, params.after, params.close)
      }
    } else {
      this.addModal(modal, button, i, params.before, params.after, params.close)
    }
  },
  // addModal prepare a modal (called by Init so you don't have use it)
  // @param modal Node
  // @param btn  ID(eg. #something)|Class(eg. .something)
  // @param i If multiple btn, points out to which btn in the array apply the event
  // @param before_callback callback called before opening a modal
  // @param after_callback callback called after opening a modal
  // @param close_callback callback called after closing a modal
  addModal: function (modal, btn, i, before_callback, after_callback, close_callback) {
    var isBtnArray = false;
    // Get the button that opens the modal
    if (!btn) {
      btn = document.getElementById("modal_btn_" + modal.id)
    } else if (typeof(btn) == "string" && btn.match(/^#/)) {
      btn = document.getElementById(btn.substr(1));
    } else if (typeof(btn) == "string" && btn.match(/^\./)) {
      btn = document.getElementsByClassName(btn.substr(1));
      isBtnArray = true;
    } else if (btn instanceof Array) {
      btn = btn.map(function(val, index) {
        if (val.match(/^#/)) {
          return document.getElementById(val.substr(1));
        } else if (val.match(/^\./)) {
          return document.getElementsByClassName(val.substr(1))[index];
        }
        return document.querySelector(val)
      })
      isBtnArray = true;
    } else {
      console.error("Couldn't find the button, please provide either a #id, a .classname or an array of #id")
      return
    }
    if ((isBtnArray) && (i > 0) && (btn.length > 0) && (btn.length > i)) {
      btn[i].addEventListener("click", function (e) {
        if (before_callback != undefined) before_callback()
        modal.style.display = "block";
        Modal.active = modal;
        if (after_callback != undefined) after_callback()
        e.preventDefault();
      });
    } else {
      btn = (isBtnArray) ? btn[0] : btn;
      // When the user clicks on the button, open the modal
      btn.addEventListener("click", function (e) {
        if (before_callback != undefined) before_callback()
        modal.style.display = "block";
        Modal.active = modal;
        if (after_callback != undefined) after_callback()
        e.preventDefault();
      });
    }
    // Get the <span> element that closes the modal
    var span = document.querySelectorAll("#" + modal.id + " .close")[0]
    // When the user clicks on <span> (x), close the modal
    span.addEventListener("click", function (e) {
      modal.style.display = "none";
      Modal.active = 0;
      if (close_callback != undefined) close_callback()
      e.preventDefault();
    });
    // When the user clicks anywhere outside of the modal, close it
    window.addEventListener("click", function (event) {
      if (event.target == modal) {
        modal.style.display = "none";
        Modal.active = 0;
        if (close_callback != undefined) close_callback()
      }
    });
  },
  // CloseActive closes the opened modal, if any
  CloseActive: function () {
    if (this.active != 0) {
      this.active.style.display = "none";
      this.active = 0;
    }
  },
  // GetActive return the opened modal div
  GetActive: function () {
    return this.active;
  },
  // Open opens a modal and closes the active one if any
  Open: function (q) {
    var activeModal = this.GetActive()
    if (activeModal != 0) {
      this.CloseActive()
    }
    var modal = document.querySelector(q);
    if (modal != undefined) {
      modal.style.display = "none";
      this.active = modal;
    }
  }
};
// @license-end
