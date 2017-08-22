# How-To Controller
Our routing system is based on [Gin-Gonic](https://github.com/gin-gonic/gin) HTTP web framework. For more information on how to use a gin.Context and what are the possibilities offered by Gin, you should take a look into their [documentation](https://github.com/gin-gonic/gin#api-examples).
This readme will only explain to you how the routing works and how to easily render a static 'hello world!' template.

## How the routing works
1. Our routing system based on gin has one main file that connect all our controllers. This file is controllers/router.go [here](https://github.com/NyaaPantsu/nyaa/blob/dev/controllers/router.go).
You can see that in the first lines, that this file import multiple sub directories. All the import that start with an underscore is a controller package. In fact, the main `router.go` only imports the different controllers and do nothing else. So if you need to create a whole new controller, you would have to **add it in this main `router.go` as an import that start with an underscore**.
2. Every controller package have another `router.go` file. It can be named differently but for homogenuity purposes, please name it like that. This file is the core of your controller package, you need it to add your controllers to the main router. All of the sub router.go files look pretty much the same.
For example, [here](https://github.com/NyaaPantsu/nyaa/blob/dev/controllers/search/router.go) is the `searchController` package `router.go` file.
* This file first import the `controllers/router` package. This step is important in order to have access to the function `router.Get()`.
* After it declares an `init()` function. This function is called automatically when the `searchController` package is imported.
* In this function we call `router.Get().XXX()`. This part is where you link the different controllers to the routes. The XXX part are gin functions. In fact, `router.Get()` return a gin.Router so you can use all the functions from the gin.Router here.
* Therefore `router.Get().Any("/", SearchHandler)` tells that every time you go to "/", it is the SearchHandler controller that should be called.
3. Every controller package have other go files. Those files should be where you put your controllers functions. Try to not put every controller functions in the same file. A controller function is the same as those designed by the Gin documentation. They should be function that accept as only argument: `c *gin.Context` and return nothing. For example, you can see [here](https://github.com/NyaaPantsu/nyaa/blob/dev/controllers/search/search.go#L19) the `SearchHandler` controller function.

## Router Nyaa Functions
Our routing system doesn't have much function that you can use besides the gin defined one. However we have one that you might find very helpful, it is `router.GetUser(c)`. This function declared in the `controllers/router` package returns the current `User` during the request. You can only use this function in a controller function which have access to the `c *gin.Context`.

## But how to display something in a controller function?
Adding a controller function should be easy now for you, and you might have seen some gin examples to display thing and understood that it works. However, except if you want to render some JSON, **please do not use the default gin rendering template for html**. As you may have seen, we have our own templating system that is based on CloudyKit Jet Template. This template system is lighter and have better performance while having more functionnalities.
To use this templating system, it is quite easy. First you need to import the templates package `templates` then you have multiple functions available to render some basic template.
For example, `templates.ModelList(c, "url/of/template/file.jet.html", ArrayOfModels, Navigation, SearchForm)` is a function displaying a template file with a local variable `Models` which is your `ArrayOfModels`.
You also have `templates.Form(c, "url/of/template/file.jet.html", form)` which is a function displaying a template file with a `Form` local variable that is equal to the `form` variable passed.
Furthermore you have `templates.Static(c, "url/of/template/file.jet.html")` which is a function displaying a static template file without any local variables
Finally, you can use a mix of `templates.CommonVariables` and `templates.Render(c, "url/of/template/file.jet.html", variables)` to render more specific templates which need special cases like [here](https://github.com/NyaaPantsu/nyaa/blob/dev/controllers/search/search.go#L70). 
When you have rendered and you don't want to continue executing the function after the display of the template, you have to return. For example, [here](https://github.com/NyaaPantsu/nyaa/blob/dev/controllers/search/search.go#L70), we want to stop the function and render the template "errors/no_results.jet.html", so we return just after calling `templates.Render()`.

## Mini Tutorial: Hello World!
To do this you will need a template file, a controller function and that's all!
### The template file
This tutorial won't go deep in the templating system, for more information about it, go to the readme in `/templates/`. First create a new file named "hello.jet.html" in `/templates/site/static/` with this in it:
```
<p>Hello <em>World</em>!</p>
```
Now let's move to the controller

### The controller package
First let's create a new directory named "helloworld" in `/controllers/`. In this directory we will create two files, a `router.go` and a `helloworld.go`.
In `helloworld.go` we will place in it the controller function `HelloWorld(c *gin.Context)` and render the static template:
```
package helloWorldController

import (
	"github.com/NyaaPantsu/nyaa/templates" // We import the templates package to use templates.* functions
	"github.com/gin-gonic/gin" // We import the gin package to use c *gin.Context
)

// HelloWorld : Controller for Hello World view page
func HelloWorld(c *gin.Context) {
  // We render the static template 
	templates.Static(c, "site/static/hello.jet.html")
  // We don't have to return since it's the end of the function
}
```

In `router.go` we will link the controller function to the route we want to use it:
```
package helloWorldController

import "github.com/NyaaPantsu/nyaa/controllers/router" // We import the router package to get the gin.Router and add the route
// When the package is called, we do what's inside
func init() {
  // we get the gin.Router an for Any method (Get/POST/PUT...) we associate /hello path url to the HelloWorld controller
	router.Get().Any("/hello", HelloWorld)
}
```

### The link to the main `router.go`
Now we need to tell to the Nyaa Programm that we have yet another controller package to add, for this you need to go to `controllers/router.go` and add the following line to the import section:
```
	_ "github.com/NyaaPantsu/nyaa/controllers/helloworld"           // helloworld controller
```

### The End
Now you can save everything and compile/run and go to http://localhost:9999/hello to see your wonderfull controller in action!
