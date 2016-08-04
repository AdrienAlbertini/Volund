package main

import (
)

type CmdFunc func()

func cmdRebuild() {
  rebuildAll = true
}

func cmdFatalError() {
  fatalError = true
}

func handleCommands(commands []string) {

  var commandsMap map[string]CmdFunc
  var commandsShortMap map[string]string
  commandsMap = make(map[string]CmdFunc)
  commandsShortMap = make(map[string]string)

  commandsMap["-rebuild"] = cmdRebuild
  commandsMap["-fatalerror"] = cmdFatalError

  commandsShortMap["-r"] = "-rebuild"

	for _, cmd := range commands {
    cmdToExec, ok := commandsMap[cmd]
    cmdShort, okShort := commandsShortMap[cmd]

    if ok == false {
      if okShort {
        cmdToExec = commandsMap[cmdShort]
        cmdToExec()
      }
    } else {
      cmdToExec()
    }
	}
}
