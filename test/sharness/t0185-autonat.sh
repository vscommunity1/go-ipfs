#!/usr/bin/env bash

test_description="Test autonat"

. lib/test-lib.sh

# NOTE: This is currently a really dumb test just to make sure this service
# starts. We need better tests but testing AutoNAT without public IP addresses
# is tricky.

test_init_ipfs

test_expect_success "enable autonat" '
<<<<<<< HEAD
<<<<<<< HEAD
=======
<<<<<<< HEAD
=======
>>>>>>> master
>>>>>>> 795845ea3e69d475f7eeab37fa155ed9964486ee
  ipfs config AutoNAT.ServiceMode enabled
'

test_launch_ipfs_daemon

test_kill_ipfs_daemon

test_expect_success "enable autonat" '
  ipfs config AutoNAT.ServiceMode disabled
<<<<<<< HEAD
=======
  ipfs config --json Swarm.EnableAutoNATService true
>>>>>>> test(sharness): make sure we can actually enable autonat
=======
<<<<<<< HEAD
=======
  ipfs config --json Swarm.EnableAutoNATService true
>>>>>>> test(sharness): make sure we can actually enable autonat
=======
>>>>>>> master
>>>>>>> 795845ea3e69d475f7eeab37fa155ed9964486ee
'

test_launch_ipfs_daemon

test_kill_ipfs_daemon

test_done
