// Package goxpress provides a fast, intuitive web framework for Go inspired by Express.js.
// This file contains built-in middleware implementations for common web application needs
// such as request logging and panic recovery.
package goxpress

import (
	"fmt"
	"log"
	"time"
)

// Logger returns a middleware that logs HTTP requests.
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
	return func(c *Context) {
		// Record start time
		start := time.Now()

		// Process request through remaining middleware/handlers
		c.Next()

		// Log request details after processing
		duration := time.Since(start)
		log.Printf("[%s] %s %s %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.RemoteAddr,
			duration,
		)
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
