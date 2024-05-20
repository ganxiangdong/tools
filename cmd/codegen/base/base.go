// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Package base defines shared basic pieces of the ent command.
package base

import (
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/util/codegen"
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/util/codegen/gen"
	"entgo.io/ent/schema/field"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

// IDType is a custom ID implementation for pflag.
type IDType field.Type

// Set implements the Set method of the flag.Value interface.
func (t *IDType) Set(s string) error {
	switch s {
	case field.TypeInt.String():
		*t = IDType(field.TypeInt)
	case field.TypeInt64.String():
		*t = IDType(field.TypeInt64)
	case field.TypeUint.String():
		*t = IDType(field.TypeUint)
	case field.TypeUint64.String():
		*t = IDType(field.TypeUint64)
	case field.TypeString.String():
		*t = IDType(field.TypeString)
	default:
		return fmt.Errorf("invalid type %q", s)
	}
	return nil
}

// Type returns the type representation of the id option for help command.
func (IDType) Type() string {
	return fmt.Sprintf("%v", []field.Type{
		field.TypeInt,
		field.TypeInt64,
		field.TypeUint,
		field.TypeUint64,
		field.TypeString,
	})
}

// String returns the default value for the help command.
func (IDType) String() string {
	return field.TypeInt.String()
}

// GenerateCmd returns the generate command for ent/c packages.
func GenerateCmd(postRun ...func(*gen.Config)) *cobra.Command {
	var (
		cfg       gen.Config
		storage   string
		table     []string
		features  []string
		templates []string
		idtype    = IDType(field.TypeInt)
		cmd       = &cobra.Command{
			Use:   "generate [flags] path",
			Short: "generate go code for the schema directory",
			Example: examples(
				"codegen generate ./ent/schema --table platform",
			),
			Args: cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, path []string) {
				opts := []codegen.Option{
					codegen.Storage(storage),
					codegen.FeatureNames(features...),
				}
				for _, tmpl := range templates {
					typ := "dir"
					if parts := strings.SplitN(tmpl, "=", 2); len(parts) > 1 {
						typ, tmpl = parts[0], parts[1]
					}
					switch typ {
					case "dir":
						opts = append(opts, codegen.TemplateDir(tmpl))
					case "file":
						opts = append(opts, codegen.TemplateFiles(tmpl))
					case "glob":
						opts = append(opts, codegen.TemplateGlob(tmpl))
					default:
						log.Fatalln("unsupported template type", typ)
					}
				}
				// If the target directory is not inferred from
				// the schema path, resolve its package path.
				if cfg.Target != "" {
					pkgPath, err := PkgPath(DefaultConfig, cfg.Target)
					if err != nil {
						log.Fatalln(err)
					}
					cfg.Package = pkgPath
				}
				cfg.IDType = &field.TypeInfo{Type: field.Type(idtype)}

				if len(table) > 0 {
					cfg.Tables = append(cfg.Tables, table...)
				}

				cfg.Dirs = []string{ // 生成代码的目录
					"internal/biz/do",
					"internal/biz/bo",
					"internal/biz/repository",
					"internal/data/repositoryimpl",
				}
				cfg.Header = `
//   ██████╗ ██████╗ ██████╗ ███████╗ ██████╗ ███████╗███╗   ██╗
//  ██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔════╝ ██╔════╝████╗  ██║
//  ██║     ██║   ██║██║  ██║█████╗  ██║  ███╗█████╗  ██╔██╗ ██║
//  ██║     ██║   ██║██║  ██║██╔══╝  ██║   ██║██╔══╝  ██║╚██╗██║
//  ╚██████╗╚██████╔╝██████╔╝███████╗╚██████╔╝███████╗██║ ╚████║
//   ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝
`

				if err := codegen.Generate(path[0], &cfg, opts...); err != nil {
					log.Fatalln(err)
				}
				for _, fn := range postRun {
					fn(&cfg)
				}
			},
		}
	)
	cmd.Flags().Var(&idtype, "idtype", "type of the id field")
	cmd.Flags().StringVar(&storage, "storage", "sql", "storage driver to support in codegen")
	cmd.Flags().StringVar(&cfg.Header, "header", "", "override codegen header")
	cmd.Flags().StringVar(&cfg.Target, "target", "", "target directory for codegen")
	cmd.Flags().StringVar(&cfg.RelativePath, "path", "", "生成代码的目录，默认从项目根目录生成")
	cmd.Flags().BoolVar(&cfg.Overwrite, "overwrite", false, "是否直接覆盖原有生成代码")
	cmd.Flags().StringSliceVarP(&table, "table", "", nil, "需要生成的表名称，多个一起生成时使用")
	cmd.Flags().StringSliceVarP(&features, "feature", "", nil, "extend codegen with additional features")
	cmd.Flags().StringSliceVarP(&templates, "template", "", nil, "external templates to execute")
	// The --idtype flag predates the field.<Type>("id") option.
	// See, https://entgo.io/docs/schema-fields#id-field.
	cobra.CheckErr(cmd.Flags().MarkHidden("idtype"))
	return cmd
}

// examples formats the given examples to the cli.
func examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}
