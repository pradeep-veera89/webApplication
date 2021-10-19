# WebApplication

## Packages
- Packages are the directory / folder in which the go file is present.
- 
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
