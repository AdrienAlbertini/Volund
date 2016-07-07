package main

import (
	"reflect"
)

type ObjectJSON struct {
	Builder   BuilderJSON     `json:"Builder"`
	Binary    CommonBuildJSON `json:"Binary"`
	StaticLib CommonBuildJSON `json:"StaticLib"`
	SharedLib CommonBuildJSON `json:"SharedLib"`
}

type BuilderJSON struct {
	Os               string                `json:"OS"`
	Toolchain        string                `json:"toolchain"`
	OutBinary        string                `json:"outBinary"`
	SharedLibsFolder string                `json:"outsharedLibsFolder"`
	OutFolder        string                `json:"outFolder"`
	Binaries         []string              `json:"binaries"`
	StaticLibs       []string              `json:"staticLibs"`
	SharedLibs       []string              `json:"sharedLibs"`
	CompilerFlags    []string              `json:"compilerFlags"`
	FullStatic       bool                  `json:"fullStatic"`
	Windows          BuilderOSSpecificJSON `json:"Windows"`
	Linux            BuilderOSSpecificJSON `json:"Linux"`
	OSX              BuilderOSSpecificJSON `json:"OSX"`
}

type BuilderOSSpecificJSON struct {
	Toolchain        string   `json:"toolchain"`
	OutBinary        string   `json:"outBinary"`
	SharedLibsFolder string   `json:"outsharedLibsFolder"`
	OutFolder        string   `json:"outFolder"`
	Binaries         []string `json:"binaries"`
	StaticLibs       []string `json:"staticLibs"`
	SharedLibs       []string `json:"sharedLibs"`
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
