package language_provider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

var _ Provider = &Go{}

type Go struct {
	fileTemplate            string
	fileTemplateWithPackage string
	fileTemplateTest        string
}

func NewGoProvider() *Go {
	return &Go{
		fileTemplate:            "language_provider/templates/golang.tmpl",
		fileTemplateWithPackage: "language_provider/templates/golang_with_package.tmpl",
		fileTemplateTest:        "language_provider/templates/golang_test.tmpl",
	}
}

func (g *Go) GetTemplate(input GetTemplateInput) (*GetTemplateOutput, error) {
	tmpl, err := template.ParseFiles(g.fileTemplate)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, input)
	if err != nil {
		return nil, err
	}
	return &GetTemplateOutput{
		Content: b.String(),
	}, nil
}

func (g *Go) RunTestcase(input RunTestcaseInput) (*RunTestcaseOutput, error) {
	id := uuid.New().String()
	err := os.Mkdir(fmt.Sprintf("storage/%s", id), 0777)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(fmt.Sprintf("storage/%s/%s.go", id, input.MethodName))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fTest, err := os.Create(fmt.Sprintf("storage/%s/%s_test.go", id, input.MethodName))
	if err != nil {
		return nil, err
	}
	defer fTest.Close()

	tmpl, err := template.ParseFiles(g.fileTemplateWithPackage)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(f, input)
	if err != nil {
		return nil, err
	}
	tmpl, err = template.ParseFiles(g.fileTemplateTest)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(fTest, input)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	mountDir := fmt.Sprintf("%s/storage/%s", currentDir, id)

	resp, err := cli.ContainerCreate(context.TODO(), &container.Config{
		Image: "picket-go-test",
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/app", mountDir)},
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}
	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}
	pass, fail, err := g.readOutput(id)
	if err != nil {
		return nil, err
	}
	log.Println("pass", pass)
	log.Println("fail", fail)
	return nil, nil
}

func (g *Go) waitForFile(filename string, timeout time.Duration) error {
	start := time.Now()
	ticker := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("Timeout waiting for file %s", filename)
		case <-ticker.C:
			if _, err := os.Stat(filename); err == nil {
				return nil
			}
			if time.Since(start) > timeout {
				return fmt.Errorf("Timeout waiting for file %s", filename)
			}
		}
	}

}

func (g *Go) readOutput(id string) ([]string, []string, error) {
	path := fmt.Sprintf("storage/%s/output.txt", id)
	err := g.waitForFile(path, time.Second*10)
	if err != nil {
		return nil, nil, err
	}
	time.Sleep(3 * time.Second)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	log := string(content)
	re := regexp.MustCompile(`--- FAIL: (\S+)`)

	matches := re.FindAllStringSubmatch(log, -1)
	failedTests := []string{}

	for _, match := range matches {
		failedTests = append(failedTests, match[1])
	}

	fmt.Printf("Failed tests: %v\n", failedTests)

	re = regexp.MustCompile(`--- PASS: (\S+)`)
	matches = re.FindAllStringSubmatch(log, -1)
	passedTests := []string{}

	for _, match := range matches {
		passedTests = append(passedTests, match[1])
	}

	fmt.Printf("Passed tests: %v\n", passedTests)

	return passedTests, failedTests, nil
}
