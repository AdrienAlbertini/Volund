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
	folderInfos     VolundBuildFolder
	allLibs         *[]*StaticLibType
	excludeSrc      []string
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

func sharedLibTypeToObjType(sharedLib SharedLibType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = sharedLib.excludeSrc
	objType.sourceFilesPath = sharedLib.sources
	objType.sourceFiles = sharedLib.sourceFileNames
	objType.headersFolders = sharedLib.headerFolders
	objType.externIncludes = sharedLib.externIncludes
	objType.externLibs = sharedLib.externLibs
	objType.staticLibs = sharedLib.staticLibs
	objType.sourceExtension = sharedLib.sourceExtension
	objType.outFolder = sharedLib.outFolder
	objType.compilerFlags = sharedLib.compilerFlags

	return
}

func staticLibTypeToObjType(staticLib StaticLibType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = staticLib.excludeSrc
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

func binaryTypeToObjType(binary BinaryType, folderInfos VolundBuildFolder,
	allLibs *[]*StaticLibType) (objType ObjFileRequirement) {

	objType.folderInfos = folderInfos
	objType.allLibs = allLibs
	objType.excludeSrc = binary.excludeSrc
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

func getObjFileArgs(objType ObjFileRequirement, fileID int, srcFilePath string) (args []string) {

	oFilePath := objType.outFolder + "/" + strings.Replace(objType.sourceFiles[fileID], objType.sourceExtension, ".o", -1)
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

	_, linkNames, linkIncludes := getStaticLibsLinks(objType.staticLibs, *objType.allLibs, objType.folderInfos.name)

	args = append(args, linkIncludes...)
	args = append(args, linkNames...)

	args = append(args, getExternIncludesArgs(objType.externIncludes)...)
	args = append(args, getExternLibsArgs(objType.externLibs)...)
	args = append(args, objType.compilerFlags...)
	return
}

func buildObjFile(objType ObjFileRequirement, srcFilePath string,
	i int, objCompleteChan chan int,
	success *bool, objectFilesPath *[]string,
	mutex *sync.Mutex, bar *pb.ProgressBar) {
	args := getObjFileArgs(objType, i, srcFilePath)
	oFilePath := objType.outFolder + "/" + strings.Replace(objType.sourceFiles[i], objType.sourceExtension, ".o", -1)
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
	//boldCyan.Printf("[%d/%d] ObjFile: ", (i + 1), len(objType.sourceFiles))
	//boldBlue.Printf("%s ", toolchain)
	//fmt.Printf("%v\n", args)
	bar.Increment()
	mutex.Unlock()
	objCompleteChan <- 1
}

func buildAndGetObjectFiles(objType ObjFileRequirement, success *bool,
	objectFilesPath *[]string) {
	objCompleteChan := make(chan int)
	var mutex = &sync.Mutex{}

	for fileID, srcFilePath := range objType.sourceFilesPath {
		if contains(objType.excludeSrc, srcFilePath) == false {
			fmt.Printf("\t%-30s %v\n", srcFilePath, getObjFileArgs(objType, fileID, srcFilePath))
		}
	}

	sourceFilesPathLen := len(objType.sourceFilesPath)
	fmt.Printf("\n")

	//boldGreen.Printf("With Args: %v\n", args)

	bar := pb.New(len(objType.sourceFilesPath))
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
	for i, srcFilePath := range objType.sourceFilesPath {

		if contains(objType.excludeSrc, srcFilePath) == false {
			go buildObjFile(objType, srcFilePath, i, objCompleteChan, success, objectFilesPath, mutex, bar)
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
	}

	bar.Finish()
	fmt.Printf("\n")
	*success = true
}
