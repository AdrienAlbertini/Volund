package main

type ObjectJSON struct {
	Builder   BuilderJSON   `json:"Builder"`
	Binary    BinaryJSON    `json:"Binary"`
	StaticLib StaticLibJSON `json:"StaticLib"`
	Plugin    PluginJSON    `json:"Plugin"`
}

type BuilderJSON struct {
	Os            string   `json:"OS"`
	Toolchain     string   `json:"toolchain"`
	OutBinary     string   `json:"outBinary"`
	Binaries      []string `json:"binaries"`
	StaticLibs    []string `json:"staticLibs"`
	Plugins       []string `json:"plugins"`
	PluginsFolder string   `json:"outPluginsFolder"`
	OutFolder     string   `json:"outFolder"`
	CompilerFlags []string `json:"compilerFlags"`
	FullStatic    bool     `json:"fullStatic"`
}

type BinaryJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	Plugins        []string `json:"plugins"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}

type StaticLibJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}

type PluginJSON struct {
	Name           string   `json:"name"`
	StaticLibs     []string `json:"staticLibs"`
	SrcFolders     []string `json:"srcFolders"`
	SrcExtension   string   `json:"srcExtension"`
	HeadersFolders []string `json:"headersFolders"`
	ExternIncludes []string `json:"externIncludes"`
	ExternLibs     []string `json:"externLibs"`
	OutFolder      string   `json:"outFolder"`
	CompilerFlags  []string `json:"compilerFlags"`
}
