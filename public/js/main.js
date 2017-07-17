// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat

// Switches between themes when a new one is selected
function switchThemes(){
  themeName = document.getElementById("theme-selector").value
  var head = document.getElementsByTagName("head")[0];
  // Remove the theme in place, it fails if one isn't set
  try{
    head.removeChild(document.getElementById("theme"));
  } catch(err){}
  // Don't add a node if we don't want extra styling
  if(themeName === ""){
    return;
  }
  // Create the new one and put it back
  var newTheme = document.createElement("link");
  newTheme.setAttribute("rel", "stylesheet");
  newTheme.setAttribute("href", "/css/"+ themeName + ".css");
  newTheme.setAttribute("id", "theme");
  head.appendChild(newTheme);
}

// Used by spoiler tags
function toggleLayer(elem) {
  if (elem.classList.contains("hide")) {
    elem.classList.remove("hide");
  } else {
    elem.classList.add("hide");
  }
}
function parseAllDates() {
  // Date formatting
  var lang = document.getElementsByTagName("html")[0].getAttribute("lang");
  var ymdOpt = { year: "numeric", month: "short", day: "numeric" };
  var hmOpt  = { hour: "numeric", minute: "numeric" };

  var list = document.getElementsByClassName("date-short");
  for(var i in list) {
    var e = list[i];
    e.title = e.innerText;
    e.innerText = new Date(e.innerText).toLocaleString(lang, ymdOpt);
  }

  var list = document.getElementsByClassName("date-full");
  for(var i in list) {
    var e = list[i];
    e.title = e.innerText;
    e.innerText = new Date(e.innerText).toLocaleString(lang);
  }
}

parseAllDates();

/*Fixed-Navbar offset fix*/
document.addEventListener("DOMContentLoaded", function(event) {
  var shiftWindow = function() { scrollBy(0, -70) };
  if (location.hash) shiftWindow();
  window.addEventListener("hashchange", shiftWindow);

  document.getElementsByClassName("search-box")[0].addEventListener("focus", function (e) {
    var w = document.getElementsByClassName("h-user")[0].offsetWidth;
    document.getElementsByClassName("h-user")[0].style.display = "none";
    document.getElementsByClassName("search-box")[0].style.width = document.getElementsByClassName("search-box")[0].offsetWidth + w + "px";
  });
  document.getElementsByClassName("search-box")[0].addEventListener("blur", function (e) {
    document.getElementsByClassName("search-box")[0].style.width = "";
    document.getElementsByClassName("h-user")[0].style.display = "inline-block";
  });

});

function playVoice() {
  var mascotAudio = document.getElementById("explosion") || document.getElementById("nyanpassu")|| document.getElementById("nyanpassu2") || document.getElementById("kawaii");
  if (mascotAudio !== undefined) {
    mascotAudio.volume = 0.2;
    mascotAudio.play();
  } else {
    console.log("Your mascot doesn't support yet audio files!")
  }
}

document.getElementsByClassName("form-input refine")[0].addEventListener("click", function (e) {
  document.getElementsByClassName("box refine")[0].style.display = document.getElementsByClassName("box refine")[0].style.display == "none" ? "block" : "none";
  if(document.getElementsByClassName("form-input refine-searchbox")[0].value != document.getElementsByClassName("form-input search-box")[0].value)
  document.getElementsByClassName("form-input refine-searchbox")[0].value = document.getElementsByClassName("form-input search-box")[0].value;
  if(document.getElementsByClassName("form-input refine-category")[0].selectedIndex != document.getElementsByClassName("form-input form-category")[0].selectedIndex)
  document.getElementsByClassName("form-input refine-category")[0].selectedIndex = document.getElementsByClassName("form-input form-category")[0].selectedIndex;
  e.preventDefault();
  if(document.getElementsByClassName("box refine")[0].style.display == "block")
  scrollTo(0, 0);
});
// @license-end
