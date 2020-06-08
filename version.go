package ipfs

// CurrentCommit is the current git commit, this is set as a ldflag in the Makefile
var CurrentCommit string

// CurrentVersionNumber is the current application's version literal
<<<<<<< HEAD
const CurrentVersionNumber = "0.5.1"
=======
const CurrentVersionNumber = "0.6.0-dev"
>>>>>>> 795845ea3e69d475f7eeab37fa155ed9964486ee

const ApiVersion = "/go-ipfs/" + CurrentVersionNumber + "/"

// UserAgent is the libp2p user agent used by go-ipfs.
//
// Note: This will end in `/` when no commit is available. This is expected.
var UserAgent = "go-ipfs/" + CurrentVersionNumber + "/" + CurrentCommit
