// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"github.com/cilium/cilium/cilium/cmd"
	_ "github.com/cilium/cilium/enterprise/cilium/cmd"
)

func main() {
	cmd.Execute()
}