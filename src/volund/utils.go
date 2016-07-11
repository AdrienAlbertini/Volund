package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type VolundBuildFolder struct {
	buildType       VolundBuildType
	path            string
	name            string
	volundBuildFile []byte
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
	fmt.Printf("Copy: %s to: %s\n", src, dst)
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(err)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getExternIncludesArgs(externIncludes []string) (args []string) {
	if externIncludes != nil {
		for _, externInclude := range externIncludes {
			args = append(args, "-I"+externInclude)
		}
	}

	return
}

func executeCommandWithPrintErr(command string, args []string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Err:" + fmt.Sprint(err) + ": " + string(output))
		return
	} else {
		fmt.Println("Out:" + string(output))
	}
}

func getLibsArgs(libs []string) (args []string) {
	if libs != nil {
		for _, lib := range libs {
			args = append(args, "-l"+lib)
		}
	}

	return
}

func getStaticLibByName(staticLibName string, allLibs []*StaticLibType) (bool, *StaticLibType) {
	for _, staticLib := range allLibs {
		if staticLib.targetName == staticLibName {
			return true, staticLib
		}
	}
	return false, nil
}

func getSourceFiles(srcFolders []string, extension string, folderInfos VolundBuildFolder) (sourceFiles []string, sourceFilesPath []string) {
	var sourceFilesInFolder []string

	for _, srcFolder := range srcFolders {
		if srcFolder == "." {
			sourceFilesInFolder = append(sourceFilesInFolder, getAllFilesFromDir(folderInfos.path+"/", extension)...)
		} else {
			sourceFilesInFolder = append(sourceFilesInFolder, getAllFilesFromDir(folderInfos.path+"/"+srcFolder, extension)...)
		}
		sourceFiles = append(sourceFiles, sourceFilesInFolder...)
		//		fmt.Printf("SourceFiles: %v\n", sourceFiles)

		if srcFolder == "." {
			sourceFilesPath = append(sourceFilesPath, sourceFilesInFolder...)
		} else {
			sourceFilesPath = append(sourceFilesPath, joinAtBegin(srcFolder+"/", sourceFilesInFolder)...)
		}
		//		fmt.Printf("SourceFilesPath: %v\n", sourceFilesPath)
		//		fmt.Printf("SourceFiles: %v\n", sourceFiles)
	}

	return
}

func getSrcAndHeadersFolderPath(folderInfos VolundBuildFolder, srcFolder string, headersFolder string) (srcPath string, headersPath string) {
	srcPath = folderInfos.path + "/" + srcFolder
	headersPath = folderInfos.path + "/" + headersFolder
	return
}

func getAllFilesFromDir(folderPath string, extension string) (finalFiles []string) {
	files, err := ioutil.ReadDir(folderPath)

	if err != nil {
		log.Fatal(err)

		return
	}

	for _, file := range files {
		filename := file.Name()

		if strings.Contains(filename, extension) {
			finalFiles = append(finalFiles, filename)
		}
	}

	return
}

func getFileJSONObj(folder VolundBuildFolder) ObjectJSON {
	var subFolderVolundJSON ObjectJSON
	json.Unmarshal(folder.volundBuildFile, &subFolderVolundJSON)
	return subFolderVolundJSON
}

func getStaticLibOSExtension() string {
	if osType == WINDOWS {
		return WINDOWS_STATIC_EXT
	}
	return LINUX_STATIC_EXT
}

func getSharedLibOsExtension() string {
	if osType == WINDOWS {
		return WINDOWS_DYNAMIC_EXT
	} else if osType == OSX {
		return OSX_DYNAMIC_EXT
	}
	return LINUX_DYNAMIC_EXT
}

func getExecutableOSExtension() string {
	if osType == WINDOWS {
		return WINDOWS_EXECUTABLE_EXT
	} else if osType == OSX {
		return OSX_EXECUTABLE_EXT
	}
	return LINUX_EXECUTABLE_EXT
}

