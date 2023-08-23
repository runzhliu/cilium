// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package egressmapha

import (
	"net"
	"testing"

	"github.com/cilium/ebpf/rlimit"
	"github.com/stretchr/testify/assert"

	"github.com/cilium/cilium/pkg/bpf"
	"github.com/cilium/cilium/pkg/ebpf"
	"github.com/cilium/cilium/pkg/hive/hivetest"
	"github.com/cilium/cilium/pkg/testutils"
)

func TestPolicyMap(t *testing.T) {
	testutils.PrivilegedTest(t)

	bpf.CheckOrMountFS("")
	assert.NoError(t, rlimit.RemoveMemlock())

	egressPolicyMap := createPolicyMap(hivetest.Lifecycle(t), DefaultPolicyConfig, ebpf.PinNone)

	sourceIP1 := net.ParseIP("1.1.1.1")
	sourceIP2 := net.ParseIP("1.1.1.2")

	_, destCIDR1, err := net.ParseCIDR("2.2.1.0/24")
	assert.NoError(t, err)
	_, destCIDR2, err := net.ParseCIDR("2.2.2.0/24")
	assert.NoError(t, err)

	egressIP1 := net.ParseIP("3.3.3.1")
	egressIP2 := net.ParseIP("3.3.3.2")

	gatewayIP1 := net.ParseIP("4.4.4.1")
	gatewayIP2 := net.ParseIP("4.4.4.2")

	// This will create 2 policies, respectively with 2 and 1 egress GWs:
	//
	// Source IP   Destination CIDR   Egress IP   Gateway
	// 1.1.1.1     2.2.1.0/24         3.3.3.1     0 => 4.4.4.1
	//                                            1 => 4.4.4.2
	// 1.1.1.2     2.2.2.0/24         3.3.3.2     0 => 4.4.4.1

	err = ApplyEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1, egressIP1, []net.IP{gatewayIP1, gatewayIP2})
	assert.NoError(t, err)

	defer RemoveEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1)

	err = ApplyEgressPolicy(egressPolicyMap, sourceIP2, *destCIDR2, egressIP2, []net.IP{gatewayIP1})
	assert.NoError(t, err)

	defer RemoveEgressPolicy(egressPolicyMap, sourceIP2, *destCIDR2)

	val, err := egressPolicyMap.Lookup(sourceIP1, *destCIDR1)
	assert.NoError(t, err)

	assert.EqualValues(t, val.Size, uint32(2))
	assert.True(t, val.EgressIP.IP().Equal(egressIP1))
	assert.True(t, val.GatewayIPs[0].IP().Equal(gatewayIP1))
	assert.True(t, val.GatewayIPs[1].IP().Equal(gatewayIP2))

	val, err = egressPolicyMap.Lookup(sourceIP2, *destCIDR2)
	assert.NoError(t, err)

	assert.EqualValues(t, val.Size, uint32(1))
	assert.True(t, val.EgressIP.IP().Equal(egressIP2))
	assert.True(t, val.GatewayIPs[0].IP().Equal(gatewayIP1))

	// Adding a policy with too many gateways should result in an error
	gatewayIPs := make([]net.IP, maxGatewayNodes+1)
	err = ApplyEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1, egressIP1, gatewayIPs)
	assert.ErrorContains(t, err, "cannot apply egress policy: too many gateways")

	// Update the first policy:
	//
	// - remove gatewayIP1 from the list of active gateways (by applying a
	//   new policy with just gatewayIP2)
	// - remove gatewayIP1 also from the list of healthy gateways
	err = ApplyEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1, egressIP1, []net.IP{gatewayIP2})
	assert.NoError(t, err)

	// Update the first policy:
	//
	// - change the active gateway from gatewayIP2 -> gatewayIP1
	//-  keep gatewayIP2 in the list of healthy gateways
	err = ApplyEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1, egressIP1, []net.IP{gatewayIP1})
	assert.NoError(t, err)

	// Update the first policy:
	//
	//-  Remove gatewayIP2 from the list of healthy gateways
	err = ApplyEgressPolicy(egressPolicyMap, sourceIP1, *destCIDR1, egressIP1, []net.IP{gatewayIP1})
	assert.NoError(t, err)

	// Remove the second policy
	err = RemoveEgressPolicy(egressPolicyMap, sourceIP2, *destCIDR2)
	assert.NoError(t, err)
}
