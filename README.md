# Volund
Volund is a c++ build system written in go. 

# Dependencies
https://github.com/fatih/color   
https://github.com/cheggaaa/pb

# Install Dependencies

```
go get github.com/fatih/color
go get gopkg.in/cheggaaa/pb.v1
```

# Description (Work in progress)

Volund is a build system like make.
We chose JSON as our file format for more flexibility.

Example:

Root/   
  |--SharedLib1/   
  |--| VolundBuild.json   
  |   
  |--StaticLib1/   
  |--| VolundBuild.json   
  |   
  |--Executable1/   
  |--| VolundBuild.json   
  |   
  |VolundBuild.json   

The root file (VolundBuild.json) represents the "Builder" type which handles the compilation of all subsystems 
(Executables, Static Libraries, Shared Libraries).   
   
Json types:   
   
```
string      OS
string      compiler
string      mainExecutable
string      mainFolder
string      mainSharedLibsFolder
[]string    executables
[]string    staticLibs
[]string    sharedLibs
[]string    externIncludes // -I (absolute)
[]string    externLibs     // -L & -l (absolute)
[]string    compilerFlags
bool        fullStatic

Executable, SharedLib, StaticLib

string   		targetName
string   		srcExtension
string   		outFolder
[]string 		staticLibsDeps
[]string 		sharedLibsDeps
[]string 		excludeSrc
[]string 		srcFolders
[]string 		headersFolders // -I (relative)
[]string 		internLibs     //  -L & -l (relative)
[]string 		externIncludes // -I (absolute)
[]string 		externLibs     // -L & -l (absolute)
[]string 		compilerFlags
```

# License 

The MIT License (MIT) see the file LICENSE in the project root.
