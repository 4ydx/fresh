package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

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

func build() (string, bool) {
	buildLog("Building...")

	cmd := exec.Command("go", "build", "-o", buildPath(), root())
	msg, ok := runBuild(cmd)
	if !ok {
		return msg, ok
	}
	cmd = exec.Command("make")
	msg, ok = runBuild(cmd)
	if !ok {
		return msg, ok
	}
	_, err := os.Stat("static/js/games/dnd/build.sh")
	if err == nil {
		buildLog("Building js chdir...")

		err = os.Chdir("static/js/games/dnd/")
		if err != nil {
			fatal(err)
		}

		buildLog("Building js...")
		cmd := exec.Command("./build.sh")
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

		buildLog("Building js retreat...")
		err = os.Chdir("../../../../")
		if err != nil {
			fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			return string(errBuf), false
		}
	}
	return "", true
}
