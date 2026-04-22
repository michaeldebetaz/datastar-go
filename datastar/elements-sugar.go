package datastar

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/valyala/bytebufferpool"
)

// ValidElementPatchModes is a list of valid element patch modes.
var ValidElementPatchModes = []ElementPatchMode{
	ElementPatchModeOuter,
	ElementPatchModeInner,
	ElementPatchModeRemove,
	ElementPatchModePrepend,
	ElementPatchModeAppend,
	ElementPatchModeBefore,
	ElementPatchModeAfter,
	ElementPatchModeReplace,
}

// ValidNamespaces is a list of valid namespaces.
var ValidNamespaces = []Namespace{
	NamespaceHTML,
	NamespaceSVG,
	NamespaceMathML,
}

// ElementPatchModeFromString converts a string to a [ElementPatchMode].
func ElementPatchModeFromString(s string) (ElementPatchMode, error) {
	switch s {
	case "outer":
		return ElementPatchModeOuter, nil
	case "inner":
		return ElementPatchModeInner, nil
	case "remove":
		return ElementPatchModeRemove, nil
	case "prepend":
		return ElementPatchModePrepend, nil
	case "append":
		return ElementPatchModeAppend, nil
	case "before":
		return ElementPatchModeBefore, nil
	case "after":
		return ElementPatchModeAfter, nil
	case "replace":
		return ElementPatchModeReplace, nil
	default:
		return "", fmt.Errorf("invalid element merge type: %s", s)
	}
}

// NamespaceFromString converts a string to a [Namespace].
func NamespaceFromString(s string) (Namespace, error) {
	switch s {
	case "html":
		return NamespaceHTML, nil
	case "svg":
		return NamespaceSVG, nil
	case "mathml":
		return NamespaceMathML, nil
	default:
		return "", fmt.Errorf("invalid namespace: %s", s)
	}
}

// WithSelectorID is a convenience wrapper for [WithSelector] option
// equivalent to calling `WithSelector("#"+id)`.
func WithSelectorID(id string) PatchElementOption {
	return WithSelector("#" + id)
}

// WithModeOuter creates a PatchElementOption that merges elements using the outer mode.
func WithModeOuter() PatchElementOption {
	return WithMode(ElementPatchModeOuter)
}

// WithModeInner creates a PatchElementOption that merges elements using the inner mode.
func WithModeInner() PatchElementOption {
	return WithMode(ElementPatchModeInner)
}

// WithModeRemove creates a PatchElementOption that removes elements from the DOM.
func WithModeRemove() PatchElementOption {
	return WithMode(ElementPatchModeRemove)
}

// WithModePrepend creates a PatchElementOption that merges elements using the prepend mode.
func WithModePrepend() PatchElementOption {
	return WithMode(ElementPatchModePrepend)
}

// WithModeAppend creates a PatchElementOption that merges elements using the append mode.
func WithModeAppend() PatchElementOption {
	return WithMode(ElementPatchModeAppend)
}

// WithModeBefore creates a PatchElementOption that merges elements using the before mode.
func WithModeBefore() PatchElementOption {
	return WithMode(ElementPatchModeBefore)
}

// WithModeAfter creates a PatchElementOption that merges elements using the after mode.
func WithModeAfter() PatchElementOption {
	return WithMode(ElementPatchModeAfter)
}

// WithModeReplace creates a PatchElementOption that replaces elements without morphing.
// This mode does not use morphing and will completely replace the element, resetting any related state.
func WithModeReplace() PatchElementOption {
	return WithMode(ElementPatchModeReplace)
}

// WithNamespaceHTML specifies the HTML namespace for the elements being patched.
func WithNamespaceHTML() PatchElementOption {
	return WithNamespace(NamespaceHTML)
}

// WithNamespace specifies the XML namespace for the elements being patched.
func WithNamespaceSVG() PatchElementOption {
	return WithNamespace(NamespaceSVG)
}

// WithNamespaceMathML specifies the MathML namespace for the elements being patched.
func WithNamespaceMathML() PatchElementOption {
	return WithNamespace(NamespaceMathML)
}

// WithViewTransitions enables the use of view transitions when merging elements.
func WithViewTransitions() PatchElementOption {
	return func(o *patchElementOptions) {
		o.UseViewTransitions = true
	}
}

// WithoutViewTransitions disables the use of view transitions when merging elements.
func WithoutViewTransitions() PatchElementOption {
	return func(o *patchElementOptions) {
		o.UseViewTransitions = false
	}
}

// PatchElementf is a convenience wrapper for [PatchElements] option
// equivalent to calling `PatchElements(fmt.Sprintf(format, args...))`.
func (sse *ServerSentEventGenerator) PatchElementf(format string, args ...any) error {
	return sse.PatchElements(fmt.Sprintf(format, args...))
}

