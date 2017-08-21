// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat

//String that will contain a far future date, used multiple times throughout multiple functions
var farFutureString 
//Array that will contain the themes that the user will switch between when triggering the function a few lines under
var UserTheme

// Switches between themes when a new one is selected
function switchThemes() {
  var themeName = document.getElementById("theme-selector").value
  var head = document.getElementsByTagName("head")[0]
  
  if (themeName === "") {
    return
  }
  
  // Remove the theme in place, it fails if one isn't set
  try {
    head.removeChild(document.getElementById("theme"))
  } catch (err) {}
  // Don't add a node if we don't want extra styling
	
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
	if(list.length == 0)
	  break;
    var e = list[i]
    var dateDifference = dateDiff(new Date(e.innerText), new Date())
    
    if(e.className != undefined && e.className.includes("scrape-date"))
      e.title = ((dateDifference.d * 24) + dateDifference.h) + " hours " + dateDifference.m + " minutes ago"
    else
      e.title = dateDifference.d + " days " + dateDifference.h + " hours ago"
	  
    e.innerText = new Date(e.innerText).toLocaleString(lang)
  }
}
function dateDiff( str1, str2 ) {
    var diff = Date.parse( str2 ) - Date.parse( str1 ); 
    return isNaN( diff ) ? NaN : {
        diff : diff,
	m  : Math.floor( diff /     60000 %   60 ),
        h  : -Math.floor( diff /  3600000 %   24 ),
        d  : -Math.floor( diff / 86400000        )
    };
}
parseAllDates()

//called if no Commit cookie is set or if the website has a newer commit than the one in cookie
function resetCookies() {
  var cookies = document.cookie.split(";")
  var excludedCookies = ["mascot", "theme", "theme2", "mascot_url", "lang", "csrf_token"]

  //Remove all cookies but exclude those in the above array
  for (var i = 0; i < cookies.length; i++) {
    var cookieName = (cookies[i].split("=")[0]).trim()
    //Remove spaces because some cookie names have it
    if (excludedCookies.includes(cookieName)) continue
    document.cookie = cookieName + "=;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
  }

  //Set new version in cookie
  document.cookie = "commit=" + commitVersion + ";expires=" + farFutureString

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
  farFutureString = new Date()
  farFutureString.setTime(farFutureString.getTime() + 50 * 36000 * 15000)
  farFutureString = farFutureString.toUTCString()
  
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
    document.getElementById("commit").className = document.getElementById("commit").innerHTML != "unknown" ? "new" : "wew";

  document.getElementById("dark-toggle").addEventListener("click", toggleTheme);

  if(document.cookie.includes("theme")) {
    var startPos = document.cookie.indexOf("theme=") + 6
    var endPos = document.cookie.substring(startPos).indexOf(";")
    UserTheme = [endPos == "-1" ? document.cookie.substring(startPos) : document.cookie.substring(startPos, endPos + startPos), "tomorrow"]
    //Get user's default theme and set the alternative one as tomorrow
  }
  else 
    UserTheme = ["g", "tomorrow"]
   //If user has no default theme, set these by default
  
  
  if(document.cookie.includes("theme2")) {
    var startPos = document.cookie.indexOf("theme2=") + 7
    var endPos = document.cookie.substring(startPos).indexOf(";")
    UserTheme[1] = endPos == "-1" ? document.cookie.substring(startPos) : document.cookie.substring(startPos, endPos + startPos)
    //If user already has ran the ToggleTheme() function in the past, we get the value of the second theme (the one the script switches to)
    if(!UserTheme.includes("tomorrow"))
      UserTheme[1] = "tomorrow"
    //If none of the theme are tomorrow, which happens if the user is on dark mode (with theme2 on g.css) and that he switches to classic or g.css in settings, we set the second one as tomorrow
    else if(UserTheme[0] == UserTheme[1])
      UserTheme[1] = "g"
    //If both theme are tomorrow, which happens if theme2 is on tomorrow (always is by default) and that the user sets tomorrow as his theme through settings page, we set secondary theme to g.css
  }
  else {
    if(UserTheme[0] == UserTheme[1])
      UserTheme[1] = "g"
    //If tomorrow is twice in UserTheme, which happens when the user already has tomorrow as his default theme and toggle the dark mode for the first time, we set the second theme as g.css
    document.cookie = "theme2=" + UserTheme[1] + ";path=/;domain=pantsu.cat;expires=" + farFutureString
    //Set cookie for future theme2 uses
  }
  
}

function toggleTheme(e) {
  var CurrentTheme = document.getElementById("theme").href
  CurrentTheme = CurrentTheme.substring(CurrentTheme.indexOf("/css/") + 5, CurrentTheme.indexOf(".css"))
  CurrentTheme = (CurrentTheme == UserTheme[0] ? UserTheme[1] : UserTheme[0])

  document.getElementById("theme").href = "/css/" + CurrentTheme + ".css";
  
  document.cookie = "theme=" + CurrentTheme + ";path=/;domain=pantsu.cat;expires=" + farFutureString
  document.cookie = "theme2=" + (CurrentTheme == UserTheme[0] ? UserTheme[1] : UserTheme[0]) + ";path=/;domain=pantsu.cat;expires=" + farFutureString
  e.preventDefault()
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
    if (document.getElementsByClassName("form-input refine-searchbox")[0].value == "")
      document.getElementsByClassName("form-input refine-searchbox")[0].value = document.getElementsByClassName("form-input search-box")[0].value
    if (document.getElementsByClassName("form-input refine-category")[0].selectedIndex == 0)
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
