var Sukebei = document.getElementById("IsUploadingToSukebei").value == "yes" ? 1 : 0;


document.getElementsByClassName("torrent-preview-table")[0].style.display = "block";
document.getElementsByClassName("table-torrent-date")[0].innerText = new Date(Date.now()).toLocaleString(document.getElementsByTagName("html")[0].getAttribute("lang"), { year: "numeric", month: "short", day: "numeric" });

document.getElementsByClassName("form-torrent-category")[0].addEventListener("change", UpdatePreviewCategory);
document.getElementsByClassName("form-torrent-name")[0].addEventListener("keyup", UpdatePreviewTorrentName);
document.getElementsByClassName("form-torrent-remake")[0].onchange = function(){
	document.getElementsByName("torrent-info tr")[0].className = "torrent-info" + (UserTrusted ? " trusted" : "") + (document.getElementsByClassName("form-torrent-remake")[0].checked ? " remake" : "");
};

function UpdatePreviewTorrentName(){
    document.getElementsByClassName("table-torrent-name")[0].innerText = document.getElementsByClassName("form-torrent-name")[0].value;
}


function UpdatePreviewCategory(){
    document.getElementsByClassName("table-torrent-category")[0].className = "nyaa-cat table-torrent-category "+ (Sukebei ? "sukebei" : "nyaa") + "-cat-" + CategoryList[Sukebei][document.getElementsByClassName("form-torrent-category")[0].selectedIndex];
}



function UpdateTorrentLang() {
	var lang_count,
		lang_value = "other";
			
		lang_count = 0;
		
	for(var lang_index = 0; lang_index < document.getElementsByName("languages").length; lang_index++) {
		if(document.getElementsByName("languages")[lang_index].checked) {
			if(++lang_count > 1){
				lang_value = "multiple";
				break;
			}
			lang_value = document.getElementsByName("languages")[lang_index].value;
		}
	}
		var lang_cat = lang_value != "other" ? (lang_value > 1 ? "multiple" : lang_value) : "other";
	document.getElementsByClassName("table-torrent-flag")[0].className = "table-torrent-flag flag flag-" + lang_cat;
}
                                                                                       
var CategoryList = [
    [5,
    12,
    5,
    13,
    6,
    3,
    4,
    7,
    14,
    8,
    9,
    10,
    18,
    11,
    15,
    16,
    1,
    2],
    [11,
    12,
    13,
    14,
    15,
    21,
    22]
];
