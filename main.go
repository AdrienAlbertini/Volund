package main

import (
	"encoding/json"
	"fmt"
	//	"gopkg.in/cheggaaa/pb.v1"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

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

	boldCyan.Printf("(%d files) Compiling Binary: %s\n", len(binary.sources), binary.name)

	if buildDependencies(binary.staticLibs, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(binaryTypeToObjType(*binary, binary.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("Build Binary: %s FAILED\n\n", binary.name)
		return false
	}

	binaryExtension := getBinaryOSExtension()

	args := []string{"-o", binary.outFolder + "/" + binary.name + binaryExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, binary.compilerFlags...)

	//args = append(args, getExternIncludesArgs(binary.externIncludes)...)
	args = append(args, getExternLibsArgs(binary.externLibs)...)

	fmt.Printf("Handle Binary args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldGreen.Printf("Binary: %s | Error: %s\n\n", binary.name, fmt.Sprint(err))
		fmt.Printf("%s\n", string(out))
		return false
	}
	fmt.Printf("%s\n\n", out)

	return true
}

func handleStatic(staticLib *StaticLibType, allLibs []*StaticLibType) bool {

	if staticLib.isBuilt == false {

		boldCyan.Printf("(%d files) Compiling StaticLib: %s\n", len(staticLib.sources), staticLib.name)

		linkPaths, linkNames, linkIncludes := getStaticLibsLinks(staticLib.staticLibs, allLibs, staticLib.name)

		if buildDependencies(staticLib.staticLibs, allLibs) == false {
			return false
		}

		var objSuccess bool
		var objectFilesPath []string
		buildAndGetObjectFiles(staticLibTypeToObjType(*staticLib, staticLib.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

		if objSuccess == false {
			boldRed.Printf("Build StaticLib: %s FAILED\n\n", staticLib.name)
			return false
		}

		staticLibExtension := getStaticLibOSExtension()

		args := []string{"rcs", staticLib.outFolder + "/" + staticLib.name + staticLibExtension}

		args = append(args, objectFilesPath...)

		args = append(args, linkIncludes...)
		args = append(args, linkPaths...)
		args = append(args, linkNames...)
		args = append(args, staticLib.compilerFlags...)

		//	args = append(args, getExternIncludesArgs(staticLib.externIncludes)...)
		args = append(args, getExternLibsArgs(staticLib.externLibs)...)

		fmt.Printf("Handle Static args: %v\n", args)
		cmd := exec.Command("ar", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			boldRed.Printf("StaticLib: %s | Error: %s\n\n", staticLib.name, fmt.Sprint(err))
			fmt.Printf("%s\n", string(out))
			return false
		}
		fmt.Printf("%s\n\n", out)

		staticLib.isBuilt = true
	}
	return true
}

func handlePlugin(plugin *PluginType, allLibs []*StaticLibType) bool {

	boldCyan.Printf("(%d files) Compiling Plugin: %s\n", len(plugin.sources), plugin.name)

	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(plugin.staticLibs, allLibs, "")

	if buildDependencies(plugin.staticLibs, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(pluginTypeToObjType(*plugin, plugin.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("Build Plugin: %s FAILED\n\n", plugin.name)
		return false
	}

	pluginLibExtension := getSharedLibOsExtension()

	args := []string{"-shared", "-o", plugin.outFolder + "/" + plugin.name + pluginLibExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, plugin.compilerFlags...)

	args = append(args, getExternLibsArgs(plugin.externLibs)...)

	fmt.Printf("Handle Plugin args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldRed.Printf("Plugin: %s | Error: %s\n\n", plugin.name, fmt.Sprint(err.Error()))
		fmt.Printf("Out: %s\n\n", out)
		return false
	}
	fmt.Printf("Out: %s\n\n", out)

	return true
}

func handleBuilder(builder BuilderJSON, binaries []*BinaryType,
	staticLibs []*StaticLibType, plugins []*PluginType) bool {
	outBinaryFound := false
	var outBinary *BinaryType

	for _, binary := range binaries {
		if contains(builder.Binaries, binary.name) && binary.isOutBinary == true {
			outBinary = binary
			outBinaryFound = true
			break
		}
	}

	if outBinaryFound == false {
		boldYellow.Printf("Out Binary not found\n")
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

	boldCyan.Printf("Copying binary files.\n")
	copy(outBinary.outFolder+"/"+outBinary.name+binaryExtension, builder.OutFolder+"/"+outBinary.name+binaryExtension)

	for _, plugin := range plugins {
		if contains(builder.Plugins, plugin.name) {
			copy(plugin.outFolder+"/"+plugin.name+pluginExtension, builder.OutFolder+"/Plugins/"+plugin.name+pluginExtension)
		}
	}

	for _, lib := range staticLibs {
		if contains(builder.StaticLibs, lib.name) {
			copy(lib.outFolder+"/"+lib.name+staticExtension, builder.OutFolder+lib.name+staticExtension)
		}
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
	if contains(obakeRootFileObj.Builder.Binaries, obakeRootFileObj.Builder.OutBinary) == false {
		obakeRootFileObj.Builder.Binaries = append(obakeRootFileObj.Builder.Binaries, obakeRootFileObj.Builder.OutBinary)
	}

	if obakeRootFileObj.Builder.FullStatic {
		obakeRootFileObj.Builder.StaticLibs = append(obakeRootFileObj.Builder.StaticLibs, obakeRootFileObj.Builder.Plugins...)
		obakeRootFileObj.Builder.Plugins = []string{}
	}

	if isValidToolchain(obakeRootFileObj.Builder.Toolchain) {
		toolchain = obakeRootFileObj.Builder.Toolchain
	} else {
		toolchain = DEFAULT_TOOLCHAIN
	}
	boldRed.Printf("Volund: OsType: %s | Toolchain: %s\n", obakeRootFileObj.Builder.Os, toolchain)
	//	fmt.Printf("SubFilesNB: %d\n", len(subFiles))

	if osType != UNKNOWN {
		for _, buildFolder := range subFiles {

			boldGreen.Print("ReadFile: ")
			fmt.Printf("%s\n", "./"+buildFolder.name+"/"+OBAKE_BS_FILENAME)
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
			if contains(obakeRootFileObj.Builder.StaticLibs, staticType.name) == false || handleStatic(staticType, staticLibs) == false {
				staticLibs = append(staticLibs[:i], staticLibs[:i+1]...)
			}
		}
		for i, pluginType := range plugins {
			if contains(obakeRootFileObj.Builder.Plugins, pluginType.name) == false || handlePlugin(pluginType, staticLibs) == false {
				plugins = append(plugins[:i], plugins[:i+1]...)
			}
		}
		for i, binaryType := range binaries {
			if contains(obakeRootFileObj.Builder.Binaries, binaryType.name) == false || handleBinary(binaryType, staticLibs) == false {
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

	initCustomColors()

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
