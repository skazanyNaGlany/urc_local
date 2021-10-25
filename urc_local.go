package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/emersion/go-autostart"
)

const AppName = "URC_LOCAL"
const AppVersion = "0.1"
const Arg0 = "sh"
const Arg1 = "-x"
const Arg2 = "rc_local.sh"

func getFullAppName() string {
	return fmt.Sprintf("%v v%v", AppName, AppVersion)
}

func printAppName() {
	log.Println(
		getFullAppName())
	log.Println()
	log.Println("Universal rc_local.sh runner.")
	log.Println()
}

func printAppInfo() {
	log.Println("Run ./rc_local.sh as root. You must set root")
	log.Println("as owner of the file along with the SUID flag.")
	log.Println()
	log.Println("Otherwise you will get error:")
	log.Println("operation not permitted")
	log.Println()
	log.Println("This app can be run only on Linux.")
	log.Println()
	log.Println("How to set:")
	log.Printf("\tsudo chown root:root %v", os.Args[0])
	log.Printf("\tsudo chmod u+s %v", os.Args[0])
	log.Println()
}

func printUsages() {
	log.Printf("Usage: %v <option>", os.Args[0])

	log.Println()
	log.Println("Options:")

	log.Println("\t--install")
	log.Println("\t\t\t autorun with the system")
	log.Println()
	log.Println("\t--uninstall")
	log.Println("\t\t\t do not autorun with the system")
	log.Println()
	log.Println("\t--run")
	log.Println("\t\t\t run ./rc_local.sh")
	log.Println()
	log.Println("\t--status")
	log.Println("\t\t\t check if app is installed")
}

func shouldPrintUsages() bool {
	len_args := len(os.Args)

	return len_args != 2 || (len_args > 1 && os.Args[1] == "--help")
}

func getGoAutostartApp() (*autostart.App, error) {
	executable, err := os.Executable()

	if err != nil {
		return nil, err
	}

	fullAppName := getFullAppName()
	app := autostart.App{
		Name:        fullAppName,
		DisplayName: fullAppName,
		Exec:        []string{executable},
	}

	return &app, nil
}

func checkInstalled() {
	app, err := getGoAutostartApp()

	if err != nil {
		log.Fatal(err)
	}

	if app.IsEnabled() {
		log.Fatal("App is installed.")
	} else {
		log.Fatal("App is not installed.")
	}
}

func printAppStatus() {
	checkInstalled()
}

func installAutorun() {
	app, err := getGoAutostartApp()

	if err != nil {
		log.Fatal(err)
	}

	if app.IsEnabled() {
		log.Fatal("App already installed.")
	}

	err = app.Enable()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("App installed.")
}

func uninstallAutorun() {
	app, err := getGoAutostartApp()

	if err != nil {
		log.Fatal(err)
	}

	if !app.IsEnabled() {
		log.Fatal("App is not installed.")
	}

	err = app.Disable()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("App uninstalled.")
}

func changeCurrentWorkingDir() {
	exeDir := filepath.Dir(os.Args[0])
	os.Chdir(exeDir)
}

func checkPlatform() {
	if runtime.GOOS != "linux" {
		log.Fatalln("This app can be used only on Linux.")
	}
}

func runRcLocal() {
	err := syscall.Setreuid(0, 0)

	if err != nil {
		log.Fatal(err)
	}

	output, err := exec.Command(Arg0, Arg1, Arg2).Output()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(output))
}

func main() {
	printAppName()
	checkPlatform()
	changeCurrentWorkingDir()

	if shouldPrintUsages() {
		printAppInfo()
		printUsages()

		os.Exit(1)
	}

	if os.Args[1] == "--status" {
		printAppStatus()
	} else if os.Args[1] == "--install" {
		installAutorun()
	} else if os.Args[1] == "--uninstall" {
		uninstallAutorun()
	} else if os.Args[1] == "--run" {
		runRcLocal()
	} else {
		printAppInfo()
		printUsages()
	}
}
