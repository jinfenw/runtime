// Copyright (c) 2014,2015,2016 Docker, Inc.
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	vc "github.com/containers/virtcontainers"
	"github.com/containers/virtcontainers/pkg/oci"
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete any resources held by one or more containers",
	ArgsUsage: `<container-id> [container-id...]

   <container-id> is the name for the instance of the container.

EXAMPLE:
   If the container id is "ubuntu01" and ` + name + ` list currently shows the
   status of "ubuntu01" as "stopped" the following will delete resources held
   for "ubuntu01" removing "ubuntu01" from the ` + name + ` list of containers:

       # ` + name + ` delete ubuntu01`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Forcibly deletes the container if it is still running (uses SIGKILL)",
		},
	},
	Action: func(context *cli.Context) error {
		args := context.Args()
		if args.Present() == false {
			return fmt.Errorf("Missing container ID, should at least provide one")
		}

		force := context.Bool("force")
		for _, cID := range []string(args) {
			if err := delete(cID, force); err != nil {
				return err
			}
		}

		return nil
	},
}

func delete(containerID string, force bool) error {
	// Checks the MUST and MUST NOT from OCI runtime specification
	if err := validContainer(containerID); err != nil {
		return err
	}

	if force == false {
		podStatus, err := vc.StatusPod(containerID)
		if err != nil {
			return err
		}

		state, err := oci.StatusToOCIState(podStatus)
		if err != nil {
			return err
		}

		running, err := processRunning(state.Pid)
		if err != nil {
			return err
		}

		if running == true {
			return fmt.Errorf("Container still running, should be stopped")
		}
	}

	pod, err := vc.StopPod(containerID)
	if err != nil {
		return err
	}

	// Retrieve OCI spec configuration.
	ociSpec, err := oci.PodToOCIConfig(*pod)
	if err != nil {
		return err
	}

	if _, err := vc.DeletePod(containerID); err != nil {
		return err
	}

	// In order to prevent any file descriptor leak related to cgroups files
	// that have been previously created, we have to remove them before this
	// function returns.
	cgroupsPathList, err := processCgroupsPath(ociSpec)
	if err != nil {
		return err
	}

	if err := removeCgroupsPath(cgroupsPathList); err != nil {
		return err
	}

	return nil
}

func removeCgroupsPath(cgroupsPathList []string) error {
	if len(cgroupsPathList) == 0 {
		ccLog.Info("Cgroups files not removed because cgroupsPath was empty")
		return nil
	}

	for _, cgroupsPath := range cgroupsPathList {
		if err := os.RemoveAll(cgroupsPath); err != nil {
			return err
		}
	}

	return nil
}
