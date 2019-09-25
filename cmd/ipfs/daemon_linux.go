// +build linux

package main

import (
<<<<<<< HEAD
	daemon "github.com/coreos/go-systemd/v22/daemon"
=======
	daemon "github.com/coreos/go-systemd/daemon"
>>>>>>> systemd: add notify support
)

func notifyReady() {
	_, _ = daemon.SdNotify(false, daemon.SdNotifyReady)
}

func notifyStopping() {
	_, _ = daemon.SdNotify(false, daemon.SdNotifyStopping)
}
