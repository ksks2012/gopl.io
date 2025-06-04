package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	i := 0
	for ; i < len(input); i++ {
		if (input[i] < '0' || input[i] > '9') && input[i] != '.' {
			break
		}
	}
	if i == 0 || i == len(input) {
		return "", fmt.Errorf("invalid input: %s", input)
	}
	value, err := strconv.ParseFloat(input[:i], 64)
	if err != nil {
		return "", err
	}
	unit := strings.ToLower(input[i:])

	switch unit {
	case "c":
		return fmt.Sprintf("%.2f°C = %.2f°F", value, value*9/5+32), nil
	case "f":
		return fmt.Sprintf("%.2f°F = %.2f°C", value, (value-32)*5/9), nil
	case "kg":
		return fmt.Sprintf("%.2fkg = %.2flb", value, value*2.20462), nil
	case "lb":
		return fmt.Sprintf("%.2flb = %.2fkg", value, value/2.20462), nil
	case "m":
		return fmt.Sprintf("%.2fm = %.2fft", value, value*3.28084), nil
	case "ft":
		return fmt.Sprintf("%.2fft = %.2fm", value, value/3.28084), nil
	default:
		return "", fmt.Errorf("unknown unit: %s", unit)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter value with unit (e.g., 100C, 10kg, 3m): ")
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			result, err := convert(line)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else if result != "" {
				fmt.Println(result)
			}
			fmt.Print("Enter value with unit (e.g., 100C, 10kg, 3m): ")
		}
	} else {
		for _, arg := range args {
			result, err := convert(arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else if result != "" {
				fmt.Println(result)
			}
		}
	}
}
