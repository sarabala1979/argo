package e2e

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/test/e2e/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

type WorkflowTemplateSuite struct {
	fixtures.E2ESuite
}

func (w *WorkflowTemplateSuite) TestNestedWorkflowTemplate() {
	w.Given().WorkflowTemplate("@smoke/workflow-template-whalesay-template.yaml").
		WorkflowTemplate("@testdata/workflow-template-nested-template.yaml").
		Workflow(`apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: workflow-template-nested-
spec:
  entrypoint: whalesay
  templates:
  - name: whalesay
    inputs:
      parameters:
      - name: message
        value: hello from nested
    templateRef:
      name: workflow-template-nested-template
      template: whalesay-template
`).When().
		CreateWorkflowTemplates().
		SubmitWorkflow().
		WaitForWorkflow(15 * time.Second).
		Then().
		Expect(func(t *testing.T, metadata *v1.ObjectMeta, status *v1alpha1.WorkflowStatus) {
			assert.Equal(t, status.Phase, v1alpha1.NodeSucceeded)
		})

}

func TestWorkflowTemplateSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTemplateSuite))
}
