$("[data-selectall='checkbox']").on("change", function(e) {
	var form = $(this).parents("form");
	$(form).find("input[type='checkbox']").prop('checked', $(this).prop('checked'));
});