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
	Input    Input      `json:"input"`
	Compress Compress   `json:"compress"`
	Output   []Output   `json:"output"`
	Notifier []Notifier `json:"notifier"`
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
