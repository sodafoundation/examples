// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ontap

import (
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/golang/glog"
	"github.com/sodafoundation/dock/pkg/utils/exec"
)

type MakeConfig struct {
	User     string
	Server   string
	Password string
	Port     string
	Timeout  time.Duration
}

func Executer() *easyssh.MakeConfig {
	ssh := &easyssh.MakeConfig{
		User:     username,
		Server:   defaultTgtBindIp,
		Password: password,
		Port:     port,
		Timeout:  timeoutForssh * time.Second,
	}
	return ssh
}

type Cli struct {
	// Command executer
	BaseExecuter exec.Executer
	// Command Root executer
	RootExecuter exec.Executer
}

func login() error {
	stdout, stderr, done, err := Executer().Run("uname", timeoutForssh*time.Second)
	if err != nil {
		glog.Errorf("unable to establish connection, stderr:%v", stderr)
		return err
	}
	glog.Infof("connection established. stdout:%v done:%v", stdout, done)
	return nil
}

func NewCli() (*Cli, error) {
	return &Cli{
		BaseExecuter: exec.NewBaseExecuter(),
		RootExecuter: exec.NewRootExecuter(),
	}, nil
}

func (c *Cli) execute(cmd ...string) (string, error) {
	return c.RootExecuter.Run(cmd[0], cmd[1:]...)
}

// this is function to check replication already exists
func (c *Cli) GetSnapmirror() error {
	snapmirrorCmd := "snapmirror" + " " + "show"

	glog.Info("the executed command:%v", snapmirrorCmd)
	stdout, stderr, done, err := Executer().Run(snapmirrorCmd, timeoutForssh*time.Second)
	if stderr != "" {
		glog.Errorf("failed to check replication status. stderr:%v", stderr)
		return err
	}
	glog.Infof("stdout:%v, done:%v", stdout, done)
	return nil
}

// this is function to creat replication between array
func (c *Cli) CreateSnapmirror(sourcevol, destvol, srcvserver, dstvserver string) error {
	// to create replication between array or volume. need to run snapmirror command
	snapmirrorCmd := "snapmirror" + " " + "create" + " -source-volume" + " " + sourcevol + " -destination-volume" +
		" " + destvol + " -source-vserver" + " " + srcvserver + " -destination-vserver" + " " + dstvserver

	glog.Info("the executed command:%v", snapmirrorCmd)

	stdout, stderr, done, err := Executer().Run(snapmirrorCmd, timeoutForssh*time.Second)
	if stderr != "" {
		glog.Errorf("failed to create replication relationship. stderr:%v", stderr)
		return err
	}
	glog.Infof("replication is created successfully. stdout:%v, done:%v", stdout, done)
	return nil
}

// this is function to delete replication between array
func (c *Cli) DeleteSnapmirror(sourcevol, destvol, srcvserver, dstvserver string) error {
	// to delete replication between array or volume. need to run snapmirror command
	snapmirrorCmd := "snapmirror" + " " + "delete" + " -source-volume" + " " + sourcevol + " -destination-volume" +
		" " + destvol + " -source-vserver" + " " + srcvserver + " -destination-vserver" + " " + dstvserver

	glog.Info("the executed command:%v", snapmirrorCmd)

	stdout, stderr, done, err := Executer().Run(snapmirrorCmd, timeoutForssh*time.Second)
	if stderr != "" {
		glog.Errorf("failed to delete replication relationship. stderr:%v", stderr)
		return err
	}
	glog.Infof("replication is delete successfully. stdout:%v, done:%v", stdout, done)
	return nil
}
