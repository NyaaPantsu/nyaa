// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat

// Switches between themes when a new one is selected

function switchThemes() {
  var themeName = document.getElementById("theme-selector").value
  var head = document.getElementsByTagName("head")[0]
  // Remove the theme in place, it fails if one isn't set
  try {
    head.removeChild(document.getElementById("theme"))
  } catch (err) {}
  // Don't add a node if we don't want extra styling
  if (themeName === "") {
    return
  }
  // Create the new one and put it back
  var newTheme = document.createElement("link")
  newTheme.setAttribute("rel", "stylesheet")
  newTheme.setAttribute("href", "/css/" + themeName + ".css")
  newTheme.setAttribute("id", "theme")
  head.appendChild(newTheme)
}

// Used by spoiler tags
function toggleLayer(elem) {
  if (elem.classList.contains("hide")) {
    elem.classList.remove("hide")
  } else {
    elem.classList.add("hide")
  }
}

function parseAllDates() {
  // Date formatting
  var lang = document.getElementsByTagName("html")[0].getAttribute("lang")
  var ymdOpt = {
    year: "numeric",
    month: "short",
    day: "numeric"
  }
  var hmOpt = {
    hour: "numeric",
    minute: "numeric"
  }

  var list = document.getElementsByClassName("date-short")
  for (var i in list) {
    var e = list[i]
    e.title = new Date(e.innerText).toLocaleString(lang)
    e.innerText = new Date(e.innerText).toLocaleString(lang, ymdOpt)
  }

  var list = document.getElementsByClassName("date-full")
  for (var i in list) {
    var e = list[i]
    var dateDifference = dateDiff(new Date(e.innerText), new Date())
	  e.title = dateDifference.d + " days " + dateDifference.h + " hours ago"
    //e.title = T.r("torrent_age", dateDifference.d, dateDifference.h)
    e.innerText = new Date(e.innerText).toLocaleString(lang)
  }
}
function dateDiff( str1, str2 ) {
    var diff = Date.parse( str2 ) - Date.parse( str1 ); 
    return isNaN( diff ) ? NaN : {
        diff : diff,
        h  : Math.floor( diff /  3600000 %   24 ),
        d  : Math.floor( diff / 86400000        )
    };
}
parseAllDates()

//called if no Commit cookie is set or if the website has a newer commit than the one in cookie
function resetCookies() {
  var cookies = document.cookie.split(";")
  var excludedCookies = ["mascot", "theme", "mascot_url", "lang", "csrf_token"]

  //Remove all cookies but exclude those in the above array
  for (var i = 0; i < cookies.length; i++) {
    var cookieName = (cookies[i].split("=")[0]).trim()
    //Remove spaces because some cookie names have it
    if (excludedCookies.includes(cookieName)) continue
    document.cookie = cookieName + "=;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
  }

  //Set new version in cookie
  var farFuture = new Date()
  farFuture.setTime(farFuture.getTime() + 50 * 36000 * 15000)
  document.cookie = "commit=" + commitVersion + ";expires=" + farFuture.toUTCString()

  var oneHour = new Date()
  oneHour.setTime(oneHour.getTime() + 1 * 3600 * 1500)
  document.cookie = "newVersion=true; expires=" + oneHour.toUTCString()
}


/*Fixed-Navbar offset fix*/
if (document.getElementsByClassName("search-box")[0] !== undefined)
  startupCode()
else
  document.addEventListener("DOMContentLoaded", function (event) {
    startupCode()
  })


function startupCode() {
  var shiftWindow = function () {
    scrollBy(0, -70)
  }
  if (location.hash) shiftWindow()
  window.addEventListener("hashchange", shiftWindow)

  if (!document.cookie.includes("commit"))
    resetCookies()
  else {
    var startPos = document.cookie.indexOf("commit") + 7,
      endPos = document.cookie.substring(startPos).indexOf(";"),
      userCommitVersion = endPos == "-1" ? document.cookie.substring(startPos) : document.cookie.substring(startPos, endPos + startPos);
    //Get start and end position of Commit string, need to start searching endPos from version cookie in case it's not the first cookie in the string
    //If endPos is equal to -1, aka if the version cookie is at the very end of the string and doesn't have an ";", the endPos is not used

    if (userCommitVersion != commitVersion)
      resetCookies()
  }

  if (document.cookie.includes("newVersion"))
    document.getElementById("commit").className = document.getElementById("commit") != "unknown" ? "new" : "wew";
}

function playVoice() {
  var mascotAudio = document.getElementById("explosion") || document.getElementById("nyanpassu") || document.getElementById("nyanpassu2") || document.getElementById("kawaii")
  if (mascotAudio !== undefined) {
    mascotAudio.volume = 0.2
    mascotAudio.play()
  } else {
    console.log("Your mascot doesn't support yet audio files!")
  }
}

document.getElementsByClassName("form-input refine")[0].addEventListener("click", function (e) {
  if(document.getElementsByClassName("form-input search-box")[0].value == "" || location.pathname != "/")
  {
    document.getElementsByClassName("box refine")[0].style.display = document.getElementsByClassName("box refine")[0].style.display == "none" ? "block" : "none"
    if (document.getElementsByClassName("form-input refine-searchbox")[0].value != document.getElementsByClassName("form-input search-box")[0].value)
      document.getElementsByClassName("form-input refine-searchbox")[0].value = document.getElementsByClassName("form-input search-box")[0].value
    if (document.getElementsByClassName("form-input refine-category")[0].selectedIndex != document.getElementsByClassName("form-input form-category")[0].selectedIndex)
      document.getElementsByClassName("form-input refine-category")[0].selectedIndex = document.getElementsByClassName("form-input form-category")[0].selectedIndex
    if (document.getElementsByClassName("box refine")[0].style.display == "block")
      scrollTo(0, 0)
    e.preventDefault()
  }
})

function humanFileSize(bytes, si) {
  var k = si ? 1000 : 1024
  var i = ~~(Math.log(bytes) / Math.log(k))
  return i == 0 ? bytes + " B" : (bytes / Math.pow(k, i)).toFixed(1) + " " + "KMGTPEZY" [i - 1] + (si ? "" : "i") + "B"
}
// @license-end
