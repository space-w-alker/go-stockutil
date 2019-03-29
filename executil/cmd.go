// Utilities that make executing commands on the local system a little bit easier.
package executil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

type CommandStatusFunc func(Status)

type Status struct {
	StartedAt  time.Time
	StoppedAt  time.Time
	Running    bool
	Successful bool
	ExitCode   int
	Error      error
	PID        int
}

func (self Status) Took() time.Duration {
	if !self.StartedAt.IsZero() {
		if !self.StoppedAt.IsZero() {
			return self.StoppedAt.Sub(self.StartedAt)
		}
	}

	return 0
}

func (self Status) String() string {
	if self.Error != nil {
		return self.Error.Error()
	} else if self.Running {
		if self.PID > 0 {
			return fmt.Sprintf("PID %d has been running for %v", self.PID, time.Since(self.StartedAt))
		} else {
			return fmt.Sprintf("Process has been running for %v, PID unknown", time.Since(self.StartedAt))
		}
	} else if !self.StoppedAt.IsZero() {
		if self.PID > 0 {
			return fmt.Sprintf("PID %d exited with status %d, took %v", self.PID, self.ExitCode, self.Took())
		} else {
			return fmt.Sprintf("Process exited with status %d, took %v, PID unknown", self.ExitCode, self.Took())
		}
	} else if self.StartedAt.IsZero() {
		return fmt.Sprintf("Process has not started yet")
	} else {
		return fmt.Sprintf("Process status is unknown")
	}
}

type Cmd struct {
	*exec.Cmd
	MonitorInterval time.Duration
	Timeout         time.Duration
	InheritEnv      bool
	OnStart         CommandStatusFunc
	OnMonitor       CommandStatusFunc
	OnComplete      CommandStatusFunc
	OnSuccess       CommandStatusFunc
	OnError         CommandStatusFunc
	status          Status
	statusLock      sync.Mutex
	reallyDone      *sync.WaitGroup
	finished        chan bool
	exitError       error
	inWriter        io.WriteCloser
}

func Wrap(cmd *exec.Cmd) *Cmd {
	return new(cmd)
}

func Command(name string, arg ...string) *Cmd {
	return new(exec.Command(name, arg...))
}

func ShellCommand(cmdline string) *Cmd {
	if shell := FindShell(); shell != `` {
		return Command(shell, `-c`, cmdline)
	}

	return nil
}

func new(wrap *exec.Cmd) *Cmd {
	return &Cmd{
		Cmd:             wrap,
		MonitorInterval: (500 * time.Millisecond),
		status:          Status{},
		finished:        make(chan bool),
		reallyDone:      &sync.WaitGroup{},
	}
}

func (self *Cmd) prestart() error {
	go self.startMonitoringCommand()

	if self.InheritEnv {
		self.Cmd.Env = append(os.Environ(), self.Cmd.Env...)
	}

	return nil
}

func (self *Cmd) CombinedOutput() ([]byte, error) {
	defer self.waitReallyDone()

	if err := self.prestart(); err == nil {
		out, err := self.Cmd.CombinedOutput()
		self.exitError = err
		self.updateStatus()
		return out, err
	} else {
		return nil, err
	}
}

func (self *Cmd) Output() ([]byte, error) {
	defer self.waitReallyDone()

	if err := self.prestart(); err == nil {
		out, err := self.Cmd.Output()
		self.exitError = err
		self.updateStatus()
		return out, err
	} else {
		return nil, err
	}
}

func (self *Cmd) SetEnv(key string, value interface{}) {
	kv := fmt.Sprintf("%v=%s", key, typeutil.String(value))

	for i, pair := range self.Cmd.Env {
		k, _ := stringutil.SplitPair(pair, `=`)

		if k == key {
			self.Cmd.Env[i] = kv
			return
		}
	}

	self.Cmd.Env = append(self.Cmd.Env, kv)
}

func (self *Cmd) Run() error {
	defer self.waitReallyDone()

	if err := self.prestart(); err == nil {
		err := self.Cmd.Run()
		self.exitError = err
		self.updateStatus()
		return err
	} else {
		return err
	}
}

