// Get the modal
var Modal = {
    Init: function (params) {
        var elements = params.elements
        var button = (params.button != undefined) ? params.button : false
        if (elements.innerHTML != undefined) {

        } else {
            var nbEl = elements.length
            for (var i=0; i < nbEl; i++) {
                var modal = elements[i];
                this.addModal(modal, button, i)
            }
        }
    },
    addModal: function(modal, btn, i) {
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
                modal.style.display = "block";
                e.preventDefault();
            });
        } else {
            btn = (isBtnArray) ? btn[0] : btn;
            // When the user clicks on the button, open the modal 
            btn.addEventListener("click", function(e) {
                modal.style.display = "block";
                e.preventDefault();
            });
        }
        // Get the <span> element that closes the modal
        var span = document.querySelectorAll("#"+modal.id+" .close")[0]
        // When the user clicks on <span> (x), close the modal
        span.addEventListener("click", function(e) {
            modal.style.display = "none";
            e.preventDefault();
        });
        // When the user clicks anywhere outside of the modal, close it
        window.addEventListener("click", function(event) {
            if (event.target == modal) {
                modal.style.display = "none";
            }
        });
    }
};