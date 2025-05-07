package functional

// Pipe composes functions left to right (f1, then f2)
func Pipe[T any](f1, f2 func(T) T) func(T) T {
	return func(x T) T {
		return f2(f1(x))
	}
}

// PipeMany composes multiple functions left to right
func PipeMany[T any](funcs ...func(T) T) func(T) T {
	return func(x T) T {
		result := x
		for _, f := range funcs {
			result = f(result)
		}
		return result
	}
}

// Compose composes functions right to left (f2, then f1)
func Compose[T any](f1, f2 func(T) T) func(T) T {
	return func(x T) T {
		return f1(f2(x))
	}
}

// ComposeMany composes multiple functions right to left
func ComposeMany[T any](funcs ...func(T) T) func(T) T {
	return func(x T) T {
		result := x
		for i := len(funcs) - 1; i >= 0; i-- {
			result = funcs[i](result)
		}
		return result
	}
}

// Apply applies a function to a value and returns the result
func Apply[T any, U any](x T, f func(T) U) U {
	return f(x)
}

// ApplyAll applies a series of functions to a value in sequence
func ApplyAll[T any](x T, funcs ...func(T) T) T {
	result := x
	for _, f := range funcs {
		result = f(result)
	}
	return result
}

// Identity returns its input unchanged
func Identity[T any](x T) T {
	return x
}

// Constant returns a function that always returns the same value
func Constant[T any, U any](x T) func(U) T {
	return func(_ U) T {
		return x
	}
}

// Memoize creates a memoized version of a function that caches its results
func Memoize[K comparable, V any](f func(K) V) func(K) V {
	cache := make(map[K]V)
	return func(k K) V {
		if v, found := cache[k]; found {
			return v
		}
		v := f(k)
		cache[k] = v
		return v
	}
}

// Chain applies a function that returns a Result if successful, otherwise propagates the error
// This is useful for composing operations that return Result types
func Chain[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsFailure() {
		return Failure[U](r.GetError())
	}
	return f(r.Unwrap())
}

// Bind applies a function to the value inside a Result, returns a new Result
// This is an alternative name for the Chain operation
func Bind[T any, U any](f func(T) Result[U]) func(Result[T]) Result[U] {
	return func(r Result[T]) Result[U] {
		return Chain(r, f)
	}
}
