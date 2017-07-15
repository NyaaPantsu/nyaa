document.getElementsByClassName("form-torrent-name")[0].onkeyup = function(){
    document.getElementsByClassName("table-torrent-name")[0].innerText = document.getElementsByClassName("form-torrent-name")[0].value;
};

function UpdatePreviewCategory(){
    document.getElementsByClassName("table-torrent-category")[0].className = "nyaa-cat table-torrent-category "+ Sukebei ? "sukebei" : "nyaa" + "-cat-" + Categorylist[Sukebei][document.getElementsByClassName("form-torrent-category")[0].selectedIndex];
}

document.getElementsByClassName("form-torrent-remake")[0].onchange = function(){
    document.getElementsByClassName("table-torrent-thead")[0].className = "torrent-info table-torrent-thead" + (UserTrusted ? " trusted" : "") + (document.getElementsByClassName("form-torrent-remake")[0].checked ? " remake" : "");
};
                                                                                                                
document.getElementsByClassName("form-torrent-hidden")[0].onchange = function(){
    document.getElementsByClassName("table-torrent-thead")[0].className = "torrent-info table-torrent-thead" + (document.getElementsByClassName("form-torrent-remake")[0].checked ? " remake" : "");
    if(UserTrusted && !document.getElementsByClassName("form-torrent-hidden")[0].checked)
        document.getElementsByClassName("table-torrent-thead")[0].className = document.getElementsByClassName("table-torrent-thead")[0].className + " trusted";
};                                                                                                            '
                                                                                       
var CategoryList = [
    [5
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
