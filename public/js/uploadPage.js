(function() {
	var torrent = $("input[name=torrent]"),
	magnet = $("input[name=magnet]"),
	name = $("input[name=name]");

	torrent.on("change", function() {
		if (torrent.val() == "") {
			enableField(magnet);
			name.attr("required", "");
		} else {
			disableField(magnet);
			// .torrent file will allow autofilling name
			name.removeAttr("required", "");
		}
	});
	magnet.on("change", function() {
		if (magnet.val() == "")
			enableField(torrent);
		else
			disableField(torrent);
	});

	function enableField(e) {
		e.attr("required", "")
			.removeAttr("disabled");
	}
	function disableField(e) {
		e.attr("disabled", "")
			.removeAttr("required");
	}
})();
