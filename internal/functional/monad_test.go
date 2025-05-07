package functional

import (
	"errors"
	"testing"
)

func TestResult_Success(t *testing.T) {
	result := Success(42)

	if !result.IsSuccess() {
		t.Errorf("Expected Success to be success")
	}

	if result.IsFailure() {
		t.Errorf("Expected Success not to be failure")
	}

	if result.GetError() != nil {
		t.Errorf("Expected Success to have no error, got %v", result.GetError())
	}

	if result.GetValue() != 42 {
		t.Errorf("Expected Success value to be 42, got %v", result.GetValue())
	}

	val := result.Unwrap()
	if val != 42 {
		t.Errorf("Expected Unwrap to return 42, got %v", val)
	}

	val = result.UnwrapOr(99)
	if val != 42 {
		t.Errorf("Expected UnwrapOr to return 42, got %v", val)
	}

	val = result.UnwrapOrElse(func(err error) int {
		return 99
	})
	if val != 42 {
		t.Errorf("Expected UnwrapOrElse to return 42, got %v", val)
	}
}

func TestResult_Failure(t *testing.T) {
	err := errors.New("test error")
	result := Failure[int](err)

	if result.IsSuccess() {
		t.Errorf("Expected Failure not to be success")
	}

	if !result.IsFailure() {
		t.Errorf("Expected Failure to be failure")
	}

	if result.GetError() != err {
		t.Errorf("Expected Failure to have error %v, got %v", err, result.GetError())
	}

	// Test default zero value
	if result.GetValue() != 0 {
		t.Errorf("Expected Failure GetValue to return zero value, got %v", result.GetValue())
	}

	// Test UnwrapOr
	val := result.UnwrapOr(99)
	if val != 99 {
		t.Errorf("Expected UnwrapOr to return fallback 99, got %v", val)
	}

	// Test UnwrapOrElse
	val = result.UnwrapOrElse(func(e error) int {
		if e == err {
			return 77
		}
		return 0
	})
	if val != 77 {
		t.Errorf("Expected UnwrapOrElse to return 77, got %v", val)
	}

	// Test Unwrap - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Unwrap to panic on Failure")
		}
	}()
	_ = result.Unwrap() // This should panic
}

func TestResult_MapResult(t *testing.T) {
	// Test with success
	success := Success(5)
	doubled := success.MapResult(func(i int) int {
		return i * 2
	})

	if !doubled.IsSuccess() || doubled.Unwrap() != 10 {
		t.Errorf("Expected MapResult on Success to return Success(10), got %v", doubled)
	}

	// Test with failure
	err := errors.New("test error")
	failure := Failure[int](err)
	mapped := failure.MapResult(func(i int) int {
		return i * 2
	})

	if !mapped.IsFailure() || mapped.GetError() != err {
		t.Errorf("Expected MapResult on Failure to keep error, got %v", mapped)
	}
}

func TestResult_FlatMap(t *testing.T) {
	// Test with success -> success
	success := Success(5)
	doubled := success.FlatMap(func(i int) Result[int] {
		return Success(i * 2)
	})

	if !doubled.IsSuccess() || doubled.Unwrap() != 10 {
		t.Errorf("Expected FlatMap on Success to Success to return Success(10), got %v", doubled)
	}

	// Test with success -> failure
	failedResult := success.FlatMap(func(i int) Result[int] {
		return Failure[int](errors.New("transform error"))
	})

	if !failedResult.IsFailure() {
		t.Errorf("Expected FlatMap on Success to Failure to return Failure")
	}

	// Test with failure
	err := errors.New("original error")
	failure := Failure[int](err)
	result := failure.FlatMap(func(i int) Result[int] {
		return Success(i * 2)
	})

	if !result.IsFailure() || result.GetError() != err {
		t.Errorf("Expected FlatMap on Failure to keep original error, got %v", result)
	}
}

