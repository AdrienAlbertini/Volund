package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func buildAndGetObjectFiles(sourceFilesPath []string, sourceFiles []string, headersFolders []string,
	externIncludes []string, externLibs []string, extension string, outFolder string, folderInfos ObakeBuildFolder,
	staticLibs []string, allLibs []*StaticLibType) (success bool, objectFilesPath []string) {
	for i, srcFilePath := range sourceFilesPath {
		oFilePath := outFolder + "/" + strings.Replace(sourceFiles[i], extension, ".o", -1)
		objectFilesPath = append(objectFilesPath, oFilePath)

		args := []string{"-c", folderInfos.path + "/" + srcFilePath, "-o", oFilePath}
		args = append(args, compilerFlags...)

		for _, headerFolder := range headersFolders {
			if headerFolder == "." {
				args = append(args, "-I"+folderInfos.path+"/")
			} else {
				args = append(args, "-I"+folderInfos.path+"/"+headerFolder+"/")
			}
		}

		_, linkNames, linkIncludes := getStaticLibsLinks(staticLibs, allLibs, folderInfos.name)

		args = append(args, linkIncludes...)
		args = append(args, linkNames...)

		args = append(args, getExternIncludesArgs(externIncludes)...)
		args = append(args, getExternLibsArgs(externLibs)...)

		fmt.Printf("Obj files: %s %v\n", toolchain, args)
		cmd := exec.Command(toolchain, args...)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("ObjFile: %s | Error: %s\n", srcFilePath, fmt.Sprint(err))
			success = false
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
	success = true
	return
}

func buildDependencies(staticLibs []string, allLibs []*StaticLibType) bool {

	for _, dependencyLib := range staticLibs {
		dependencyFound, dependency := getStaticLibByName(dependencyLib, allLibs)

		if dependencyFound && dependency.isBuilt == false {
			if handleStatic(dependency, allLibs) == false {
				return false
			}
		}
	}
	return true
}

func handleBinary(binary *BinaryType, allLibs []*StaticLibType) bool {

	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(binary.staticLibs, allLibs, binary.name)

	if buildDependencies(binary.staticLibs, allLibs) == false {
		return false
	}

	objSuccess, objectFilesPath := buildAndGetObjectFiles(binary.sources, binary.sourceFileNames, binary.headerFolders,
		binary.externIncludes, binary.externLibs, binary.sourceExtension, binary.outFolder, binary.folderInfos,
		binary.staticLibs, allLibs)

	if objSuccess == false {
		fmt.Printf("Build Binary: %s FAILED\n", binary.name)
		return false
	}

	binaryExtension := getBinaryOSExtension()

	args := []string{"-o", binary.outFolder + "/" + binary.name + binaryExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)

	//args = append(args, getExternIncludesArgs(binary.externIncludes)...)
	args = append(args, getExternLibsArgs(binary.externLibs)...)

	fmt.Printf("Handle Binary args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Binary: %s | Error: %s\n", binary.name, fmt.Sprint(err))
		fmt.Printf("Out: %s\n\n", out)
		return false
	}
	fmt.Printf("Out: %s\n\n", out)

	return true
}

