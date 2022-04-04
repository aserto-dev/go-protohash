//go:build mage
// +build mage

package main

import (
	"github.com/aserto-dev/mage-loot/buf"
	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/magefile/mage/mg"
)

// All executes all build targets in dependency order.
func All() {
	mg.Deps(Deps, Gen, Lint, Test)
}

// Deps installs build dependencies.
func Deps() {
	deps.GetAllDeps()
}

// Gen executes code generators.
func Gen() error {
	if err := bufGenerate(); err != nil {
		return err
	}

	if err := bufLint(); err != nil {
		return err
	}
	return nil
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

func bufGenerate() error {
	return buf.Run(
		buf.AddArg("generate"),
	)
}

func bufLint() error {
	return buf.Run(
		buf.AddArg("lint"),
	)
}
