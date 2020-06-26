// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package chartrepo manages a list of configured helm repositories.
package chartrepo

import (
	"context"
	"fmt"
	"os/exec"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
)

type registry interface {
	Helm(ctx context.Context, version string) (string, bool, error)
}

// Add installs all specified Helm Chart repositories.
// https://helm.sh/docs/topics/chart_repository/
// helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com
// helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com --username my-username --password my-password
func Add(ctx context.Context, repos []config.HelmChartRepository, reg registry, logger *zap.Logger) error {
	helm, _, err := reg.Helm(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to find helm to add repos (%v)", err)
	}
	for _, repo := range repos {
		args := []string{"repo", "add", repo.Address}
		if repo.Username != "" || repo.Password != "" {
			args = append(args, "--username", repo.Username, "--password", repo.Password)
		}
		cmd := exec.CommandContext(ctx, helm, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to add chart repository %s: %s (%w)", repo.Name, string(out), err)
		}
		logger.Info(fmt.Sprintf("successfully added chart repository %s", repo.Name))
	}
	return nil
}