func handleStatic(staticLib *StaticLibType, allLibs []*StaticLibType) bool {

	if staticLib.isBuilt == false {

		linkPaths, linkNames, linkIncludes := getStaticLibsLinks(staticLib.staticLibs, allLibs, staticLib.name)

		if buildDependencies(staticLib.staticLibs, allLibs) == false {
			return false
		}

		objSuccess, objectFilesPath := buildAndGetObjectFiles(staticLib.sources, staticLib.sourceFileNames, staticLib.headerFolders,
			staticLib.externIncludes, staticLib.externLibs, staticLib.sourceExtension, staticLib.outFolder,
			staticLib.folderInfos, staticLib.staticLibs, allLibs)

		if objSuccess == false {
			fmt.Printf("Build StaticLib: %s FAILED\n", staticLib.name)
			return false
		}

		staticLibExtension := getStaticLibOSExtension()

		args := []string{"rcs", staticLib.outFolder + "/" + staticLib.name + staticLibExtension}

		args = append(args, objectFilesPath...)

		args = append(args, linkIncludes...)
		args = append(args, linkPaths...)
		args = append(args, linkNames...)

		//	args = append(args, getExternIncludesArgs(staticLib.externIncludes)...)
		args = append(args, getExternLibsArgs(staticLib.externLibs)...)

		fmt.Printf("Handle Static args: %v\n", args)

		//	executeCommandWithPrintErr("ar", args)

		cmd := exec.Command("ar", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("StaticLib: %s | Error: %s\n", staticLib.name, fmt.Sprint(err))
			fmt.Printf("Out: %s\n\n", out)
			return false
		}
		fmt.Printf("Out: %s\n\n", out)

		/*out, err := exec.Command("ar", args...).Output()

		if err != nil {
			fmt.Printf("StaticLib: %s | Error: %s\n\n", staticLib.name, err)
			return false
		} else {
			fmt.Printf("Out: %s\n\n", out)
		}
		*/
		staticLib.isBuilt = true
	}
	return true
}

func handlePlugin(plugin *PluginType, allLibs []*StaticLibType) bool {

	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(plugin.staticLibs, allLibs, "")

	if buildDependencies(plugin.staticLibs, allLibs) == false {
		return false
	}
	objSuccess, objectFilesPath := buildAndGetObjectFiles(plugin.sources, plugin.sourceFileNames, plugin.headerFolders,
		plugin.externIncludes, plugin.externLibs, plugin.sourceExtension, plugin.outFolder, plugin.folderInfos,
		plugin.staticLibs, allLibs)

	if objSuccess == false {
		fmt.Printf("Build Plugin: %s FAILED\n", plugin.name)
		return false
	}

	pluginLibExtension := getSharedLibOsExtension()

	args := []string{"-shared", "-o", plugin.outFolder + "/" + plugin.name + pluginLibExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)

	//	args = append(args, getExternIncludesArgs(plugin.externIncludes)...)
	args = append(args, getExternLibsArgs(plugin.externLibs)...)

	fmt.Printf("Handle Plugin args: %v\n", args)

	executeCommandWithPrintErr(toolchain, args)
	/*
		out, err := exec.Command(toolchain, args...).Output()

		if err != nil {
			fmt.Printf("PluginLib: %s | Error: %s\n\n", plugin.name, err)
			return false
		}
		fmt.Printf("Out: %s\n\n", out)
	*/
	return true
}

func handleBuilder(builder BuilderJSON, binaries []*BinaryType,
	staticLibs []*StaticLibType, plugins []*PluginType) bool {
	outBinaryFound := false
	var outBinary *BinaryType

	for _, binary := range binaries {
		if binary.isOutBinary == true {
			outBinary = binary
			outBinaryFound = true
			break
		}
	}

	if outBinaryFound == false {
		return false
	}

	binaryExtension := getBinaryOSExtension()
	staticExtension := getStaticLibOSExtension()
	pluginExtension := getSharedLibOsExtension()

	success, _ := exists(builder.OutFolder)
	if !success {
		os.MkdirAll(builder.OutFolder, os.ModePerm)
	}
	success, _ = exists(builder.OutFolder + "Plugins")
	if !success {
		os.MkdirAll(builder.OutFolder+"/Plugins", os.ModePerm)
	}

	copy(outBinary.outFolder+"/"+outBinary.name+binaryExtension, builder.OutFolder+"/"+outBinary.name+binaryExtension)

	for _, plugin := range plugins {
		copy(plugin.outFolder+"/"+plugin.name+pluginExtension, builder.OutFolder+"/Plugins/"+plugin.name+pluginExtension)
	}

	for _, lib := range staticLibs {
		copy(lib.outFolder+"/"+lib.name+staticExtension, builder.OutFolder+lib.name+staticExtension)
	}

	return true
}

