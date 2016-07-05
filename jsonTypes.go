package main

type ObjectJSON struct {
	Builder   BuilderJSON     `json:"Builder"`
	Binary    CommonBuildJSON `json:"Binary"`
	StaticLib CommonBuildJSON `json:"StaticLib"`
	SharedLib CommonBuildJSON `json:"SharedLib"`
}

type BuilderJSON struct {
	Os               string   `json:"OS"`
	Toolchain        string   `json:"toolchain"`
	OutBinary        string   `json:"outBinary"`
	Binaries         []string `json:"binaries"`
	StaticLibs       []string `json:"staticLibs"`
	SharedLibs       []string `json:"sharedLibs"`
	SharedLibsFolder string   `json:"outsharedLibsFolder"`
	OutFolder        string   `json:"outFolder"`
	CompilerFlags    []string `json:"compilerFlags"`
	FullStatic       bool     `json:"fullStatic"`
}

type CommonBuildJSON struct {
	Name           string               `json:"name"`
	SrcExtension   string               `json:"srcExtension"`
	OutFolder      string               `json:"outFolder"`
	StaticLibs     []string             `json:"staticLibs"`
	SharedLibs     []string             `json:"sharedLibs"`
	SrcFolders     []string             `json:"srcFolders"`
	HeadersFolders []string             `json:"headersFolders"`
	ExternIncludes []string             `json:"externIncludes"`
	ExternLibs     []string             `json:"externLibs"`
	CompilerFlags  []string             `json:"compilerFlags"`
	Windows        OSSpecificParamsJSON `json:"Windows"`
	Linux          OSSpecificParamsJSON `json:"Linux"`
	OSX            OSSpecificParamsJSON `json:"OSX"`
}

type OSSpecificParamsJSON struct {
	Name           string   `json:"name"`
	SrcExtension   string   `json:"srcExtension"`
	OutFolder      string   `json:"outFolder"`
	StaticLibs     []string `json:"staticLibs"`
	SharedLibs     []string `json:"sharedLibs"`
	SrcFolders     []string `json:"srcFolders"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	CompilerFlags  []string `json:"compilerFlags"`
}

/*
type BinaryJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	sharedLibs        []string `json:"sharedLibs"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}

type StaticLibJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}

type sharedLibJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}
*/