// TemplComponent satisfies the component rendering interface for HTML template engine [Templ].
// This separate type ensures compatibility with [Templ] without imposing a dependency requirement
// on those who prefer to use a different template engine.
//
// [Templ]: https://templ.guide/
type TemplComponent interface {
	Render(ctx context.Context, w io.Writer) error
}

// PatchElementTempl is a convenience adaptor of [sse.PatchElements] for [TemplComponent].
func (sse *ServerSentEventGenerator) PatchElementTempl(c TemplComponent, opts ...PatchElementOption) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := c.Render(sse.Context(), buf); err != nil {
		return fmt.Errorf("failed to patch element: %w", err)
	}
	if err := sse.PatchElements(buf.String(), opts...); err != nil {
		return fmt.Errorf("failed to patch element: %w", err)
	}
	return nil
}

// GoStarElementRenderer satisfies the component rendering interface for HTML template engine [GoStar].
// This separate type ensures compatibility with [GoStar] without imposing a dependency requirement
// on those who prefer to use a different template engine.
//
// [GoStar]: https://github.com/delaneyj/gostar
type GoStarElementRenderer interface {
	Render(w io.Writer) error
}

// PatchElementGostar is a convenience adaptor of [sse.PatchElements] for [GoStarElementRenderer].
func (sse *ServerSentEventGenerator) PatchElementGostar(child GoStarElementRenderer, opts ...PatchElementOption) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := child.Render(buf); err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}
	if err := sse.PatchElements(buf.String(), opts...); err != nil {
		return fmt.Errorf("failed to patch element: %w", err)
	}
	return nil
}

type contentTypeChoice string

const (
	ContentTypeJSON contentTypeChoice = "json"
	ContentTypeForm contentTypeChoice = "form"
)

type retryChoice string

const (
	// Default, retries on network errors only
	RetryAuto retryChoice = "auto"
	// Retries on 4xx and 5xx responses
	RetryError retryChoice = "error"
	// Retries on all non-204 responses except redirects
	RetryAlways retryChoice = "always"
	// Disables retries
	RetryNever retryChoice = "never"
)

// see: https://data-star.dev/reference/actions#options
type actionSSEOptions struct {
	ContentType         *contentTypeChoice
	FilterSignals       *string
	FormSelector        *string
	Headers             *string
	OpenWhenHidden      *bool
	Payload             *string
	RequestCancellation *string
	Retry               *retryChoice
	RetryInterval       *int
	RetryScaler         *int
	RetryMaxWait        *int
	RetryMaxCount       *int
}

// If the option provided is the default datastar value, it will not be inserted.
type ActionSSEOption func(o *actionSSEOptions)

// The type of content to send. A value of json sends all signals in a JSON request. A value of form tells the action to look for the closest form to the element on which it is placed (unless a selector option is provided), perform validation on the form elements, and send them to the backend using a form request (no signals are sent). Defaults to json.
func WithContentType(contentType contentTypeChoice) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.ContentType = &contentType
	}
}

// A filter object with an include property that accepts a regular expression to match signal paths (defaults to all signals: /.*/), and an optional exclude property to exclude specific signal paths (defaults to all signals that do not have a _ prefix: /(^_|\._).*/).
func WithFilterSignals(filterSignals string) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.FilterSignals = &filterSignals
	}
}

// Optionally specifies a form to send when the contentType option is set to form. If the value is null, the closest form is used. Defaults to null.
func WithFormSelector(selector string) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.FormSelector = &selector
	}
}

// An object containing headers to send with the request.
func WithHeaders(headers string) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.Headers = &headers
	}
}

// Whether to keep the connection open when the page is hidden. Useful for dashboards but can cause a drain on battery life and other resources when enabled. Defaults to false.
func WithOpenWhenHidden(openWhenHidden bool) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.OpenWhenHidden = &openWhenHidden
	}
}

// Allows the fetch payload to be overridden with a custom object.
func WithPayload(payload string) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.Payload = &payload
	}
}

// Determines when to retry requests. Can be 'auto' (default, retries on network errors only), 'error' (retries on 4xx and 5xx responses), 'always' (retries on all non-204 responses except redirects), or 'never' (disables retries). Defaults to 'auto'.
func WithRetry(retry retryChoice) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.Retry = &retry
	}
}

// The retry interval in milliseconds. Defaults to 1000 (one second).
func WithRetryInterval(retryIntervalMs int) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.RetryInterval = &retryIntervalMs
	}
}

// A numeric multiplier applied to scale retry wait times. Defaults to 2.
func WithRetryScaler(retryScaler int) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.RetryScaler = &retryScaler
	}
}

