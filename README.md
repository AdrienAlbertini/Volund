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
  |  | VolundBuild.json
  |
  |--StaticLib1/
  |  | VolundBuild.json
  |
  |--Executable1/
  |  | VolundBuild.json
  |
  | VolundBuild.json

The root file (VolundBuild.json) represents the "Builder" type which handles the compilation of all subsystems 
(Executables, Static Libraries, Shared Libraries).

Json types:

```
Builder

string                `json:"OS"` 
string                `json:"compiler"`
string                `json:"mainExecutable"`
string                `json:"mainFolder"`
string                `json:"mainSharedLibsFolder"`
[]string              `json:"executables"`
[]string              `json:"staticLibs"`
[]string              `json:"sharedLibs"`
[]string              `json:"externIncludes"` // -I (absolute)
[]string              `json:"externLibs"`     // -L & -l (absolute)
[]string              `json:"compilerFlags"`
bool                  `json:"fullStatic"`

Executable, SharedLib, StaticLib

string   `json:"targetName"`
string   `json:"srcExtension"`
string   `json:"outFolder"`
[]string `json:"staticLibsDeps"`
[]string `json:"sharedLibsDeps"`
[]string `json:"excludeSrc"`
[]string `json:"srcFolders"`
[]string `json:"headersFolders"` // -I (relative)
[]string `json:"internLibs"`     //  -L & -l (relative)
[]string `json:"externIncludes"` // -I (absolute)
[]string `json:"externLibs"`     // -L & -l (absolute)
[]string `json:"compilerFlags"`
```

# License 

The MIT License (MIT) see the file LICENSE in the project root.
