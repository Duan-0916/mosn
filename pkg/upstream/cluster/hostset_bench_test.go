/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cluster

import (
	"testing"

	"sofastack.io/sofa-mosn/pkg/api/v2"
	"sofastack.io/sofa-mosn/pkg/types"
)

func benchAddHost(b *testing.B, count int) {
	pool := makePool(2 * count)
	oldHosts := pool.MakeHosts(count, nil)
	newHosts := pool.MakeHosts(count, nil)
	newHosts = append(newHosts, oldHosts...)
	for i := 0; i < b.N; i++ {
		hs := &hostSet{
			allHosts: oldHosts,
		}
		hs.UpdateHosts(newHosts)
	}
}

// add and delete
func benchUpdateHost(b *testing.B, count int) {
	pool := makePool(3 * count)
	oldHosts := pool.MakeHosts(2*count, nil)
	newHosts := pool.MakeHosts(count, nil)
	newHosts = append(newHosts, oldHosts[:count]...)
	for i := 0; i < b.N; i++ {
		hs := &hostSet{
			allHosts: oldHosts,
		}
		hs.UpdateHosts(newHosts)

	}
}

// Test HostSet Update Host
func BenchmarkHostSetUpdateHost(b *testing.B) {

	b.Run("AddHost10", func(b *testing.B) {
		benchAddHost(b, 10)
	})

	b.Run("AddHost100", func(b *testing.B) {
		benchAddHost(b, 100)
	})

	b.Run("AddHost500", func(b *testing.B) {
		benchAddHost(b, 500)
	})

	b.Run("Update50", func(b *testing.B) {
		benchUpdateHost(b, 50)
	})

	b.Run("Update500", func(b *testing.B) {
		benchUpdateHost(b, 500)
	})

}
func BenchmarkRemoveHosts(b *testing.B) {
	pool := makePool(100)
	totalHosts := pool.MakeHosts(100, nil)
	removedAddrs := []string{}
	for _, h := range totalHosts[:50] {
		removedAddrs = append(removedAddrs, h.AddressString())
	}
	b.Run("RemoveHosts", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hs := &hostSet{
				allHosts: totalHosts,
			}
			hs.RemoveHosts(removedAddrs)
		}
	})
}

func BenchmarkRefreshHost(b *testing.B) {
	pool := makePool(100)
	totalHosts := pool.MakeHosts(50, nil)
	totalHosts = append(totalHosts, pool.MakeHosts(50, v2.Metadata{
		"zone": "a",
	})...)
	hs := &hostSet{}
	hs.UpdateHosts(totalHosts)
	hs.createSubset(func(h types.Host) bool {
		if h.Metadata() != nil && h.Metadata()["zone"] == "a" {
			return true
		}
		return false
	})
	host := hs.Hosts()[55]
	b.Run("RefreshHost", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				host.SetHealthFlag(types.FAILED_ACTIVE_HC)
			} else {
				host.ClearHealthFlag(types.FAILED_ACTIVE_HC)
			}
			hs.refreshHealthHosts(host)
		}
	})
}
