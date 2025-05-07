// Package functional provides functional programming utilities and patterns
package functional

// TYPES

// Result is a monadic type representing either a success value or an error
// Similar to Either in functional programming
type Result[T any] struct {
	value T
	err   error
}

// IO is a monadic type representing a computation with side effects
type IO[T any] struct {
	run func() T
}

// Option is a monadic type representing an optional value (Some or None)
type Option[T any] struct {
	value *T
}

// RESULT METHODS

// Unwrap returns the value if successful or panics if there's an error
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr returns the value if successful or the provided default if there's an error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value if successful or calls the provided function if there's an error
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.err != nil {
		return f(r.err)
	}
	return r.value
}

// IsSuccess returns true if the Result contains a success value
func (r Result[T]) IsSuccess() bool {
	return r.err == nil
}

// IsFailure returns true if the Result contains an error
func (r Result[T]) IsFailure() bool {
	return r.err != nil
}

// MapResult applies a function to the value if successful, otherwise propagates the error
func (r Result[T]) MapResult(f func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Success(f(r.value))
}

// FlatMap (or Bind) applies a function that returns a Result if successful, otherwise propagates the error
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return f(r.value)
}

// GetError returns the error if present, or nil if the result is successful
func (r Result[T]) GetError() error {
	return r.err
}

// GetValue returns the value regardless of whether the result is successful
// Check IsSuccess() first to determine if the value is valid
func (r Result[T]) GetValue() T {
	return r.value
}

// IO METHODS

// Perform executes the IO operation and returns the result
func (io IO[T]) Perform() T {
	return io.run()
}

// Map applies a function to the result of an IO operation
func (io IO[T]) Map(f func(T) T) IO[T] {
	return NewIO(func() T {
		return f(io.Perform())
	})
}

// OPTION METHODS

// IsSome returns true if the Option contains a value
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone returns true if the Option contains no value
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Unwrap returns the value if it exists or panics
func (o Option[T]) Unwrap() T {
	if o.value == nil {
		panic("attempted to unwrap a None value")
	}
	return *o.value
}

// UnwrapOr returns the value if it exists or the provided default
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.value == nil {
		return defaultValue
	}
	return *o.value
}

// Map applies a function to the value if it exists
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.value == nil {
		return o
	}
	result := f(*o.value)
	return Some(result)
}

// Bind (or FlatMap) applies a function that returns an Option if the value exists
func (o Option[T]) Bind(f func(T) Option[T]) Option[T] {
	if o.value == nil {
		return o
	}
	return f(*o.value)
}

// ToResult converts an Option to a Result with the provided error for None
func (o Option[T]) ToResult(err error) Result[T] {
	if o.value == nil {
		return Failure[T](err)
	}
	return Success(*o.value)
}

// EXPORTED FUNCTIONS

// Success creates a new Result with a success value
func Success[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		err:   nil,
	}
}

// Failure creates a new Result with an error
func Failure[T any](err error) Result[T] {
	var zero T
	return Result[T]{
		value: zero,
		err:   err,
	}
}

// MapResultTo applies a function that changes the type if successful, otherwise propagates the error
func MapResultTo[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.err != nil {
		return Failure[U](r.err)
	}
	return Success(f(r.value))
}

// BindResult is a generic version of FlatMap that allows changing the type
func BindResult[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Failure[U](r.err)
	}
	return f(r.value)
}

// Pure lifts a value into a Result context (Applicative pattern)
func Pure[T any](value T) Result[T] {
	return Success(value)
}

// ApplyResult applies a function inside a Result to a value inside another Result (Applicative pattern)
func ApplyResult[T any, U any](rf Result[func(T) U], ra Result[T]) Result[U] {
	if rf.IsFailure() {
		return Failure[U](rf.GetError())
	}
	if ra.IsFailure() {
		return Failure[U](ra.GetError())
	}
	f := rf.Unwrap()
	a := ra.Unwrap()
	return Success(f(a))
}

// NewIO creates a new IO with a function that produces a value
func NewIO[T any](f func() T) IO[T] {
	return IO[T]{run: f}
}

// MapIO applies a function that changes the type to the result of an IO operation
func MapIO[T any, U any](io IO[T], f func(T) U) IO[U] {
	return NewIO(func() U {
		return f(io.Perform())
	})
}

// FlatMapIO (or Bind) chains IO operations
func FlatMapIO[T any, U any](io IO[T], f func(T) IO[U]) IO[U] {
	return NewIO(func() U {
		t := io.Perform()
		return f(t).Perform()
	})
}

// IOPure lifts a value into an IO context
func IOPure[T any](value T) IO[T] {
	return NewIO(func() T {
		return value
	})
}

// Some creates a new Option with a value
func Some[T any](value T) Option[T] {
	return Option[T]{&value}
}

// None creates a new Option with no value
func None[T any]() Option[T] {
	return Option[T]{nil}
}

// MapOption applies a function that changes the type if the value exists
func MapOption[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return Some(f(*o.value))
}

// BindOption is a generic version of Bind that allows changing the type
func BindOption[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return f(*o.value)
}

// OptPure lifts a value into an Option context (Applicative pattern)
func OptPure[T any](value T) Option[T] {
	return Some(value)
}

// OptApply applies a function inside an Option to a value inside another Option (Applicative pattern)
func OptApply[T any, U any](of Option[func(T) U], oa Option[T]) Option[U] {
	if of.IsNone() || oa.IsNone() {
		return None[U]()
	}
	f := of.Unwrap()
	a := oa.Unwrap()
	return Some(f(a))
}

// TryCatch executes a function that may return an error and converts it to a Result
func TryCatch[T any](f func() (T, error)) Result[T] {
	value, err := f()
	if err != nil {
		return Failure[T](err)
	}
	return Success(value)
}
