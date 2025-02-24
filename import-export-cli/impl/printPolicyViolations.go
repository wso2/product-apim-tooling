/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package impl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Violation represents a violation
type Violation struct {
	Policy   string    `json:"policy"`
	Rulesets []Ruleset `json:"rulesets"`
}

// Ruleset represents a ruleset
type Ruleset struct {
	Ruleset        string          `json:"ruleset"`
	Type           string          `json:"type"`
	RuleViolations []RuleViolation `json:"rule-violations"`
}

// RuleViolation represents a rule violation
type RuleViolation struct {
	Path     string `json:"path"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// PrintViolations prints the violations in the given format
// If apiLoggingCmdFormat is "json", it prints the violations in JSON format
// If apiLoggingCmdFormat is "list", it prints the violations in list format
// If apiLoggingCmdFormat is "table" or empty, it prints the violations in table format
// @param violations: List of violations
// @param apiLoggingCmdFormat: Format to print the violations
// @return None
func PrintViolations(violations []Violation, apiLoggingCmdFormat string) {

	if apiLoggingCmdFormat == "json" {
		violationsJSON, err := json.MarshalIndent(violations, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling violations to JSON:", err)
			return
		}
		fmt.Println(string(violationsJSON))
		return

	} else if apiLoggingCmdFormat == "list" {
		// Print violations in list format
		fmt.Println("\nViolations:")
		for _, violation := range violations {
			fmt.Println("- Policy:", violation.Policy)
			for _, ruleset := range violation.Rulesets {
				fmt.Println("  Ruleset:", ruleset.Ruleset)
				for _, rule := range ruleset.RuleViolations {
					fmt.Printf("    - Path: %s | Message: %s | Severity: %s\n", rule.Path, rule.Message, rule.Severity)
				}
			}
		}
	} else if apiLoggingCmdFormat == "" || apiLoggingCmdFormat == "table" {
		// Print violations in table format
		for _, violation := range violations {

			for _, ruleset := range violation.Rulesets {
				fmt.Printf("\nPolicy: %s\nRuleset: %s\n\n", violation.Policy, ruleset.Ruleset)
				// Create table for each ruleset
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Path", "Message", "Severity"})
				// Append violations to the table
				for _, violation := range ruleset.RuleViolations {
					table.Append([]string{violation.Path, violation.Message, violation.Severity})
				}
				// Render the table
				table.Render()
			}
		}
	}
}