func handleFiles(rootOBSFile []byte, subFiles []ObakeBuildFolder) {
	var binaries []*BinaryType
	var staticLibs []*StaticLibType
	var plugins []*PluginType
	var obakeRootFileObj ObjectJSON

	json.Unmarshal(rootOBSFile, &obakeRootFileObj)
	//fmt.Printf("RootOBSFile: %v\n", obakeRootFileObj)

	osType = getOsType(obakeRootFileObj.Builder.Os)
	compilerFlags = obakeRootFileObj.Builder.CompilerFlags

	if isValidToolchain(obakeRootFileObj.Builder.Toolchain) {
		toolchain = obakeRootFileObj.Builder.Toolchain
	} else {
		toolchain = DEFAULT_TOOLCHAIN
	}
	fmt.Printf("OsType: %s | Toolchain: %s\n", obakeRootFileObj.Builder.Os, toolchain)
	//	fmt.Printf("SubFilesNB: %d\n", len(subFiles))

	if osType != UNKNOWN {
		for _, buildFolder := range subFiles {

			fmt.Printf("ReadFile: %s\n", "./"+buildFolder.name+"/"+OBAKE_BS_FILENAME)
			buildFolder.obakeBuildFile, _ = ioutil.ReadFile("./" + buildFolder.name + "/" + OBAKE_BS_FILENAME)
			obakeCurrentFile := getBuildFileJSONObj(buildFolder)

			if obakeCurrentFile.Binary.Name != "" {
				buildFolder.buildType = BINARY
				binaries = append(binaries, makeBinaryType(buildFolder, obakeRootFileObj.Builder.OutBinary))
			} else if obakeCurrentFile.Plugin.Name != "" {
				buildFolder.buildType = PLUGIN
				plugins = append(plugins, makePluginType(buildFolder))
			} else if obakeCurrentFile.StaticLib.Name != "" {
				buildFolder.buildType = STATIC_LIB
				staticLibs = append(staticLibs, makeStaticLibType(buildFolder))
			}
		}
		fmt.Printf("\n")
		for i, staticType := range staticLibs {
			if handleStatic(staticType, staticLibs) == false {
				staticLibs = append(staticLibs[:i], staticLibs[:i+1]...)
			}
		}
		for i, pluginType := range plugins {
			if handlePlugin(pluginType, staticLibs) == false {
				plugins = append(plugins[:i], plugins[:i+1]...)
			}
		}
		for i, binaryType := range binaries {
			if handleBinary(binaryType, staticLibs) == false {
				binaries = append(binaries[:i], binaries[:i+1]...)
			}
		}

		handleBuilder(obakeRootFileObj.Builder, binaries, staticLibs, plugins)
	}
}

func main() {
	//	var argsWithProg []string = os.Args
	var subFiles []ObakeBuildFolder
	var rootOBSFile []byte
	var subOBSFile []byte

	// argsWithProg = os.Args[1:]

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		if file.IsDir() {
			filetest := "./" + filename
			//	fmt.Printf("DirFound: %s\n", filetest)
			files, err := ioutil.ReadDir(filetest)
			if err == nil {
				subfolderName := filename
				for _, file := range files {
					filename = file.Name()
					if filename == OBAKE_BS_FILENAME {

						//	fmt.Printf("	Obake SubBuild File Found\n")
						subOBSFile, _ = ioutil.ReadFile(filename)
						var subFolderInfo ObakeBuildFolder
						subFolderInfo.buildType = NONE
						subFolderInfo = ObakeBuildFolder{path: "./" + subfolderName, name: subfolderName, obakeBuildFile: subOBSFile}
						subFiles = append(subFiles, subFolderInfo)
					}
				}
			} else {
				fmt.Printf("ERR: %s\n", err)
				log.Fatal(err)
			}
		} else if filename == OBAKE_BS_FILENAME {
			rootOBSFile, _ = ioutil.ReadFile(filename)
			//	fmt.Printf("Obake RootBuild File Found\n\n")

		}
	}

	handleFiles(rootOBSFile, subFiles)
}
