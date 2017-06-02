document.querySelector("[data-selectall='checkbox']").addEventListener("change", function(e) {
	var cbs = document.querySelectorAll("input[type='checkbox'].selectable");
	var l = cbs.length;
	for (var i=0; i<l; i++) cbs[i].checked = e.target.checked;
});