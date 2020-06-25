package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/mocks"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
)

const workflow string =`
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  name: hello-world
  namespace: default
spec:
  entrypoint: whalesay
  templates:
  - name: whalesay
    container:
      image: docker/whalesay:latest
      command: [cowsay]
      args: ["hello world"]
`

func TestSubmitFromResource(t *testing.T) {
	client := mocks.Client{}
	wfClient := mocks.WorkflowServiceClient{}
	wfClient.On("SubmitFrom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&wfv1.Workflow{}, nil)
	client.On("NewWorkflowServiceClient").Return(&wfClient)
	CLIOpt.client = &client
	CLIOpt.ctx = context.TODO()
	output := CaptureOutput(func(){submitWorkflowFromResource("workflowtemplatetest",&wfv1.SubmitOpts{},&cliSubmitOpts{})})
	assert.Contains(t, output, "Created:")
}

func TestSubmitWorkflows(t *testing.T) {
	client := mocks.Client{}
	wfClient := mocks.WorkflowServiceClient{}
	var wf wfv1.Workflow
	wfClient.On("CreateWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&wf, nil)
	client.On("NewWorkflowServiceClient").Return(&wfClient)
	CLIOpt.client = &client
	CLIOpt.ctx = context.TODO()

	err:=yaml.Unmarshal([]byte(workflow), &wf)
	assert.NoError(t, err)
	workflows := []wfv1.Workflow{wf}
	output := CaptureOutput(func(){submitWorkflows(workflows,&wfv1.SubmitOpts{},&cliSubmitOpts{})})
	fmt.Println(output)
	assert.Contains(t, output, "Created:")
}

func CaptureOutput(f func()) string {
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr
	var buf bytes.Buffer
	log.SetOutput(&buf)
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr
	return string(out)+buf.String()
}