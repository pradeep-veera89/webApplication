# WebApplication
## CMD (cmd)
#### web (cmd/web)
- main.go
  - Contains main.go
- routes.go
  - Contains all the routes of the application
  - uses chi router for building the routes.
  - chi router has middleware which is not present in standard http router
  - Here we have used Recoverer middleware from chi to provide stacktrace and error log in case of panic situation
  - Used NoSurf Package 'github.com/justinas/nosurf' for CSRF Token Protection
  - CSRFToken is generated using middleware of the CHI Router.
  - SessionLoad to load and save sessions from SCS Session Package
- middleware.go
  - new middleware to add CSRF Protection
  - new middleware to load and save Sessions from SCS Session package. 
  
## PKG (pkg)
- Packages are the directory / folder in which the go file is present.
  
#### Config (pkg/config/config.go)
- Application wide config.
- This cofiguration is accesed by any kind of application
- This package is imported by other parts of the application but the config package doesnot import any other packages from the application.
- Add Template Cache to the appConfig.

#### Models (pkg/models/config.go)

- config.go (pkg/models/config.go)
  - Added templateData struct to store all the infromation which are passed to template
   
#### Handlers
-  Package : "net/http"
-  HandleFunc
-  ListenAndServe

#### Renders
- Creates template cache and executes the template based on route/handler
- Also adds the default data while executing the template.


## Templates
- Always use helper (renderTemplate) to parse the template
- Go uses 'template.ParseFiles' to parse the html template and 'Execute' for executing the template.
- Generate template cache to automatically fetch the templates.
  
#### Base layout
- Add a base layout by defining base '{{define "base"}}'
- Add content block followed by "."
  
#### Child layout
- extends the base layout using '{{template "base"}}'
- define the custom body by using '{{define "content"}}'
  
#### StadardFunctions
- Package : "html/template"
-  ParseFiles
-  Execute

## Points
- Function name starting with Capital letter are visible to other packages inside the project.
- Function name starting with small letter are not visible to other packages inside the project
- Comments for the function should start with function name

#### Import Cycle Problem
- There are package with name A and B
- This problem occurs when A imports B and B also imports A.

## External Golang Packages
- NoSurf Package 'github.com/justinas/nosurf' for CSRF Token Protection
- CHI Router Package 'go get -u github.com/go-chi/chi/v5' for Routing instead of http router.
- SCS (Session) Package 'https://github.com/alexedwards/scs' for session management.