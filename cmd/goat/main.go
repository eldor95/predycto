// Copyright 2022 gorse Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"errors"
	"fmt"
	"github.com/scylladb/go-set/strset"
	"github.com/spf13/cobra"
	"modernc.org/cc/v3"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var supportedTypes = strset.New("int64_t")

type TranslateUnit struct {
	Source     string
	Assembly   string
	Object     string
	GoAssembly string
	Go         string
	Package    string
	Options    []string
}

func NewTranslateUnit(source string, outputDir string, options ...string) TranslateUnit {
	sourceExt := filepath.Ext(source)
	noExtSourcePath := source[:len(source)-len(sourceExt)]
	noExtSourceBase := filepath.Base(noExtSourcePath)
	return TranslateUnit{
		Source:     source,
		Assembly:   noExtSourcePath + ".s",
		Object:     noExtSourcePath + ".o",
		GoAssembly: filepath.Join(outputDir, noExtSourceBase+".s"),
		Go:         filepath.Join(outputDir, noExtSourceBase+".go"),
		Package:    filepath.Base(outputDir),
		Options:    options,
	}
}

// parseSource parse C source file and extract functions declarations.
func (t *TranslateUnit) parseSource() ([]Function, error) {
	// TODO: Add include paths
	ast, err := cc.Parse(&cc.Config{}, nil, []string{
		"/usr/include",
		"/usr/local/lib/gcc/x86_64-pc-linux-gnu/12.1.1/include"}, []cc.Source{{Name: t.Source}})
	if err != nil {
		return nil, err
	}
	var functions []Function
	for _, nodes := range ast.Scope {
		if len(nodes) != 1 || nodes[0].Position().Filename != t.Source {
			continue
		}
		node := nodes[0]
		if declarator, ok := node.(*cc.Declarator); ok {
			funcIdent := declarator.DirectDeclarator
			if funcIdent.Case != cc.DirectDeclaratorFuncParam {
				continue
			}
			if function, err := convertFunction(funcIdent); err != nil {
				return nil, err
			} else {
				functions = append(functions, function)
			}
		}
	}
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].Position < functions[j].Position
	})
	return functions, nil
}

func (t *TranslateUnit) generateGoStubs(functions []Function) error {
	// generate code
	var builder strings.Builder
	builder.WriteString("//go:build !noasm\n")
	builder.WriteString("// AUTO-GENERATED BY GOAT -- DO NOT EDIT\n\n")
	builder.WriteString(fmt.Sprintf("package %v\n\n", t.Package))
	builder.WriteString("import \"unsafe\"\n")
	for _, function := range functions {
		builder.WriteString("\n//go:noescape\n")
		builder.WriteString(fmt.Sprintf("func %v(%s unsafe.Pointer)\n",
			function.Name, strings.Join(function.Parameters, ", ")))
	}

	// write file
	f, err := os.Create(t.Go)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}(f)
	_, err = f.WriteString(builder.String())
	return err
}

func (t *TranslateUnit) compile(args ...string) error {
	args = append(args, "-mno-red-zone", "-mstackrealign", "-mllvm", "-inline-threshold=1000",
		"-fno-asynchronous-unwind-tables", "-fno-exceptions", "-fno-rtti")
	_, err := runCommand("clang", append([]string{"-S", "-c", t.Source, "-o", t.Assembly}, args...)...)
	if err != nil {
		return err
	}
	_, err = runCommand("clang", append([]string{"-c", t.Source, "-o", t.Object}, args...)...)
	return err
}

func (t *TranslateUnit) Translate() error {
	functions, err := t.parseSource()
	if err != nil {
		return err
	}
	if err = t.generateGoStubs(functions); err != nil {
		return err
	}
	if err = t.compile(t.Options...); err != nil {
		return err
	}
	assembly, err := parseAssembly(t.Assembly)
	if err != nil {
		return err
	}
	dump, _ := runCommand("objdump", "-d", t.Object)
	err = parseObjectDump(dump, assembly)
	if err != nil {
		return err
	}
	for i, name := range functions {
		functions[i].Lines = assembly[name.Name]
	}
	return generateGoAssembly(t.GoAssembly, functions)
}

type Function struct {
	Name       string
	Position   int
	Parameters []string
	Lines      []Line
}

// convertFunction extracts the function definition from cc.DirectDeclarator.
func convertFunction(declarator *cc.DirectDeclarator) (Function, error) {
	params, err := convertFunctionParameters(declarator.ParameterTypeList.ParameterList)
	if err != nil {
		return Function{}, err
	}
	return Function{
		Name:       declarator.DirectDeclarator.Token.Value.String(),
		Position:   declarator.Position().Line,
		Parameters: params,
	}, nil
}

func convertFunctionParameters(params *cc.ParameterList) ([]string, error) {
	declaration := params.ParameterDeclaration
	paramName := declaration.Declarator.DirectDeclarator.Token.Value
	paramType := declaration.DeclarationSpecifiers.TypeSpecifier.Token.Value
	isPointer := declaration.Declarator.Pointer != nil
	if !isPointer && !supportedTypes.Has(paramType.String()) {
		position := declaration.Position()
		return nil, fmt.Errorf("%v:%v:%v: error: unsupported type: %v\n",
			position.Filename, position.Line, position.Column, paramType)
	}
	paramNames := []string{paramName.String()}
	if params.ParameterList != nil {
		if nextParamNames, err := convertFunctionParameters(params.ParameterList); err != nil {
			return nil, err
		} else {
			paramNames = append(paramNames, nextParamNames...)
		}
	}
	return paramNames, nil
}

// runCommand runs a command and extract its output.
func runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			return "", errors.New(string(output))
		} else {
			return "", err
		}
	}
	return string(output), nil
}

var command = &cobra.Command{
	Use:  "goat source [-o output_directory]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.PersistentFlags().GetString("output")
		if output == "" {
			var err error
			if output, err = os.Getwd(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		var options []string
		machineOptions, _ := cmd.PersistentFlags().GetStringSlice("machine-option")
		for _, m := range machineOptions {
			options = append(options, "-m"+m)
		}
		optimizeLevel, _ := cmd.PersistentFlags().GetInt("optimize-level")
		options = append(options, fmt.Sprintf("-O%d", optimizeLevel))
		file := NewTranslateUnit(args[0], output, options...)
		if err := file.Translate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	command.PersistentFlags().StringP("output", "o", "", "output directory of generated files")
	command.PersistentFlags().StringSliceP("machine-option", "m", nil, "machine option for clang")
	command.PersistentFlags().IntP("optimize-level", "O", 0, "optimization level for clang")
}

func main() {
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
