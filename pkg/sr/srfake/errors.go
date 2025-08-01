package srfake

import (
	"fmt"
	"net/http"

	"github.com/twmb/franz-go/pkg/sr"
)

// registryError is a single, unified error type for the mock registry.
type registryError struct {
	// The user-facing message for the HTTP response.
	Message string
	// The HTTP status code to return.
	HTTPStatus int
	// The Schema Registry error code to return.
	SRCode int
}

// Error implements the standard error interface.
func (e *registryError) Error() string {
	return e.Message
}

// newErr creates a new registryError with the given parameters.
func newErr(httpStatus, srCode int, format string, a ...any) *registryError {
	return &registryError{
		HTTPStatus: httpStatus,
		SRCode:     srCode,
		Message:    fmt.Sprintf(format, a...),
	}
}

// Specific error constructors for common error scenarios

func errSubjectNotFound(subject string) *registryError {
	return newErr(http.StatusNotFound, sr.ErrSubjectNotFound.Code, "subject %q not found", subject)
}

func errVersionNotFound(subject string, version int) *registryError {
	return newErr(http.StatusNotFound, sr.ErrVersionNotFound.Code, "version %d not found for %q", version, subject)
}

func errSchemaNotFound() *registryError {
	return newErr(http.StatusNotFound, sr.ErrSchemaNotFound.Code, "schema not found")
}

func errSchemaIsReferenced(subject string, version int, by []int) *registryError {
	return newErr(http.StatusConflict, sr.ErrInvalidSchema.Code, "Cannot delete schema %s:%d as it is still referenced by schema IDs: %v", subject, version, by)
}

func errInvalidReference(ref sr.SchemaReference) *registryError {
	return newErr(http.StatusUnprocessableEntity, sr.ErrInvalidSchema.Code, "reference %q subject %q version %d not found", ref.Name, ref.Subject, ref.Version)
}

func errInvalidSchema(msg string) *registryError {
	return newErr(http.StatusUnprocessableEntity, sr.ErrInvalidSchema.Code, msg)
}

func errInvalidSchemaWithCause(cause error, msg string) *registryError {
	return newErr(http.StatusUnprocessableEntity, sr.ErrInvalidSchema.Code, msg)
}

func errInvalidVersion(msg string) *registryError {
	return newErr(http.StatusBadRequest, sr.ErrInvalidVersion.Code, msg)
}

func errInvalidCompatLevel(msg string) *registryError {
	return newErr(http.StatusBadRequest, sr.ErrInvalidCompatibilityLevel.Code, msg)
}

func errCircularDependency(subject string) *registryError {
	return newErr(http.StatusUnprocessableEntity, sr.ErrInvalidSchema.Code, "circular dependency detected: subject %s is referenced in a cycle", subject)
}
