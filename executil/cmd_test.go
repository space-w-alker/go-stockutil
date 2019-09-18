package executil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func incr(value *int) CommandStatusFunc {
	return func(s Status) {
		if value != nil {
			*value += 1
		}
	}
}

func TestExecTrue(t *testing.T) {
	var starts int
	var completes int
	var successes int
	var failures int

	assert := require.New(t)

	cmd := Command(`/bin/true`)
	cmd.OnStart = incr(&starts)
	cmd.OnComplete = incr(&completes)
	cmd.OnSuccess = incr(&successes)
	cmd.OnError = incr(&failures)

	assert.False(cmd.Status().Running)
	assert.False(cmd.Status().Successful)
	assert.Zero(cmd.Status().PID)
	assert.True(cmd.Status().StartedAt.IsZero())
	assert.True(cmd.Status().StoppedAt.IsZero())
	assert.Zero(cmd.Status().ExitCode)
	assert.Nil(cmd.Status().Error)

	cmd.Run()

	assert.False(cmd.Status().Running)
	assert.True(cmd.Status().Successful)
	assert.False(cmd.Status().PID == 0)
	assert.False(cmd.Status().StartedAt.IsZero())
	assert.False(cmd.Status().StoppedAt.IsZero())
	assert.True(cmd.Status().Took() > 0)
	assert.Zero(cmd.Status().ExitCode)
	assert.Nil(cmd.Status().Error)

	assert.Equal(1, starts)
	assert.Equal(1, completes)
	assert.Equal(1, successes)
	assert.Equal(0, failures)
}

func TestExecFalse(t *testing.T) {
	var starts int
	var completes int
	var successes int
	var failures int

	assert := require.New(t)

	cmd := Command(`/bin/false`)
	cmd.OnStart = incr(&starts)
	cmd.OnComplete = incr(&completes)
	cmd.OnSuccess = incr(&successes)
	cmd.OnError = incr(&failures)

	assert.False(cmd.Status().Running)
	assert.False(cmd.Status().Successful)
	assert.Zero(cmd.Status().PID)
	assert.True(cmd.Status().StartedAt.IsZero())
	assert.True(cmd.Status().StoppedAt.IsZero())
	assert.Zero(cmd.Status().ExitCode)
	assert.Nil(cmd.Status().Error)

	cmd.Run()

	assert.False(cmd.Status().Running)
	assert.False(cmd.Status().Successful)
	assert.False(cmd.Status().PID == 0)
	assert.False(cmd.Status().StartedAt.IsZero())
	assert.False(cmd.Status().StoppedAt.IsZero())
	assert.True(cmd.Status().Took() > 0)
	assert.True(cmd.Status().ExitCode == 1)
	assert.EqualError(cmd.Status().Error, `Process exited with status 1`)

	assert.Equal(1, starts)
	assert.Equal(1, completes)
	assert.Equal(0, successes)
	assert.Equal(1, failures)
}
