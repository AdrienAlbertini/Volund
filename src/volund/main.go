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

func handleExecutable(executable *ExecutableType, allLibs []*StaticLibType) bool {

	_, linkNames, linkIncludes := getStaticLibsLinks(executable.staticLibsDeps, allLibs, executable.targetName)

	boldCyan.Printf("(%d files) Compiling Executable: %s\n", len(executable.src), executable.targetName)

	if buildDependencies(executable.staticLibsDeps, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(executableTypeToObjType(*executable, executable.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("ERROR: Build Executable: %s FAILED\n\n", executable.targetName)
		return false
	}

	executableExtension := getExecutableOSExtension()

	args := []string{"-o", executable.outFolder + "/" + executable.targetName + executableExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	//	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, executable.compilerFlags...)

	args = append(args, getExternIncludesArgs(executable.externIncludes)...)
	args = append(args, getLibsArgs(executable.externLibs)...)

	fmt.Printf("Handle Executable args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldRed.Printf("ERROR: Executable: %s | Error: %s\n\n", executable.targetName, fmt.Sprint(err))
		fmt.Printf("%s\n", string(out))
		return false
	}

	if osType == OSX || osType == LINUX {
		args = []string{"+x", executable.outFolder + "/" + executable.targetName + executableExtension}
		exec.Command("chmod", args...).Run()
	}

	fmt.Printf("%s\n\n", out)
	executable.isBuilt = true

	return true
}

func handleStatic(staticLib *StaticLibType, allLibs []*StaticLibType) bool {

	if staticLib.isBuilt == false {

		boldCyan.Printf("(%d files) Compiling StaticLib: %s\n", len(staticLib.src), staticLib.targetName)

		//	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(staticLib.staticLibs, allLibs, staticLib.name)

		if buildDependencies(staticLib.staticLibsDeps, allLibs) == false {
			return false
		}

		var objSuccess bool
		var objectFilesPath []string
		buildAndGetObjectFiles(staticLibTypeToObjType(*staticLib, staticLib.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

		if objSuccess == false {
			boldRed.Printf("ERROR: Build StaticLib: %s FAILED\n\n", staticLib.targetName)
			return false
		}

		staticLibExtension := getStaticLibOSExtension()

		args := []string{"rcs", staticLib.outFolder + "/" + staticLib.targetName + staticLibExtension}

		args = append(args, objectFilesPath...)

		//	args = append(args, linkIncludes...)
		//	args = append(args, linkPaths...)
		//	args = append(args, linkNames...)
		//	args = append(args, staticLib.compilerFlags...)

		//	args = append(args, getExternIncludesArgs(staticLib.externIncludes)...)
		//	args = append(args, getExternLibsArgs(staticLib.externLibs)...)

		fmt.Printf("Handle Static args: %v\n", args)
		cmd := exec.Command("ar", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			boldRed.Printf("ERROR: StaticLib: %s | Linker Error: %s\n\n", staticLib.targetName, fmt.Sprint(err))
			fmt.Printf("%s\n", string(out))
			return false
		}
		fmt.Printf("%s\n\n", out)

		staticLib.isBuilt = true
	}
	return true
}

func handleSharedLib(sharedLib *SharedLibType, allLibs []*StaticLibType) bool {

	boldCyan.Printf("(%d files) Compiling SharedLib: %s\n", len(sharedLib.src), sharedLib.targetName)

	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(sharedLib.staticLibsDeps, allLibs, "")

	if buildDependencies(sharedLib.staticLibsDeps, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(sharedLibTypeToObjType(*sharedLib, sharedLib.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("ERROR: Build SharedLib: %s FAILED\n\n", sharedLib.targetName)
		return false
	}

	sharedLibExtension := getSharedLibOsExtension()

	osSharedFlag := "-shared"

	switch osType {
	case LINUX:
		osSharedFlag = "-fPIC"
	}

	args := []string{osSharedFlag, "-o", sharedLib.outFolder + "/" + sharedLib.targetName + sharedLibExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, sharedLib.compilerFlags...)

	args = append(args, getLibsArgs(sharedLib.externLibs)...)

	fmt.Printf("Handle sharedLib args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldRed.Printf("ERROR: SharedLib: %s | Error: %s\n\n", sharedLib.targetName, fmt.Sprint(err.Error()))
		fmt.Printf("Out: %s\n\n", out)
		return false
	}
	fmt.Printf("Out: %s\n\n", out)
	sharedLib.isBuilt = true

	return true
}

func handleBuilder(mainBinaryError bool, builder BuilderJSON, executables []*ExecutableType,
	staticLibs []*StaticLibType, sharedLibs []*SharedLibType) bool {
	mainExecutableFound := false
	var mainExecutable *ExecutableType

	for _, executable := range executables {
		fmt.Printf("OutBinary: %s | CurrentExecutable: %s\n", builder.MainExecutable, executable.targetName)
		if builder.MainExecutable == executable.targetName {
			mainExecutable = executable
			mainExecutableFound = true
			break
		}
	}

	if mainExecutableFound == false {
		if mainBinaryError {
			boldRed.Printf("ERROR: Out Binary Build FAILED\n")
		} else {
			boldYellow.Printf("Out Binary not found\n")
		}
		return false
	}

	binaryExtension := getExecutableOSExtension()
	staticExtension := getStaticLibOSExtension()
	sharedLibExtension := getSharedLibOsExtension()

	success, _ := exists(builder.MainFolder)
	if !success {
		os.MkdirAll(builder.MainFolder, os.ModePerm)
	}
	success, _ = exists(builder.MainFolder + "sharedLibs")
	if !success {
		os.MkdirAll(builder.MainFolder+"/sharedLibs", os.ModePerm)
	}

	boldCyan.Printf("\nCopying out binary files.\n")
	copy(mainExecutable.outFolder+"/"+mainExecutable.targetName+binaryExtension, builder.MainFolder+"/"+mainExecutable.targetName+binaryExtension)

	for _, sharedLib := range sharedLibs {
		if contains(mainExecutable.sharedLibsDeps, sharedLib.targetName) {
			copy(sharedLib.outFolder+"/"+sharedLib.targetName+sharedLibExtension, builder.MainFolder+"/sharedLibs/"+sharedLib.targetName+sharedLibExtension)
		}
	}

	for _, lib := range staticLibs {
		if contains(mainExecutable.staticLibsDeps, lib.targetName) {
			copy(lib.outFolder+"/"+lib.targetName+staticExtension, builder.MainFolder+lib.targetName+staticExtension)
		}
	}

	return true
}

func handleFiles(rootVolundFile []byte, subFiles []VolundBuildFolder) {
	var executables []*ExecutableType
	var staticLibs []*StaticLibType
	var sharedLibs []*SharedLibType
	var volundRootFileObj ObjectJSON
	mainExecutableError := false

	json.Unmarshal(rootVolundFile, &volundRootFileObj)

	if volundRootFileObj.IsEmpty() || volundRootFileObj.Builder.IsEmpty() {
		boldRed.Printf("ERROR : Can't parse builder json\n")
		return
	}

	osType = getOsType(volundRootFileObj.Builder.Os)

	if osType == UNKNOWN {
		osType = getRuntimeOS()
	}
	//	volundRootFileObj.Builder = resolveBuilderOSParams(volundRootFileObj.Builder)

	if osType == UNKNOWN {
		boldRed.Printf("ERROR: OS not supported\n")
		return
	}

	compilerFlags = volundRootFileObj.Builder.CompilerFlags
	if contains(volundRootFileObj.Builder.Executables, volundRootFileObj.Builder.MainExecutable) == false {
		volundRootFileObj.Builder.Executables = append(volundRootFileObj.Builder.Executables, volundRootFileObj.Builder.MainExecutable)
	}

	if volundRootFileObj.Builder.FullStatic {
		volundRootFileObj.Builder.StaticLibs = append(volundRootFileObj.Builder.StaticLibs, volundRootFileObj.Builder.SharedLibs...)
		volundRootFileObj.Builder.SharedLibs = []string{}
	}

	if isValidToolchain(volundRootFileObj.Builder.Compiler) {
		toolchain = volundRootFileObj.Builder.Compiler
	} else {
		toolchain = DEFAULT_TOOLCHAIN
	}

	boldRed.Printf("Volund: OS: %s | Toolchain: %s\n\n", osType.ToString(), toolchain)
	//	fmt.Printf("SubFilesNB: %d\n", len(subFiles))

	if osType != UNKNOWN {
		for _, buildFolder := range subFiles {

			boldGreen.Print("ReadFile: ")
			fmt.Printf("%s\n", "./"+buildFolder.name+"/"+VOLUND_BUILD_FILENAME)
			buildFolder.volundBuildFile, _ = ioutil.ReadFile("./" + buildFolder.name + "/" + VOLUND_BUILD_FILENAME)
			volundCurrentFile := getFileJSONObj(buildFolder)

			if resolveBuildType(&volundCurrentFile, &buildFolder, &volundRootFileObj.Builder.Executables,
				&volundRootFileObj.Builder.StaticLibs, &volundRootFileObj.Builder.SharedLibs) {
				switch buildFolder.buildType {
				case EXECUTABLE:
					executables = append(executables, makeExecutableType(buildFolder, volundRootFileObj.Builder.MainExecutable))
					if contains(volundRootFileObj.Builder.Executables, volundCurrentFile.Executable.TargetName) == false {
						boldYellow.Printf("WARNING: %s will not be build (not present in Builder file).\n", volundCurrentFile.Executable.TargetName)
					}
				case SHARED_LIB:
					sharedLibs = append(sharedLibs, makeSharedLibType(buildFolder))
					if contains(volundRootFileObj.Builder.SharedLibs, volundCurrentFile.SharedLib.TargetName) == false {
						boldYellow.Printf("WARNING: %s will not be build (not present in Builder file).\n", volundCurrentFile.SharedLib.TargetName)
					}
				case STATIC_LIB:
					staticLibs = append(staticLibs, makeStaticLibType(buildFolder))
					if contains(volundRootFileObj.Builder.StaticLibs, volundCurrentFile.StaticLib.TargetName) == false {
						boldYellow.Printf("WARNING: %s will not be build (not present in Builder file).\n", volundCurrentFile.StaticLib.TargetName)
					}
				case NONE:
					boldRed.Printf("ERROR: can't find build type for this file.\n")
				}
			}
			/*
				if volundCurrentFile.Binary.IsEmpty() == false {
					buildFolder.buildType = BINARY
					binaries = append(binaries, makeBinaryType(buildFolder, volundRootFileObj.Builder.OutBinary))
				} else if volundCurrentFile.SharedLib.IsEmpty() == false {
					buildFolder.buildType = SHARED_LIB
					sharedLibs = append(sharedLibs, makeSharedLibType(buildFolder))
				} else if volundCurrentFile.StaticLib.IsEmpty() == false {
					buildFolder.buildType = STATIC_LIB
					staticLibs = append(staticLibs, makeStaticLibType(buildFolder))
				} else {
					buildFolder.buildType = NONE
					boldYellow.Printf("WARNING : Can't parse json: %s\n", "./"+buildFolder.name+"/"+VOLUND_BUILD_FILENAME)
				}
			*/
		}

		var mainExecutable *ExecutableType
		for _, checkExecutable := range executables {
			if checkExecutable.targetName == volundRootFileObj.Builder.MainExecutable {
				mainExecutable = checkExecutable
				break
			}
		}

		fmt.Printf("\n")
		for i := 0; i < len(staticLibs); i++ {
			staticType := staticLibs[i]
			if (contains(volundRootFileObj.Builder.StaticLibs, staticType.targetName) == false && contains(mainExecutable.staticLibsDeps, staticType.targetName) == false) || handleStatic(staticType, staticLibs) == false {
				staticLibs = append(staticLibs[:i], staticLibs[i+1:]...)
				i = -1
			}
		}
		for i := 0; i < len(sharedLibs); i++ {
			sharedLibType := sharedLibs[i]
			if (contains(volundRootFileObj.Builder.SharedLibs, sharedLibType.targetName) == false && contains(mainExecutable.sharedLibsDeps, sharedLibType.targetName) == false) || handleSharedLib(sharedLibType, staticLibs) == false {
				sharedLibs = append(sharedLibs[:i], sharedLibs[i+1:]...)
				i = -1
			}
		}

		for i := 0; i < len(executables); i++ {
			executableType := executables[i]
			if contains(volundRootFileObj.Builder.Executables, executableType.targetName) == false || handleExecutable(executableType, staticLibs) == false {
				executables = append(executables[:i], executables[i+1:]...)
				if volundRootFileObj.Builder.MainExecutable == executableType.targetName {
					mainExecutableError = true
				}
				i = -1
			}
		}

		handleBuilder(mainExecutableError, volundRootFileObj.Builder, executables, staticLibs, sharedLibs)
	}
}

func main() {
	//	var argsWithProg []string = os.Args
	var subFiles []VolundBuildFolder
	var rootVolundFile []byte
	var subVolundFile []byte

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
			//	fmt.Printf("DirFound: %s\n", filetest)3
			files, err := ioutil.ReadDir(filetest)
			if err == nil {
				subfolderName := filename
				for _, file := range files {
					filename = file.Name()
					if filename == VOLUND_BUILD_FILENAME {

						//	fmt.Printf("	Volund SubBuild File Found\n")
						subVolundFile, _ = ioutil.ReadFile(filename)
						var subFolderInfo VolundBuildFolder
						subFolderInfo.buildType = NONE
						subFolderInfo = VolundBuildFolder{path: "./" + subfolderName, name: subfolderName, volundBuildFile: subVolundFile}
						subFiles = append(subFiles, subFolderInfo)
					}
				}
			} else {
				fmt.Printf("ERR: %s\n", err)
				log.Fatal(err)
			}
		} else if filename == VOLUND_BUILD_FILENAME {
			rootVolundFile, _ = ioutil.ReadFile(filename)
			//	fmt.Printf("Volund RootBuild File Found\n\n")

		}
	}

	handleFiles(rootVolundFile, subFiles)
}
