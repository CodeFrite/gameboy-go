package gameboy

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/codefrite/gameboy-go/gameboy"
	"github.com/cucumber/godog"
)

// Define a custom type for the context key
type contextKey string

const timerKey contextKey = "timer"

// Initialize the context with a *Timer
func initializeContext(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	return toContext(ctx, gameboy.NewTimer(1)), nil
}

func toContext(ctx context.Context, timer *gameboy.Timer) context.Context {
	return context.WithValue(ctx, timerKey, timer)
}

func fromContext(ctx context.Context) *gameboy.Timer {
	timer, _ := ctx.Value(timerKey).(*gameboy.Timer)
	return timer
}

// Step Definitions
func iAddASubscriber(ctx context.Context) (context.Context, error) {
	fmt.Println("iAddASubscriber>")
	timer := fromContext(ctx)
	if timer == nil {
		return ctx, fmt.Errorf("timer is nil")
	}
	cpu := gameboy.NewCPU(nil)
	timer.Subscribe(cpu)
	return ctx, nil
}

func iAddSubscribers(ctx context.Context, arg1 int) (context.Context, error) {
	fmt.Println("iAddASubscriber>")
	timer := fromContext(ctx)
	if timer == nil {
		return ctx, fmt.Errorf("timer is nil")
	}
	cpu := gameboy.NewCPU(nil)
	for i := 0; i < arg1; i++ {
		timer.Subscribe(cpu)
	}
	return ctx, nil
}

func iInstantiateANewStructAs(ctx context.Context, _type string, _name string) (context.Context, error) {
	fmt.Println("iInstantiateANewStructAs>")
	if _type == "Timer" {
		timer := gameboy.NewTimer(1)
		ctx = toContext(ctx, timer)
	}
	return ctx, nil
}

func iShouldHaveAMapWithTheFollowingKeyValuePairs(ctx context.Context, arg1 string, table *godog.Table) (context.Context, error) {
	fmt.Println("iShouldHaveAMapWithTheFollowingKeyValuePairs>")

	// Extract headers
	headers := table.Rows[0].Cells
	if len(headers) != 2 {
		return ctx, fmt.Errorf("expected 2 columns in the table, got %d", len(headers))
	}

	// Create a map to store the key-value pairs
	expectedMap := make(map[string]string)

	// Iterate over the rows to extract key-value pairs
	for _, row := range table.Rows[1:] {
		if len(row.Cells) != 2 {
			return ctx, fmt.Errorf("expected 2 columns in the row, got %d", len(row.Cells))
		}
		key := row.Cells[0].Value
		value := row.Cells[1].Value
		expectedMap[key] = value
	}

	// Retrieve the actual map from the context
	timer := fromContext(ctx)
	if timer == nil {
		return ctx, fmt.Errorf("timer is nil")
	}

	// Use reflection to access the Registers field
	val := reflect.ValueOf(timer).Elem().FieldByName("Registers")
	if !val.IsValid() {
		return ctx, fmt.Errorf("field Registers not found in Timer struct")
	}

	// Ensure the field is a map
	if val.Kind() != reflect.Map {
		return ctx, fmt.Errorf("Registers field is not a map")
	}

	// Compare the actual map with the expected map
	for key, expectedValue := range expectedMap {
		actualValue := val.MapIndex(reflect.ValueOf(key))
		if !actualValue.IsValid() {
			return ctx, fmt.Errorf("key %s not found in the actual map", key)
		}
		if fmt.Sprintf("0x%X", actualValue.Interface()) != expectedValue {
			return ctx, fmt.Errorf("for key %s, expected value %s, got %s", key, expectedValue, fmt.Sprintf("0x%X", actualValue.Interface()))
		}
	}

	return ctx, nil
}

func iShouldHaveAVariableWithTheFollowingFields(ctx context.Context, arg1 string, arg2 *godog.Table) (context.Context, error) {
	fmt.Println("iShouldHaveAVariableWithTheFollowingFields>")
	return ctx, godog.ErrPending
}

func theTimerShouldHaveSubscribers(ctx context.Context, arg1 int) (context.Context, error) {
	fmt.Println("theTimerShouldHaveSubscribers>")
	timer := fromContext(ctx)
	if timer == nil {
		return ctx, fmt.Errorf("timer is nil")
	}
	if len(timer.Subscribers) != arg1 {
		return ctx, fmt.Errorf("expected %d subscribers, got %d", arg1, len(timer.Subscribers))
	}
	return ctx, nil
}

// Initialize the scenario
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Define the preconditions for the scenario: share the test suite context
	ctx.Before(initializeContext)

	// Define the steps for the scenario
	ctx.Step(`^i add a subscriber$`, iAddASubscriber)
	ctx.Step(`^i instantiate a new struct "([^"]*)" as "([^"]*)"$`, iInstantiateANewStructAs)
	ctx.Step(`^i should have a map "([^"]*)" with the following key value pairs:$`, iShouldHaveAMapWithTheFollowingKeyValuePairs)
	ctx.Step(`^i should have a variable "([^"]*)" with the following fields:$`, iShouldHaveAVariableWithTheFollowingFields)
	ctx.Step(`^the timer should have (\d+) subscribers$`, theTimerShouldHaveSubscribers)
	ctx.Step(`^i add (\d+) subscribers$`, iAddSubscribers)

}

// Test runner
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/timer.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
