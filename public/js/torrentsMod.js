{{ range _, cat := GetCategories(false, true) }}
T.Add("{{ cat.ID }}", "{{ T(cat.Name) }}")
{{end}}
{{ range _, language := GetTorrentLanguages() }}
T.Add("{{ language.Code }}", "{{ LanguageName(language, T) }}")
{{ if language.Tag != language.Code }}
T.Add("{{ language.Tag }}", "{{ LanguageName(language, T) }}")
{{end}}
{{end}}
Templates.Add("torrents.item", function(torrent) {
    return "<tr id=\"torrent_" + torrent.id + "\" class=\"torrent-info"+ ((torrent.status == 2) ? " remake" : ((torrent.status == 3) ? " trusted" : ((torrent.status == 3) ? " aplus" : "" )))+"\">"+
    {{ if User.HasAdmin() }}
		(( TorrentsMod.enabled ) ? "<td class=\"tr-cb\">" :  "<td class=\"tr-cb hide\"" + ((TorrentsMod.enabled) ? "style=\"display:table-cell;\"" : "") +">")+
        "<input data-name=\""+Templates.EncodeEntities(torrent.name)+"\" type=\"checkbox\" id=\"torrent_cb_"+torrent.id+"\" name=\"torrent_id\" value=\""+torrent.id+"\">"+
        "</td>"+
    {{ end }}
    "<td class=\"tr-cat home-td\">"+
    {{ if Sukebei() }}
        "<div class=\"nyaa-cat sukebei-cat-"+ torrent.category + torrent.sub_category +"\">"+
    {{else}}
        "<div class=\"nyaa-cat nyaa-cat-"+ torrent.sub_category +"\">"+
    {{end}}
            "<a href=\"{{URL.Parse("/search?c=") }}"+ torrent.category + "_" + torrent.sub_category +"\" title=\""+ T.r(torrent.category+"_"+torrent.sub_category)+"\" class=\"category\">"+
                ((torrent.languages[0] != "") ? "<a href=\"{{URL.Parse("/search?c=") }}"+ torrent.category + "_" + torrent.sub_category +"&lang=" + torrent.languages.join(",") +"\"><img src=\"img/blank.gif\" class=\"flag flag-"+ ((torrent.languages.length == 1) ? flagCode(torrent.languages[0]) : "multiple") +"\" title=\""+torrent.languages.map(function (el, i) { return T.r(el)}).join(",")+"\"></a>" : "") +
            "</a>"+
        "</div></td>"+
        "<td class=\"tr-name home-td\""+  (torrent.comments.length  == 0 ? "colspan=\"2\"" : "" ) +"><a href=\"/view/"+torrent.id+"\">"+Templates.EncodeEntities(torrent.name) +"</a></td>"+
        ((torrent.comments.length > 0) ? "<td class=\"tr-cs home-td\"><span>"+ torrent.comments.length + "</span></td>" : "")+
        "<td class=\"tr-links home-td\">"+
            "<a href=\""+torrent.magnet +"\" title=\"{{ T("magnet_link") }}\">"+
                "<div class=\"icon-magnet\"></div>"+
            "</a>"+(torrent.torrent != "" ? " <a href=\""+torrent.torrent+"\" title=\"{{ T("torrent_file") }}\"><div class=\"icon-floppy\"></div></a>" : "") +
        "</td>"+
        "<td class=\"tr-size home-td hide-xs\">"+humanFileSize(torrent.filesize)+"</td>"+
        "<td class=\"tr-se home-td hide-xs\">"+torrent.seeders+"</td>"+
        "<td class=\"tr-le home-td hide-xs\">"+torrent.leechers+"</td>"+
        "<td class=\"tr-dl home-td hide-xs\">"+torrent.completed+"</td>"+
        "<td class=\"tr-date home-td date-short hide-xs\">"+torrent.date+"</td>"+
    "</tr>";
});
