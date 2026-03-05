// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/output"
)

// FetchRunPolicyDetails retrieves all policy check and evaluation details for a run.
// When fetchLogs is true, raw Sentinel log output is fetched for each policy check.
func FetchRunPolicyDetails(c *client.TfxClient, runID string, fetchLogs bool) (*view.RunPolicyResult, error) {
	log := output.Get().Logger()
	log.Debug("Fetching policy details for run", "runID", runID)

	// Fetch the run with task stages included
	run, err := c.Client.Runs.ReadWithOptions(c.Context, runID, &tfe.RunReadOptions{
		Include: []tfe.RunIncludeOpt{tfe.RunTaskStages},
	})
	if err != nil {
		log.Error("Failed to fetch run", "runID", runID, "error", err)
		return nil, err
	}

	result := &view.RunPolicyResult{
		RunID:     run.ID,
		RunStatus: string(run.Status),
	}

	// Fetch legacy policy checks
	policyChecks, err := fetchLegacyPolicyChecks(c, runID, fetchLogs)
	if err != nil {
		log.Error("Failed to fetch policy checks", "runID", runID, "error", err)
		return nil, err
	}
	result.PolicyChecks = policyChecks

	// Fetch newer policy evaluations from task stages
	evaluations, err := fetchPolicyEvaluations(c, run.TaskStages, fetchLogs)
	if err != nil {
		log.Error("Failed to fetch policy evaluations", "runID", runID, "error", err)
		return nil, err
	}
	result.PolicyEvaluations = evaluations

	result.HasPolicies = len(result.PolicyChecks) > 0 || len(result.PolicyEvaluations) > 0

	log.Debug("Policy details fetched",
		"runID", runID,
		"policyChecks", len(result.PolicyChecks),
		"policyEvaluations", len(result.PolicyEvaluations),
	)
	return result, nil
}

// fetchLegacyPolicyChecks lists and reads all legacy Sentinel policy checks for a run
func fetchLegacyPolicyChecks(c *client.TfxClient, runID string, fetchLogs bool) ([]view.PolicyCheckDetail, error) {
	log := output.Get().Logger()

	allChecks, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.PolicyCheck, *client.Pagination, error) {
		opts := &tfe.PolicyCheckListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}
		res, err := c.Client.PolicyChecks.List(c.Context, runID, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
	if err != nil {
		return nil, err
	}

	var details []view.PolicyCheckDetail
	for _, pc := range allChecks {
		// Read full details to get the Result struct
		full, err := c.Client.PolicyChecks.Read(c.Context, pc.ID)
		if err != nil {
			log.Error("Failed to read policy check", "id", pc.ID, "error", err)
			continue
		}

		detail := view.PolicyCheckDetail{
			ID:     full.ID,
			Status: string(full.Status),
			Scope:  string(full.Scope),
		}
		if full.Actions != nil {
			detail.Overridable = full.Actions.IsOverridable
		}
		if full.Result != nil {
			detail.Passed = full.Result.Passed
			detail.AdvisoryFailed = full.Result.AdvisoryFailed
			detail.SoftFailed = full.Result.SoftFailed
			detail.HardFailed = full.Result.HardFailed
			detail.TotalFailed = full.Result.TotalFailed
			detail.Duration = full.Result.Duration
		}

		if fetchLogs {
			logReader, err := c.Client.PolicyChecks.Logs(c.Context, full.ID)
			if err != nil {
				log.Error("Failed to fetch policy check logs", "id", full.ID, "error", err)
			} else if logReader != nil {
				logBytes, err := io.ReadAll(logReader)
				if err != nil {
					log.Error("Failed to read policy check logs", "id", full.ID, "error", err)
				} else {
					detail.Logs = string(logBytes)
				}
			}
		}

		details = append(details, detail)
	}

	return details, nil
}

// fetchPolicyEvaluations fetches policy evaluations from task stages
func fetchPolicyEvaluations(c *client.TfxClient, taskStages []*tfe.TaskStage, fetchLogs bool) ([]view.PolicyEvaluationDetail, error) {
	log := output.Get().Logger()

	var details []view.PolicyEvaluationDetail

	for _, ts := range taskStages {
		allEvals, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.PolicyEvaluation, *client.Pagination, error) {
			opts := &tfe.PolicyEvaluationListOptions{
				ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			}
			res, err := c.Client.PolicyEvaluations.List(c.Context, ts.ID, opts)
			if err != nil {
				return nil, nil, err
			}
			return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
		})
		if err != nil {
			log.Error("Failed to list policy evaluations", "taskStageID", ts.ID, "error", err)
			return nil, err
		}

		for _, eval := range allEvals {
			detail := view.PolicyEvaluationDetail{
				ID:         eval.ID,
				Status:     string(eval.Status),
				PolicyKind: string(eval.PolicyKind),
				Stage:      string(ts.Stage),
			}
			if eval.ResultCount != nil {
				detail.Passed = eval.ResultCount.Passed
				detail.AdvisoryFailed = eval.ResultCount.AdvisoryFailed
				detail.MandatoryFailed = eval.ResultCount.MandatoryFailed
				detail.Errored = eval.ResultCount.Errored
			}

			// Fetch policy set outcomes for this evaluation
			policySets, err := fetchPolicySetOutcomes(c, eval.ID, fetchLogs)
			if err != nil {
				log.Error("Failed to fetch policy set outcomes", "evaluationID", eval.ID, "error", err)
				return nil, err
			}
			detail.PolicySets = policySets

			details = append(details, detail)
		}
	}

	return details, nil
}

