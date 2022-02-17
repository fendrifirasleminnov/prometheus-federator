/*
Copyright 2022 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type PrometheusRuleHandler func(string, *v1.PrometheusRule) (*v1.PrometheusRule, error)

type PrometheusRuleController interface {
	generic.ControllerMeta
	PrometheusRuleClient

	OnChange(ctx context.Context, name string, sync PrometheusRuleHandler)
	OnRemove(ctx context.Context, name string, sync PrometheusRuleHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() PrometheusRuleCache
}

type PrometheusRuleClient interface {
	Create(*v1.PrometheusRule) (*v1.PrometheusRule, error)
	Update(*v1.PrometheusRule) (*v1.PrometheusRule, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.PrometheusRule, error)
	List(namespace string, opts metav1.ListOptions) (*v1.PrometheusRuleList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PrometheusRule, err error)
}

type PrometheusRuleCache interface {
	Get(namespace, name string) (*v1.PrometheusRule, error)
	List(namespace string, selector labels.Selector) ([]*v1.PrometheusRule, error)

	AddIndexer(indexName string, indexer PrometheusRuleIndexer)
	GetByIndex(indexName, key string) ([]*v1.PrometheusRule, error)
}

type PrometheusRuleIndexer func(obj *v1.PrometheusRule) ([]string, error)

type prometheusRuleController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewPrometheusRuleController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) PrometheusRuleController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &prometheusRuleController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromPrometheusRuleHandlerToHandler(sync PrometheusRuleHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.PrometheusRule
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.PrometheusRule))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *prometheusRuleController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.PrometheusRule))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdatePrometheusRuleDeepCopyOnChange(client PrometheusRuleClient, obj *v1.PrometheusRule, handler func(obj *v1.PrometheusRule) (*v1.PrometheusRule, error)) (*v1.PrometheusRule, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *prometheusRuleController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *prometheusRuleController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *prometheusRuleController) OnChange(ctx context.Context, name string, sync PrometheusRuleHandler) {
	c.AddGenericHandler(ctx, name, FromPrometheusRuleHandlerToHandler(sync))
}

func (c *prometheusRuleController) OnRemove(ctx context.Context, name string, sync PrometheusRuleHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromPrometheusRuleHandlerToHandler(sync)))
}

func (c *prometheusRuleController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *prometheusRuleController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *prometheusRuleController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *prometheusRuleController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *prometheusRuleController) Cache() PrometheusRuleCache {
	return &prometheusRuleCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *prometheusRuleController) Create(obj *v1.PrometheusRule) (*v1.PrometheusRule, error) {
	result := &v1.PrometheusRule{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *prometheusRuleController) Update(obj *v1.PrometheusRule) (*v1.PrometheusRule, error) {
	result := &v1.PrometheusRule{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *prometheusRuleController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *prometheusRuleController) Get(namespace, name string, options metav1.GetOptions) (*v1.PrometheusRule, error) {
	result := &v1.PrometheusRule{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *prometheusRuleController) List(namespace string, opts metav1.ListOptions) (*v1.PrometheusRuleList, error) {
	result := &v1.PrometheusRuleList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *prometheusRuleController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *prometheusRuleController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.PrometheusRule, error) {
	result := &v1.PrometheusRule{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type prometheusRuleCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *prometheusRuleCache) Get(namespace, name string) (*v1.PrometheusRule, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.PrometheusRule), nil
}

func (c *prometheusRuleCache) List(namespace string, selector labels.Selector) (ret []*v1.PrometheusRule, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PrometheusRule))
	})

	return ret, err
}

func (c *prometheusRuleCache) AddIndexer(indexName string, indexer PrometheusRuleIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.PrometheusRule))
		},
	}))
}

func (c *prometheusRuleCache) GetByIndex(indexName, key string) (result []*v1.PrometheusRule, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.PrometheusRule, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.PrometheusRule))
	}
	return result, nil
}
