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
	OnStart         CommandStatusFunc
	OnMonitor       CommandStatusFunc
	OnComplete      CommandStatusFunc
	OnSuccess       CommandStatusFunc
	OnError         CommandStatusFunc
	status          Status
	statusLock      sync.Mutex
	monitor         *time.Ticker
	done            chan Status
	reallyDone      chan bool
	start           chan bool
	exitError       error
	inWriter        io.WriteCloser
}

func Wrap(cmd *exec.Cmd) *Cmd {
	return new(cmd)
}

func Command(name string, arg ...string) *Cmd {
	return new(exec.Command(name, arg...))
}

func new(wrap *exec.Cmd) *Cmd {
	cmd := &Cmd{
		Cmd:             wrap,
		MonitorInterval: (500 * time.Millisecond),
		status:          Status{},
		done:            make(chan Status),
		start:           make(chan bool),
		reallyDone:      make(chan bool),
	}

	go cmd.startMonitoringCommand()
	return cmd
}

func (self *Cmd) prestart() error {
	self.start <- true
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

// Wait for the process to complete, then return the last status.
// Process must have been started using the Start() function.
func (self *Cmd) WaitStatus() Status {
	self.exitError = self.Wait()
	self.updateStatus()
	return self.status
}

func (self *Cmd) waitReallyDone() {
	<-self.reallyDone
}

func (self *Cmd) startMonitoringCommand() {
	<-self.start

	self.statusLock.Lock()
	self.status.StartedAt = time.Now()
	self.statusLock.Unlock()

	self.updateStatus()

	if fn := self.OnStart; fn != nil {
		fn(self.status)
	}

	self.monitor = time.NewTicker(self.MonitorInterval)

MonitorLoop:
	for {
		select {
		case <-self.monitor.C:
			self.updateStatus()

			if fn := self.OnMonitor; fn != nil {
				fn(self.Status())
			}
		case <-self.done:
			self.monitor.Stop()
			break MonitorLoop
		}
	}

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

	self.reallyDone <- true
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
			self.done <- self.status
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

				self.done <- self.status
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
