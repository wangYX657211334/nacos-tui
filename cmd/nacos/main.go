package main

import (
	"database/sql"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/storage"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/internal/ui/home"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	logPath := filepath.Join(userHome, ".nacos-tui", "nacos-tui.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		log.Panic(err)
	}
	file, err := os.Create(logPath)
	if err != nil {
		log.Panic(err)
	}
	log.Default().SetOutput(file)
}

func main() {
	db, err := storage.NewConnection()
	if err != nil {
		log.Panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)
	homeModel, err := home.NewHomeModel(db)
	if err != nil {
		log.Panic(err)
	}
	program := tea.NewProgram(homeModel, tea.WithAltScreen())
	event.RegisterSubscribe(event.RefreshScreenEvent, func(a ...any) {
		program.Send(base.RefreshScreenMsg)
	})
	tea.EnterAltScreen()
	event.RegisterSubscribe(event.QuitEvent, func(a ...any) {
		err := program.ReleaseTerminal()
		if err != nil {
			log.Panic(err)
		}
		os.Exit(0)
	})

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
