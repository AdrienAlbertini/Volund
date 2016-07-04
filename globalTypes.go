package main

type ObakeToolchainType int64
type ObakeOSType int64
type ObakeBuildType int64

const OBAKE_BS_FILENAME string = "ObakeBuild.json"
const DEFAULT_TOOLCHAIN string = "clang++"

const (
	BINARY ObakeBuildType = iota
	STATIC_LIB
	PLUGIN
	Builder
	NONE
)

const (
	WINDOWS ObakeOSType = iota
	LINUX
	OSX
	UNKNOWN
)

const (
	WINDOWS_STATIC_EXT  = ".lib"
	WINDOWS_DYNAMIC_EXT = ".dll"
	WINDOWS_BINARY_EXT  = ".exe"

	LINUX_STATIC_EXT  = ".a"
	LINUX_DYNAMIC_EXT = ".so"
	LINUX_BINARY_EXT  = ""

	OSX_STATIC_EXT  = ".a"
	OSX_DYNAMIC_EXT = ".dylib"
	OSX_BINARY_EXT  = ""
)

var compilerFlags []string
var osType ObakeOSType
var toolchain string
var builder BuilderType
