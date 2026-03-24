package main

import (
	"errors"
	"fmt"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	state       State
	failures    int           // error counter
	maxFailures int           // limit faliures before Open
	timeout     time.Duration // Open state duration
	lastFaliure time.Time     // when was a last failure
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	fmt.Println("init state:", cb.state)
	if cb.state == StateOpen {
		if time.Since(cb.lastFaliure) >= cb.timeout {
			cb.state = StateHalfOpen
			cb.lastFaliure = time.Now()
			fmt.Println("state now:", cb.state)
			if err := fn(); err != nil {
				cb.state = StateOpen
				cb.lastFaliure = time.Now()
				fmt.Printf("state %v: %v\n", cb.state, err)
				return err
			}
			cb.failures = 0
			return nil

		}
		return errors.New("circuit is open")
	}

	fmt.Println("state now:", cb.state)
	if err := fn(); err != nil {
		cb.failures++
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			cb.lastFaliure = time.Now()
			return err
		}
		return err
	}

	cb.failures = 0
	cb.state = StateClosed
	fmt.Println("state now:", cb.state)

	return nil
}

func NewCircuitBraker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		failures:    0,
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

func main() {

	cb := NewCircuitBraker(3, 2*time.Second)

	failingService := func() error {
		return fmt.Errorf("service unavailable")
	}

	okService := func() error {
		fmt.Println("запит успішний!")
		return nil
	}

	// if err := cb.Execute(failingService); err != nil {
	// 	fmt.Println(err)
	// }
	// if err := cb.Execute(okService); err != nil {
	// 	fmt.Println(err)
	// }

	// 3 помилки → переходимо в Open
	cb.Execute(failingService)
	cb.Execute(failingService)
	cb.Execute(failingService) // ← тут відкривається

	// одразу після — має повернути "circuit is open"
	cb.Execute(okService)

	// чекаємо таймаут
	time.Sleep(3 * time.Second)

	// тепер HalfOpen — пробуємо знову
	cb.Execute(okService) // ← має спрацювати і перейти в Closed

}
