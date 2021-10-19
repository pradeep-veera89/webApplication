
## Points to remember
- use time.Time for Datetime datatype.
- int16 or int32 or int64 is choosen automatically based whether 16bit or 32bit or 64bit computer 
- Struct with Capital letter allows to use struct in another package.
  - Go does not have public or private functionality
  - func or struct or variable with Capital letters is similar to Public.
  - func or struct or variable with small is similar to private,

## Why Go?
- Compiles to single binary file.(Complied Language)
- No runtime to worry about (PHP or Python should have correct version to be installed.)
- Statically typed.
- Object oriented(sort of , has some of the elements of oop)
- Concurrency (can run application using differenet cors at the same time)
- Cross platform(Linux, Windows, Max)
- Excellent package management(similar to composer in PHP) & testing built in(similar to PHP Unit).

## Datastructures
### Map
- we use make keyword to create.
- specify index type inside the square brackets, following the square brackets is the type stored inside the index.
- syntax to create  
  - 'myMap := make(map[string]string)'  -> StringMap with Key as String
  - 'myMap := make(map[int]string)'     -> StringMap with Key as Int.
  - 'myMap := make(map[int]int)'        -> IntMap
  

[MapWithStruct](/home/pradeep/Pictures/Golang/tutorial/GolangMaps.png)

- Maps are build into system but not sorted, so should be looked by using key.
- Maps can be passed by using pointer.
- When the return type of Map is unknown interface can be used, but its not recommended all the time. 'myMap := make(map[string]interface{})'

### Slices
- Similar to array in PHP 
- Considering for strings, if we need to store more than on string in a variable.
- syntax : 'var myString[] string'

[slice](/home/pradeep/Pictures/Golang/tutorial/slice.png)

### Variable Shawdowing
- When there are two variable with same name 
  - outside the function after import statement (public for the package)#
  - inside the function.
- priority is given to the variable which is inside the function 