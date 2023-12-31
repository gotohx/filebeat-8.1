// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		version string
		err     bool
		result  Version
	}{
		{
			version: "1.2.3",
			err:     false,
			result:  Version{Major: 1, Minor: 2, Bugfix: 3, version: "1.2.3"},
		},
		{
			version: "1.3.3",
			err:     false,
			result:  Version{Major: 1, Minor: 3, Bugfix: 3, version: "1.3.3"},
		},
		{
			version: "1.3.2-alpha1",
			err:     false,
			result:  Version{Major: 1, Minor: 3, Bugfix: 2, version: "1.3.2-alpha1", Meta: "alpha1"},
		},
		{
			version: "alpha1",
			err:     true,
		},
	}

	for _, test := range tests {
		v, err := NewVersion(test.version)
		if test.err {
			assert.Error(t, err)
			continue
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, *v, test.result)
	}
}

func TestLessThanOrEqual(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		version1 string
		meta     bool
		result   bool
	}{
		{
			name:     "5.4.3 < 6.4.3",
			version:  "5.4.3",
			version1: "6.4.3",
			meta:     true,
			result:   true,
		},
		{
			name:     "5.7.3 < 5.8.3",
			version:  "5.7.3",
			version1: "5.8.2",
			meta:     true,
			result:   true,
		},
		{
			name:     "5.4.3 > 5.4.2",
			version:  "5.4.3",
			version1: "5.4.2",
			meta:     true,
			result:   false,
		},
		{
			name:     "5.4.3 = 5.4.3",
			version:  "5.4.3",
			version1: "5.4.3",
			meta:     true,
			result:   true,
		},
		{
			name:     "1.2.3 > 1.2.3-alpha1",
			version:  "1.2.3",
			version1: "1.2.3-alpha1",
			meta:     true,
			result:   false,
		},
		{
			name:     "1.2.3.alpha < 1.2.3",
			version:  "1.2.3-alpha",
			version1: "1.2.3",
			meta:     true,
			result:   true,
		},
		{
			name:     "1.2.3.alpha < 1.2.3.beta",
			version:  "1.2.3-alpha",
			version1: "1.2.3-beta",
			meta:     true,
			result:   true,
		},
		{
			name:     "1.2.3.rc1 < 1.2.3.rc2",
			version:  "1.2.3-rc1",
			version1: "1.2.3-rc2",
			meta:     true,
			result:   true,
		},
		{
			name:     "5.4.3 = 5.4.3",
			version:  "5.4.3",
			version1: "5.4.3",
			meta:     false,
			result:   true,
		},
		{
			name:     "1.2.3 = 1.2.3-alpha1",
			version:  "1.2.3",
			version1: "1.2.3-alpha1",
			meta:     false,
			result:   true,
		},
		{
			name:     "1.2.3.beta1 = 1.2.3-beta2",
			version:  "1.2.3-beta1",
			version1: "1.2.3-beta2",
			meta:     false,
			result:   true,
		},
	}

	for _, test := range tests {
		v, err := NewVersion(test.version)
		assert.NoError(t, err)
		v1, err := NewVersion(test.version1)
		assert.NoError(t, err)

		assert.Equal(t, test.result, v.LessThanOrEqual(test.meta, v1), test.name)
	}
}

func TestVersionLessThan(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		version1 string
		result   bool
	}{
		{
			name:     "1.2.3 < 2.0.0",
			version:  "1.2.3",
			version1: "2.0.0",
			result:   true,
		},
		{
			name:     "1.2.3 = 1.2.3-beta1",
			version:  "1.2.3",
			version1: "1.2.3-beta1",
			result:   false,
		},
		{
			name:     "5.4.1 < 5.4.2",
			version:  "5.4.1",
			version1: "5.4.2",
			result:   true,
		},
		{
			name:     "5.5.1 > 5.4.2",
			version:  "5.5.1",
			version1: "5.4.2",
			result:   false,
		},
		{
			name:     "6.1.1-alpha3 < 6.2.0",
			version:  "6.1.1-alpha3",
			version1: "6.2.0",
			result:   true,
		},
	}

	for _, test := range tests {
		v, err := NewVersion(test.version)
		assert.NoError(t, err)
		v1, err := NewVersion(test.version1)
		assert.NoError(t, err)

		assert.Equal(t, v.LessThan(v1), test.result, test.name)
	}
}

func TestVersionLessThanMajorMinor(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		version1 string
		result   bool
	}{
		{
			name:     "1.2.x < 2.0.x",
			version:  "1.2.3",
			version1: "2.0.0",
			result:   true,
		},
		{
			name:     "1.2.3 = 1.2.3-beta1",
			version:  "1.2.3",
			version1: "1.2.3-beta1",
			result:   false,
		},
		{
			name:     "5.4.x == 5.4.x",
			version:  "5.4.1",
			version1: "5.4.2",
			result:   false,
		},
		{
			name:     "5.5.x > 5.4.x",
			version:  "5.5.1",
			version1: "5.4.2",
			result:   false,
		},
		{
			name:     "6.1.x < 6.2.x",
			version:  "6.1.1-alpha3",
			version1: "6.2.0",
			result:   true,
		},
	}

	for _, test := range tests {
		v, err := NewVersion(test.version)
		assert.NoError(t, err)
		v1, err := NewVersion(test.version1)
		assert.NoError(t, err)

		assert.Equal(t, v.LessThanMajorMinor(v1), test.result, test.name)
	}
}
