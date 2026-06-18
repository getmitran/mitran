package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type process struct {
	name string
	cmd  *exec.Cmd
	port string
}

var (
	serveDashboard bool
	enginePort     string
	workerPort     string
	dashboardPort  string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Mitran engine, worker, and optionally dashboard",
	Long:  "Starts the Go engine, Python agent worker, and optionally the React dashboard dev server. Watches for crashes and prints a status table.",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().BoolVar(&serveDashboard, "dashboard", false, "Also start the dashboard dev server")
	serveCmd.Flags().StringVar(&enginePort, "engine-port", "7780", "Engine HTTP port")
	serveCmd.Flags().StringVar(&workerPort, "worker-port", "8888", "Agent worker port")
	serveCmd.Flags().StringVar(&dashboardPort, "dashboard-port", "5173", "Dashboard dev server port")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	var procs []*process
	var mu sync.Mutex

	// Engine
	engineCmd := exec.Command("go", "run", ".")
	engineCmd.Dir = "../server"
	engineCmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", enginePort))
	engineCmd.Stdout = os.Stdout
	engineCmd.Stderr = os.Stderr
	procs = append(procs, &process{name: "engine", cmd: engineCmd, port: enginePort})

	// Worker
	workerCmd := exec.Command("python3", "-m", "worker.main")
	workerCmd.Dir = "../agent-worker"
	workerCmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", workerPort))
	workerCmd.Stdout = os.Stdout
	workerCmd.Stderr = os.Stderr
	procs = append(procs, &process{name: "worker", cmd: workerCmd, port: workerPort})

	// Dashboard (optional)
	if serveDashboard {
		dashCmd := exec.Command("npm", "run", "dev", "--", "--port", dashboardPort)
		dashCmd.Dir = "../dashboard"
		dashCmd.Stdout = os.Stdout
		dashCmd.Stderr = os.Stderr
		procs = append(procs, &process{name: "dashboard", cmd: dashCmd, port: dashboardPort})
	}

	// Start all processes
	for _, p := range procs {
		if err := p.cmd.Start(); err != nil {
			// Kill already-started processes
			for _, started := range procs {
				if started.cmd.Process != nil {
					_ = started.cmd.Process.Kill()
				}
			}
			return fmt.Errorf("failed to start %s: %w", p.name, err)
		}
	}

	// Print status table
	time.Sleep(500 * time.Millisecond)
	printStatusTable(procs)

	// Watch for crashes
	crashCh := make(chan *process, len(procs))
	for _, p := range procs {
		go func(proc *process) {
			_ = proc.cmd.Wait()
			mu.Lock()
			defer mu.Unlock()
			crashCh <- proc
		}(p)
	}

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		fmt.Printf("\n⚡ Received %s, shutting down...\n", sig)
	case crashed := <-crashCh:
		fmt.Printf("\n💥 %s crashed (exit: %s)\n", crashed.name, crashed.cmd.ProcessState)
	}

	// Shutdown all
	for _, p := range procs {
		if p.cmd.Process != nil {
			_ = p.cmd.Process.Signal(syscall.SIGTERM)
		}
	}
	time.Sleep(time.Second)
	for _, p := range procs {
		if p.cmd.Process != nil {
			_ = p.cmd.Process.Kill()
		}
	}

	return nil
}

func printStatusTable(procs []*process) {
	fmt.Println()
	fmt.Println("┌─────────────┬────────┬───────────────────────┐")
	fmt.Println("│ Service     │ Port   │ Status                │")
	fmt.Println("├─────────────┼────────┼───────────────────────┤")
	for _, p := range procs {
		status := "✅ running"
		if p.cmd.ProcessState != nil {
			status = "❌ exited"
		}
		fmt.Printf("│ %-11s │ %-6s │ %-21s │\n", p.name, p.port, status)
	}
	fmt.Println("└─────────────┴────────┴───────────────────────┘")
	fmt.Println()
}
