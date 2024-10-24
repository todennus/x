package xerror

import "fmt"

type ErrorWraperConfigs struct {
	goodErrors   []error
	defaultError RichError
}

func NewWrapperConfigs(defaultErr RichError, goodErrors ...error) *ErrorWraperConfigs {
	if len(goodErrors) == 0 {
		panic("invalid good errors")
	}

	return &ErrorWraperConfigs{
		goodErrors:   goodErrors,
		defaultError: defaultErr,
	}
}

func (config *ErrorWraperConfigs) Event(err error, event string, attributes ...any) *errorWrapper {
	return &errorWrapper{err: err, config: config, event: event, attributes: attributes}
}

type errorWrapper struct {
	config     *ErrorWraperConfigs
	err        error
	event      string
	attributes []any
	final      error
}

func (w *errorWrapper) Enrich(code error) *errorWrapperReplaceBy {
	return &errorWrapperReplaceBy{w: w, code: code}
}

func (w *errorWrapper) EnrichWith(code error, description string, format ...any) *errorWrapperReplaceBy {
	return &errorWrapperReplaceBy{w: w, code: code, description: fmt.Sprintf(description, format...)}
}

func (w *errorWrapper) Error() error {
	if w.final == nil {
		w.final = w.config.defaultError.Hide(w.err, w.event, w.attributes...)
	}

	return w.final
}

type errorWrapperReplaceBy struct {
	w           *errorWrapper
	code        error
	description string
}

func (b *errorWrapperReplaceBy) If(target error) *errorWrapper {
	if b.description == "" {
		return replaceBy(b.w, Enrich(b.code, b.w.err.Error()), target)
	} else {
		return replaceBy(b.w, Enrich(b.code, b.description).Hide(b.w.err, b.w.event, b.w.attributes...), target)
	}
}

func (b *errorWrapperReplaceBy) Error() error {
	return b.If(nil).Error()
}

func replaceBy(w *errorWrapper, err error, conditionErr error) *errorWrapper {
	if w.err == nil || w.final != nil {
		return w
	}

	if conditionErr != nil && !Iss(conditionErr, w.config.goodErrors...) {
		panic("condition error must be a good error")
	}

	if conditionErr == nil || Is(w.err, conditionErr) {
		w.final = err
	}

	return w
}
