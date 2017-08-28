# How-To Template 
Templating system is based on [CloudyKit Jet Template](https://github.com/CloudyKit/jet). Therefore it is very much like the Golang basic template system but with some improvements.
For all syntaxic question, it is recommanded to look into Jet Template documentation [here](https://github.com/CloudyKit/jet/wiki/Jet-template-syntax).

## File naming
You can, pretty much, name your files however you want. But, to keep some homogenuity, it would be preferable to keep the name in lowercase, you can use underscores and you have to use the suffix ".jet.html"

## Global Variables
In every template file, besides the built-in functions from Jet Template, you have also some Nyaa Global Variables that you can access to. Those variables are set in `template.go` in CommonVariable function. If you want to add a global variable, it is [there](https://github.com/NyaaPantsu/nyaa/blob/dev/templates/template.go#L58). Be aware to set only important variables here, not page specific one.
### What do they do?
Here we will try to look into each variable and explain how do they work.
* `Config` variable is the website configuration set in config/config.yml. You can access every exported Properties and functions from `config/struct.go`. For example, Config.Torrents.Tags.Default return the default tag type and Config.Port return the running port.    
* `CsrfToken` variable return a string for the csrf input field. You have to use it in a hidden input every time you want to make a POST request. 
* `Errors` variable is a map[string][] with all the errors listed, from an internal error to a more simple form input error. *Better to use the yield function `errors()`.*  
* `Infos` variable is a map[string][] with all the information messages listed. *Better to use the yield function `infos()`.*
* `Mascot` variable is the choosen mascot.
* `MascotURL` variable is the url of the mascot set by the user.
* `Search` variable is link to the search form state. For example, you can access to the set category with `Search.Category`. This variables also inherit all exported Properties from `TorrentParam` (utils/search/torrentParam.go).
* `Theme` variable is the name of the choosen theme.
* `URL` variable is the current URL of the page. This variable is a net/url.Url struct. Therefore you can use all exported properties and functions from it. More information in the [golang doc](https://golang.org/pkg/net/url/#URL).
* `User` variable is the current User. This variable is defined in models/user.go. Every exported properties and functions are available. For example, `User.Username` gives you the user's username and `User.IsModerator()` tells you if a user is a moderator or not.

### How to use them?
Pretty simple, just type: `{{ NameOfVariable }}`. For example, {{ URL.String() }} to get the current URL.

## Global Functions
Same as global variables, there are also global functions. They are all defined [here](https://github.com/NyaaPantsu/nyaa/blob/dev/templates/template_functions.go#L24).
* `contains(language, string)` tells you if a language corresponds to the language code provided
* `genActivityContent(activity, T function)` returns the translated activty.
* `genUploaderLink(uploaderID, uploaderName, hidden)` return a `<a href="">Username</a>` for the user provided based on the elements provided. For example, if you provide a username and no userID or a hidden bool to true. Then it won't return a link but will return "renchon" username.
* `getCategory(MainCategory, keepParent)` return an array of `Category` (struct set in utils/categories/categories.go) based on the MainCategory string provided and the bool keepParent. If keepParent is true, the MainCategory is included.
* `categoryName(maincat, subcat)` returns the category name
* `Sukebei()` returns a boolean, true if the website is sukebei, false if not
And many more...
### How to use them?
Same as global variables, just do `{{ function() }}`. For example, `flagCode("fr-fr")` returns the flag code "fr".

## Helpers Functions
Beside global functions accessible whenever you want, you can also import template functions or make them.
All template functions beside the global ones, should be set in separate files in `/templates/layouts/partials/helpers/`.
List of template functions available (non exhaustive, please look for them in helpers folder):
* `badge_user()` is a function used in menu to display the user badge
* `captcha(captchaid)` is a function used in forms to display a captcha input, it needs a captchaid to work
* `csrf_field()` is a function used in forms to add a hidden input for the csrf check. The token is directly taken from the global variables
* `errors(name)` is a function used to display the error messages. You just need to specify what errors do you want.
* `infos(name)` is a function used to display the informations messages. You just to specify what infos do you want.
### How to use them?
A bit less easy, you need to use this format: `{{ yield function(param="something") }}`. For example, `{{ yield csrf_field() }}` will add a csrf hidden input or `{{ errors(name="Description") }}` will display all the errors for Description.
Furthermore, you need to import (relative to /templates/) the file where the function is defined at the **start** of template! For example, `{{ import "layouts/partials/helpers/csrf" }}` to import the csrf_field() function and use it.
