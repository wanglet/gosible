package gosible

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/go-basic/uuid"
)

func NewJob(params map[string]interface{}) (*Job, error) {
	job := &Job{}
	job.Environments = Environments{ForceColor: true}

	res, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(res, job)
	if job.UUID == "" {
		job.UUID = uuid.New()
	}

	return job, nil
}

type Job struct {
	AnsibleOptions             `json:",inline"`
	Playbook                   string `json:"playbook"`
	Pattern                    string `json:"pattern"`
	PlaybookOptions            `json:",inline"`
	ConnectionOptions          `json:",inline"`
	PrivilegeEscalationOptions `json:",inline"`
	Environments               `json:"environments"`

	Ctx          context.Context
	cmd          *exec.Cmd
	UUID         string
	EventHandler io.Writer
}

func (j *Job) validate() error {
	if j.Pattern != "" && j.Playbook != "" {
		return errors.New("one of playbook and pattern is allowed")
	}
	if j.Pattern == "" && j.Playbook == "" {
		return errors.New("one of playbook and pattern is required")
	}
	if j.Pattern != "" {
		if j.AnsibleOptions.ModuleName == "" {
			return errors.New("moduleName is required when running ansible")
		}
	}
	return nil
}

func (j *Job) prepare() error {
	var commandName string
	args := []string{}
	if j.Pattern != "" {
		commandName = "ansible"
		args = append(args, j.Pattern)
		args = append(args, marshal(j.AnsibleOptions, "argument", " ")...)
	}
	if j.Playbook != "" {
		commandName = "ansible-playbook"
		args = append(args, j.Playbook)
		args = append(args, marshal(j.PlaybookOptions, "argument", " ")...)
	}
	args = append(args, marshal(j.ConnectionOptions, "argument", " ")...)
	args = append(args, marshal(j.PrivilegeEscalationOptions, "argument", " ")...)

	if j.Ctx != nil {
		j.cmd = exec.CommandContext(j.Ctx, commandName, args...)
	} else {
		j.cmd = exec.Command(commandName, args...)
	}

	j.cmd.Env = append(j.cmd.Env, marshal(j.Environments, "environment", "=")...)
	//a.cmd.Env = append(a.cmd.Env, os.Environ()...)

	if j.EventHandler != nil {
		j.cmd.Stdout = j.EventHandler
	} else {
		j.cmd.Stdout = EventHandler{
			UUID: j.UUID,
		}
		j.cmd.Stderr = j.cmd.Stdout
	}
	return nil
}

func (j *Job) Run() error {
	if err := j.Start(); err != nil {
		return err
	}
	return j.cmd.Wait()
}

func (j *Job) Kill() error {
	if j.cmd == nil {
		return errors.New("no cmd")
	}
	return j.cmd.Process.Kill()
}

func (j *Job) Start() error {
	var err error
	if err = j.validate(); err != nil {
		return err
	}
	if err = j.prepare(); err != nil {
		return err
	}
	fmt.Println(j.cmd)
	return j.cmd.Start()
}

// func (j *Job) Wait() error {
// 	if j.cmd == nil {
// 		return errors.New("no cmd")
// 	}
// 	return j.cmd.Wait()
// }
