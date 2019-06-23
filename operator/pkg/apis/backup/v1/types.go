/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Backup is a specification for a Backup resource
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupSpec `json:"spec"`
}

// BackupSpec is the spec for a Backup resource
type BackupSpec struct {
	Name     string     `json:"name"`
	Cron     string     `json:"cron"`
	Type     string     `json:"type"`
	Input    Input      `json:"input"`
	Compress Compress   `json:"compress"`
	Output   []Output   `json:"output"`
	Notifier []Notifier `json:"notifier"`
}

func (b Backup) ConstrucArguments() []string {
	var args []string

	if b.Spec.Input.Type != "" {
		for key, value := range b.Spec.Input.Config {
			args = append(args, fmt.Sprintf("-i %s::%s=%s ", b.Spec.Input.Type, key, value))
		}
	}

	for _, output := range b.Spec.Output {
		var arg string
		if output.Type != "" {
			for key, value := range output.Config {
				arg = fmt.Sprintf("-o %s::%s=%s ", output.Type, key, value)
			}
		}
		args = append(args, arg)
	}

	if b.Spec.Compress.Type != "" {
		for key, value := range b.Spec.Compress.Config {
			args = append(args, fmt.Sprintf("-z %s::%s=%v ", b.Spec.Compress.Type, key, value))
		}
	}

	for _, notifier := range b.Spec.Notifier {
		var arg string
		if notifier.Type != "" {
			for key, value := range notifier.Config {
				arg = fmt.Sprintf("-n %s::%s=%s ", notifier.Type, key, value)
			}
		}
		args = append(args, arg)
	}

	return args
}

type Input struct {
	Type   string
	Config map[string]string
}
type Compress struct {
	Type   string
	Config map[string]int
}
type Output struct {
	Type   string
	Config map[string]string
}
type Notifier struct {
	Type   string
	Config map[string]string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupList is a list of Backup resources
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Backup `json:"items"`
}
