package handlers

import (
	"fantomrocks-api/internal/services"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	corsAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	corsAccessControlRequestMethod    = "Access-Control-Request-Method"
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsAccessControlMaxAge           = "Access-Control-Max-Age"
	corsHeaderOrigin                  = "Origin"
	maxHeaderElements                 = 25
)

// Define configuration container with CORS settings.
type CORSOptions struct {
	// A list of origins a cross-domain request can be executed from.
	// The special "*" value enables all origins to be accepted.
	// An origin may contain a single wildcard (*) (i.e.: *.dummy.net).
	AllowOrigins []string

	// A list of methods the client is allowed to use with cross-domain requests.
	AllowMethods []string

	// A list of non simple headers allowed with cross-domain requests.
	AllowHeaders []string

	// Is client allowed to include user credentials like cookies, HTTP authentication or client side SSL certificates?
	AllowCredentials bool

	//  How long (in seconds) can client cache result of test OPTIONS request.
	MaxAge int

	// List of compiled wildcard origins to be used to match incoming requests.
	wildcardOrigins []*regexp.Regexp

	// Pre-compiled expression for header fields split.
	splitHeaderRegex *regexp.Regexp
}

// Get default cors options with every origin allowed and basic methods enabled.
func CORSDefault() *CORSOptions {
	return &CORSOptions{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}
}

// Get new CORS handler.
func CORSHandler(log services.Logger, opt *CORSOptions, h http.Handler) http.Handler {
	// compile wildcard origins
	opt.wildcardOrigins = compileWildcardOrigins(opt.AllowOrigins, log)
	opt.splitHeaderRegex = regexp.MustCompile(`(,|\s)+`)

	// make new handler using closure
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// is this an OPTIONS test request?
		if r.Method == http.MethodOptions && "" != r.Header.Get(corsAccessControlRequestMethod) {
			// log the action
			log.Debugf("CORS(): Test OPTIONS request from remote client.")

			// handle options
			opt.handleOptions(w, r, log)

			// close the chain with OK response
			w.WriteHeader(http.StatusOK)
			return
		}

		// handle regular request
		opt.handleRequest(w, r, log)

		// pass the request down the chain
		h.ServeHTTP(w, r)
	})
}

// Compile origins with wildcards to regular expressions we can use to match incoming origins.
func compileWildcardOrigins(origins []string, log services.Logger) []*regexp.Regexp {
	// prep container
	rs := make([]*regexp.Regexp, 0)

	// loop for origins for compilation
	for _, o := range origins {
		// do we have a wildcard in this origin?
		if strings.ContainsRune(o, '*') {
			var exp strings.Builder

			// write expression start
			exp.WriteString("^")

			// replace dots with slash-dots; replace wildcard with dot-wildcard
			exp.WriteString(strings.Replace(strings.ReplaceAll(o, ".", "\\."), "*", ".*", 1))

			// write expression end
			exp.WriteString("$")

			// compile
			r, err := regexp.Compile(exp.String())
			if err != nil {
				log.Errorf("Invalid regular expression '%s'. %s", exp.String(), err.Error())
			} else {
				// add the expression
				rs = append(rs, r)
			}
		}
	}

	return rs
}

// Handle test OPTIONS request from remote client.
func (opt *CORSOptions) handleOptions(w http.ResponseWriter, r *http.Request, log services.Logger) {
	// make sure we handle what we should
	if r.Method != http.MethodOptions {
		return
	}

	// prep headers to be written to
	headers := w.Header()

	// apply Vary headers to inform browser about our intentions
	headers.Add("Vary", corsHeaderOrigin)
	headers.Add("Vary", corsAccessControlRequestMethod)
	headers.Add("Vary", corsAccessControlRequestHeaders)

	// get the actual origin; if not set we can not continue
	origin := r.Header.Get(corsHeaderOrigin)
	if "" == origin {
		log.Debugf("Empty origin, request is aborted.")
		return
	}

	// validate the origin
	if !opt.isOriginAllowed(origin) {
		log.Debugf("Origin '%s' not allowed, request is aborted.", origin)
		return
	}

	// validate the Method
	rMethod := r.Header.Get(corsAccessControlRequestMethod)
	if !opt.isMethodAllowed(rMethod) {
		log.Debugf("Method '%s' not allowed, request is aborted.", rMethod)
		return
	}

	// validate headers
	rHeaders := opt.splitHeaderRegex.Split(r.Header.Get(corsAccessControlRequestHeaders), maxHeaderElements)
	if !opt.isHeadersAllowed(rHeaders) {
		log.Debugf("Certain headers not allowed, request is aborted. %v", rHeaders)
		return
	}

	// indicate we allow the origin
	headers.Set(corsAccessControlAllowOrigin, origin)

	// indicate we allow the request rMethod
	headers.Set(corsAccessControlAllowMethods, strings.ToUpper(rMethod))

	// indicate we allow requested headers
	if 0 < len(rHeaders) {
		headers.Set(corsAccessControlAllowHeaders, strings.Join(rHeaders, ", "))
	}

	// indicate we allow credentials to be passed
	if opt.AllowCredentials {
		headers.Set(corsAccessControlAllowCredentials, "true")
	}

	// indicate max age for cache if set
	if 0 < opt.MaxAge {
		headers.Set(corsAccessControlMaxAge, strconv.Itoa(opt.MaxAge))
	}

	// we are done
	log.Debugf("Response headers %v.", headers)
}

// Handle actual request from remote client.
func (opt *CORSOptions) handleRequest(w http.ResponseWriter, r *http.Request, log services.Logger) {
	// prep headers to be written to
	headers := w.Header()

	// apply Vary headers to inform browser about our intentions
	headers.Add("Vary", corsHeaderOrigin)

	// get the actual origin; if not set we can not continue
	origin := r.Header.Get(corsHeaderOrigin)
	if "" == origin {
		log.Debugf("Empty origin, request is aborted.")
		return
	}

	// validate the origin
	if !opt.isOriginAllowed(origin) {
		log.Debugf("Origin '%s' not allowed, request is aborted.", origin)
		return
	}

	// validate the method (this is not required by specs)
	if !opt.isMethodAllowed(r.Method) {
		log.Debugf("Method '%s' not allowed, request should not happen.", r.Method)
		return
	}

	// indicate we allow the origin
	headers.Set(corsAccessControlAllowOrigin, origin)

	// indicate we allow credentials to be passed
	if opt.AllowCredentials {
		headers.Set(corsAccessControlAllowCredentials, "true")
	}

	// we are done
	log.Debugf("Response headers %v.", headers)
}

// Check if the given Origin is allowed by the options.
func (opt *CORSOptions) isOriginAllowed(origin string) bool {
	origin = strings.ToLower(origin)

	// check simple origins
	for _, o := range opt.AllowOrigins {
		if o == origin {
			return true
		}
	}

	// check compiled regex origins
	for _, rx := range opt.wildcardOrigins {
		if rx.MatchString(origin) {
			return true
		}
	}

	return false
}

// Check if the given Method is allowed by the options.
func (opt *CORSOptions) isMethodAllowed(method string) bool {
	method = strings.ToUpper(method)
	for _, m := range opt.AllowMethods {
		if m == method {
			return true
		}
	}
	return false
}

// Check if the given list of headers is allowed.
func (opt *CORSOptions) isHeadersAllowed(headers []string) bool {
	// loop request headers
	for _, header := range headers {
		header = http.CanonicalHeaderKey(header)
		found := false

		// loop allowed headers
		for _, ah := range opt.AllowHeaders {
			if ah == header {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
