// Copyright © 2016 Asteris, LLC
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

package rpm_test

import (
	"testing"

	"github.com/asteris-llc/converge/helpers/fakerenderer"
	"github.com/asteris-llc/converge/resource"
	"github.com/asteris-llc/converge/resource/package/rpm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPreparerInterfaces ensures that the correct interfaces are implemented by
// the preparer
func TestPreparerInterfaces(t *testing.T) {
	t.Parallel()
	assert.Implements(t, (*resource.Resource)(nil), new(rpm.Preparer))
}

// TestPrepareCreatesPackage tests rpm.Package creation from the preparerer
// ensuring that the state field is respected.
func TestPreparerCreatesPackage(t *testing.T) {
	t.Parallel()
	t.Run("when-state-present", func(t *testing.T) {
		p := &rpm.Preparer{Name: "test1", State: "present"}
		task, err := p.Prepare(fakerenderer.New())
		require.NoError(t, err)
		asRPM, ok := task.(*rpm.Package)
		require.True(t, ok)
		assert.Equal(t, "present", string(asRPM.State))
	})
	t.Run("when-state-absent", func(t *testing.T) {
		p := &rpm.Preparer{Name: "test1", State: "absent"}
		task, err := p.Prepare(fakerenderer.New())
		require.NoError(t, err)
		asRPM, ok := task.(*rpm.Package)
		require.True(t, ok)
		assert.Equal(t, "absent", string(asRPM.State))
	})
	t.Run("when-state-missing", func(t *testing.T) {
		p := &rpm.Preparer{Name: "test1"}
		task, err := p.Prepare(fakerenderer.New())
		require.NoError(t, err)
		asRPM, ok := task.(*rpm.Package)
		require.True(t, ok)
		assert.Equal(t, "present", string(asRPM.State))
	})
}
