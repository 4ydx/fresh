package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type BuildOk struct {
	*sync.RWMutex
	Val bool
}

func runBuild(cmd *exec.Cmd) (string, bool) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		fatal(err)
	}
	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}
	return "", true
}

func buildHelper(wg *sync.WaitGroup, ok *BuildOk, project string) {
	defer wg.Done()

	description := strings.Split(project, " ")
	project = description[0]

	buildLog("Building %s", project)

	relativePath := description[1]
	base := filepath.Base(project)
	target := relativePath + "/" + base + ".js"

	cmd := exec.Command("gopherjs", "build", project, "-o", target)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		fatal(err)
	}
	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		mainLog("Build Failed: \n %s", string(errBuf))
		createBuildErrorsLog(string(errBuf))

		ok.Lock()
		ok.Val = false
		ok.Unlock()
	}
}

func build(started bool) bool {
	buildLog("Building...")

	cmd := exec.Command("go", "build", "-o", buildPath(), root())
	msg, ok := runBuild(cmd)
	if !ok {
		mainLog("Build Failed: \n %s", msg)
		createBuildErrorsLog(msg)
		return ok
	}
	cmd = exec.Command("make")
	msg, ok = runBuild(cmd)
	if !ok {
		mainLog("Build Failed: \n %s", msg)
		createBuildErrorsLog(msg)
		return ok
	}

	b, err := ioutil.ReadFile("build_projects")
	if err != nil {
		fatal(err)
	}
	projects := strings.Split(string(b), "\n")

	wg := &sync.WaitGroup{}
	buildOk := &BuildOk{Val: true, RWMutex: &sync.RWMutex{}}
	for _, project := range projects {
		if strings.TrimSpace(project) == "" {
			continue
		}
		wg.Add(1)
		go buildHelper(wg, buildOk, project)
	}
	wg.Wait()

	return buildOk.Val
}
