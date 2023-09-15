package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	workDir := "/home/holyzing/Xtalpi/Golang/HolyGo/Python/base"

	err := os.Chdir(workDir)
	if err != nil {
		panic(err)
	}
	/**
	process, err := os.StartProcess("/home/holyzing/Softwares/miniconda3/condabin/conda", []string{"activate crawler-env"}, &os.ProcAttr{
		Dir:   workDir,
		Env:   []string{},
		Files: []*os.File{},
		Sys:   nil,
	})
	if err != nil {
		panic(err)
	}
	wait, err := process.Wait()
	if err != nil {
		panic(err)
	}
	code := wait.ExitCode()
	fmt.Println(code)
	*/

	s, err2 := exec.LookPath("export")
	fmt.Println(s, err2)
	s, err2 = exec.LookPath("sh")
	fmt.Println(s, err2)

	command := "/bin/bash"
	args := []string{"-c", "MYPATH=$PWD:/home/code:$APP_ENV python3 mulitprocess.py"}

	// stdout, err := os.Create("stdout.log")
	// if err != nil {
	// 	panic(err)
	// }
	// stderr, err := os.Create("stderr.log")
	// if err != nil {
	// 	panic(err)
	// }

	// MYPATH=$PWD:/home/code:$APP_ENV

	// te := exec.Command("bash", "-c", "export", "MYPATH=$PWD:/home/code:$APP_ENV")
	te := exec.Command("bash", "-c", "MYPATH=$PWD:/home/code:$APP_ENV python3 --version")
	te.Stdout = os.Stdout
	te.Stderr = os.Stderr

	if err = te.Start(); err != nil {
		panic(err)
	}
	if err = te.Wait(); err != nil {
		fmt.Println(err.Error())
		return
	}

	executor := exec.Command(command, args...)
	executor.Stdout = os.Stdout
	executor.Stderr = os.Stderr

	err = executor.Start()
	if err != nil {
		panic(err)
	}
	err = executor.Wait()
	if err != nil {
		panic(err)
	}
	// all, err := ioutil.ReadAll(stdout)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("--------------- stdout --------------\n", string(all))
	// all, err = ioutil.ReadAll(stderr)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("--------------- stderr --------------\n", string(all))
}
