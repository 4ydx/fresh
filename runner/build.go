package runner

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
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

	target := description[1] + "/" + filepath.Base(project) + ".js"
	buildLog("Building %s", target)

	running := []string{"go", "gopherjs"}

	for _, run := range running {
		var cmd *exec.Cmd
		switch run {
		case "go":
			cmd = exec.Command(run, "build", project)
		case "gopherjs":
			cmd = exec.Command(run, "build", project, "-o", target)
		}
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
			mainLog("Build Failed: \n %s \n %s", cmd.String(), string(errBuf))
			createBuildErrorsLog(string(errBuf))

			ok.Lock()
			ok.Val = false
			ok.Unlock()
		}
	}
}

func build(started bool) bool {
	buildLog("Now building path %s root %s", buildPath(), root())

	cmd := exec.Command("make")
	msg, ok := runBuild(cmd)
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
	mainLog("Build Waiting")

	wg.Wait()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://localhost:10000/refresh")
	if err != nil {
		mainLog("Refresh Failed: %s", err)
		buildOk.Val = false
	}
	defer resp.Body.Close()

	if ok, err := ioutil.ReadAll(resp.Body); err != nil || string(ok) != "OK" {
		if err != nil {
			mainLog("Refresh Failed: %s", err)
		} else {
			mainLog("Refresh Failed: %s", string(ok))
		}
		buildOk.Val = false
	}

	return buildOk.Val
}
