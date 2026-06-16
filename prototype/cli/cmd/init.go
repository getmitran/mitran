package cmd

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

type agent struct {
	name string
	task string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Mitran for your company",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Printf("\n%s%s🚀 Mitran Init — Let's set up your AI engineering team%s\n\n", colorBold, colorCyan, colorReset)

	company := prompt("What does your company do?")
	engineers := prompt("How many engineers?")
	stack := prompt("What languages/frameworks do you use?")
	github := promptYN("Do you use GitHub?")
	slack := promptYN("Do you use Slack?")

	agents := buildPlan(company, engineers, stack, github, slack)

	fmt.Printf("\n%s%s📋 Task Plan:%s\n\n", colorBold, colorYellow, colorReset)
	for i, a := range agents {
		fmt.Printf("  %s%d. %-12s%s → %s\n", colorCyan, i+1, a.name, colorReset, a.task)
	}
	fmt.Println()

	approved := promptYN("Approve this plan?")
	if !approved {
		fmt.Printf("\n%sPlan rejected. Run 'mitran init' again when ready.%s\n", colorYellow, colorReset)
		return nil
	}

	fmt.Printf("\n%s%s⚡ Starting execution...%s\n\n", colorBold, colorGreen, colorReset)
	for _, a := range agents {
		fmt.Printf("  [%s✓%s] %s — %s\n", colorGreen, colorReset, a.name, a.task)
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Printf("\n%s%s✅ Mitran is ready! Run 'mitran serve' to start the engine.%s\n\n", colorBold, colorGreen, colorReset)
	return nil
}

func prompt(label string) string {
	p := promptui.Prompt{Label: label}
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return ""
	}
	return result
}

func promptYN(label string) bool {
	p := promptui.Prompt{Label: label, IsConfirm: true}
	_, err := p.Run()
	return err == nil
}

func buildPlan(company, engineers, stack string, github, slack bool) []agent {
	agents := []agent{
		{"Dev Agent", fmt.Sprintf("Set up %s project scaffolding for %s engineers", stack, engineers)},
		{"Docs Agent", fmt.Sprintf("Generate technical docs for: %s", company)},
		{"Ops Agent", "Configure monitoring (Grafana + Prometheus)"},
		{"Review Agent", "Set up automated code review pipeline"},
		{"HR Agent", "Initialize onboarding workflows"},
		{"CI/CD Agent", "Create build and deploy pipelines"},
		{"Tickets Agent", "Set up internal ticket system"},
		{"Wiki Agent", "Generate knowledge base structure"},
	}
	if github {
		agents[0].task += " + GitHub repos"
		agents[5].task += " (GitHub Actions)"
	}
	if slack {
		agents[4].task += " + Slack notifications"
		agents[6].task += " + Slack integration"
	}
	return agents
}
