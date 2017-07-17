var Kilo = function (params) {
  // self reference
  var self = this
  // variables
  this.sukebei = (params.sukebei !== undefined) ? params.sukebei : 0
  this.userTrusted = (params.userTrusted !== undefined) ? params.userTrusted : false
  this.isMember = (params.isMember !== undefined) ? params.isMember : false
  this.langSelect = (params.langSelect !== undefined) ? params.langSelect : 'languages'
  this.locale = (params.locale !== undefined) ? params.locale : ''
  this.formatDate = {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }
  this.categories = []
  if (this.locale == '' && document.getElementsByTagName('html')[0].getAttribute('lang') !== null) {
    this.locale = document.getElementsByTagName('html')[0].getAttribute('lang')
  }

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
	
	//Adding the torrent under and above the previewed one. Akuma, you do this
	var torrentHTML = ["", ""];
	document.getElementById("torrentListResults").innerHTML = torrentHTML[0] + document.getElementById("torrentListResults").innerHTML + torrentHTML[1];

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
	
  }
  // Helpers function for events and render
  this.setRemake = function (b) {
    if (b) {
      document.getElementsByName('torrent-info tr')[0].classList.add('remake')
    } else {
      document.getElementsByName('torrent-info tr')[0].classList.remove('remake')
    }
  }
  this.setHidden = function (b) {
    if (!b) {
      document.getElementsByName('torrent-info tr')[0].classList.remove('trusted')
    } else if (this.userTrusted) {
      document.getElementsByName('torrent-info tr')[0].classList.add('trusted')
    }
  }
  this.setName = function (value) {
    document.getElementsByClassName('table-torrent-name')[0].innerText = value
  }
  this.setCategory = function (index) {
    var tableCategory = document.getElementsByClassName('table-torrent-category')[0]
    tableCategory.className = 'nyaa-cat table-torrent-category ' + (this.sukebei ? 'sukebei' : 'nyaa') + '-cat-' + this.categories[index]
    tableCategory.title = document.getElementsByClassName('form-torrent-category')[0].querySelectorAll("option")[index].textContent
  }

  // Event handlers
  var updatePreviewRemake = function (e) {
    var el = e.target
    self.setRemake(el.checked)
  }
  var updatePreviewTorrentName = function (e) {
    var el = e.target
    self.setName(el.value)

    var Keywords_flags= [
	    ["vostfr","[ita]"],
	    ["fr", "it"] ];
	  
    var torrentLowerCaseName = el.value.toLowerCase(),
	updateLang = false;  
	//we don't want to be running the updateLang function for every time the loop loops
	//same for lower case
	  
    for(var KeywordIndex = 0; KeywordIndex < Keywords_flags[0].length; KeywordIndex++)
	if(torrentLowerCaseName.includes(Keywords_flags[0][KeywordIndex])) {
	   document.getElementById("upload-lang-" + Keywords_flags[1][KeywordIndex]).checked = true;
	   updateLang = true;
	}
  
     if(updateLang) updateTorrentLang();
  }

  var updateHidden = function (e) {
    var el = e.target
    self.setHidden(el.checked)
  }
  var updatePreviewCategory = function (e) {
    var el = e.target
    self.setCategory(el.selectedIndex)
  }
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
    document.getElementsByClassName('table-torrent-flag')[0].className = 'table-torrent-flag flag flag-' + langCat
    document.getElementsByClassName('table-torrent-flag')[0].title = langTitle
  }
}
