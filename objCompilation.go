package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type ObjFileRequirement struct {
	folderInfos     ObakeBuildFolder
	allLibs         *[]*StaticLibType
	sourceFilesPath []string
	sourceFiles     []string
	headersFolders  []string
	externIncludes  []string
	externLibs      []string
	staticLibs      []string
	compilerFlags   []string
	sourceExtension string
	outFolder       string
}

func pluginTypeToObjType(plugin PluginType, folderInfos ObakeBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.sourceFilesPath = plugin.sources
	objType.sourceFiles = plugin.sourceFileNames
	objType.headersFolders = plugin.headerFolders
	objType.externIncludes = plugin.externIncludes
	objType.externLibs = plugin.externLibs
	objType.staticLibs = plugin.staticLibs
	objType.sourceExtension = plugin.sourceExtension
	objType.outFolder = plugin.outFolder
	objType.compilerFlags = plugin.compilerFlags

	return
}

func staticLibTypeToObjType(staticLib StaticLibType, folderInfos ObakeBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.sourceFilesPath = staticLib.sources
	objType.sourceFiles = staticLib.sourceFileNames
	objType.headersFolders = staticLib.headerFolders
	objType.externIncludes = staticLib.externIncludes
	objType.externLibs = staticLib.externLibs
	objType.staticLibs = staticLib.staticLibs
	objType.sourceExtension = staticLib.sourceExtension
	objType.outFolder = staticLib.outFolder
	objType.compilerFlags = staticLib.compilerFlags

	return
}

func binaryTypeToObjType(binary BinaryType, folderInfos ObakeBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.sourceFilesPath = binary.sources
	objType.sourceFiles = binary.sourceFileNames
	objType.headersFolders = binary.headerFolders
	objType.externIncludes = binary.externIncludes
	objType.externLibs = binary.externLibs
	objType.staticLibs = binary.staticLibs
	objType.sourceExtension = binary.sourceExtension
	objType.outFolder = binary.outFolder
	objType.compilerFlags = binary.compilerFlags

	return
}

func buildAndGetObjectFiles(objType ObjFileRequirement, success *bool, objectFilesPath *[]string) {
	for i, srcFilePath := range objType.sourceFilesPath {
		oFilePath := objType.outFolder + "/" + strings.Replace(objType.sourceFiles[i], objType.sourceExtension, ".o", -1)
		*objectFilesPath = append(*objectFilesPath, oFilePath)

		//		fmt.Printf("SrcFilePath: %s\n", srcFilePath)
		args := []string{"-c", objType.folderInfos.path + "/" + srcFilePath, "-o", oFilePath}
		args = append(args, compilerFlags...)

		for _, headerFolder := range objType.headersFolders {
			if headerFolder == "." {
				args = append(args, "-I"+objType.folderInfos.path+"/")
			} else {
				args = append(args, "-I"+objType.folderInfos.path+"/"+headerFolder+"/")
			}
		}

		_, linkNames, linkIncludes := getStaticLibsLinks(objType.staticLibs, *objType.allLibs, objType.folderInfos.name)

		args = append(args, linkIncludes...)
		args = append(args, linkNames...)

		args = append(args, getExternIncludesArgs(objType.externIncludes)...)
		args = append(args, getExternLibsArgs(objType.externLibs)...)
		args = append(args, objType.compilerFlags...)

		boldCyan.Printf("[%d/%d] ObjFile: ", (i + 1), len(objType.sourceFiles))
		boldBlue.Printf("%s ", toolchain)
		fmt.Printf("%v\n", args)
		cmd := exec.Command(toolchain, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			boldRed.Printf("ObjFile: %s | Error: %s:\n\n", srcFilePath, fmt.Sprint(err))
			fmt.Printf("%s\n", string(out))
			//fmt.Printf("ObjFile: %s | Error: %s\n", srcFilePath, fmt.Sprint(err))
			*success = false
			return
		}
		/*
			fmt.Printf("Obj files: %s %v\n", toolchain, args)
			_, err := exec.Command(toolchain, args...).Output()
			if err != nil {
				fmt.Printf("ObjFile: %s | Error: %s\n", srcFilePath, err)
			}
		*/
	}
	*success = true
}
