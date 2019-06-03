package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/developgo/skm"
	cli "gopkg.in/urfave/cli.v1"
)

type resticConfig struct {
	Repository   string `json:"repository"`
	PasswordFile string `json:"password_file"`
}

func mustHaveRestic(env *skm.Environment) {
	if env.ResticPath == "" {
		skm.Fatalf("Restic not available. See https://restic.net/ for installation instructions.\n")
	}
}

func createAndOpenResticConfig(path string) (*os.File, error) {
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	cfg := resticConfig{
		Repository:   filepath.Join(os.Getenv("HOME"), ".skm-backups"),
		PasswordFile: filepath.Join(os.Getenv("HOME"), ".skm-backups.passwd"),
	}
	if err := json.NewEncoder(fp).Encode(&cfg); err != nil {
		fp.Close()
		return nil, err
	}
	fp.Close()
	return os.Open(path)
}

func mustLoadOrCreateResticSettings(env *skm.Environment, ctx *cli.Context) *resticConfig {
	configPath := filepath.Join(env.StorePath, "restic.json")
	fp, err := os.Open(configPath)
	if os.IsNotExist(err) {
		fp, err = createAndOpenResticConfig(configPath)
	}
	if err != nil {
		skm.Fatalf("Failed to open %s: %s\n", configPath, err.Error())
		return nil
	}
	defer fp.Close()
	cfg := resticConfig{}
	if err := json.NewDecoder(fp).Decode(&cfg); err != nil {
		skm.Fatalf("Failed to parse %s: %s\n", configPath, err.Error())
		return nil
	}
	return &cfg
}

func ensureInitializedResticRepo(cfg *resticConfig, env *skm.Environment) {
	if _, err := os.Stat(cfg.PasswordFile); err != nil {
		if !os.IsNotExist(err) {
			skm.Fatalf("Failed to check restic password file: %s\n", err.Error())
		}
		skm.Fatalf("Please create %s with the password you want to use for your restic backups.\n", cfg.PasswordFile)
	}

	if _, err := os.Stat(filepath.Join(cfg.Repository, "config")); err != nil {
		if !os.IsNotExist(err) {
			skm.Fatalf("Failed to check restic repository: %s\n", err.Error())
			return
		}
		log.Printf("Restic repository (%s) doesn't exist yet. Creating it...", cfg.Repository)
		if err := exec.Command(env.ResticPath, "init", "--password-file", cfg.PasswordFile, "--repo", cfg.Repository).Run(); err != nil {
			skm.Fatalf("Failed to initialize restic repository (%s): %s\n", cfg.Repository, err.Error())
			return
		}
	}
}
