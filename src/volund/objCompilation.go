package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
	//	"sync/atomic"
	"gopkg.in/cheggaaa/pb.v1"
)

type ObjFileRequirement struct {
	folderInfos    VolundBuildFolder
	allLibs        *[]*StaticLibType
	excludeSrc     []string
	srcFilesPaths  []string
	srcFiles       []string
	headersFolders []string
	internLibs     []string
	externIncludes []string
	externLibs     []string
	staticLibsDeps []string
	sharedLibsDeps []string
	compilerFlags  []string
	srcExtension   string
	outFolder      string
}

func sharedLibTypeToObjType(sharedLib SharedLibType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = sharedLib.excludeSrc
	objType.srcFilesPaths = sharedLib.src
	objType.srcFiles = sharedLib.srcFileNames
	objType.headersFolders = sharedLib.headersFolders
	objType.internLibs = sharedLib.internLibs
	objType.externIncludes = sharedLib.externIncludes
	objType.externLibs = sharedLib.externLibs
	objType.staticLibsDeps = sharedLib.staticLibsDeps
	objType.srcExtension = sharedLib.srcExtension
	objType.outFolder = sharedLib.outFolder
	objType.compilerFlags = sharedLib.compilerFlags

	return
}

func staticLibTypeToObjType(staticLib StaticLibType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = staticLib.excludeSrc
	objType.srcFilesPaths = staticLib.src
	objType.srcFiles = staticLib.srcFileNames
	objType.headersFolders = staticLib.headersFolders
	objType.internLibs = staticLib.internLibs
	objType.externIncludes = staticLib.externIncludes
	objType.externLibs = staticLib.externLibs
	objType.staticLibsDeps = staticLib.staticLibsDeps
	objType.srcExtension = staticLib.srcExtension
	objType.outFolder = staticLib.outFolder
	objType.compilerFlags = staticLib.compilerFlags

	return
}

func executableTypeToObjType(executable ExecutableType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = executable.excludeSrc
	objType.srcFilesPaths = executable.src
	objType.srcFiles = executable.srcFileNames
	objType.headersFolders = executable.headersFolders
	objType.internLibs = executable.internLibs
	objType.externIncludes = executable.externIncludes
	objType.externLibs = executable.externLibs
	objType.staticLibsDeps = executable.staticLibsDeps
	objType.srcExtension = executable.srcExtension
	objType.outFolder = executable.outFolder
	objType.compilerFlags = executable.compilerFlags

	return
}

func getObjFileArgs(objType ObjFileRequirement, fileID int, srcFilePath string) (args []string) {

	oFilePath := objType.outFolder + "/" + strings.Replace(objType.srcFiles[fileID], objType.srcExtension, ".o", -1)
	//	*objectFilesPath = append(*objectFilesPath, oFilePath)

	//		fmt.Printf("SrcFilePath: %s\n", srcFilePath)
	args = []string{"-c", objType.folderInfos.path + "/" + srcFilePath, "-o", oFilePath}
	args = append(args, compilerFlags...)

	for _, headerFolder := range objType.headersFolders {
		if headerFolder == "." {
			args = append(args, "-I"+objType.folderInfos.path+"/")
		} else {
			args = append(args, "-I"+objType.folderInfos.path+"/"+headerFolder+"/")
		}
	}

	_, linkNames, linkIncludes := getStaticLibsLinks(objType.staticLibsDeps, *objType.allLibs, objType.folderInfos.name)

	args = append(args, linkIncludes...)
	args = append(args, linkNames...)

	args = append(args, getExternIncludesArgs(objType.externIncludes)...)
	args = append(args, getLibsArgs(objType.externLibs)...)
	args = append(args, objType.compilerFlags...)
	return
}

func buildObjFile(objType ObjFileRequirement, srcFilePath string,
	i int, objCompleteChan chan int,
	success *bool, objectFilesPath *[]string,
	externLibs []string, externIncludes []string,
	mutex *sync.Mutex, bar *pb.ProgressBar) {
	args := getObjFileArgs(objType, i, srcFilePath)
	args = append(args, externIncludes...)
	args = append(args, externLibs...)

	oFilePath := objType.outFolder + "/" + strings.Replace(objType.srcFiles[i], objType.srcExtension, ".o", -1)
	*objectFilesPath = append(*objectFilesPath, oFilePath)

	cmd := exec.Command(toolchain, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		mutex.Lock()
		bar.Finish()
		boldRed.Printf("ERROR: ObjFile: %s | Error: %s:\nArgs: %v \n\n", srcFilePath, fmt.Sprint(err), args)
		fmt.Printf("%s\n", string(out))
		mutex.Unlock()
		//fmt.Printf("ObjFile: %s | Error: %s\n", srcFilePath, fmt.Sprint(err))
		*success = false
		objCompleteChan <- 1
		return
	}
	mutex.Lock()
	*success = true
	//boldCyan.Printf("[%d/%d] ObjFile: ", (i + 1), len(objType.sourceFiles))
	//boldBlue.Printf("%s ", toolchain)
	//fmt.Printf("%v\n", args)
	bar.Increment()
	mutex.Unlock()
	objCompleteChan <- 1
}

func buildAndGetObjectFiles(objType ObjFileRequirement, externLibs []string, externIncludes []string,
	success *bool, objectFilesPath *[]string) {

	if len(objType.srcFilesPaths) > 0 {
		objCompleteChan := make(chan int)
		var mutex = &sync.Mutex{}

		for fileID, srcFilePath := range objType.srcFilesPaths {
			if contains(objType.excludeSrc, srcFilePath) == false {
				args := getObjFileArgs(objType, fileID, srcFilePath)
				args = append(args, externIncludes...)
				args = append(args, externLibs...)
				fmt.Printf("\t%s \n\t%v\n\n", srcFilePath, args)
			}
		}

		sourceFilesPathLen := len(objType.srcFilesPaths)
		fmt.Printf("\n")

		//boldGreen.Printf("With Args: %v\n", args)

		bar := pb.New(len(objType.srcFilesPaths))
		bar.SetRefreshRate(time.Millisecond)
		bar.Prefix("Compile ObjFiles:  ")
		bar.ShowBar = true
		bar.ShowPercent = true
		bar.ShowCounters = true
		bar.ShowSpeed = true
		bar.SetWidth(80)
		bar.SetMaxWidth(80)
		bar.SetUnits(pb.U_NO)
		bar.Start()

		for i, srcFilePath := range objType.srcFilesPaths {

			if contains(objType.excludeSrc, srcFilePath) == false {
				go buildObjFile(objType, srcFilePath, i, objCompleteChan, success, objectFilesPath, externLibs, externIncludes, mutex, bar)
			} else {
				sourceFilesPathLen--
			}
			/*
					fmt.Printf("Obj files: %s %v\n", toolchain, args)
					_, err := exec.Command(toolchain, args...).Output()
					if err != nil {
					fmt.Printf("ObjFile: %s | Error: %s\n", srcFilePath, err)
				}
			*/
		}

		for i := 0; i < sourceFilesPathLen; i++ {
			<-objCompleteChan
			if *success == false {
				return
			}
		}

		bar.Finish()

	} else {
		boldRed.Printf("No source files found.\n")
		*success = false
		return
	}

	fmt.Printf("\n")
	*success = true
}
