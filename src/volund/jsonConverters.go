package main

import (
	"os"
)

type BinaryType struct {
	folderInfos     VolundBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	excludeSrc      []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	compilerFlags   []string
	outFolder       string
	isOutBinary     bool
	isBuilt         bool
}

type StaticLibType struct {
	folderInfos     VolundBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	excludeSrc      []string
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
	folderInfos     VolundBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	sharedLibs      []string
	excludeSrc      []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	compilerFlags   []string
	outFolder       string
	isBuilt         bool
}

type BuilderType struct {
	os               VolundOSType
	outBinary        BinaryType
	binaries         []BinaryType
	staticLibs       []StaticLibType
	sharedLibs       []SharedLibType
	sharedLibsFolder string
	outFolder        string
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
	resolvedJson.ExcludeSrc = append(jsonObj.ExcludeSrc, resolveJsonObj.ExcludeSrc...)
	resolvedJson.SrcFolders = append(jsonObj.SrcFolders, resolveJsonObj.SrcFolders...)
	resolvedJson.HeadersFolders = append(jsonObj.HeadersFolders, resolveJsonObj.HeadersFolders...)
	resolvedJson.ExternIncludes = append(jsonObj.ExternIncludes, resolveJsonObj.ExternIncludes...)
	resolvedJson.ExternLibs = append(jsonObj.ExternLibs, resolveJsonObj.ExternLibs...)
	resolvedJson.CompilerFlags = append(jsonObj.CompilerFlags, resolveJsonObj.CompilerFlags...)
	return
}

func resolveBuilderOSParams(jsonObj BuilderJSON) (resolvedJson BuilderJSON) {
	var resolveJsonObj BuilderOSSpecificJSON

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

	resolvedJson.Toolchain = returnDefaultIfEmpty(resolveJsonObj.Toolchain, jsonObj.Toolchain)
	resolvedJson.OutBinary = returnDefaultIfEmpty(resolveJsonObj.OutBinary, jsonObj.OutBinary)
	resolvedJson.SharedLibsFolder = returnDefaultIfEmpty(resolveJsonObj.SharedLibsFolder, jsonObj.SharedLibsFolder)
	resolvedJson.OutFolder = returnDefaultIfEmpty(resolveJsonObj.OutFolder, jsonObj.OutFolder)
	resolvedJson.Binaries = append(jsonObj.Binaries, resolveJsonObj.Binaries...)
	resolvedJson.StaticLibs = append(jsonObj.StaticLibs, resolveJsonObj.StaticLibs...)
	resolvedJson.SharedLibs = append(jsonObj.SharedLibs, resolveJsonObj.SharedLibs...)
	resolvedJson.CompilerFlags = append(jsonObj.CompilerFlags, resolveJsonObj.CompilerFlags...)
	resolvedJson.FullStatic = jsonObj.FullStatic
	return
}

func makeStaticLibType(folderInfos VolundBuildFolder) *StaticLibType {
	staticLib := new(StaticLibType)
	jsonObj := getFileJSONObj(folderInfos)

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
	staticLib.excludeSrc = commonBuildObj.ExcludeSrc
	staticLib.isBuilt = false

	success, _ := exists(staticLib.outFolder)
	if !success {
		os.MkdirAll(staticLib.outFolder, os.ModePerm)
	}

	staticLib.sourceFileNames, staticLib.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return staticLib
}

func makeSharedLibType(folderInfos VolundBuildFolder) *SharedLibType {
	sharedLib := new(SharedLibType)
	jsonObj := getFileJSONObj(folderInfos)

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
	sharedLib.excludeSrc = commonBuildObj.ExcludeSrc
	sharedLib.isBuilt = false

	success, _ := exists(sharedLib.outFolder)
	if !success {
		os.MkdirAll(sharedLib.outFolder, os.ModePerm)
	}

	sharedLib.sourceFileNames, sharedLib.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return sharedLib
}

func makeBinaryType(folderInfos VolundBuildFolder, outBinary string) *BinaryType {
	binary := new(BinaryType)
	jsonObj := getFileJSONObj(folderInfos)

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
	binary.excludeSrc = commonBuildObj.ExcludeSrc
	binary.isBuilt = false

	success, _ := exists(binary.outFolder)
	if !success {
		os.MkdirAll(binary.outFolder, os.ModePerm)
	}

	binary.sourceFileNames, binary.sources = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return binary
}
