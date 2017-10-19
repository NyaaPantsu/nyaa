// @source https://github.com/NyaaPantsu/nyaa/tree/dev/public/js
// @license magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt Expat

//String that will contain a far future date, used multiple times throughout multiple functions
var farFutureString 
//Array that will contain the themes that the user will switch between when triggering the function a few lines under
var UserTheme

var Mirror = false

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
  newTheme.setAttribute("href", "/css/themes/" + themeName + ".css")
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

  var list = document.getElementsByClassName("date-short")
  for(var i = 0; i < list.length; i++) {
    var e = list[i]
    if(e.className.includes("date-converted")) continue
    e.innerText = new Date(e.title).toLocaleString(lang, ymdOpt)
    e.title = new Date(e.title).toLocaleString(lang)
    e.className = e.className + " date-converted"
  }

  var list = document.getElementsByClassName("date-full")
  for(var i = 0; i < list.length; i++) {
    var e = list[i]
    var dateDifference = dateDiff(new Date(e.innerText), new Date())
    

    e.title = (dateDifference.d == 0 ? "" : dateDifference.d+" days ")
    e.title = e.title + (dateDifference.h == 0 ? "" : dateDifference.h+" hours ")
    e.title = e.title + (dateDifference.m == 0 ? "" : dateDifference.m+" minutes ")
	if(e.title == "") 
		e.title = dateDifference.s + " seconds "
    e.title = e.title + "ago"
	  
    e.innerText = new Date(e.innerText).toLocaleString(lang)
  }
}
function dateDiff( str1, str2 ) {
    var diff = Date.parse( str2 ) - Date.parse( str1 ); 
    return isNaN( diff ) ? NaN : {
        diff : diff,
		s  : Math.floor( diff /     1000          ),
		m  : Math.floor( diff /    60000 %     60 ),
        h  : Math.floor( diff /  3600000 %     24 ),
        d  : Math.floor( diff / 86400000          )
    };
}
parseAllDates()

//called if no Commit cookie is set or if the website has a newer commit than the one in cookie
function resetCookies() {
  var cookies = document.cookie.split(";")
  var excludedCookies = ["session", "mascot", "version", "theme", "theme2", "mascot_url", "lang", "csrf_token", "altColors", "EU_Cookie", "oldNav"]
  //Excluded cookies are either left untouched or deleted then re-created
  //Ignored Cookies are constantly left untouched
  
  //Get HostName without subDomain
  var hostName = window.location.host

  if(!Mirror) {
	  var lastDotIndex = hostName.lastIndexOf(".")
	  var secondLast = -1
	  
	  for(var index = 0; index < lastDotIndex; index++) {
		if(hostName[index] == '.')
		  secondLast = index
	  }
	  hostName = hostName.substr(secondLast == -1 ? 0 : secondLast)
  }

  for (var i = 0; i < cookies.length; i++) {
    var cookieName = (cookies[i].split("=")[0]).trim()
    //Trim spaces because some cookie names have them at times
    if (excludedCookies.includes(cookieName)) {
      if(domain == hostName) {
	//only execute if cookie are supposed to be shared between nyaa & sukebei, aka on host name without subdomain
        var cookieValue = getCookieValue(cookieName)
        document.cookie = cookieName + "=;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
        document.cookie = cookieName + "=;path=/;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
        if(cookieName != "session")
	  document.cookie = cookieName + "=" + cookieValue + ";path=/;expires=" + farFutureString + ";domain=" + domain
	else document.cookie = cookieName + "=" + cookieValue + ";path=/;expires=" + farFutureString
        //Remove cookie from both current & general path, then re-create it to ensure domain is correct
	//Force current domain for session cookie
        }
      continue
    }
    document.cookie = cookieName + "=;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
    document.cookie = cookieName + "=;path=/;expires=Thu, 01 Jan 1970 00:00:00 UTC;"
  }

  //Set new version in cookie
  document.cookie = "commit=" + commitVersion + ";path=/;expires=" + farFutureString + ";domain=" + domain
  document.cookie = "version=" + websiteVersion + ";path=/;expires=" + farFutureString + ";domain=" + domain

  var oneHour = new Date()
  oneHour.setTime(oneHour.getTime() + 1 * 3600 * 1500)
  document.cookie = "newVersion=true;path=/;expires=" + oneHour.toUTCString() + ";domain=" + domain
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


  if(!window.location.host.includes(domain)) {
	  domain = window.location.host
	  Mirror = true
  }

  if (!document.cookie.includes("commit") && !document.cookie.includes("version"))
    resetCookies()
  else {
    var userCommitVersion = getCookieValue("commit"), userWebsiteVersion = getCookieValue("version");
    if (userCommitVersion != commitVersion || userWebsiteVersion != websiteVersion)
      resetCookies()
  }
  
  if(document.getElementById("cookie-warning-close") != null) {
	document.getElementById("cookie-warning-close").addEventListener("click", function (e) {
      document.getElementById("cookie-warning").outerHTML = "";
      document.cookie = "EU_Cookie=true;path=/;expires=" + farFutureString + ";domain=" + domain
    })
  }
	  

  if (document.cookie.includes("newVersion"))
    document.getElementById("commit").className = document.getElementById("commit").innerHTML != "unknown" ? "new" : "wew";

  document.getElementById("dark-toggle").addEventListener("click", toggleTheme);

  if(document.cookie.includes("theme=")) {
    UserTheme = [getCookieValue("theme"), darkTheme]
    //Get user's default theme and set the alternative one as dark theme
  }
  else 
    UserTheme = ["g", darkTheme]
   //If user has no default theme, set these by default
  
  
  if(document.cookie.includes("theme2=")) {
    UserTheme[1] = getCookieValue("theme2")
    //If user already has ran the ToggleTheme() function in the past, we get the value of the second theme (the one the script switches to)
    if(!UserTheme.includes(darkTheme))
      UserTheme[1] = darkTheme
    //If none of the theme are darkTheme, which happens if the user is on dark mode (with theme2 on g.css) and that he switches to classic or g.css in settings, we set the second one as darkTheme
    else if(UserTheme[0] == UserTheme[1])
      UserTheme[1] = "g"
    //If both theme are darkTheme, which happens if theme2 is on darkTheme (always is by default) and that the user sets darkTheme as his theme through settings page, we set secondary theme to g.css
  }
  else if(UserTheme[0] == UserTheme[1])
    UserTheme[1] = "g"
    //If darkTheme is twice in UserTheme, which happens when the user already has darkTheme as his default theme and toggle the dark mode for the first time, we set the second theme as g.css
  
}

