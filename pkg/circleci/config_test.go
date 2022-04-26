package circleci

import (
	_ "embed"
	"github.com/go-test/deep"
	"testing"
)

//go:embed "tests/config.yml"
var s string

func Test_unmarshalConfig(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "should marshall a correctly formatted circleci configuration",
			args: args{
				b: []byte(s),
			},
			want: &Config{
				Version: "2.1",
				Parameters: map[string]Parameter{
					"src-repo-url": {
						Name:        "src-repo-url",
						Description: "The repository url",
						Type:        "string",
						Default:     "https://github.com/esnet/iperf.git",
					},
					"branch-name": {
						Name:        "branch-name",
						Description: "The branch name",
						Type:        "string",
						Default:     "3.8.1",
					},
					"common-build-params": {
						Name:        "common-build-params",
						Description: "The common build params",
						Type:        "string",
						Default:     "--disable-shared --disable-static",
					},
				},
				Jobs: map[string]Job{
					"build-linux": {
						Name: "build-linux",
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-linux",
							},
						},
					},
					"build-windows": {
						Name: "build-windows",
						Environment: map[string]string{
							"FOO": "BAR",
						},
						Parallelism: 2,
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-cygwin64",
							},
						},
					},
					"build-macos": {
						Name: "build-macos",
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-macos",
							},
						},
					},
					"test-linux": {
						Name: "test-linux",
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-linux",
							},
						},
					},
					"test-windows": {
						Name: "test-windows",
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-cygwin64",
							},
						},
					},
					"test-macos": {
						Name:        "test-macos",
						Parallelism: 3,
						Parameters: map[string]Parameter{
							"label": {
								Name:    "label",
								Type:    "string",
								Default: "iperf3-macos",
							},
						},
					},
					"release": {
						Name: "release",
					},
				},
				Workflows: Workflows{
					Version: "2",
					Workflows: map[string]Workflow{
						"build-test-release": {
							Name: "build-test-release",
							Jobs: []WorkflowJob{
								{
									Name: "build-linux",
								},
								{
									Name: "build-windows",
								},
								{
									Name: "build-macos",
								},
								{
									Name: "test-linux",
									Requires: []string{
										"build-linux",
									},
								},
								{
									Name: "test-windows",
									Requires: []string{
										"build-windows",
									},
								},
								{
									Name: "macos",
									Requires: []string{
										"build-macos",
									},
								},
								{
									Name: "release",
									Requires: []string{
										"test-linux",
										"test-windows",
										"test-macos",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalConfig(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("unmarshalConfig() diff = %v", diff)
			}
		})
	}
}
