/*
Copyright 2017 Jetstack Ltd.

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

// This file was automatically generated by informer-gen

package v1beta1

import (
	internalinterfaces "github.com/jetstack/navigator/third_party/k8s.io/client-go/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// ClusterRoles returns a ClusterRoleInformer.
	ClusterRoles() ClusterRoleInformer
	// ClusterRoleBindings returns a ClusterRoleBindingInformer.
	ClusterRoleBindings() ClusterRoleBindingInformer
	// Roles returns a RoleInformer.
	Roles() RoleInformer
	// RoleBindings returns a RoleBindingInformer.
	RoleBindings() RoleBindingInformer
}

type version struct {
	factory internalinterfaces.SharedInformerFactory
	filter  internalinterfaces.FilterFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, filter internalinterfaces.FilterFunc) Interface {
	return &version{factory: f, filter: filter}
}

// ClusterRoles returns a ClusterRoleInformer.
func (v *version) ClusterRoles() ClusterRoleInformer {
	return &clusterRoleInformer{factory: v.factory, filter: v.filter}
}

// ClusterRoleBindings returns a ClusterRoleBindingInformer.
func (v *version) ClusterRoleBindings() ClusterRoleBindingInformer {
	return &clusterRoleBindingInformer{factory: v.factory, filter: v.filter}
}

// Roles returns a RoleInformer.
func (v *version) Roles() RoleInformer {
	return &roleInformer{factory: v.factory, filter: v.filter}
}

// RoleBindings returns a RoleBindingInformer.
func (v *version) RoleBindings() RoleBindingInformer {
	return &roleBindingInformer{factory: v.factory, filter: v.filter}
}
