// Package goxpress provides a fast, intuitive web framework for Go inspired by Express.js.
// This file contains built-in middleware implementations for common web application needs
// such as request logging and panic recovery.
package goxpress

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LoggerConfig defines configuration options for the logger middleware
type LoggerConfig struct {
	// SkipPaths is a list of URL paths to skip logging for.
	// Supports exact matches and simple wildcard patterns with *.
	// Examples: "/health", "/metrics", "/api/*/health"
	SkipPaths []string

	// Output specifies where to write the log output.
	// If nil, defaults to os.Stdout.
	Output io.Writer

	// Formatter specifies a function to format log entries.
	// If nil, defaults to DefaultLogFormatter.
	Formatter LogFormatter
}

// LogFormatter is a function type for custom log formatting
type LogFormatter func(c *Context, start time.Time, duration time.Duration) string

// DefaultLogFormatter returns the default log format
func DefaultLogFormatter(c *Context, start time.Time, duration time.Duration) string {
	return fmt.Sprintf("[%s] %s %s %v\n",
		c.Request.Method,
		c.Request.URL.Path,
		c.Request.RemoteAddr,
		duration,
	)
}

// matchPath checks if a path matches any of the skip patterns
func matchPath(path string, skipPaths []string) bool {
	for _, pattern := range skipPaths {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		// Also support simple wildcard matching
		if strings.Contains(pattern, "*") {
			if simpleWildcardMatch(path, pattern) {
				return true
			}
		} else if path == pattern {
			return true
		}
	}
	return false
}

// simpleWildcardMatch performs simple wildcard matching
func simpleWildcardMatch(path, pattern string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return path == pattern
	}

	// Check if path starts with first part
	if !strings.HasPrefix(path, parts[0]) {
		return false
	}

	// Check if path ends with last part
	if !strings.HasSuffix(path, parts[len(parts)-1]) {
		return false
	}

	// Check middle parts
	remaining := path
	for i, part := range parts {
		if i == 0 {
			remaining = remaining[len(part):]
			continue
		}
		if i == len(parts)-1 {
			break
		}

		idx := strings.Index(remaining, part)
		if idx == -1 {
			return false
		}
		remaining = remaining[idx+len(part):]
	}

	return true
}

// Logger returns a middleware that logs HTTP requests using default configuration.
// It records the HTTP method, URL path, client address, and processing time
// for each request. The log output goes to the standard logger.
//
// The middleware logs after the request is processed, allowing it to
// measure the actual processing time including all other middleware
// and the final handler.
//
// Example:
//
//	app := goxpress.New()
//	app.Use(Logger()) // Enable request logging
//	app.GET("/", handler)
//
// Output format: [METHOD] path clientAddr duration
// Example output: [GET] /api/users 127.0.0.1:54321 1.2ms
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig returns a middleware that logs HTTP requests with custom configuration.
// It allows you to configure skip paths, output destination, and log format.
//
// Example:
//
//	config := goxpress.LoggerConfig{
//		SkipPaths: []string{"/health", "/metrics", "/api/*/internal"},
//		Output:    logFile, // io.Writer
//		Formatter: goxpress.DefaultLogFormatter,
//	}
//	app.Use(goxpress.LoggerWithConfig(config))
func LoggerWithConfig(config LoggerConfig) HandlerFunc {
	// Set defaults
	if config.Output == nil {
		config.Output = os.Stdout
	}
	
	if config.Formatter == nil {
		config.Formatter = DefaultLogFormatter
	}

	return func(c *Context) {
		// Check if this path should be skipped
		if matchPath(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// Record start time
		start := time.Now()

		// Process request through remaining middleware/handlers
		c.Next()

		// Log request details after processing
		duration := time.Since(start)
		logEntry := config.Formatter(c, start, duration)
		log.Println(logEntry)

		// Write to configured output
		config.Output.Write([]byte(logEntry))
	}
}

// Recover returns a middleware that recovers from panics that occur
// during request processing. When a panic is caught, it is converted
// to an error and passed to the error handling middleware chain.
//
// This middleware prevents panics from crashing the entire server
// and allows for graceful error handling and logging.
//
// The middleware handles both error-type panics and arbitrary value panics,
// converting them to appropriate error instances.
//
// Example:
//
//	app := goxpress.New()
//	app.Use(Recover()) // Enable panic recovery
//
//	// Add error handler to process recovered panics
//	app.UseError(func(err error, c *Context) {
//		log.Printf("Recovered from panic: %v", err)
//		c.JSON(500, map[string]string{"error": "Internal Server Error"})
//	})
//
//	app.GET("/panic", func(c *Context) {
//		panic("Something went wrong!") // Will be recovered
//	})
func Recover() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic for debugging
				log.Printf("Panic recovered: %v", r)

				// Abort further processing
				c.Abort()

				// Convert panic to error and pass to error handlers
				var err error
				if e, ok := r.(error); ok {
					// Panic value is already an error
					err = e
				} else {
					// Convert arbitrary panic value to error
					err = fmt.Errorf("%v", r)
				}

				// Pass error to error handling middleware
				c.Next(err)
			}
		}()

		// Continue with normal request processing
		c.Next()
	}
}
