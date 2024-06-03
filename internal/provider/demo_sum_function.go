package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &DemoSumFunction{}

type DemoSumFunction struct{}

func NewDemoSumFunction() function.Function {
	return &DemoSumFunction{}
}

func (f *DemoSumFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "demo_sum"
}

func (f *DemoSumFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Demo sum function.",
		Description: "It adds 2 numbers and return the sum.",
		Parameters: []function.Parameter{
			function.Int64Parameter{
				Name:        "num1",
				Description: "First number.",
			},
			function.Int64Parameter{
				Name:        "num2",
				Description: "Second number.",
			},
		},
		Return: function.Int64Return{},
	}
}

func (f *DemoSumFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var num1 int64
	var num2 int64
	var sum int64

	// Read Terraform argument data into the variables
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &num1, &num2))

	sum = num1 + num2

	// Set the result
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, sum))
}
