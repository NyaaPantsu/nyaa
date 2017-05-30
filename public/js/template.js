// Templates variable
var Templates = {
	tmpl: [],
	Add: function(templateName, template) {
		this.tmpl[templateName] = template
	},
	Render: function(templateName, model) {
		return this.tmpl[templateName](model)
	},
	ApplyItemListRenderer: function(params) {
		return function(models) {
			for (var i=models.length-1; i >= 0; i--) {
				var object = Templates.Render(params.templateName, models[i]);
				if (params.method == "append") {
					params.element.innerHTML = params.element.innerHTML + object
				} else if (params.method == "prepend") {
					params.element.innerHTML = object + params.element.innerHTML
				}
			}
		};
	},
	EncodeEntities: function(value) {
		return value.
		replace(/&/g, '&amp;').
		replace(/[\uD800-\uDBFF][\uDC00-\uDFFF]/g, function(value) {
		var hi = value.charCodeAt(0);
		var low = value.charCodeAt(1);
		return '&#' + (((hi - 0xD800) * 0x400) + (low - 0xDC00) + 0x10000) + ';';
		}).
		replace(/([^\#-~| |!])/g, function(value) {
		return '&#' + value.charCodeAt(0) + ';';
		}).
		replace(/</g, '&lt;').
		replace(/>/g, '&gt;');
	}
};
