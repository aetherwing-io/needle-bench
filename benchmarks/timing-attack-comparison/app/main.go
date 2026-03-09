package main

import (
	"fmt"
	"os"

	"timing-attack/auth"
	"timing-attack/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: authserver <command> [args...]")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  serve             Start the auth server (blocks)")
		fmt.Fprintln(os.Stderr, "  verify <password>  Verify a password against stored hash")
		fmt.Fprintln(os.Stderr, "  test              Run security tests")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		server.Start()
	case "verify":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: authserver verify <password>")
			os.Exit(1)
		}
		store := auth.NewUserStore()
		store.AddUser("admin", "correct-horse-battery-staple")
		if store.Authenticate("admin", os.Args[2]) {
			fmt.Println("Authentication successful")
		} else {
			fmt.Println("Authentication failed")
			os.Exit(1)
		}
	case "test":
		result := auth.RunSecurityTest()
		if result {
			fmt.Println("PASS: All security tests passed")
		} else {
			fmt.Println("FAIL: Security test failed")
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
