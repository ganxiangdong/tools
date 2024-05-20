// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := exec.Command("go", "run", "codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/cmd/codegen", "generate", "./ent/schema", "target", "./xxx")
	stderr := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	require.NoError(t, cmd.Run(), stderr.String())
}
