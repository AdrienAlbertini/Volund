package main

type VolundToolchainType int64
type VolundOSType int64
type VolundBuildType int64

const VOLUND_BUILD_FILENAME string = "VolundBuild.json"
const DEFAULT_TOOLCHAIN string = "clang++"

const (
	BINARY VolundBuildType = iota
	STATIC_LIB
	SHARED_LIB
	BUILDER
	NONE
)

const (
	WINDOWS VolundOSType = iota
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
var osType VolundOSType
var toolchain string
var builder BuilderType