func (self *Cmd) Start() error {
	if err := self.prestart(); err == nil {
		if err := self.Cmd.Start(); err != nil {
			self.exitError = err
			self.updateStatus()
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}

// Return the current status of the process.
func (self *Cmd) Status() Status {
	self.statusLock.Lock()
	defer self.statusLock.Unlock()

	return self.status
}

// Kill the running command.
func (self *Cmd) Kill() error {
	if p := self.Process; p != nil {
		return p.Kill()
	} else {
		return nil
	}
}

// Wait for the process to complete, then return the last status.
// Process must have been started using the Start() function.
func (self *Cmd) WaitStatus() Status {
	self.exitError = self.Wait()
	self.updateStatus()
	return self.status
}

func (self *Cmd) waitReallyDone() {
	// this quits the startMonitoringCommand() monitoring loop
	select {
	case self.finished <- true:
	default:
	}

	// this waits for startMonitoringCommand() to fire callbacks and actually exit
	self.reallyDone.Wait()
}

// a goroutine that is launched whenever a command
func (self *Cmd) startMonitoringCommand() {
	self.reallyDone.Add(1)
	self.statusLock.Lock()
	self.status.StartedAt = time.Now()
	self.status.Running = true
	self.statusLock.Unlock()

	self.updateStatus()

	if fn := self.OnStart; fn != nil {
		fn(self.status)
	}

	// if there's a timeout set, and its shorter than our monitor interval,
	// reduce the monitor interval to that
	if self.Timeout > 0 && self.Timeout < self.MonitorInterval {
		self.MonitorInterval = self.Timeout
	}

	ticker := time.NewTicker(self.MonitorInterval)

MonitorLoop:
	for self.Timeout == 0 || time.Since(self.status.StartedAt) < self.Timeout {
		select {
		case <-ticker.C:
			self.updateStatus()

			if !self.status.Running {
				break MonitorLoop
			}

			if fn := self.OnMonitor; fn != nil {
				fn(self.Status())
			}

		case <-self.finished:
			ticker.Stop()
			break MonitorLoop
		}
	}

	ticker.Stop()
	self.Kill()
	self.updateStatus()
	finalStatus := self.Status()

	// fire off callbacks
	if fn := self.OnComplete; fn != nil {
		fn(finalStatus)
	}

	if finalStatus.Successful {
		if fn := self.OnSuccess; fn != nil {
			fn(finalStatus)
		}
	} else if fn := self.OnError; fn != nil {
		fn(finalStatus)
	}

	self.reallyDone.Done()
}

func (self *Cmd) updateStatus() {
	self.statusLock.Lock()
	defer self.statusLock.Unlock()

	var pstate *os.ProcessState

	// handle known exit errors
	if self.exitError != nil {
		self.status.Error = self.exitError
		self.status.Running = false
		self.status.StoppedAt = time.Now()

		if xerr, ok := self.exitError.(*exec.ExitError); ok {
			pstate = xerr.ProcessState
		} else {
			return
		}
	} else if s := self.Cmd.ProcessState; s != nil {
		pstate = s
	}

	// check ProcessState
	if pstate != nil {
		self.status.PID = pstate.Pid()
		self.status.Successful = pstate.Success()

		if status, ok := pstate.Sys().(syscall.WaitStatus); ok {
			if status.Exited() {
				self.status.Running = false
				self.status.ExitCode = status.ExitStatus()
				self.status.StoppedAt = time.Now()

				if status.Signaled() {
					if sig := status.StopSignal(); sig != 0 {
						self.status.Error = errors.New(sig.String())
					} else if sig := status.Signal(); sig != 0 {
						self.status.Error = errors.New(sig.String())
					}
				} else if e := self.status.ExitCode; e != 0 {
					self.status.Error = fmt.Errorf("Process exited with status %d", e)
				}
			} else if status.Stopped() {
				self.status.Running = false
			} else {
				self.status.Running = true
			}
		}
	} else if process := self.Cmd.Process; process != nil {
		self.status.PID = process.Pid
	}
}

func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	return &Cmd{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
}