// fetchPolicySetOutcomes fetches all policy set outcomes for a policy evaluation.
// When fetchLogs is true, a raw API call is made to the list endpoint to retrieve
// the output.print field that go-tfe does not expose.
func fetchPolicySetOutcomes(c *client.TfxClient, evaluationID string, fetchLogs bool) ([]view.PolicySetDetail, error) {
	log := output.Get().Logger()

	allOutcomes, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.PolicySetOutcome, *client.Pagination, error) {
		opts := &tfe.PolicySetOutcomeListOptions{
			ListOptions: &tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}
		res, err := c.Client.PolicySetOutcomes.List(c.Context, evaluationID, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
	if err != nil {
		return nil, err
	}

	// TODO: go-tfe's Outcome struct is missing the output field (contains print logs).
	// Once go-tfe adds Output to tfe.Outcome, remove this raw API call and
	// read o.Output directly in the loop below. Also delete fetchEvaluationOutputs
	// and policySetOutcomeListRaw.
	// Ref: https://github.com/hashicorp/go-tfe/issues/XXX
	var outcomeOutputs map[string]map[string]string // outcomeID -> policyName -> output
	if fetchLogs {
		var err error
		outcomeOutputs, err = fetchEvaluationOutputs(c, evaluationID)
		if err != nil {
			log.Error("Failed to fetch evaluation outputs", "evaluationID", evaluationID, "error", err)
		}
	}

	var details []view.PolicySetDetail
	for _, pso := range allOutcomes {
		detail := view.PolicySetDetail{
			ID:              pso.ID,
			PolicySetName:   pso.PolicySetName,
			Error:           pso.Error,
			Passed:          pso.ResultCount.Passed,
			AdvisoryFailed:  pso.ResultCount.AdvisoryFailed,
			MandatoryFailed: pso.ResultCount.MandatoryFailed,
			Errored:         pso.ResultCount.Errored,
		}
		if pso.Overridable != nil {
			detail.Overridable = *pso.Overridable
		}

		for _, o := range pso.Outcomes {
			pd := view.PolicyDetail{
				PolicyName:       o.PolicyName,
				EnforcementLevel: string(o.EnforcementLevel),
				Status:           o.Status,
				Description:      o.Description,
				Query:            o.Query,
			}
			// TODO: Replace with direct o.Output read once go-tfe exposes it
			if outcomeOutputs != nil {
				if outputs, ok := outcomeOutputs[pso.ID]; ok {
					pd.Output = outputs[o.PolicyName]
				}
			}
			detail.Policies = append(detail.Policies, pd)
		}

		details = append(details, detail)
	}

	return details, nil
}

// policySetOutcomeListRaw represents the raw JSON API list response for policy set outcomes.
// This is needed because go-tfe's Outcome struct does not include the output field.
type policySetOutcomeListRaw struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Outcomes []struct {
				PolicyName string `json:"policy_name"`
				Output     []struct {
					Print string `json:"print"`
				} `json:"output"`
			} `json:"outcomes"`
		} `json:"attributes"`
	} `json:"data"`
}

// fetchEvaluationOutputs makes a direct API call to the policy-set-outcomes list
// endpoint to retrieve the output.print field for each policy in each outcome.
// Returns a nested map: outcomeID -> policyName -> output.
func fetchEvaluationOutputs(c *client.TfxClient, evaluationID string) (map[string]map[string]string, error) {
	log := output.Get().Logger()
	apiURL := fmt.Sprintf("https://%s/api/v2/policy-evaluations/%s/policy-set-outcomes", c.Hostname, evaluationID)

	req, err := http.NewRequestWithContext(c.Context, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching policy set outcomes for evaluation %s", resp.StatusCode, evaluationID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug("Raw policy set outcomes list response", "evaluationID", evaluationID, "body", string(body))

	var raw policySetOutcomeListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	result := make(map[string]map[string]string)
	for _, item := range raw.Data {
		outputs := make(map[string]string)
		for _, o := range item.Attributes.Outcomes {
			var prints []string
			for _, out := range o.Output {
				if strings.TrimSpace(out.Print) != "" {
					prints = append(prints, out.Print)
				}
			}
			if len(prints) > 0 {
				outputs[o.PolicyName] = strings.Join(prints, "")
				log.Debug("Found policy output", "outcomeID", item.ID, "policyName", o.PolicyName)
			}
		}
		if len(outputs) > 0 {
			result[item.ID] = outputs
		}
	}

	return result, nil
}
