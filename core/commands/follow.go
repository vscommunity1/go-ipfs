package commands

import (
	"errors"
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	e "github.com/ipfs/go-ipfs/core/commands/e"

	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
)

type ipnsFollowResult struct {
	OK bool
}

var IpnsFollowCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Follow IPNS names.",
		ShortDescription: `
Periodically resolve and optionally pin IPNS names in the background.
`,
	},
	Subcommands: map[string]*cmds.Command{
		"add":    ipnsFollowAddCmd,
		"pin":    ipnsFollowPinCmd,
		"list":   ipnsFollowListCmd,
		"cancel": ipnsFollowCancelCmd,
	},
}

var ipnsFollowAddCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Follow a name without pinning",
		ShortDescription: `
Follows an IPNS name by periodically resolving in the backround.
`,
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Namecache == nil {
			res.SetError(errors.New("IPNS Namecache is not available"), cmdkit.ErrClient)
			return
		}

		n.Namecache.Follow(req.Arguments()[0], false)

		res.SetOutput(&ipnsFollowResult{true})
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("name", true, false, "IPNS Name to follow."),
	},
	Type: ipnsFollowResult{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: marshalFollowResult,
	},
}

var ipnsFollowPinCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Follows and pins a name",
		ShortDescription: `
Follows an IPNS name by periodically resolving and recursively
pinning in the backround.
`,
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Namecache == nil {
			res.SetError(errors.New("IPNS Namecache is not available"), cmdkit.ErrClient)
			return
		}

		n.Namecache.Follow(req.Arguments()[0], true)

		res.SetOutput(&ipnsFollowResult{true})
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("name", true, false, "IPNS Name to follow."),
	},
	Type: ipnsFollowResult{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: marshalFollowResult,
	},
}

var ipnsFollowListCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "List names followed by the daemon",
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Namecache == nil {
			res.SetError(errors.New("IPNS Namecache is not available"), cmdkit.ErrClient)
			return
		}

		res.SetOutput(&stringList{n.Namecache.ListFollows()})
	},
	Type: stringList{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: stringListMarshaler,
	},
}

var ipnsFollowCancelCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Cancels a follow",
	},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Namecache == nil {
			res.SetError(errors.New("IPNS Namecache is not available"), cmdkit.ErrClient)
			return
		}

		n.Namecache.Unfollow(req.Arguments()[0])

		res.SetOutput(&ipnsFollowResult{true})
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("name", true, false, "Name follow to cancel."),
	},
	Type: ipnsFollowResult{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: marshalFollowResult,
	},
}

func marshalFollowResult(res cmds.Response) (io.Reader, error) {
	v, err := unwrapOutput(res.Output())
	if err != nil {
		return nil, err
	}

	output, ok := v.(*ipnsFollowResult)
	if !ok {
		return nil, e.TypeErr(output, v)
	}

	var state string
	if output.OK {
		state = "ok"
	} else {
		state = "wtf"
	}

	return strings.NewReader(state + "\n"), nil
}
