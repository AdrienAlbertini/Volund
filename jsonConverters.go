package main

import (
	"os"
)

type BinaryType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	plugins         []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	outFolder       string
	isOutBinary     bool
}

type StaticLibType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	plugins         []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	outFolder       string
	isBuilt         bool
}

type PluginType struct {
	folderInfos     ObakeBuildFolder
	name            string
	sourceExtension string
	staticLibs      []string
	plugins         []string
	sources         []string
	sourceFileNames []string
	headerFolders   []string
	externIncludes  []string
	externLibs      []string
	outFolder       string
}

type BuilderType struct {
	os            ObakeOSType
	outBinary     BinaryType
	binaries      []BinaryType
	staticLibs    []StaticLibType
	plugins       []PluginType
	pluginsFolder string
	outFolder     string
}

func makeStaticLibType(folderInfos ObakeBuildFolder) *StaticLibType {
	staticLib := new(StaticLibType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	staticLib.folderInfos = folderInfos
	staticLib.name = jsonObj.StaticLib.Name
	staticLib.outFolder = staticLib.folderInfos.path + "/" + jsonObj.StaticLib.OutFolder
	staticLib.staticLibs = jsonObj.StaticLib.StaticLibs
	staticLib.headerFolders = jsonObj.StaticLib.HeadersFolders
	staticLib.sourceExtension = jsonObj.StaticLib.SrcExtension
	staticLib.externIncludes = jsonObj.StaticLib.ExternIncludes
	staticLib.externLibs = jsonObj.StaticLib.ExternLibs
	staticLib.isBuilt = false

	success, _ := exists(staticLib.outFolder)
	if !success {
		os.MkdirAll(staticLib.outFolder, os.ModePerm)
	}

	staticLib.sourceFileNames, staticLib.sources = getSourceFiles(jsonObj.StaticLib.SrcFolders, jsonObj.StaticLib.SrcExtension, folderInfos)

	return staticLib
}

func makePluginType(folderInfos ObakeBuildFolder) *PluginType {
	plugin := new(PluginType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	plugin.folderInfos = folderInfos
	plugin.name = jsonObj.Plugin.Name
	plugin.outFolder = plugin.folderInfos.path + "/" + jsonObj.Plugin.OutFolder
	plugin.staticLibs = jsonObj.Plugin.StaticLibs
	plugin.headerFolders = jsonObj.Plugin.HeadersFolders
	plugin.sourceExtension = jsonObj.Plugin.SrcExtension
	plugin.externIncludes = jsonObj.Plugin.ExternIncludes
	plugin.externLibs = jsonObj.Plugin.ExternLibs

	success, _ := exists(plugin.outFolder)
	if !success {
		os.MkdirAll(plugin.outFolder, os.ModePerm)
	}

	plugin.sourceFileNames, plugin.sources = getSourceFiles(jsonObj.Plugin.SrcFolders, jsonObj.Plugin.SrcExtension, folderInfos)

	return plugin
}

func makeBinaryType(folderInfos ObakeBuildFolder, outBinary string) *BinaryType {
	binary := new(BinaryType)
	jsonObj := getBuildFileJSONObj(folderInfos)

	binary.folderInfos = folderInfos
	binary.name = jsonObj.Binary.Name
	binary.isOutBinary = outBinary == binary.name
	binary.outFolder = binary.folderInfos.path + "/" + jsonObj.Binary.OutFolder
	binary.staticLibs = jsonObj.Binary.StaticLibs
	binary.headerFolders = jsonObj.Binary.HeadersFolders
	binary.sourceExtension = jsonObj.Binary.SrcExtension
	binary.externIncludes = jsonObj.Binary.ExternIncludes
	binary.externLibs = jsonObj.Binary.ExternLibs

	success, _ := exists(binary.outFolder)
	if !success {
		os.MkdirAll(binary.outFolder, os.ModePerm)
	}

	binary.sourceFileNames, binary.sources = getSourceFiles(jsonObj.Binary.SrcFolders, jsonObj.Binary.SrcExtension, folderInfos)

	return binary
}
