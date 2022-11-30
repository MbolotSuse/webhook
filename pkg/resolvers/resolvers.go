// Package resolvers resolves what rules different users and roleTemplates our bound to
package resolvers

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
)

// ErrUnimplemented is an error returned when a function is not implemented.
var ErrUnimplemented = fmt.Errorf("not implemented")

// ruleAccumulator based off kubernetes struct
// https://github.com/kubernetes/kubernetes/blob/d5fdf3135e7c99e5f81e67986ae930f6a2ffb047/pkg/registry/rbac/validation/rule.go#L124#L137
type ruleAccumulator struct {
	rules  []rbacv1.PolicyRule
	errors []error
}

func (r *ruleAccumulator) visit(source fmt.Stringer, rule *rbacv1.PolicyRule, err error) bool {
	if rule != nil {
		r.rules = append(r.rules, *rule)
	}
	if err != nil {
		r.errors = append(r.errors, err)
	}
	return true
}

// getError will combine all of the recorded errors into a single error.
func (r *ruleAccumulator) getError() error {
	if len(r.errors) == 0 {
		return nil
	}
	if len(r.errors) == 1 {
		return r.errors[0]
	}
	var errorStr string
	for _, err := range r.errors {
		errorStr += fmt.Sprintf(", %s", err.Error())
	}
	const leadingChars = 2
	return fmt.Errorf("[%s]", errorStr[leadingChars:])
}

// visitRules calls visitor on each rule in the list with the given Stringer and error.
func visitRules(source fmt.Stringer, rules []rbacv1.PolicyRule, err error, visitor func(source fmt.Stringer, rule *rbacv1.PolicyRule, err error) bool) bool {
	if err != nil && rules == nil {
		return visitor(source, nil, err)
	}
	for i := range rules {
		// we do not care about the return here
		if !visitor(source, &rules[i], err) {
			return false
		}
	}
	return true
}
