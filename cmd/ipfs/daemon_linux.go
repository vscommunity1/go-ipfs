// +build linux

package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
	daemon "github.com/coreos/go-systemd/v22/daemon"
=======
	daemon "github.com/coreos/go-systemd/daemon"
>>>>>>> systemd: add notify support
=======
	daemon "github.com/coreos/go-systemd/v22/daemon"
>>>>>>> master
)

func notifyReady() {
	_, _ = daemon.SdNotify(false, daemon.SdNotifyReady)
}

func notifyStopping() {
	_, _ = daemon.SdNotify(false, daemon.SdNotifyStopping)
}
