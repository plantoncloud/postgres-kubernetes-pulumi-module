package postgres

import (
	"context"
	postgresblueprintoutputs "github.com/plantoncloud/postgres-kubernetes-pulumi-blueprint/pkg/postgres/outputs"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	code2cloudv1deploypgk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/postgreskubernetes/model"
	code2cloudv1deploypgk8sstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/postgreskubernetes/stack/model"
	"github.com/plantoncloud/postgres-kubernetes-pulumi-blueprint/pkg/postgres/network/hostname"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *code2cloudv1deploypgk8sstackmodel.PostgresKubernetesStackInput) (*code2cloudv1deploypgk8smodel.PostgresKubernetesStatus, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *code2cloudv1deploypgk8sstackmodel.PostgresKubernetesStackInput) *code2cloudv1deploypgk8smodel.PostgresKubernetesStatus {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &code2cloudv1deploypgk8smodel.PostgresKubernetesStatus{}
	}
	return &code2cloudv1deploypgk8smodel.PostgresKubernetesStatus{
		Namespace:               backend.GetVal(stackOutput, postgresblueprintoutputs.GetNamespaceNameOutputName()),
		Service:                 backend.GetVal(stackOutput, postgresblueprintoutputs.GetKubeServiceNameOutputName()),
		PortForwardCommand:      backend.GetVal(stackOutput, postgresblueprintoutputs.GetKubePortForwardCommandOutputName()),
		KubeEndpoint:            backend.GetVal(stackOutput, postgresblueprintoutputs.GetKubeEndpointOutputName()),
		ExternalClusterHostname: backend.GetVal(stackOutput, hostname.GetExternalClusterHostnameOutputName()),
		InternalClusterHostname: backend.GetVal(stackOutput, hostname.GetInternalClusterHostnameOutputName()),
	}
}
