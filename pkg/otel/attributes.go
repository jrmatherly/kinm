package otel

import (
	"fmt"

	"github.com/obot-platform/kinm/pkg/types"

	"go.opentelemetry.io/otel/attribute"
	"k8s.io/apiserver/pkg/storage"
)

// stringerOrEmpty returns an attribute with the Stringer's value, or an empty string if nil.
// This prevents panics when Label or Field selectors are nil.
func stringerOrEmpty(name string, s fmt.Stringer) attribute.KeyValue {
	if s == nil {
		return attribute.String(name, "")
	}
	return attribute.Stringer(name, s)
}

func ListOptionsToAttributes(opts storage.ListOptions, otherAttributes ...attribute.KeyValue) []attribute.KeyValue {
	return append(otherAttributes,
		attribute.String("resourceVersion", opts.ResourceVersion),
		attribute.String("continue", opts.Predicate.Continue),
		attribute.Int64("limit", opts.Predicate.Limit),
		attribute.Bool("allowWatchBookmarks", opts.Predicate.AllowWatchBookmarks),
		attribute.StringSlice("indexLabels", opts.Predicate.IndexLabels),
		attribute.StringSlice("indexFields", opts.Predicate.IndexFields),
		stringerOrEmpty("labelSelector", opts.Predicate.Label),
		stringerOrEmpty("fieldSelector", opts.Predicate.Field),
		attribute.String("resourceVersionMatch", string(opts.ResourceVersionMatch)),
		attribute.Bool("progressNotify", opts.ProgressNotify),
		attribute.Bool("recursive", opts.Recursive),
		attribute.Bool("sendInitialEvents", opts.SendInitialEvents == nil || *opts.SendInitialEvents),
	)
}

func ObjectToAttributes(obj types.Object, otherAttributes ...attribute.KeyValue) []attribute.KeyValue {
	return append(otherAttributes,
		attribute.String("name", obj.GetName()),
		attribute.String("namespace", obj.GetNamespace()),
	)
}
