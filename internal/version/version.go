// Copyright 2026 Christophe Pallier <christophe@pallier.org>
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
//
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"runtime"
)

// These variables are populated at build time using -ldflags
var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// These are typically static constants
const (
	Author  = "Christophe Pallier <christophe@pallier.org>"
	License = "Apache-2.0"
)

// Info returns a formatted string containing all version metadata
func Info() string {
	return fmt.Sprintf(
		`Version:    %s
Git Commit: %s
Build Time: %s
Go Version: %s
OS/Arch:    %s/%s
Author:     %s
`,
		Version, GitCommit, BuildTime, runtime.Version(), runtime.GOOS, runtime.GOARCH, Author,
	)
}
