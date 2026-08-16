package ui

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogRequest(method, path string, statusCode int, note string) {
	timeStr := time.Now().Format("15:04:05")

	methodColor := color.New(color.FgCyan, color.Bold).SprintfFunc()
	switch method {
	case "GET":
		methodColor = color.New(color.FgBlue, color.Bold).SprintfFunc()
	case "POST":
		methodColor = color.New(color.FgGreen, color.Bold).SprintfFunc()
	case "PUT", "PATCH":
		methodColor = color.New(color.FgYellow, color.Bold).SprintfFunc()
	case "DELETE":
		methodColor = color.New(color.FgRed, color.Bold).SprintfFunc()
	}

	statusColor := color.New(color.FgGreen).SprintfFunc()
	if statusCode >= 400 && statusCode < 500 {
		statusColor = color.New(color.FgYellow).SprintfFunc()
	} else if statusCode >= 500 {
		statusColor = color.New(color.FgRed, color.Bold).SprintfFunc()
	}

	timeFormatted := color.New(color.FgHiBlack).Sprintf(timeStr)
	out := fmt.Sprintf("  %s %-7s %-32s -> %s", timeFormatted, methodColor(method), path, statusColor("%d", statusCode))
	if note != "" {
		out += fmt.Sprintf(" %s", color.New(color.FgHiBlack).Sprintf("(%s)", note))
	}

	fmt.Println(out)
}

func (l *Logger) Info(msg string) {
	fmt.Println(color.New(color.FgCyan).Sprintf("[INFO] %s", msg))
}

func (l *Logger) Success(msg string) {
	fmt.Println(color.New(color.FgGreen).Sprintf("[SUCCESS] %s", msg))
}

func (l *Logger) Error(msg string) {
	fmt.Println(color.New(color.FgRed, color.Bold).Sprintf("[ERROR] %s", msg))
}
