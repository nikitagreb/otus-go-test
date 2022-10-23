package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	envStr := envStrings(env)

	cmdToExec, cmdArgs := executeCmd(cmd)

	execCmd, err := doRun(cmdToExec, cmdArgs, envStr)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", "error read file", err))
	}

	return execCmd.ProcessState.ExitCode()
}

func doRun(cmdToExec string, cmdArgs []string, envStr []string) (*exec.Cmd, error) {
	execCmd := exec.Command(cmdToExec, cmdArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = os.Environ()
	execCmd.Env = append(execCmd.Env, envStr...)
	return execCmd, execCmd.Run()
}

func executeCmd(cmds []string) (string, []string) {
	return cmds[0], cmds[1:]
}

func envStrings(envs Environment) []string {
	envStr := make([]string, 0, len(envs))
	for key, val := range envs {
		envStr = append(envStr, key+"="+val.Value)
	}
	return envStr
}
