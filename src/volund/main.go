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

	_, linkNames, linkIncludes := getStaticLibsLinks(binary.staticLibs, allLibs, binary.name)

	boldCyan.Printf("(%d files) Compiling Binary: %s\n", len(binary.sources), binary.name)

	if buildDependencies(binary.staticLibs, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(binaryTypeToObjType(*binary, binary.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("ERROR: Build Binary: %s FAILED\n\n", binary.name)
		return false
	}

	binaryExtension := getBinaryOSExtension()

	args := []string{"-o", binary.outFolder + "/" + binary.name + binaryExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	//	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, binary.compilerFlags...)

	//args = append(args, getExternIncludesArgs(binary.externIncludes)...)
	args = append(args, getExternLibsArgs(binary.externLibs)...)

	fmt.Printf("Handle Binary args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldRed.Printf("ERROR: Binary: %s | Error: %s\n\n", binary.name, fmt.Sprint(err))
		fmt.Printf("%s\n", string(out))
		return false
	}
	fmt.Printf("%s\n\n", out)
	binary.isBuilt = true

	return true
}

func handleStatic(staticLib *StaticLibType, allLibs []*StaticLibType) bool {

	if staticLib.isBuilt == false {

		boldCyan.Printf("(%d files) Compiling StaticLib: %s\n", len(staticLib.sources), staticLib.name)

		//	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(staticLib.staticLibs, allLibs, staticLib.name)

		if buildDependencies(staticLib.staticLibs, allLibs) == false {
			return false
		}

		var objSuccess bool
		var objectFilesPath []string
		buildAndGetObjectFiles(staticLibTypeToObjType(*staticLib, staticLib.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

		if objSuccess == false {
			boldRed.Printf("ERROR: Build StaticLib: %s FAILED\n\n", staticLib.name)
			return false
		}

		staticLibExtension := getStaticLibOSExtension()

		args := []string{"rcs", staticLib.outFolder + "/" + staticLib.name + staticLibExtension}

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
			boldRed.Printf("ERROR: StaticLib: %s | Linker Error: %s\n\n", staticLib.name, fmt.Sprint(err))
			fmt.Printf("%s\n", string(out))
			return false
		}
		fmt.Printf("%s\n\n", out)

		staticLib.isBuilt = true
	}
	return true
}

func handlesharedLib(sharedLib *SharedLibType, allLibs []*StaticLibType) bool {

	boldCyan.Printf("(%d files) Compiling SharedLib: %s\n", len(sharedLib.sources), sharedLib.name)

	linkPaths, linkNames, linkIncludes := getStaticLibsLinks(sharedLib.staticLibs, allLibs, "")

	if buildDependencies(sharedLib.staticLibs, allLibs) == false {
		return false
	}

	var objSuccess bool
	var objectFilesPath []string
	buildAndGetObjectFiles(sharedLibTypeToObjType(*sharedLib, sharedLib.folderInfos, &allLibs), &objSuccess, &objectFilesPath)

	if objSuccess == false {
		boldRed.Printf("ERROR: Build SharedLib: %s FAILED\n\n", sharedLib.name)
		return false
	}

	sharedLibLibExtension := getSharedLibOsExtension()

	osSharedFlag := "-shared"

	switch osType {
	case LINUX:
		osSharedFlag = "-fPIC"
	}

	args := []string{osSharedFlag, "-o", sharedLib.outFolder + "/" + sharedLib.name + sharedLibLibExtension}

	args = append(args, compilerFlags...)
	args = append(args, objectFilesPath...)

	args = append(args, linkIncludes...)
	args = append(args, linkPaths...)
	args = append(args, linkNames...)
	args = append(args, sharedLib.compilerFlags...)

	args = append(args, getExternLibsArgs(sharedLib.externLibs)...)

	fmt.Printf("Handle sharedLib args: %v\n", args)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		boldRed.Printf("ERROR: SharedLib: %s | Error: %s\n\n", sharedLib.name, fmt.Sprint(err.Error()))
		fmt.Printf("Out: %s\n\n", out)
		return false
	}
	fmt.Printf("Out: %s\n\n", out)
	sharedLib.isBuilt = true

	return true
}

func handleBuilder(outBinaryError bool, builder BuilderJSON, binaries []*BinaryType,
	staticLibs []*StaticLibType, sharedLibs []*SharedLibType) bool {
	outBinaryFound := false
	var outBinary *BinaryType

	for _, binary := range binaries {
		fmt.Printf("OutBinary: %s | CurrentBinary: %s\n", builder.OutBinary, binary.name)
		if builder.OutBinary == binary.name {
			outBinary = binary
			outBinaryFound = true
			break
		}
	}

	if outBinaryFound == false {
		if outBinaryError {
			boldRed.Printf("ERROR: Out Binary Build FAILED\n")
		} else {
			boldYellow.Printf("Out Binary not found\n")
		}
		return false
	}

	binaryExtension := getBinaryOSExtension()
	staticExtension := getStaticLibOSExtension()
	sharedLibExtension := getSharedLibOsExtension()

	success, _ := exists(builder.OutFolder)
	if !success {
		os.MkdirAll(builder.OutFolder, os.ModePerm)
	}
	success, _ = exists(builder.OutFolder + "sharedLibs")
	if !success {
		os.MkdirAll(builder.OutFolder+"/sharedLibs", os.ModePerm)
	}

	boldCyan.Printf("Copying out binary files.\n")
	copy(outBinary.outFolder+"/"+outBinary.name+binaryExtension, builder.OutFolder+"/"+outBinary.name+binaryExtension)

	for _, sharedLib := range sharedLibs {
		if contains(outBinary.sharedLibs, sharedLib.name) {
			copy(sharedLib.outFolder+"/"+sharedLib.name+sharedLibExtension, builder.OutFolder+"/sharedLibs/"+sharedLib.name+sharedLibExtension)
		}
	}

	for _, lib := range staticLibs {
		if contains(outBinary.staticLibs, lib.name) {
			copy(lib.outFolder+"/"+lib.name+staticExtension, builder.OutFolder+lib.name+staticExtension)
		}
	}

	return true
}

func handleFiles(rootVolundFile []byte, subFiles []VolundBuildFolder) {
	var binaries []*BinaryType
	var staticLibs []*StaticLibType
	var sharedLibs []*SharedLibType
	var volundRootFileObj ObjectJSON
	outBinaryError := false

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
	if contains(volundRootFileObj.Builder.Binaries, volundRootFileObj.Builder.OutBinary) == false {
		volundRootFileObj.Builder.Binaries = append(volundRootFileObj.Builder.Binaries, volundRootFileObj.Builder.OutBinary)
	}

	if volundRootFileObj.Builder.FullStatic {
		volundRootFileObj.Builder.StaticLibs = append(volundRootFileObj.Builder.StaticLibs, volundRootFileObj.Builder.SharedLibs...)
		volundRootFileObj.Builder.SharedLibs = []string{}
	}

	if isValidToolchain(volundRootFileObj.Builder.Toolchain) {
		toolchain = volundRootFileObj.Builder.Toolchain
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

			if resolveBuildType(&volundCurrentFile, &buildFolder, &volundRootFileObj.Builder.Binaries,
				&volundRootFileObj.Builder.StaticLibs, &volundRootFileObj.Builder.SharedLibs) {
				switch buildFolder.buildType {
				case BINARY:
					binaries = append(binaries, makeBinaryType(buildFolder, volundRootFileObj.Builder.OutBinary))
				case SHARED_LIB:
					sharedLibs = append(sharedLibs, makeSharedLibType(buildFolder))
				case STATIC_LIB:
					staticLibs = append(staticLibs, makeStaticLibType(buildFolder))
				}
			} else {
				boldRed.Printf("ERROR: can't find build type for this file.\n")
				continue
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

		var outBinary *BinaryType
		for _, checkBinary := range binaries {
			if checkBinary.name == volundRootFileObj.Builder.OutBinary {
				outBinary = checkBinary
				break
			}
		}

		fmt.Printf("\n")
		for i := 0; i < len(staticLibs); i++ {
			staticType := staticLibs[i]
			if (contains(volundRootFileObj.Builder.StaticLibs, staticType.name) == false && contains(outBinary.staticLibs, staticType.name) == false) || handleStatic(staticType, staticLibs) == false {
				staticLibs = append(staticLibs[:i], staticLibs[i+1:]...)
				i = -1
			}
		}
		for i := 0; i < len(sharedLibs); i++ {
			sharedLibType := sharedLibs[i]
			if (contains(volundRootFileObj.Builder.SharedLibs, sharedLibType.name) == false && contains(outBinary.sharedLibs, sharedLibType.name) == false) || handlesharedLib(sharedLibType, staticLibs) == false {
				sharedLibs = append(sharedLibs[:i], sharedLibs[i+1:]...)
				i = -1
			}
		}

		for i := 0; i < len(binaries); i++ {
			binaryType := binaries[i]
			fmt.Printf("Binary[%d]: %s\n", i, binaryType.name)
			if contains(volundRootFileObj.Builder.Binaries, binaryType.name) == false || handleBinary(binaryType, staticLibs) == false {
				binaries = append(binaries[:i], binaries[i+1:]...)
				if volundRootFileObj.Builder.OutBinary == binaryType.name {
					outBinaryError = true
				}
				i = -1
			}
		}

		handleBuilder(outBinaryError, volundRootFileObj.Builder, binaries, staticLibs, sharedLibs)
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
