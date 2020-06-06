package ipfs

// CurrentCommit is the current git commit, this is set as a ldflag in the Makefile
var CurrentCommit string

// CurrentVersionNumber is the current application's version literal
<<<<<<< HEAD
const CurrentVersionNumber = "0.5.0-dev"
=======
const CurrentVersionNumber = "0.4.23-rc1"
>>>>>>> release: bump to v0.4.23

const ApiVersion = "/go-ipfs/" + CurrentVersionNumber + "/"

// UserAgent is the libp2p user agent used by go-ipfs.
//
// Note: This will end in `/` when no commit is available. This is expected.
var UserAgent = "go-ipfs/" + CurrentVersionNumber + "/" + CurrentCommit
