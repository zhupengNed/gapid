// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package executor provides an interface for executing compiled API programs.
package executor

import (
	"unsafe"

	"github.com/google/gapid/core/codegen"
	"github.com/google/gapid/gapil/compiler"
)

// Executor is used to create execution environments for a compiled program.
// Use New() to create Executors, do not create directly.
type Executor struct {
	program        *compiler.Program
	exec           *codegen.Executor
	createContext  unsafe.Pointer
	destroyContext unsafe.Pointer
	globalsSize    uint64
	cmdFunctions   map[string]unsafe.Pointer
}

// New returns a new and initialized Executor for the given program.
func New(prog *compiler.Program, optimize bool) *Executor {
	e, err := prog.Module.Executor(optimize)
	if err != nil {
		panic(err)
	}

	if prog.CreateContext.IsNull() || prog.DestroyContext.IsNull() {
		panic("Program has no context functions. Was EmitContext not set to true?")
	}

	exec := &Executor{
		program:        prog,
		exec:           e,
		createContext:  e.FunctionAddress(prog.CreateContext),
		destroyContext: e.FunctionAddress(prog.DestroyContext),
		globalsSize:    uint64(e.SizeOf(prog.Globals.Type)),
		cmdFunctions:   map[string]unsafe.Pointer{},
	}

	for name, info := range prog.Commands {
		exec.cmdFunctions[name] = e.FunctionAddress(info.Execute)
	}

	return exec
}
