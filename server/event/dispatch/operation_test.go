package dispatch

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"gopkg.in/square/go-jose.v2/jwt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	"github.com/argoproj/argo/server/auth"
	"github.com/argoproj/argo/server/auth/types"
	"github.com/argoproj/argo/util/instanceid"
	"github.com/argoproj/argo/workflow/common"
)

func Test_metaData(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		data := metaData(context.TODO())
		assert.Empty(t, data)
	})
	t.Run("Headers", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
			"x-valid": []string{"true"},
			"ignored": []string{"false"},
		})
		data := metaData(ctx)
		if assert.Len(t, data, 1) {
			assert.Equal(t, []string{"true"}, data["x-valid"])
		}
	})
}

func TestNewOperation(t *testing.T) {
	// set-up
	client := fake.NewSimpleClientset(
		&wfv1.ClusterWorkflowTemplate{
			ObjectMeta: metav1.ObjectMeta{Name: "my-cwft", Labels: map[string]string{common.LabelKeyControllerInstanceID: "my-instanceid"}},
		},
		&wfv1.WorkflowTemplate{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wft", Namespace: "my-ns", Labels: map[string]string{common.LabelKeyControllerInstanceID: "my-instanceid"}},
		},
		&wfv1.WorkflowTemplate{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wft-2", Namespace: "my-ns", Labels: map[string]string{common.LabelKeyControllerInstanceID: "my-instanceid"}},
		},
		&wfv1.WorkflowTemplate{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wft-3", Namespace: "my-ns"},
		},
	)
	ctx := context.WithValue(context.WithValue(context.Background(), auth.WfKey, client), auth.ClaimsKey, &types.Claims{Claims: jwt.Claims{Subject: "my-sub"}})
	recorder := record.NewFakeRecorder(6)

	// act
	operation, err := NewOperation(ctx, instanceid.NewService("my-instanceid"), recorder, []wfv1.WorkflowEventBinding{
		// test a malformed binding
		{
			ObjectMeta: metav1.ObjectMeta{Name: "malformed", Namespace: "my-ns"},
		},
		// test a binding that cannot match
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-0", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "false"},
			},
		},
		// test a binding with a cluster template
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-1", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-cwft", ClusterScope: true},
					Arguments:           &wfv1.Arguments{Parameters: []wfv1.Parameter{{Name: "my-param", ValueFrom: &wfv1.ValueFrom{Event: `"bar"`}}}},
				},
			},
		},
		// test a binding with a namespace template
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-2", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					Arguments:           &wfv1.Arguments{Parameters: []wfv1.Parameter{{Name: "my-param", ValueFrom: &wfv1.ValueFrom{Event: `"bar"`}}}},
				},
			},
		},
		// test a bind with a payload and expression with a map in it
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-3", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "payload.foo.bar == 'baz'"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft-2"},
					Arguments:           &wfv1.Arguments{Parameters: []wfv1.Parameter{{Name: "my-param", ValueFrom: &wfv1.ValueFrom{Event: "payload.foo"}}}},
				},
			},
		},
		// test a binding that errors due to missing workflow template
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-4", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "not-found"},
				},
			},
		},
		// test a binding that errors match due to wrong instance ID on template
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-5", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft-3"},
				},
			},
		},
		// test a binding with an invalid event selector
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-6", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "garbage!!!!!!"},
			},
		},
		// test a binding with a non-bool event selector
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-6", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: `"garbage"`},
			},
		},
		// test a binding with an invalid param expression
		{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-7", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					Arguments:           &wfv1.Arguments{Parameters: []wfv1.Parameter{{Name: "my-param", ValueFrom: &wfv1.ValueFrom{Event: "rubbish!!!"}}}},
				},
			},
		},
	}, "my-ns", "my-discriminator", &wfv1.Item{Value: json.RawMessage(`{"foo": {"bar": "baz"}}`)})
	assert.NoError(t, err)
	operation.Dispatch()

	expectedParamValues := []string{"bar", "bar", `{"bar":"baz"}`}
	// assert
	list, err := client.ArgoprojV1alpha1().Workflows("my-ns").List(metav1.ListOptions{})
	if assert.NoError(t, err) && assert.Len(t, list.Items, 3) {
		for i, wf := range list.Items {
			assert.Equal(t, "my-instanceid", wf.Labels[common.LabelKeyControllerInstanceID])
			assert.Equal(t, "my-sub", wf.Labels[common.LabelKeyCreator])
			assert.Contains(t, wf.Labels, common.LabelKeyWorkflowEventBinding)
			assert.Equal(t, []wfv1.Parameter{{Name: "my-param", Value: wfv1.AnyStringPtr(expectedParamValues[i])}}, wf.Spec.Arguments.Parameters)
		}
	}
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow template expression: unexpected token EOF (1:1)", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to get workflow template: workflowtemplates.argoproj.io \"not-found\" not found", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to validate workflow template instanceid: 'my-wft-3' is not managed by the current Argo Server", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow template expression: unexpected token Operator(\"!\") (1:8)\n | garbage!!!!!!\n | .......^", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: malformed workflow template expression: did not evaluate to boolean", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow template parameter \"my-param\" expression: unexpected token Operator(\"!\") (1:8)\n | rubbish!!!\n | .......^", <-recorder.Events)
}

