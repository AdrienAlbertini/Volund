package main

import (
	"reflect"
)

type ObjectJSON struct {
	Builder    BuilderJSON     `json:"Builder"`
	Executable CommonBuildJSON `json:"Executable"`
	StaticLib  CommonBuildJSON `json:"StaticLib"`
	SharedLib  CommonBuildJSON `json:"SharedLib"`
}

type BuilderJSON struct {
	Os                   string                `json:"OS"`
	Compiler             string                `json:"compiler"`
	MainExecutable       string                `json:"mainExecutable"`
	MainFolder           string                `json:"mainFolder"`
	MainSharedLibsFolder string                `json:"mainSharedLibsFolder"`
	Executables          []string              `json:"executables"`
	StaticLibs           []string              `json:"staticLibs"`
	SharedLibs           []string              `json:"sharedLibs"`
	ExternIncludes       []string              `json:"externIncludes"` // -I (absolute)
	ExternLibs           []string              `json:"externLibs"`     // -L & -l (absolute)
	CompilerFlags        []string              `json:"compilerFlags"`
	FullStatic           bool                  `json:"fullStatic"`
	Windows              BuilderOSSpecificJSON `json:"Windows"`
	Linux                BuilderOSSpecificJSON `json:"Linux"`
	OSX                  BuilderOSSpecificJSON `json:"OSX"`
	PS4                  BuilderOSSpecificJSON `json:"PS4"`
}

type BuilderOSSpecificJSON struct {
	Compiler             string   `json:"compiler"`
	MainExecutable       string   `json:"mainExecutable"`
	MainFolder           string   `json:"mainFolder"`
	MainSharedLibsFolder string   `json:"mainSharedLibsFolder"`
	Executables          []string `json:"executables"`
	StaticLibs           []string `json:"staticLibs"`
	SharedLibs           []string `json:"sharedLibs"`
	ExternIncludes       []string `json:"externIncludes"` // -I (absolute)
	ExternLibs           []string `json:"externLibs"`     // -L & -l (absolute)
	CompilerFlags        []string `json:"compilerFlags"`
	FullStatic           bool     `json:"fullStatic"`
}

type CommonBuildJSON struct {
	TargetName     string               `json:"targetName"`
	SrcExtension   string               `json:"srcExtension"`
	OutFolder      string               `json:"outFolder"`
	StaticLibsDeps []string             `json:"staticLibsDeps"`
	SharedLibsDeps []string             `json:"sharedLibsDeps"`
	ExcludeSrc     []string             `json:"excludeSrc"`
	SrcFolders     []string             `json:"srcFolders"`
	HeadersFolders []string             `json:"headersFolders"` // -I (relative)
	InternLibs     []string             `json:"internLibs"`     //  -L & -l (relative)
	ExternIncludes []string             `json:"externIncludes"` // -I (absolute)
	ExternLibs     []string             `json:"externLibs"`     // -L & -l (absolute)
	CompilerFlags  []string             `json:"compilerFlags"`
	Windows        OSSpecificParamsJSON `json:"Windows"`
	Linux          OSSpecificParamsJSON `json:"Linux"`
	OSX            OSSpecificParamsJSON `json:"OSX"`
	PS4            OSSpecificParamsJSON `json:"PS4"`
}

type OSSpecificParamsJSON struct {
	TargetName     string   `json:"targetName"`
	SrcExtension   string   `json:"srcExtension"`
	OutFolder      string   `json:"outFolder"`
	StaticLibsDeps []string `json:"staticLibsDeps"`
	SharedLibsDeps []string `json:"sharedLibsDeps"`
	ExcludeSrc     []string `json:"excludeSrc"`
	SrcFolders     []string `json:"srcFolders"`
	HeadersFolders []string `json:"headersFolders"` // -I (relative)
	InternLibs     []string `json:"internLibs"`     //  -L & -l (relative)
	ExternIncludes []string `json:"externIncludes"` // -I (absolute)
	ExternLibs     []string `json:"externLibs"`     // -L & -l (absolute)
	CompilerFlags  []string `json:"compilerFlags"`
}

func (s ObjectJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, ObjectJSON{})
}

func (s BuilderJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, BuilderJSON{})
}

func (s BuilderOSSpecificJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, BuilderOSSpecificJSON{})
}

func (s CommonBuildJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, CommonBuildJSON{})
}

func (s OSSpecificParamsJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, OSSpecificParamsJSON{})
}
