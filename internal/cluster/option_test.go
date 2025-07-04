/*
 * MIT License
 *
 * Copyright (c) 2022-2025  Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cluster

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tochemey/goakt/v3/internal/size"
	"github.com/tochemey/goakt/v3/log"
	testkit "github.com/tochemey/goakt/v3/mocks/hash"
)

func TestOptions(t *testing.T) {
	mockHasher := new(testkit.Hasher)
	// nolint
	tlsConfig := &tls.Config{}
	size := uint64(1 * size.MB)
	testCases := []struct {
		name     string
		option   Option
		expected Engine
	}{
		{
			name:     "WithPartitionsCount",
			option:   WithPartitionsCount(2),
			expected: Engine{partitionsCount: 2},
		},
		{
			name:     "WithLogger",
			option:   WithLogger(log.DefaultLogger),
			expected: Engine{logger: log.DefaultLogger},
		},
		{
			name:     "WithWriteTimeout",
			option:   WithWriteTimeout(2 * time.Minute),
			expected: Engine{writeTimeout: 2 * time.Minute},
		},
		{
			name:     "WithReadTimeout",
			option:   WithReadTimeout(2 * time.Minute),
			expected: Engine{readTimeout: 2 * time.Minute},
		},
		{
			name:     "WithShutdownTimeout",
			option:   WithShutdownTimeout(2 * time.Minute),
			expected: Engine{shutdownTimeout: 2 * time.Minute},
		},
		{
			name:     "WithHasher",
			option:   WithHasher(mockHasher),
			expected: Engine{hasher: mockHasher},
		},
		{
			name:     "WithReplicaCount",
			option:   WithReplicaCount(2),
			expected: Engine{replicaCount: 2},
		},
		{
			name:     "WithMinimumPeersQuorum",
			option:   WithMinimumPeersQuorum(3),
			expected: Engine{minimumPeersQuorum: 3},
		},
		{
			name:     "WithWriteQuorum",
			option:   WithWriteQuorum(3),
			expected: Engine{writeQuorum: 3},
		},
		{
			name:     "WithReadQuorum",
			option:   WithReadQuorum(3),
			expected: Engine{readQuorum: 3},
		},
		{
			name:     "WithTLS",
			option:   WithTLS(tlsConfig, tlsConfig),
			expected: Engine{serverTLS: tlsConfig, clientTLS: tlsConfig},
		},
		{
			name:     "WithStorageSize",
			option:   WithTableSize(size),
			expected: Engine{tableSize: size},
		},
		{
			name:     "WithBootstrapTimeout",
			option:   WithBootstrapTimeout(3),
			expected: Engine{bootstrapTimeout: 3},
		},
		{
			name:     "WithCacheSyncInterval",
			option:   WithCacheSyncInterval(3),
			expected: Engine{cacheSyncInterval: 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cl Engine
			tc.option.Apply(&cl)
			assert.Equal(t, tc.expected, cl)
		})
	}
}
