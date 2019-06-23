/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/copybird/copybird/operator/pkg/apis/backup/v1"
	scheme "github.com/copybird/copybird/operator/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BackupsGetter has a method to return a BackupInterface.
// A group's client should implement this interface.
type BackupsGetter interface {
	Backups(namespace string) BackupInterface
}

// BackupInterface has methods to work with Backup resources.
type BackupInterface interface {
	Create(*v1.Backup) (*v1.Backup, error)
	Update(*v1.Backup) (*v1.Backup, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Backup, error)
	List(opts metav1.ListOptions) (*v1.BackupList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Backup, err error)
	BackupExpansion
}

// backups implements BackupInterface
type backups struct {
	client rest.Interface
	ns     string
}

// newBackups returns a Backups
func newBackups(c *BackupsV1Client, namespace string) *backups {
	return &backups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the backup, and returns the corresponding backup object, and an error if there is any.
func (c *backups) Get(name string, options metav1.GetOptions) (result *v1.Backup, err error) {
	result = &v1.Backup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Backups that match those selectors.
func (c *backups) List(opts metav1.ListOptions) (result *v1.BackupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.BackupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested backups.
func (c *backups) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a backup and creates it.  Returns the server's representation of the backup, and an error, if there is any.
func (c *backups) Create(backup *v1.Backup) (result *v1.Backup, err error) {
	result = &v1.Backup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("backups").
		Body(backup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a backup and updates it. Returns the server's representation of the backup, and an error, if there is any.
func (c *backups) Update(backup *v1.Backup) (result *v1.Backup, err error) {
	result = &v1.Backup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backups").
		Name(backup.Name).
		Body(backup).
		Do().
		Into(result)
	return
}

// Delete takes name of the backup and deletes it. Returns an error if one occurs.
func (c *backups) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *backups) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched backup.
func (c *backups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Backup, err error) {
	result = &v1.Backup{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("backups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}