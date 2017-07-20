var Kilo = function (params) {
  // self reference
  var self = this

  // public variables
  // Boolean defining if we are in sukebei
  this.sukebei = (params.sukebei !== undefined) ? params.sukebei : 0
  // Boolean defining if a user is trusted
  this.userTrusted = (params.userTrusted !== undefined) ? params.userTrusted : false
  // Boolean defining if a user is logged
  this.isMember = (params.isMember !== undefined) ? params.isMember : false
  // Boolean enabling the AJAX load of torrents
  this.listContext = (params.listContext !== undefined) ? params.listContext : false
  // Variable defining the <select> of languages
  this.langSelect = (params.langSelect !== undefined) ? params.langSelect : 'languages'
  // Variable defining the language of the user
  this.locale = (params.locale !== undefined) ? params.locale : ''
  // Format of the date in the torrent listing
  this.formatDate = {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }
  // Array of categories (loaded from the select html tag categories)
  this.categories = []

  // if no locale provided as a parameter, fallback to the language set by html tag
  if (this.locale == '' && document.getElementsByTagName('html')[0].getAttribute('lang') !== null) {
    this.locale = document.getElementsByTagName('html')[0].getAttribute('lang')
  }

  // Private variables
  var Keywords_flags= [
	["vostfr","vosfr", "[ita]", "[eng]", " eng ","[english]","[english sub]", "[jp]","[jpn]","[japanese]"],
  ["fr","fr", "it", "en","en","en","en", "ja","ja","ja"]
  ]
  var Keywords_categories = [
		[ ["[jav]"], [7] ], 
		[ [""], [0] ]
  ]

  // Parsing categories
  document.querySelectorAll(".form-torrent-category option").forEach(function(el) {
    var subcat
    if (self.sukebei) {
      subcat = el.value.replace("_", "")
    } else {
      subcat = el.value.split("_")[1]
    }
    subcat = (subcat === undefined) ? 0 : subcat
    self.categories.push(subcat)
  })

  this.render = function () {
    console.log(this.locale)
    // Displaying the block and set the locale timestamp
    document.getElementsByClassName('torrent-preview-table')[0].style.display = 'block'
    document.getElementsByClassName('table-torrent-date')[0].innerText = new Date(Date.now()).toISOString()
	
    // Adding listener events
    for (var langIndex = 0; langIndex < document.getElementsByName(this.langSelect).length; langIndex++) {
      document.getElementsByName(this.langSelect)[langIndex].addEventListener('change', updateTorrentLang)
    }
    var formCategory = document.getElementsByClassName('form-torrent-category')[0]
    formCategory.addEventListener('change', updatePreviewCategory)
    var formName = document.getElementsByClassName('form-torrent-name')[0]
    formName.addEventListener('keyup', updatePreviewTorrentName)
    var formRemake = document.getElementsByClassName('form-torrent-remake')[0]
    formRemake.addEventListener('change', updatePreviewRemake)
    var formHidden = document.getElementsByClassName('form-torrent-hidden')[0]
    if (this.isMember) {
      formHidden.addEventListener('change', updateHidden)
    }

    // Setting default values
    this.setRemake(formRemake.checked)
    if (this.isMember) {
      this.setHidden(formHidden.checked)
    }
    this.setName(formName.value)
    this.setCategory(formCategory.selectedIndex)
    updateTorrentLang()

    //Adding the torrent under and above the previewed one. 
    if (this.listContext) {
      Query.Get('/api/search?limit=2', function (data) {
        torrents = data.torrents
        var torrentHTML = []
        var l = torrents.length
        for (var i = 0; i < l; i++) {
          torrentHTML.push(Templates.Render('torrents.item', torrents[i]))
        }
        document.getElementById("torrentListResults").innerHTML = torrentHTML[0] + document.getElementsByName("torrent-info tr")[0].outerHTML + torrentHTML[1]
      })
    }
  }

  // Helpers function for events and render
  // set the class remake with b a boolean
  this.setRemake = function (b) {
    if (b) {
      document.getElementsByName('torrent-info tr')[0].classList.add('remake')
    } else {
      document.getElementsByName('torrent-info tr')[0].classList.remove('remake')
    }
  }
  // set the class hidden with b a boolean
  this.setHidden = function (b) {
    if (!b) {
      document.getElementsByName('torrent-info tr')[0].classList.remove('trusted')
    } else if (this.userTrusted) {
      document.getElementsByName('torrent-info tr')[0].classList.add('trusted')
    }
  }
  // set the name of the torrent according to value string
  this.setName = function (value) {
    document.getElementsByClassName('table-torrent-name')[0].innerText = value
  }
  // set the category of the torrent according to index int
  this.setCategory = function (index) {
    var tableCategory = document.getElementsByClassName('table-torrent-category')[0]
    tableCategory.className = 'nyaa-cat table-torrent-category ' + (this.sukebei ? 'sukebei' : 'nyaa') + '-cat-' + this.categories[index]
    tableCategory.title = document.getElementsByClassName('form-torrent-category')[0].querySelectorAll("option")[index].textContent
  }
  // 
  this.addKeywordFlags = function(value) {
    var torrentLowerCaseName = value.toLowerCase()
    var updateLang = false
    
    for(var KeywordIndex = 0; KeywordIndex < Keywords_flags[0].length; KeywordIndex++)
      if(torrentLowerCaseName.includes(Keywords_flags[0][KeywordIndex])) {
		   document.getElementById("upload-lang-" + Keywords_flags[1][KeywordIndex]).checked = true
		   updateLang = true
      }
  
    if(updateLang) updateTorrentLang()
  }
  //
  this.addKeywordCategories = function(value) {
    if(document.getElementsByClassName('form-torrent-category')[0].selectedIndex != 0)
      return
		
    var torrentLowerCaseName = value.toLowerCase(),
      IsOnSukebei = params.sukebei ? 0 : 1;

    for(var KeywordIndex = 0; KeywordIndex < Keywords_categories[IsOnSukebei][0].length; KeywordIndex++)
      if(torrentLowerCaseName.includes(Keywords_categories[IsOnSukebei][0][KeywordIndex])) {
        document.getElementsByClassName('form-torrent-category')[0].selectedIndex = Keywords_categories[IsOnSukebei][1][KeywordIndex];
        this.setCategory(document.getElementsByClassName('form-torrent-category')[0].selectedIndex)
        break
      }	
  }
  // Helper to prevent the functions on keyup/keydown to slow the user typing
  this.debounce = function (func, wait, immediate) {
    var timeout
    return function() {
      var context = this, args = arguments
      var later = function() {
        timeout = null
        if (!immediate) func.apply(context, args)
      }
      var callNow = immediate && !timeout
      clearTimeout(timeout)
      timeout = setTimeout(later, wait)
      if (callNow) func.apply(context, args)
    }
  }

  // Event handlers
  // Event on remake checkbox
  var updatePreviewRemake = function (e) {
    var el = e.target
    self.setRemake(el.checked)
  }
  // Event on torrent name keyup
  var updatePreviewTorrentName = function (e) {
    var el = e.target
    self.setName(el.value)
    self.debounce(function(value) {
      self.addKeywordFlags(value)
      self.addKeywordCategories(value)
    }, 300)(el.value)
  }
  // Event on hidden checkbox
  var updateHidden = function (e) {
    var el = e.target
    self.setHidden(el.checked)
  }
  // Event on cateogry change
  var updatePreviewCategory = function (e) {
    var el = e.target
    self.setCategory(el.selectedIndex)
  }
  // Event on languages checkbox
  var updateTorrentLang = function () {
    var langCount = 0
    var langValue = 'other'
    var langTitle = ''

    for (var langIndex = 0; langIndex < document.getElementsByName('languages').length; langIndex++) {
      if (document.getElementsByName('languages')[langIndex].checked) {
        langTitle = langTitle + document.getElementsByName('upload-lang-languagename')[langIndex].innerText + ','
        if (++langCount > 1) {
          langValue = 'multiple'
          continue
        }
        langValue = document.getElementsByName('languages')[langIndex].value
      }
    }
    var langCat = langValue !== 'other' ? (langValue > 1 ? 'multiple' : langValue) : 'other'
    document.getElementsByClassName('table-torrent-flag')[0].className = 'table-torrent-flag flag flag-' + flagCode(langCat)
    document.getElementsByClassName('table-torrent-flag')[0].title = langTitle
  }
}
