// Package main provides the entry point for the secret-env-manager (sem) CLI application.
// sem manages secrets from cloud providers (AWS Secrets Manager, Google Cloud Secret Manager)
// and exports them as environment variables.
package main

import (
	"fmt"
	"os"

	"github.com/gumi-tsd/secret-env-manager/cmd"
	"github.com/gumi-tsd/secret-env-manager/internal/formatting"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/logging"
	"github.com/urfave/cli/v2"
)

// These variables are set during build time using -ldflags
var (
	version    = "dev"
	buildTime  = "unknown"
	commitHash = "unknown"
)

// CLI flag definitions
var (
	inputFlag = &cli.StringFlag{
		Name:    "input",
		Aliases: []string{"i"},
		Usage:   "Input configuration file name",
		Value:   ".env",
	}
	exportFlag = &cli.BoolFlag{
		Name:    "with-export",
		Aliases: []string{"e"},
		Usage:   "Add 'export' prefix when displaying environment variables to standard output.",
	}
	endpointURLFlag = &cli.StringFlag{
		Name:  "endpoint-url",
		Usage: "Custom endpoint URL for AWS Secrets Manager",
		Value: "",
	}
	noQuotesFlag = &cli.BoolFlag{
		Name:    "no-quotes",
		Aliases: []string{"q"},
		Usage:   "Don't wrap values in single quotes when generating the cache file",
		Value:   false,
	}
	noExpandJsonFlag = &cli.BoolFlag{
		Name:  "no-expand-json",
		Usage: "Don't automatically expand JSON values into separate environment variables",
		Value: false,
	}
)

// Logger instance
var logger = logging.DefaultLogger()

func main() {
	// Create app configuration
	app := createApp()

	// Run app and handle potential errors in a functional way
	runResult := runApp(app, os.Args)
	if runResult.IsFailure() {
		logFatalError(runResult.GetError().Error())
	}
}

// createApp sets up the CLI application configuration
// Pure function: Always returns the same output for the same inputs
func createApp() *cli.App {
	return &cli.App{
		Name:    "secret-env-manager (sem)",
		Usage:   "manage secret environment variables",
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, commitHash),
		Commands: []*cli.Command{
			{
				Name: "init",
				Usage: "This command interactively lists available secrets from cloud providers and generates URIs for use in environment files.\n" +
					"It supports both AWS Secrets Manager and Google Cloud Secret Manager.\n" +
					"For AWS, it requires AWS_PROFILE and AWS_REGION environment variables to be set.\n" +
					"For Google Cloud, it requires GOOGLE_CLOUD_PROJECT environment variable to be set.\n" +
					"Custom AWS endpoints can be specified using the --endpoint-url flag for local development with services like LocalStack.\n",
				Action: cmd.Init,
				Flags: []cli.Flag{
					endpointURLFlag,
				},
			},
			{
				Name: "load",
				Usage: "This command reads the cached env file specified in the input and displays it to standard output.\n" +
					"By using the -e option, you can output each line with 'export ' prefixed as 'export ENV=VALUE'.\n" +
					"If the cache file does not exist, you will be prompted to run update.\n",
				Action: cmd.Load,
				Flags: []cli.Flag{
					inputFlag,
					exportFlag,
				},
			},
			{
				Name: "update",
				Usage: "This command retrieves secrets from cloud providers based on the specified env file and caches them in a file named .cache.$input.\n" +
					"If the cache file is not excluded from version control (git tracking) at runtime, a warning will be displayed and the process will exit with code 1 before generating the file.\n" +
					"Please ensure the file is added to gitignore before running this command.\n",
				Action: cmd.Update,
				Flags: []cli.Flag{
					inputFlag,
					endpointURLFlag,
					noQuotesFlag,
					noExpandJsonFlag,
				},
			},
		},
	}
}

// runApp runs the CLI application with the given arguments
// Returns a Result monad to handle errors in a functional way
func runApp(app *cli.App, args []string) functional.Result[bool] {
	if err := app.Run(args); err != nil {
		return functional.Failure[bool](err)
	}
	return functional.Success(true)
}

// logFatalError logs a fatal error and exits the program
// This function wraps the side effect of logging and exiting
func logFatalError(message string) {
	errorMsg := formatting.Error("Error running the application: %s", message)
	logger.Error("%s", errorMsg)
	os.Exit(1)
}
