package engine

import "net/http"

// PageLoader is a function that returns a populated PageContext for a page
type PageLoader func(r *http.Request) *PageContext

var pageLoaders = map[string]PageLoader{}

// RegisterPage registers a data loader for a named page.
// The loader is called before every render of that page to inject live data.
func RegisterPage(name string, loader PageLoader) {
	pageLoaders[name] = loader
}

// GetPageLoader returns the registered loader for a page, or nil
func GetPageLoader(name string) PageLoader {
	return pageLoaders[name]
}
