// Package strategy provides the core strategy layer defining interfaces for Kubernetes-like API resource operations.
package strategy

import "k8s.io/apimachinery/pkg/runtime"

type CompleteStrategy interface {
	Creater
	Updater
	StatusUpdater
	Getter
	Lister
	Deleter
	Watcher

	Destroy()
	Scheme() *runtime.Scheme
}
