package circleci

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config describes a circleci configuration. Get doumentation about it here
// https://circleci.com/docs/2.0/configuration-reference/#version
type Config struct {
	// The version field is intended to be used in order to issue warnings for
	// deprecation or breaking changes.
	Version string `yaml:"version"`

	// Pipeline parameters declared for use in the configuration.
	Parameters map[string]Parameter `yaml:"parameters"`

	// A Workflow is comprised of one or more uniquely named jobs. Jobs are
	// specified in the jobs map. The name of the job is the key in the map, and
	// the value is a map describing the job.
	Jobs map[string]Job `yaml:"jobs"`

	// Used for orchestrating all jobs. Each workflow consists of the workflow
	// name as a key and a map as a value.
	Workflows Workflows `yaml:"workflows"`
}

// fillNames set name in objects issued from a map (with the map key)
func (c *Config) fillNames() {
	for k, v := range c.Parameters {
		v.Name = k
		c.Parameters[k] = v
	}
	for k, v := range c.Jobs {
		v.Name = k
		c.Jobs[k] = v

		for k2, v2 := range v.Parameters {
			v2.Name = k2
			c.Jobs[k].Parameters[k2] = v2
		}
	}
	for k, v := range c.Workflows.Workflows {
		v.Name = k
		c.Workflows.Workflows[k] = v
	}
}

// Pipeline parameter declared for use in the configuration.
type Parameter struct {
	// The parameter name.
	Name string

	// Optional. Used to generate documentation for orbs.
	Description string `yaml:"description"`

	// The parameter types.
	Type string `yaml:"type"`

	// The default value for the parameter. If not present, the parameter is
	// implied to be required.
	Default any `yaml:"default"`
}

// Pipeline job declared for use in the configuration.
type Job struct {
	// The job name.
	Name string

	// A map of environment variable names and values.
	Environment map[string]string `yaml:"environment"`

	// Number of parallel instances of this job to run.
	Parallelism int `yaml:"parallelism"`

	// Parameters for making a job explicitly configurable in a workflow.
	Parameters map[string]Parameter `yaml:"parameters"`
}

// Pipeline workflows declared for use in the configuration.
type Workflows struct {
	// The Workflows version field is used to issue warnings for deprecation or
	// breaking changes.
	Version string `yaml:"version"`

	// Used for orchestrating all jobs. Each workflow consists of the workflow
	// name as a key and a map as a value.
	Workflows map[string]Workflow `yaml:",inline"`
}

// Pipeline workflow declared for use in the configuration.
type Workflow struct {
	// The workflow name.
	Name string
	// The workflow jobs.
	Jobs []WorkflowJob `yaml:"jobs"`
}

// A job declaration inside a workflow.
type WorkflowJob struct {
	// The job name.
	Name string
	// A list of jobs that must succeed for the job to start.
	Requires []string
}

func (j *WorkflowJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var job map[string]struct {
		Name     string   `yaml:"name"`
		Requires []string `yaml:"requires"`
	}

	if err := unmarshal(&job); err == nil {
		for k, v := range job {
			if v.Name != "" {
				j.Name = v.Name
			} else {
				j.Name = k
			}
			j.Requires = v.Requires
			return nil
		}
	}

	var name string
	if err := unmarshal(&name); err == nil {
		j.Name = name
	}
	return nil
}

// NewConfig returns a newly instanciated config initialized from config file
// path.
func NewConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(fmt.Errorf("Read file error: %w", err))
		return nil, err
	}

	return unmarshalConfig(f)
}

func unmarshalConfig(b []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		log.Errorf("Cannot unmarshal yaml: %v", err)
		return nil, err
	}
	config.fillNames()
	return config, nil
}
