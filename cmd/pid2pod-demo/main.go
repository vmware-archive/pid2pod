// Copyright 2017 by the contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main runs "pid2pod-demo", which lists Kubernetes pod metadata for all proceses.
package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/go-ps"

	"github.com/heptiolabs/pid2pod"
)

func main() {
	processes, err := ps.Processes()
	if err != nil {
		log.Fatalf("could not list processes: %v", err)
	}
	for _, proc := range processes {
		pid := proc.Pid()
		id, err := pid2pod.LookupPod(pid)
		if err != nil {
			log.Fatalf("could not get ID of process %d: %v", pid, err)
		}
		if id != nil {
			fmt.Printf("PID %d (%s): %+#v\n", pid, proc.Executable(), id)
		}

	}
}
