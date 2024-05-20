// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/cmd/codegen/base"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	log.SetFlags(0)
	cmd := &cobra.Command{Use: "codegen"}
	cmd.AddCommand(
		base.GenerateCmd(),
	)
	_ = cmd.Execute()
}
