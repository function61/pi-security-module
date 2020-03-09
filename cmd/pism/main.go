package main

import (
	"fmt"
	"github.com/function61/gokit/dynversion"
	"github.com/function61/gokit/logex"
	"github.com/function61/gokit/ossignal"
	"github.com/function61/gokit/stopper"
	"github.com/function61/gokit/systemdinstaller"
	"github.com/function61/pi-security-module/pkg/httpserver"
	"github.com/function61/pi-security-module/pkg/keepassimport"
	"github.com/function61/pi-security-module/pkg/sshagent"
	"github.com/function61/pi-security-module/pkg/state"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func serverEntrypoint() *cobra.Command {
	server := &cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			rootLogger := log.New(os.Stderr, "", log.LstdFlags)

			logl := logex.Levels(logex.Prefix("serverEntrypoint", rootLogger))

			logl.Info.Printf("%s starting", dynversion.Version)
			defer logl.Info.Printf("Stopped")

			workers := stopper.NewManager()

			go func() {
				logl.Info.Printf("Received signal %s; stopping", <-ossignal.InterruptOrTerminate())

				workers.StopAllWorkersAndWait()
			}()

			exitIfError(httpserver.Run(workers.Stopper(), rootLogger))
		},
	}

	server.AddCommand(&cobra.Command{
		Use:   "init-config [adminUsername] [adminPassword]",
		Short: "Initializes configuration file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			exitIfError(state.InitConfig(args[0], args[1]))
		},
	})

	server.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Installs systemd unit file to make pi-security-module start on system boot",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			service := systemdinstaller.SystemdServiceFile(
				"pi-security-module",
				"Pi security module",
				systemdinstaller.Args("server"))

			exitIfError(systemdinstaller.Install(service))

			fmt.Println(systemdinstaller.GetHints(service))
		},
	})

	return server
}

func main() {
	rootCmd := &cobra.Command{
		Use:     os.Args[0],
		Short:   "Software for a hardware security module",
		Version: dynversion.Version,
	}

	rootCmd.AddCommand(serverEntrypoint())

	rootCmd.AddCommand(sshagent.Entrypoint())

	rootCmd.AddCommand(keepassimport.Entrypoint())

	exitIfError(rootCmd.Execute())
}

func exitIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