func getStaticLibsLinks(libsToLink []string, libs []*StaticLibType, avoidLib string) (linkPaths []string, linkNames []string,
	linkIncludes []string) {

	//	fmt.Printf("GetStaticLibsLinks LibsToLink: %v\n", libsToLink)
	////	fmt.Printf("GetStaticLibsLinks AvoidLib: %s\n", avoidLib)

	for _, staticLib := range libs {

		//		fmt.Printf("GetStaticLibsLinks Libs: %s\n", staticLib.name)
		if (staticLib.targetName != avoidLib) && (contains(libsToLink, staticLib.targetName)) {
			path := "-L" + staticLib.outFolder
			name := staticLib.outFolder + "/" + staticLib.targetName + getStaticLibOSExtension()

			linkIncludes = append(linkIncludes, "-I"+"./"+staticLib.targetName+"/.")
			for _, includeHeader := range staticLib.headersFolders {
				linkIncludes = append(linkIncludes, "-I"+"./"+staticLib.targetName+"/"+includeHeader)
			}

			linkPaths = append(linkPaths, path)
			linkNames = append(linkNames, name)

			//	fmt.Printf("GetStaticLibsLinks Names: %v\n", linkNames)
		}
	}
	//	fmt.Printf("GetStaticLibsLinks LinkPaths: %v\n", linkPaths)
	//	fmt.Printf("GetStaticLibsLinks linkNames: %v\n", linkNames)
	//	fmt.Printf("GetStaticLibsLinks linkIncludes: %v\n", linkIncludes)

	return
}

func getExternalDependencies(externLibs []string, externIncludes []string) (linkLibs []string, linkIncludes []string) {

	for _, externalLib := range externLibs {
		linkLibs = append(linkLibs, "-l"+externalLib)
	}
	for _, externalInclude := range externIncludes {
		linkIncludes = append(linkIncludes, "-I"+externalInclude)
	}

	return
}

func isValidToolchain(testToolchain string) bool {
	validToolchains := []string{"clang", "", "gcc", "g++"}

	for _, toolchain := range validToolchains {
		if testToolchain == toolchain {
			return true
		}
	}

	return false
}

func resolveBuildType(builderJSON *BuilderJSON, jsonOBJ *ObjectJSON, buildFolder *VolundBuildFolder, executables *[]string,
	staticLibs *[]string, sharedLibs *[]string) bool {

	if jsonOBJ.Executable.TargetName != "" && contains(*executables, jsonOBJ.Executable.TargetName) == false && builderJSON.MainExecutable == jsonOBJ.Executable.TargetName {
		buildFolder.buildType = EXECUTABLE
	} else if jsonOBJ.Executable.TargetName != "" && contains(*executables, jsonOBJ.Executable.TargetName) {
		buildFolder.buildType = EXECUTABLE
	} else if jsonOBJ.SharedLib.TargetName != "" && contains(*sharedLibs, jsonOBJ.SharedLib.TargetName) {
		buildFolder.buildType = SHARED_LIB
	} else if jsonOBJ.StaticLib.TargetName != "" && contains(*staticLibs, jsonOBJ.StaticLib.TargetName) {
		buildFolder.buildType = STATIC_LIB
	} else {
		buildFolder.buildType = NONE
		return false
	}
	return true
}

func getRuntimeOS() VolundOSType {
	switch runtime.GOOS {
	case "windows":
		return WINDOWS
	case "linux":
		return LINUX
	case "darwin":
		return OSX
	}
	return UNKNOWN
}

func getOsType(osStr string) VolundOSType {
	switch osStr {
	case "Auto":
		return getRuntimeOS()
	case "Windows":
		return WINDOWS
	case "Linux":
		return LINUX
	case "OSX":
		return OSX
	}
	return UNKNOWN
}

func (osConst VolundOSType) ToString() string {
	switch osConst {
	case WINDOWS:
		return "Windows"
	case LINUX:
		return "Linux"
	case OSX:
		return "OSX"
	}
	return "Unknown OS"
}

func returnDefaultIfEmpty(toCheck string, defaultStr string) string {
	if toCheck == "" {
		return defaultStr
	}
	return toCheck
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func joinAtBegin(toJoin string, strs []string) (finalStrs []string) {
	for _, str := range strs {
		finalStrs = append(finalStrs, toJoin+str)
	}

	return
}