func TestIO(t *testing.T) {
	// Test creating and performing IO
	counter := 0
	io := NewIO(func() int {
		counter++
		return 42
	})

	// Perform should execute the function
	val := io.Perform()
	if val != 42 {
		t.Errorf("Expected IO.Perform to return 42, got %v", val)
	}
	if counter != 1 {
		t.Errorf("Expected function to be called once, got %d calls", counter)
	}

	// Performing again should execute the function again
	val = io.Perform()
	if counter != 2 {
		t.Errorf("Expected function to be called twice, got %d calls", counter)
	}

	// Test IO Map
	mappedIO := io.Map(func(i int) int {
		return i * 2
	})

	mappedVal := mappedIO.Perform()
	if mappedVal != 84 {
		t.Errorf("Expected mapped IO to return 84, got %v", mappedVal)
	}
	if counter != 3 {
		t.Errorf("Expected original function to be called again, got %d calls", counter)
	}
}

func TestOption_Some(t *testing.T) {
	opt := Some(42)

	if !opt.IsSome() {
		t.Errorf("Expected Some to be Some")
	}

	if opt.IsNone() {
		t.Errorf("Expected Some not to be None")
	}

	val := opt.Unwrap()
	if val != 42 {
		t.Errorf("Expected Some.Unwrap to return 42, got %v", val)
	}

	val = opt.UnwrapOr(99)
	if val != 42 {
		t.Errorf("Expected Some.UnwrapOr to return 42, got %v", val)
	}
}

func TestOption_None(t *testing.T) {
	opt := None[int]()

	if opt.IsSome() {
		t.Errorf("Expected None not to be Some")
	}

	if !opt.IsNone() {
		t.Errorf("Expected None to be None")
	}

	val := opt.UnwrapOr(99)
	if val != 99 {
		t.Errorf("Expected None.UnwrapOr to return fallback 99, got %v", val)
	}

	// Test Unwrap - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected None.Unwrap to panic")
		}
	}()
	_ = opt.Unwrap() // This should panic
}

func TestOption_Map(t *testing.T) {
	// Test Some
	some := Some(5)
	doubled := some.Map(func(i int) int {
		return i * 2
	})

	if !doubled.IsSome() || doubled.Unwrap() != 10 {
		t.Errorf("Expected Some.Map to return Some(10), got %v", doubled)
	}

	// Test None
	none := None[int]()
	result := none.Map(func(i int) int {
		return i * 2
	})

	if !result.IsNone() {
		t.Errorf("Expected None.Map to return None")
	}
}

func TestOption_Bind(t *testing.T) {
	// Test Some -> Some
	some := Some(5)
	doubled := some.Bind(func(i int) Option[int] {
		return Some(i * 2)
	})

	if !doubled.IsSome() || doubled.Unwrap() != 10 {
		t.Errorf("Expected Some.Bind to Some to return Some(10), got %v", doubled)
	}

	// Test Some -> None
	none := some.Bind(func(i int) Option[int] {
		return None[int]()
	})

	if !none.IsNone() {
		t.Errorf("Expected Some.Bind to None to return None")
	}

	// Test None
	originalNone := None[int]()
	stillNone := originalNone.Bind(func(i int) Option[int] {
		return Some(i * 2)
	})

	if !stillNone.IsNone() {
		t.Errorf("Expected None.Bind to return None regardless of function")
	}
}

func TestOption_ToResult(t *testing.T) {
	// Test Some to Success
	some := Some(42)
	err := errors.New("test error")
	result := some.ToResult(err)

	if !result.IsSuccess() || result.Unwrap() != 42 {
		t.Errorf("Expected Some.ToResult to return Success(42)")
	}

	// Test None to Failure
	none := None[int]()
	failResult := none.ToResult(err)

	if !failResult.IsFailure() || failResult.GetError() != err {
		t.Errorf("Expected None.ToResult to return Failure with provided error")
	}
}

func TestMapResultTo(t *testing.T) {
	// Test Success
	success := Success(5)
	mapped := MapResultTo(success, func(i int) string {
		return "value: " + string(rune('0'+i))
	})

	if !mapped.IsSuccess() || mapped.Unwrap() != "value: 5" {
		t.Errorf("Expected MapResultTo on Success to return mapped Success")
	}

	// Test Failure
	err := errors.New("test error")
	failure := Failure[int](err)
	mappedFailure := MapResultTo(failure, func(i int) string {
		return "value: " + string(rune('0'+i))
	})

	if !mappedFailure.IsFailure() || mappedFailure.GetError() != err {
		t.Errorf("Expected MapResultTo on Failure to propagate error")
	}
}