// The maximum allowable wait time in milliseconds between retries. Defaults to 30000 (30 seconds).
func WithRetryMaxWait(retryMaxWaitMs int) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.RetryMaxWait = &retryMaxWaitMs
	}
}

// The maximum number of retry attempts. Defaults to 10.
func WithRetryMaxCount(retryMaxCount int) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.RetryMaxCount = &retryMaxCount
	}
}

// Controls request cancellation behavior. Can be 'auto' (default, cancels existing requests on the same element), 'disabled' (allows concurrent requests), or an AbortController instance for custom control. Defaults to 'auto'.
func WithRequestCancellation(requestCancellation string) ActionSSEOption {
	return func(o *actionSSEOptions) {
		o.RequestCancellation = &requestCancellation
	}
}

func SSEGet(url string, opts ...ActionSSEOption) string {
	return sseAction("@get", url, opts...)
}

func SSEPost(url string, opts ...ActionSSEOption) string {
	return sseAction("@post", url, opts...)
}

func SSEPut(url string, opts ...ActionSSEOption) string {
	return sseAction("@put", url, opts...)
}

func SSEPatch(url string, opts ...ActionSSEOption) string {
	return sseAction("@patch", url, opts...)
}

func SSEDelete(url string, opts ...ActionSSEOption) string {
	return sseAction("@delete", url, opts...)
}

func sseAction(action string, url string, opts ...ActionSSEOption) string {
	var sb strings.Builder
	sb.Grow(256)

	sb.WriteString(action)
	sb.WriteString("('")
	sb.WriteString(url)
	sb.WriteString("'")

	if len(opts) == 0 {
		sb.WriteString(")")
		return sb.String()
	}

	options := &actionSSEOptions{}
	for _, opt := range opts {
		opt(options)
	}

	sb.WriteString(", { ")
	first := true

	writeKey := func(key string) {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(key)
		first = false
	}

	if options.ContentType != nil {
		writeKey("contentType: '")
		sb.WriteString(string(*options.ContentType))
		sb.WriteString("'")
	}

	if options.FilterSignals != nil {
		writeKey("filterSignals: ")
		sb.WriteString(*options.FilterSignals)
	}

	if options.FormSelector != nil {
		writeKey("selector: '")
		sb.WriteString(*options.FormSelector)
		sb.WriteString("'")
	}

	if options.Headers != nil {
		writeKey("headers: '")
		sb.WriteString(*options.Headers)
		sb.WriteString("'")
	}

	if options.OpenWhenHidden != nil {
		writeKey("openWhenHidden: ")
		sb.WriteString(strconv.FormatBool(*options.OpenWhenHidden))
	}

	if options.Payload != nil {
		writeKey("payload: ")
		sb.WriteString(*options.Payload)

	}

	if options.RequestCancellation != nil {
		writeKey("requestCancellation: '")
		sb.WriteString(*options.RequestCancellation)
		sb.WriteString("'")
	}

	if options.Retry != nil {
		writeKey("retry: '")
		sb.WriteString(string(*options.Retry))
		sb.WriteString("'")
	}

	if options.RetryInterval != nil {
		writeKey("retryInterval: ")
		sb.WriteString(strconv.Itoa(*options.RetryInterval))
	}

	if options.RetryMaxCount != nil {
		writeKey("retryMaxCount: ")
		sb.WriteString(strconv.Itoa(*options.RetryMaxCount))
	}

	if options.RetryMaxWait != nil {
		writeKey("retryMaxWait: ")
		sb.WriteString(strconv.Itoa(*options.RetryMaxWait))
	}

	if options.RetryScaler != nil {
		writeKey("retryScaler: ")
		sb.WriteString(strconv.Itoa(*options.RetryScaler))
	}

	sb.WriteString(" })")

	return sb.String()
}

// RemoveElement is a convenience method for removing elements from the DOM.
// It uses PatchElements with the remove mode and the specified selector.
func (sse *ServerSentEventGenerator) RemoveElement(selector string, opts ...PatchElementOption) error {
	// Prepend the remove mode option
	allOpts := append([]PatchElementOption{WithModeRemove(), WithSelector(selector)}, opts...)
	return sse.PatchElements("", allOpts...)
}

// RemoveElementf is a convenience wrapper for RemoveElement that formats the selector string
// using the provided format and arguments similar to fmt.Sprintf.
func (sse *ServerSentEventGenerator) RemoveElementf(selectorFormat string, args ...any) error {
	selector := fmt.Sprintf(selectorFormat, args...)
	return sse.RemoveElement(selector)
}

// RemoveElementByID is a convenience wrapper for RemoveElement that removes an element by its ID.
// Equivalent to calling RemoveElement("#"+id).
func (sse *ServerSentEventGenerator) RemoveElementByID(id string) error {
	return sse.RemoveElement("#" + id)
}
