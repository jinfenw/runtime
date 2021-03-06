//
// Copyright (c) 2017 Intel Corporation
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
//

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// package variables set in TestMain
var testDir = ""
var testDirMode = os.FileMode(0750)

// TestMain is the common main function used by ALL the test functions
// for this package.
func TestMain(m *testing.M) {
	var err error

	testDir, err = ioutil.TempDir("", fmt.Sprintf("%s-", name))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(testDir, testDirMode)
	if err != nil {
		fmt.Printf("Could not create test directory %s: %s\n", testDir, err)
		os.Exit(1)
	}

	ret := m.Run()

	os.RemoveAll(testDir)

	os.Exit(ret)
}