func TestBindResult(t *testing.T) {
	// Test Success -> Success
	success := Success(5)
	result := BindResult(success, func(i int) Result[string] {
		return Success("value: " + string(rune('0'+i)))
	})

	if !result.IsSuccess() || result.Unwrap() != "value: 5" {
		t.Errorf("Expected BindResult on Success to Success to return Success")
	}

	// Test Success -> Failure
	result = BindResult(success, func(i int) Result[string] {
		return Failure[string](errors.New("bind error"))
	})

	if !result.IsFailure() {
		t.Errorf("Expected BindResult on Success to Failure to return Failure")
	}

	// Test Failure
	err := errors.New("original error")
	failure := Failure[int](err)
	result = BindResult(failure, func(i int) Result[string] {
		return Success("This shouldn't be reached")
	})

	if !result.IsFailure() || result.GetError() != err {
		t.Errorf("Expected BindResult on Failure to propagate original error")
	}
}

func TestPure(t *testing.T) {
	result := Pure(42)

	if !result.IsSuccess() || result.Unwrap() != 42 {
		t.Errorf("Expected Pure to create a Success Result")
	}
}

func TestApplyResult(t *testing.T) {
	// Test Success + Success
	rf := Success(func(i int) string { return "value: " + string(rune('0'+i)) })
	ra := Success(5)

	result := ApplyResult(rf, ra)
	if !result.IsSuccess() || result.Unwrap() != "value: 5" {
		t.Errorf("Expected ApplyResult(Success, Success) to be Success with applied function")
	}

	// Test Failure in function
	errF := errors.New("function error")
	failF := Failure[func(int) string](errF)

	result = ApplyResult(failF, ra)
	if !result.IsFailure() || result.GetError() != errF {
		t.Errorf("Expected ApplyResult with failed function to return function's error")
	}

	// Test Failure in argument
	errA := errors.New("argument error")
	failA := Failure[int](errA)

	result = ApplyResult(rf, failA)
	if !result.IsFailure() || result.GetError() != errA {
		t.Errorf("Expected ApplyResult with failed argument to return argument's error")
	}
}

func TestTryCatch(t *testing.T) {
	// Test successful execution
	success := TryCatch(func() (int, error) {
		return 42, nil
	})

	if !success.IsSuccess() || success.Unwrap() != 42 {
		t.Errorf("Expected TryCatch with no error to return Success")
	}

	// Test with error
	err := errors.New("test error")
	failure := TryCatch(func() (int, error) {
		return 0, err
	})

	if !failure.IsFailure() || failure.GetError() != err {
		t.Errorf("Expected TryCatch with error to return Failure with error")
	}
}

func TestChain(t *testing.T) {
	// Test Success -> Success
	success := Success(5)
	result := Chain(success, func(i int) Result[string] {
		return Success("value: " + string(rune('0'+i)))
	})

	if !result.IsSuccess() || result.Unwrap() != "value: 5" {
		t.Errorf("Expected Chain on Success to Success to return Success")
	}

	// Test Success -> Failure
	result = Chain(success, func(i int) Result[string] {
		return Failure[string](errors.New("chain error"))
	})

	if !result.IsFailure() {
		t.Errorf("Expected Chain on Success to Failure to return Failure")
	}

	// Test Failure
	err := errors.New("original error")
	failure := Failure[int](err)
	result = Chain(failure, func(i int) Result[string] {
		return Success("This shouldn't be reached")
	})

	if !result.IsFailure() || result.GetError() != err {
		t.Errorf("Expected Chain on Failure to propagate original error")
	}
}

func TestBind(t *testing.T) {
	// Create a bind function
	bindFn := Bind(func(i int) Result[string] {
		return Success("value: " + string(rune('0'+i)))
	})

	// Test with Success
	success := Success(5)
	result := bindFn(success)

	if !result.IsSuccess() || result.Unwrap() != "value: 5" {
		t.Errorf("Expected Bind with Success to return Success with transformed value")
	}

	// Test with Failure
	err := errors.New("original error")
	failure := Failure[int](err)
	result = bindFn(failure)

	if !result.IsFailure() || result.GetError() != err {
		t.Errorf("Expected Bind with Failure to propagate original error")
	}
}
