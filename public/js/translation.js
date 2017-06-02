function Translations() {
    var translations = {};
    this.Add =  function(tr, val) {
        if (val != undefined) {
            tr[tr] = val;
        }
        Object.assign(translations, tr);
    };
    this.r = function(string, ...args) {
        if ((string != undefined) && (translations[string] != undefined)) {
            if (args != undefined) {
                return this.format(translations[string], ...args);
            }
            return translations[string];
        }
        console.error("No translation string for %s! Please check!", string);
        return "";
    };
    this.format = function(format, ...args) {
        return format.replace(/{(\d+)}/g, function(match, number) { 
        return typeof args[number] != 'undefined'
            ? args[number] 
            : match
        ;
        });
    };
}

var T = new Translations();