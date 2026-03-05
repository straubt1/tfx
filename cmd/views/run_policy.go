// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

import (
	"fmt"
)

// RunPolicyResult holds all policy information for a run
type RunPolicyResult struct {
	RunID     string
	RunStatus string

	// Legacy Sentinel policy checks
	PolicyChecks []PolicyCheckDetail

	// Newer OPA/Sentinel policy evaluations (via task stages)
	PolicyEvaluations []PolicyEvaluationDetail

	// True if any policy data exists
	HasPolicies bool
}

// PolicyCheckDetail represents a single legacy Sentinel policy check
type PolicyCheckDetail struct {
	ID             string
	Status         string
	Scope          string
	Overridable    bool
	Passed         int
	AdvisoryFailed int
	SoftFailed     int
	HardFailed     int
	TotalFailed    int
	Duration       int
	Logs           string // Raw Sentinel log output (populated when --logs flag is used)
}

// PolicyEvaluationDetail represents a newer policy evaluation (OPA or Sentinel)
type PolicyEvaluationDetail struct {
	ID              string
	Status          string
	PolicyKind      string // "opa" or "sentinel"
	Stage           string // "post_plan", "pre_apply", etc.
	Passed          int
	AdvisoryFailed  int
	MandatoryFailed int
	Errored         int
	PolicySets      []PolicySetDetail
}

// PolicySetDetail represents a policy set within an evaluation
type PolicySetDetail struct {
	ID              string
	PolicySetName   string
	Overridable     bool
	Error           string
	Passed          int
	AdvisoryFailed  int
	MandatoryFailed int
	Errored         int
	Policies        []PolicyDetail
}

// PolicyDetail represents an individual policy outcome
type PolicyDetail struct {
	PolicyName       string
	EnforcementLevel string
	Status           string
	Description      string
	Query            string
	Output           string // Raw log output from output.print (populated when --logs flag is used)
}

// RunPolicyView handles rendering for run policy command
type RunPolicyView struct{ *BaseView }

func NewRunPolicyView() *RunPolicyView { return &RunPolicyView{NewBaseView()} }

// JSON output types
type runPolicyOutput struct {
	RunID             string                   `json:"runId"`
	RunStatus         string                   `json:"runStatus"`
	HasPolicies       bool                     `json:"hasPolicies"`
	PolicyChecks      []policyCheckOutput      `json:"policyChecks"`
	PolicyEvaluations []policyEvaluationOutput `json:"policyEvaluations"`
}

type policyCheckOutput struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Scope          string `json:"scope"`
	Overridable    bool   `json:"overridable"`
	Passed         int    `json:"passed"`
	AdvisoryFailed int    `json:"advisoryFailed"`
	SoftFailed     int    `json:"softFailed"`
	HardFailed     int    `json:"hardFailed"`
	TotalFailed    int    `json:"totalFailed"`
	Duration       int    `json:"duration"`
	Logs           string `json:"logs,omitempty"`
}

type policyEvaluationOutput struct {
	ID              string            `json:"id"`
	Status          string            `json:"status"`
	PolicyKind      string            `json:"policyKind"`
	Stage           string            `json:"stage"`
	Passed          int               `json:"passed"`
	AdvisoryFailed  int               `json:"advisoryFailed"`
	MandatoryFailed int               `json:"mandatoryFailed"`
	Errored         int               `json:"errored"`
	PolicySets      []policySetOutput `json:"policySets"`
}

type policySetOutput struct {
	ID              string               `json:"id"`
	PolicySetName   string               `json:"policySetName"`
	Overridable     bool                 `json:"overridable"`
	Error           string               `json:"error,omitempty"`
	Passed          int                  `json:"passed"`
	AdvisoryFailed  int                  `json:"advisoryFailed"`
	MandatoryFailed int                  `json:"mandatoryFailed"`
	Errored         int                  `json:"errored"`
	Policies        []policyDetailOutput `json:"policies"`
}

type policyDetailOutput struct {
	PolicyName       string `json:"policyName"`
	EnforcementLevel string `json:"enforcementLevel"`
	Status           string `json:"status"`
	Description      string `json:"description,omitempty"`
	Query            string `json:"query,omitempty"`
	Output           string `json:"output,omitempty"`
}

func (v *RunPolicyView) Render(result *RunPolicyResult) error {
	if v.IsJSON() {
		return v.renderJSON(result)
	}
	return v.renderTerminal(result)
}

