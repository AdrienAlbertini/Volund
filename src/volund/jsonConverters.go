package main

import (
	"os"
)

type ExecutableType struct {
	folderInfos      VolundBuildFolder
	targetName       string
	srcExtension     string
	staticLibsDeps   []string
	sharedLibsDeps   []string
	excludeSrc       []string
	src              []string
	srcFileNames     []string
	headersFolders   []string
	internLibs       []string
	externIncludes   []string
	externLibs       []string
	compilerFlags    []string
	outFolder        string
	isMainExecutable bool
	isBuilt          bool
}

type StaticLibType struct {
	folderInfos    VolundBuildFolder
	targetName     string
	srcExtension   string
	staticLibsDeps []string
	sharedLibsDeps []string
	excludeSrc     []string
	src            []string
	srcFileNames   []string
	headersFolders []string
	internLibs     []string
	externIncludes []string
	externLibs     []string
	compilerFlags  []string
	outFolder      string
	isBuilt        bool
}

type SharedLibType struct {
	folderInfos    VolundBuildFolder
	targetName     string
	srcExtension   string
	staticLibsDeps []string
	sharedLibsDeps []string
	excludeSrc     []string
	src            []string
	srcFileNames   []string
	headersFolders []string
	internLibs     []string
	externIncludes []string
	externLibs     []string
	compilerFlags  []string
	outFolder      string
	isBuilt        bool
}

type BuilderType struct {
	os                   VolundOSType
	mainExecutable       ExecutableType
	executables          []ExecutableType
	staticLibs           []StaticLibType
	sharedLibs           []SharedLibType
	externIncludes       []string
	externLibs           []string
	mainSharedLibsFolder string
	mainFolder           string
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

	resolvedJson.TargetName = returnDefaultIfEmpty(resolveJsonObj.TargetName, jsonObj.TargetName)
	resolvedJson.SrcExtension = returnDefaultIfEmpty(resolveJsonObj.SrcExtension, jsonObj.SrcExtension)
	resolvedJson.OutFolder = returnDefaultIfEmpty(resolveJsonObj.OutFolder, jsonObj.OutFolder)
	resolvedJson.StaticLibsDeps = append(jsonObj.StaticLibsDeps, resolveJsonObj.StaticLibsDeps...)
	resolvedJson.SharedLibsDeps = append(jsonObj.SharedLibsDeps, resolveJsonObj.SharedLibsDeps...)
	resolvedJson.ExcludeSrc = append(jsonObj.ExcludeSrc, resolveJsonObj.ExcludeSrc...)
	resolvedJson.SrcFolders = append(jsonObj.SrcFolders, resolveJsonObj.SrcFolders...)
	resolvedJson.HeadersFolders = append(jsonObj.HeadersFolders, resolveJsonObj.HeadersFolders...)
	resolvedJson.InternLibs = append(jsonObj.InternLibs, resolveJsonObj.InternLibs...)
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

	resolvedJson.Compiler = returnDefaultIfEmpty(resolveJsonObj.Compiler, jsonObj.Compiler)
	resolvedJson.MainExecutable = returnDefaultIfEmpty(resolveJsonObj.MainExecutable, jsonObj.MainExecutable)
	resolvedJson.MainFolder = returnDefaultIfEmpty(resolveJsonObj.MainFolder, jsonObj.MainFolder)
	resolvedJson.MainSharedLibsFolder = returnDefaultIfEmpty(resolveJsonObj.MainSharedLibsFolder, jsonObj.MainSharedLibsFolder)
	resolvedJson.Executables = append(jsonObj.Executables, resolveJsonObj.Executables...)
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
	staticLib.targetName = commonBuildObj.TargetName
	staticLib.outFolder = staticLib.folderInfos.path + "/" + commonBuildObj.OutFolder
	staticLib.staticLibsDeps = commonBuildObj.StaticLibsDeps
	staticLib.headersFolders = commonBuildObj.HeadersFolders
	staticLib.srcExtension = commonBuildObj.SrcExtension
	staticLib.externIncludes = commonBuildObj.ExternIncludes
	staticLib.externLibs = commonBuildObj.ExternLibs
	staticLib.compilerFlags = commonBuildObj.CompilerFlags
	staticLib.excludeSrc = commonBuildObj.ExcludeSrc
	staticLib.isBuilt = false

	success, _ := exists(staticLib.outFolder)
	if !success {
		os.MkdirAll(staticLib.outFolder, os.ModePerm)
	}

	staticLib.srcFileNames, staticLib.src = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return staticLib
}

func makeSharedLibType(folderInfos VolundBuildFolder) *SharedLibType {
	sharedLib := new(SharedLibType)
	jsonObj := getFileJSONObj(folderInfos)

	commonBuildObj := resolveOSParams(jsonObj.SharedLib)
	sharedLib.folderInfos = folderInfos
	sharedLib.targetName = commonBuildObj.TargetName
	sharedLib.outFolder = sharedLib.folderInfos.path + "/" + commonBuildObj.OutFolder
	sharedLib.staticLibsDeps = commonBuildObj.StaticLibsDeps
	sharedLib.headersFolders = commonBuildObj.HeadersFolders
	sharedLib.srcExtension = commonBuildObj.SrcExtension
	sharedLib.externIncludes = commonBuildObj.ExternIncludes
	sharedLib.externLibs = commonBuildObj.ExternLibs
	sharedLib.compilerFlags = commonBuildObj.CompilerFlags
	sharedLib.excludeSrc = commonBuildObj.ExcludeSrc
	sharedLib.isBuilt = false

	success, _ := exists(sharedLib.outFolder)
	if !success {
		os.MkdirAll(sharedLib.outFolder, os.ModePerm)
	}

	sharedLib.srcFileNames, sharedLib.src = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return sharedLib
}

func makeExecutableType(folderInfos VolundBuildFolder, mainExecutable string) *ExecutableType {
	executable := new(ExecutableType)
	jsonObj := getFileJSONObj(folderInfos)

	commonBuildObj := resolveOSParams(jsonObj.Executable)
	executable.folderInfos = folderInfos
	executable.targetName = commonBuildObj.TargetName
	executable.isMainExecutable = mainExecutable == executable.targetName
	executable.outFolder = executable.folderInfos.path + "/" + commonBuildObj.OutFolder
	executable.staticLibsDeps = commonBuildObj.StaticLibsDeps
	executable.sharedLibsDeps = commonBuildObj.SharedLibsDeps
	executable.headersFolders = commonBuildObj.HeadersFolders
	executable.srcExtension = commonBuildObj.SrcExtension
	executable.externIncludes = commonBuildObj.ExternIncludes
	executable.externLibs = commonBuildObj.ExternLibs
	executable.compilerFlags = commonBuildObj.CompilerFlags
	executable.excludeSrc = commonBuildObj.ExcludeSrc
	executable.isBuilt = false

	success, _ := exists(executable.outFolder)
	if !success {
		os.MkdirAll(executable.outFolder, os.ModePerm)
	}

	executable.srcFileNames, executable.src = getSourceFiles(commonBuildObj.SrcFolders, commonBuildObj.SrcExtension, folderInfos)

	return executable
}
