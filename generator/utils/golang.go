package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type GolangCodeFormatter struct {
}

func (f GolangCodeFormatter) Format(input string) (*string, error) {
	filePath := f.randomFilePath()
	if err := f.writeContentsToFile(filePath, input); err != nil {
		return nil, fmt.Errorf("writing contents to %q: %+v", filePath, err)
	}
	defer f.deleteFileContents(filePath)

	if err := f.runGoFmt(filePath); err != nil {
		return nil, fmt.Errorf("running gofmt on %q: %+v", filePath, err)
	}

	contents, err := f.readFileContents(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading contents from %q: %+v", filePath, err)
	}

	return contents, nil
}

func (f GolangCodeFormatter) randomFilePath() string {
	time := time.Now().Unix()
	return fmt.Sprintf("%stemp-%d.go", os.TempDir(), time)
}

func (f GolangCodeFormatter) runGoFmt(filePath string) error {
	cmd := exec.Command("gofmt", "-w", filePath)
	// intentionally not using these errors since the exit codes are kinda uninteresting
	cmd.Start()
	cmd.Wait()
	return nil
}

func (f GolangCodeFormatter) deleteFileContents(filePath string) error {
	return os.Remove(filePath)
}

func (f GolangCodeFormatter) readFileContents(filePath string) (*string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contents := string(data)
	return &contents, nil
}

func (GolangCodeFormatter) writeContentsToFile(filePath, contents string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(contents); err != nil {
		return err
	}

	return nil
}
