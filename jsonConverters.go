package main

import (
	"os"
	"reflect"
)

type BinaryType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	compilerFlags   []string
	outFolder       string
	isOutBinary     bool
}

type StaticLibType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	compilerFlags   []string
	outFolder       string
	isBuilt         bool
}

type SharedLibType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	compilerFlags   []string
	outFolder       string
}

type BuilderType struct {
	os               ObakeOSType
	outBinary        BinaryType
	binaries         []BinaryType
	staticLibs       []StaticLibType
	sharedLibs       []SharedLibType
	sharedLibsFolder string
	outFolder        string
}

func (s OSSpecificParamsJSON) IsEmpty() bool {
	return reflect.DeepEqual(s, OSSpecificParamsJSON{})
}

func resolveOSParams(jsonObj CommonBuildJSON) (resolvedJson CommonBuildJSON) {
	var resolveJsonObj OSSpecificParamsJSON

	resolvedJson = jsonObj
	switch osType {
	case WINDOWS:
		if jsonObj.Windows.IsEmpty() {
			return
		}
		resolveJsonObj = jsonObj.Windows
	case LINUX:
		if jsonObj.Linux.IsEmpty() {
			return
		}
		resolveJsonObj = jsonObj.Linux
	case OSX:
		if jsonObj.OSX.IsEmpty() {
			return
		}
		resolveJsonObj = jsonObj.OSX
	}

	resolvedJson.Name = returnDefaultIfEmpty(resolveJsonObj.Name, jsonObj.Name)
	resolvedJson.SrcExtension = returnDefaultIfEmpty(resolveJsonObj.SrcExtension, jsonObj.SrcExtension)
	resolvedJson.OutFolder = returnDefaultIfEmpty(resolveJsonObj.OutFolder, jsonObj.OutFolder)
	resolvedJson.StaticLibs = append(jsonObj.StaticLibs, resolveJsonObj.StaticLibs...)
	resolvedJson.SharedLibs = append(jsonObj.SharedLibs, resolveJsonObj.SharedLibs...)
	resolvedJson.SrcFolders = append(jsonObj.SrcFolders, resolveJsonObj.SrcFolders...)
	resolvedJson.HeadersFolders = append(jsonObj.HeadersFolders, resolveJsonObj.HeadersFolders...)
	resolvedJson.ExternIncludes = append(jsonObj.ExternIncludes, resolveJsonObj.ExternIncludes...)
	resolvedJson.ExternLibs = append(jsonObj.ExternLibs, resolveJsonObj.ExternLibs...)
	resolvedJson.CompilerFlags = append(jsonObj.CompilerFlags, resolveJsonObj.CompilerFlags...)
	return
}

func makeStaticLibType(folderInfos ObakeBuildFolder) *StaticLibType {
	staticLib := new(StaticLibType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	commonBuildObj := resolveOSParams(jsonObj.StaticLib)
	staticLib.folderInfos = folderInfos
	staticLib.name = commonBuildObj.Name
	staticLib.outFolder = staticLib.folderInfos.path + "/" + commonBuildObj.OutFolder
	staticLib.staticLibs = commonBuildObj.StaticLibs
	staticLib.headerFolders = commonBuildObj.HeadersFolders
	staticLib.sourceExtension = commonBuildObj.SrcExtension
	staticLib.externIncludes = commonBuildObj.ExternIncludes
	staticLib.externLibs = commonBuildObj.ExternLibs
	staticLib.compilerFlags = commonBuildObj.CompilerFlags
	staticLib.isBuilt = false

	success, _ := exists(staticLib.outFolder)
	if !success {
		os.MkdirAll(staticLib.outFolder, os.ModePerm)
	}

	staticLib.sourceFileNames, staticLib.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return staticLib
}

func makeSharedLibType(folderInfos ObakeBuildFolder) *SharedLibType {
	sharedLib := new(SharedLibType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	commonBuildObj := resolveOSParams(jsonObj.SharedLib)
	sharedLib.folderInfos = folderInfos
	sharedLib.name = commonBuildObj.Name
	sharedLib.outFolder = sharedLib.folderInfos.path + "/" + commonBuildObj.OutFolder
	sharedLib.staticLibs = commonBuildObj.StaticLibs
	sharedLib.headerFolders = commonBuildObj.HeadersFolders
	sharedLib.sourceExtension = commonBuildObj.SrcExtension
	sharedLib.externIncludes = commonBuildObj.ExternIncludes
	sharedLib.externLibs = commonBuildObj.ExternLibs
	sharedLib.compilerFlags = commonBuildObj.CompilerFlags

	success, _ := exists(sharedLib.outFolder)
	if !success {
		os.MkdirAll(sharedLib.outFolder, os.ModePerm)
	}

	sharedLib.sourceFileNames, sharedLib.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return sharedLib
}

func makeBinaryType(folderInfos ObakeBuildFolder, outBinary string) *BinaryType {
	binary := new(BinaryType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	commonBuildObj := resolveOSParams(jsonObj.Binary)
	binary.folderInfos = folderInfos
	binary.name = commonBuildObj.Name
	binary.isOutBinary = outBinary == binary.name
	binary.outFolder = binary.folderInfos.path + "/" + commonBuildObj.OutFolder
	binary.staticLibs = commonBuildObj.StaticLibs
	binary.sharedLibs = commonBuildObj.SharedLibs
	binary.headerFolders = commonBuildObj.HeadersFolders
	binary.sourceExtension = commonBuildObj.SrcExtension
	binary.externIncludes = commonBuildObj.ExternIncludes
	binary.externLibs = commonBuildObj.ExternLibs
	binary.compilerFlags = commonBuildObj.CompilerFlags

	success, _ := exists(binary.outFolder)
	if !success {
		os.MkdirAll(binary.outFolder, os.ModePerm)
	}

	binary.sourceFileNames, binary.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return binary
}
