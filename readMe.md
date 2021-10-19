# WebApplication

## Packages
- Packages are the directory / folder in which the go file is present.
  
### Config (pkg/config/config.go)
- Application wide config.
- This cofiguration is accesed by any kind of application
- This package is imported by other parts of the application but the config package doesnot import any other packages from the application.
- Add Template Cache to the appConfig.

### Models
- Added templateData struct to store all the infromation which are passed to template
   
## Handlers
-  Package : "net/http"
-  HandleFunc
-  ListenAndServe


## Templates
- Always use helper (renderTemplate) to parse the template
- Go uses 'template.ParseFiles' to parse the html template and 'Execute' for executing the template.
- Generate template cache to automatically fetch the templates.
  
### Base layout
- Add a base layout by defining base '{{define "base"}}'
- Add content block followed by "."
  
### Child layout
- extends the base layout using '{{template "base"}}'
- define the custom body by using '{{define "content"}}'
  
### StadardFunctions
- Package : "html/template"
-  ParseFiles
-  Execute

## Points
- Function name starting with Capital letter are visible to other packages inside the project.
- Function name starting with small letter are not visible to other packages inside the project
- Comments for the function should start with function name

### Import Cycle Problem
- There are package with name A and B
- This problem occurs when A imports B and B also imports A.
