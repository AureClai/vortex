//go:build js && wasm

package animation

import (
	"log"
	"time"
)

// ParameterCondition checks animation parameter values
type ParameterCondition struct {
	ParameterName string
	Operator      ConditionOperator
	Value         interface{}
}

func (pc *ParameterCondition) Evaluate(sm *Graph) bool {
	param, exists := sm.GetParameter(pc.ParameterName)
	if !exists {
		return false
	}

	switch pc.Operator {
	case OperatorEquals:
		return param == pc.Value
	case OperatorGreater:
		return compareValues(param, pc.Value) > 0
	case OperatorLess:
		return compareValues(param, pc.Value) < 0
	case OperatorGreaterOrEqual:
		return compareValues(param, pc.Value) >= 0
	case OperatorLessOrEqual:
		return compareValues(param, pc.Value) <= 0
	case OperatorNotEqual:
		return param != pc.Value
	}
	return false
}

func compareValues(a, b interface{}) int {
	switch a.(type) {
	// Compare int values
	case int:
		return a.(int) - b.(int)
	// Compare float64 values
	case float64:
		return int(a.(float64) - b.(float64))
		// Default to 0
	default:
		// Log the error
		log.Printf("Unsupported value type: %T", a)
		return 0
	}
}

// TimeCondition checks animation time
type TimeCondition struct {
	MinTime, MaxTime time.Duration
}

func (tc *TimeCondition) Evaluate(sm *Graph) bool {
	if sm.CurrentNode == nil {
		return false
	}

	elapsed := sm.CurrentNode.LocalTime
	return elapsed >= tc.MinTime && (tc.MaxTime == 0 || elapsed <= tc.MaxTime)
}

// CompleteCondition checks if current animation is complete
type CompleteCondition struct{}

func (cc *CompleteCondition) Evaluate(sm *Graph) bool {
	if sm.CurrentNode == nil || sm.CurrentNode.Clip == nil {
		return false
	}

	return sm.CurrentNode.LocalTime >= sm.CurrentNode.Clip.Duration
}

// CombinedCondition checks multiple conditions
type CombinedCondition struct {
	Conditions []EdgeCondition
	Operator   LogicalOperator
}

func (cc *CombinedCondition) Evaluate(sm *Graph) bool {
	switch cc.Operator {
	case LogicalAnd:
		// Check all conditions
		for _, condition := range cc.Conditions {
			if !condition.Evaluate(sm) {
				return false
			}
		}
		return true
	case LogicalOr:
		// Check any condition
		for _, condition := range cc.Conditions {
			if condition.Evaluate(sm) {
				return true
			}
		}
		return false
	}
	return false
}

// Enums for conditions
type ConditionOperator int

const (
	OperatorEquals ConditionOperator = iota
	OperatorGreater
	OperatorLess
	OperatorGreaterOrEqual
	OperatorLessOrEqual
	OperatorNotEqual
)

type LogicalOperator int

const (
	LogicalAnd LogicalOperator = iota
	LogicalOr
)
