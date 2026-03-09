package main

import (
	"fmt"
	"os"

	"split-brain/cluster"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: split-brain <simulate|check>")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  simulate  - Run leader election simulation with partition")
		fmt.Fprintln(os.Stderr, "  check     - Check if split-brain is possible in current implementation")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "simulate":
		result := cluster.RunSimulation()
		if result.SplitBrain {
			fmt.Printf("SPLIT-BRAIN DETECTED: nodes %v all believe they are leader\n", result.Leaders)
			fmt.Printf("Partition: %v | %v\n", result.PartitionA, result.PartitionB)
			os.Exit(1)
		}
		fmt.Printf("OK: single leader %d elected\n", result.Leaders[0])
	case "check":
		safe := cluster.CheckSafety()
		if !safe {
			fmt.Println("UNSAFE: split-brain is possible under partition")
			os.Exit(1)
		}
		fmt.Println("SAFE: leader election handles partitions correctly")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
