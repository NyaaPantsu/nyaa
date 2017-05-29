// Templates variable
var Templates = {
	tmpl: [],
	Add: function(templateName, template) {
		this.tmpl[templateName] = template
	},
	Render: function(templateName, model) {
		console.log(model)
		return this.tmpl[templateName](model)
	},
	ApplyItemListRenderer: function(params) {
		return function(models) {
			console.log("Parsing results...")
			for (var i=0; i < models.length; i++) {
				var object = Templates.Render(params.templateName, models[i]);
				if (params.method == "append") {
					params.element.innerHTML = params.element.innerHTML + object
				} else if (params.method == "prepend") {
					params.element.innerHTML = object + params.element.innerHTML
				}
			}
		};
	}
};