function toggleTheme(e) {
  var CurrentTheme = document.getElementById("theme").href
  CurrentTheme = CurrentTheme.substring(CurrentTheme.indexOf("/themes/") + 8, CurrentTheme.indexOf(".css"))
  CurrentTheme = (CurrentTheme == UserTheme[0] ? UserTheme[1] : UserTheme[0])

  document.getElementById("theme").href = "/css/themes/" + CurrentTheme + ".css";
  
  if(UserID > 0 ){
    Query.Get("/dark", function(data) {})
    //If user logged in, we're forced to go through this page in order to save the new user theme
  }
  else {
    document.cookie = "theme=" + CurrentTheme + ";path=/;expires=" + farFutureString + ";domain=" + domain
    document.cookie = "theme2=" + (CurrentTheme == UserTheme[0] ? UserTheme[1] : UserTheme[0]) + ";path=/;expires=" + farFutureString + ";domain=" + domain
    //Otherwise, we can just set the theme through cookies
  }
  
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

document.getElementsByClassName("form-input refine-btn")[0].addEventListener("click", function (e) {
  var inputs = document.querySelectorAll(".box.refine form input")
  var select = document.querySelectorAll(".box.refine form select")
	
  for(var i = 0; i < inputs.length; i++) 
    if(inputs[i].value == "") inputs[i].disabled = true;
	
  for(var i = 0; i < select.length; i++) 
    if(select[i].selectedIndex == 0) select[i].disabled = true;
	
	
  if(document.querySelector(".box.refine form input[name='limit']").value == "50")
    document.querySelector(".box.refine form input[name='limit']").disabled = true

  if(document.querySelector(".box.refine form select[name='sort']").selectedIndex == 5)
    document.querySelector(".box.refine form select[name='sort']"). disabled = true;
  else  document.querySelector(".box.refine form select[name='sort']"). disabled = false;
	
  if(document.querySelector(".box.refine form select[name='order']").selectedIndex == 1)
    document.querySelector(".box.refine form select[name='order']"). disabled = true;
  else  document.querySelector(".box.refine form select[name='order']"). disabled = false;
	
  if(document.querySelector(".box.refine form select[name='sizeType']").selectedIndex == 2 &&
    document.querySelector(".box.refine form input[name='minSize']").value == "" &&
	document.querySelector(".box.refine form input[name='maxSize']").value == "")
	document.querySelector(".box.refine form select[name='sizeType']"). disabled = true
  else document.querySelector(".box.refine form select[name='sizeType']"). disabled = false
  
  if(document.querySelector(".box.refine form select[name='order']").selectedIndex == 1)
    document.querySelector(".box.refine form select[name='order']"). disabled = true;
  else  document.querySelector(".box.refine form select[name='order']"). disabled = false;
})

function humanFileSize(bytes, si) {
  if (bytes == 0) 
    return "Unknown"
  var k = si ? 1000 : 1024
  var i = ~~(Math.log(bytes) / Math.log(k))
  return i == 0 ? bytes + " B" : (bytes / Math.pow(k, i)).toFixed(1) + " " + "KMGTPEZY" [i - 1] + (si ? "" : "i") + "B"
}

function getCookieValue(cookieName) {
    var startPos = document.cookie.indexOf(cookieName + "=") 
    if(startPos == -1) return ""
    startPos +=  cookieName.length + 1
    var endPos = document.cookie.substring(startPos).indexOf(";")
    return endPos == -1 ? document.cookie.substring(startPos) : document.cookie.substring(startPos, endPos + startPos)
}

// @license-end