func Test_populateWorkflowMetadata(t *testing.T) {
	// set-up
	client := fake.NewSimpleClientset(
		&wfv1.WorkflowTemplate{
			ObjectMeta: metav1.ObjectMeta{Name: "my-wft", Namespace: "my-ns", Labels: map[string]string{common.LabelKeyControllerInstanceID: "my-instanceid"}},
		},
	)
	ctx := context.WithValue(context.WithValue(context.Background(), auth.WfKey, client), auth.ClaimsKey, &types.Claims{Claims: jwt.Claims{Subject: "my-sub"}})
	recorder := record.NewFakeRecorder(10)

	// act
	operation, err := NewOperation(ctx, instanceid.NewService("my-instanceid"), recorder, []wfv1.WorkflowEventBinding{
		{
			// No name specified
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-1", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
				},
			},
		},
		{
			// Fixed name, label and annotation given
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-2", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta: metav1.ObjectMeta{
						Name:        "\"my-wfeb-2\"",
						Labels:      map[string]string{"aLabel": "\"someValue\""},
						Annotations: map[string]string{"anAnnotation": "\"otherValue\""},
					},
				},
			},
		},
		{
			// Name, label and annotations using expressions
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-3", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta: metav1.ObjectMeta{
						Name:        "\"my-wfeb-\" + payload.foo.bar",
						Labels:      map[string]string{"aLabel": "payload.list[0]"},
						Annotations: map[string]string{"anAnnotation": "payload.list[1]"},
					},
				},
			},
		},
		{
			// Invalid name expression
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-4", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.......foo[.numeric]"},
				},
			},
		},
		{
			// Invalid label expression
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-4", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"invalidLabel": "foo...bar"}},
				},
			},
		},
		{
			// Invalid annotation expression
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-4", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{"invalidAnnotation": "foo.[..]bar"}},
				},
			},
		},
		{
			// Name expression evaluates to invalid type (float)
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-5", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.foo.numeric"},
				},
			},
		},
		{
			// Name expression evaluates to invalid type (boolean)
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-6", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.foo.bool"},
				},
			},
		},
		{
			// Name expression evaluates to invalid type (map)
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-7", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.foo"},
				},
			},
		},
		{
			// Name expression evaluates to invalid type (list)
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-8", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.list"},
				},
			},
		},
		{
			// Name expression evaluates to non existent value (nil)
			ObjectMeta: metav1.ObjectMeta{Name: "my-wfeb-9", Namespace: "my-ns"},
			Spec: wfv1.WorkflowEventBindingSpec{
				Event: wfv1.Event{Selector: "true"},
				Submit: &wfv1.Submit{
					WorkflowTemplateRef: wfv1.WorkflowTemplateRef{Name: "my-wft"},
					ObjectMeta:          metav1.ObjectMeta{Name: "payload.nothing"},
				},
			},
		},
	}, "my-ns", "my-discriminator",
		&wfv1.Item{Value: json.RawMessage(`{"foo": {"bar": "baz", "numeric": 8675309, "bool": true}, "list": ["one", "two"]}`)})

	assert.NoError(t, err)
	operation.Dispatch()

	list, err := client.ArgoprojV1alpha1().Workflows("my-ns").List(metav1.ListOptions{})

	assert.NoError(t, err)
	assert.Len(t, list.Items, 3)

	assert.Contains(t, list.Items[0].Name, "my-wft")

	assert.Equal(t, list.Items[1].Name, "my-wfeb-2")
	assert.Contains(t, list.Items[1].Labels, "aLabel")
	assert.Equal(t, "someValue", list.Items[1].Labels["aLabel"])
	assert.Contains(t, list.Items[1].Annotations, "anAnnotation")
	assert.Equal(t, "otherValue", list.Items[1].Annotations["anAnnotation"])

	assert.Equal(t, list.Items[2].Name, "my-wfeb-baz")
	assert.Contains(t, list.Items[2].Labels, "aLabel")
	assert.Equal(t, "one", list.Items[2].Labels["aLabel"])
	assert.Contains(t, list.Items[2].Annotations, "anAnnotation")
	assert.Equal(t, "two", list.Items[2].Annotations["anAnnotation"])

	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow name expression: unexpected token Operator(\"..\") (1:10)\n | payload.......foo[.numeric]\n | .........^", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow label \"invalidLabel\" expression: cannot use pointer accessor outside closure (1:6)\n | foo...bar\n | .....^", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: failed to evaluate workflow annotation \"invalidAnnotation\" expression: expected name (1:6)\n | foo.[..]bar\n | .....^", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: workflow name expression must evaluate to a string, not a float64", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: workflow name expression must evaluate to a string, not a bool", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: workflow name expression must evaluate to a string, not a map[string]interface {}", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: workflow name expression must evaluate to a string, not a []interface {}", <-recorder.Events)
	assert.Equal(t, "Warning WorkflowEventBindingError failed to dispatch event: workflow name expression must evaluate to a string, not a <nil>", <-recorder.Events)
}

func Test_expressionEnvironment(t *testing.T) {
	env, err := expressionEnvironment(context.TODO(), "my-ns", "my-d", &wfv1.Item{Value: []byte(`{"foo":"bar"}`)})
	if assert.NoError(t, err) {
		assert.Equal(t, "my-ns", env["namespace"])
		assert.Equal(t, "my-d", env["discriminator"])
		assert.Contains(t, env, "metadata")
		assert.Equal(t, map[string]interface{}{"foo": "bar"}, env["payload"], "make sure we parse an object as a map")
	}
}