func (v *RunPolicyView) renderJSON(result *RunPolicyResult) error {
	out := runPolicyOutput{
		RunID:             result.RunID,
		RunStatus:         result.RunStatus,
		HasPolicies:       result.HasPolicies,
		PolicyChecks:      make([]policyCheckOutput, len(result.PolicyChecks)),
		PolicyEvaluations: make([]policyEvaluationOutput, len(result.PolicyEvaluations)),
	}

	for i, pc := range result.PolicyChecks {
		out.PolicyChecks[i] = policyCheckOutput{
			ID:             pc.ID,
			Status:         pc.Status,
			Scope:          pc.Scope,
			Overridable:    pc.Overridable,
			Passed:         pc.Passed,
			AdvisoryFailed: pc.AdvisoryFailed,
			SoftFailed:     pc.SoftFailed,
			HardFailed:     pc.HardFailed,
			TotalFailed:    pc.TotalFailed,
			Duration:       pc.Duration,
			Logs:           pc.Logs,
		}
	}

	for i, eval := range result.PolicyEvaluations {
		evalOut := policyEvaluationOutput{
			ID:              eval.ID,
			Status:          eval.Status,
			PolicyKind:      eval.PolicyKind,
			Stage:           eval.Stage,
			Passed:          eval.Passed,
			AdvisoryFailed:  eval.AdvisoryFailed,
			MandatoryFailed: eval.MandatoryFailed,
			Errored:         eval.Errored,
			PolicySets:      make([]policySetOutput, len(eval.PolicySets)),
		}

		for j, ps := range eval.PolicySets {
			psOut := policySetOutput{
				ID:              ps.ID,
				PolicySetName:   ps.PolicySetName,
				Overridable:     ps.Overridable,
				Error:           ps.Error,
				Passed:          ps.Passed,
				AdvisoryFailed:  ps.AdvisoryFailed,
				MandatoryFailed: ps.MandatoryFailed,
				Errored:         ps.Errored,
				Policies:        make([]policyDetailOutput, len(ps.Policies)),
			}

			for k, p := range ps.Policies {
				psOut.Policies[k] = policyDetailOutput{
					PolicyName:       p.PolicyName,
					EnforcementLevel: p.EnforcementLevel,
					Status:           p.Status,
					Description:      p.Description,
					Query:            p.Query,
					Output:           p.Output,
				}
			}

			evalOut.PolicySets[j] = psOut
		}

		out.PolicyEvaluations[i] = evalOut
	}

	return v.Output().RenderJSON(out)
}

func (v *RunPolicyView) renderTerminal(result *RunPolicyResult) error {
	// Run-level summary
	props := []PropertyPair{
		{Key: "Run ID", Value: result.RunID},
		{Key: "Run Status", Value: result.RunStatus},
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}

	if !result.HasPolicies {
		v.Output().Message("")
		v.Output().Message("No policy checks or evaluations found for this run.")
		return nil
	}

	// Legacy Policy Checks (Sentinel)
	if len(result.PolicyChecks) > 0 {
		v.Output().Message("")
		v.Output().Message("=== Policy Checks (Sentinel) ===")

		for _, pc := range result.PolicyChecks {
			if err := v.Output().RenderTags(fmt.Sprintf("Policy Check: %s", pc.ID), []PropertyPair{
				{Key: "Status", Value: pc.Status},
				{Key: "Scope", Value: pc.Scope},
				{Key: "Overridable", Value: pc.Overridable},
				{Key: "Passed", Value: pc.Passed},
				{Key: "Hard Failed", Value: pc.HardFailed},
				{Key: "Soft Failed", Value: pc.SoftFailed},
				{Key: "Advisory Failed", Value: pc.AdvisoryFailed},
				{Key: "Total Failed", Value: pc.TotalFailed},
			}); err != nil {
				return err
			}

			if pc.Logs != "" {
				v.Output().Message("")
				v.Output().Message("  --- Logs ---")
				v.Output().Message("%s", pc.Logs)
			}
		}
	}

	// Newer Policy Evaluations (OPA/Sentinel)
	if len(result.PolicyEvaluations) > 0 {
		v.Output().Message("")
		v.Output().Message("=== Policy Evaluations ===")

		for _, eval := range result.PolicyEvaluations {
			if err := v.Output().RenderTags(fmt.Sprintf("Evaluation: %s", eval.ID), []PropertyPair{
				{Key: "Kind", Value: eval.PolicyKind},
				{Key: "Stage", Value: eval.Stage},
				{Key: "Status", Value: eval.Status},
				{Key: "Passed", Value: eval.Passed},
				{Key: "Advisory Failed", Value: eval.AdvisoryFailed},
				{Key: "Mandatory Failed", Value: eval.MandatoryFailed},
				{Key: "Errored", Value: eval.Errored},
			}); err != nil {
				return err
			}

			// Render each policy set
			for _, ps := range eval.PolicySets {
				v.Output().Message("")
				v.Output().Message("  Policy Set: %s", ps.PolicySetName)
				v.Output().Message("  Overridable: %v", ps.Overridable)
				if ps.Error != "" {
					v.Output().Message("  Error: %s", ps.Error)
				}

				if len(ps.Policies) > 0 {
					headers := []string{"Policy Name", "Enforcement", "Status", "Description"}
					var rows [][]interface{}
					for _, p := range ps.Policies {
						rows = append(rows, []interface{}{
							p.PolicyName,
							p.EnforcementLevel,
							p.Status,
							p.Description,
						})
					}
					if err := v.Output().RenderTable(headers, rows); err != nil {
						return err
					}

					// Render per-policy logs if present
					for _, p := range ps.Policies {
						if p.Output != "" {
							v.Output().Message("")
							v.Output().Message("  --- %s Logs ---", p.PolicyName)
							v.Output().Message("%s", p.Output)
						}
					}
				}
			}
		}
	}

	return nil
}
