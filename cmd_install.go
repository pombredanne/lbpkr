package main

import (
	"fmt"
	"regexp"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

func pkr_make_cmd_install() *commander.Command {
	cmd := &commander.Command{
		Run:       pkr_run_cmd_install,
		UsageLine: "install [options] <rpmname> [<version> [<release>]]",
		Short:     "install a RPM from the yum repository",
		Long: `
install installs a RPM from the yum repository.

ex:
 $ pkr install LHCb
`,
		Flag: *flag.NewFlagSet("pkr-install", flag.ExitOnError),
	}
	cmd.Flag.Bool("v", false, "enable verbose mode")
	cmd.Flag.String("type", "lhcb", "config type (lhcb|atlas)")
	return cmd
}

func pkr_run_cmd_install(cmd *commander.Command, args []string) error {
	var err error

	cfgtype := cmd.Flag.Lookup("type").Value.Get().(string)
	debug := cmd.Flag.Lookup("v").Value.Get().(bool)

	cfg := NewConfig(cfgtype)
	ctx, err := New(cfg, debug)
	if err != nil {
		return err
	}

	ctx.msg.Infof("hello: %v\n", cfg.Prefix())

	rpmname := ""
	version := ""
	release := ""
	switch len(args) {
	case 0:
		ctx.msg.Errorf("please specify at least the name of the RPM to install\n")
		cmd.Usage()
		return fmt.Errorf("pkr: invalid number of arguments (got=%d)", len(args))
	case 1:
		rpmname = args[0]
	case 2:
		rpmname = args[0]
		version = args[1]
	case 3:
		rpmname = args[0]
		version = args[1]
		release = args[2]
	default:
		return fmt.Errorf("pkr: invalid number of arguments. expected n=1|2|3. got=%d (%v)",
			len(args),
			args,
		)
	}

	re := regexp.MustCompile(`(.*)-([\d\.]+)-(\d)$`).FindAllStringSubmatch(rpmname, -1)
	if len(re) == 1 {
		m := re[0]
		switch len(m) {
		case 2:
			rpmname = m[1]
		case 3:
			rpmname = m[1]
			version = m[2]
		case 4:
			rpmname = m[1]
			version = m[2]
			release = m[3]
		}
	}
	ctx.msg.Infof("installing RPM %s %s %s\n", rpmname, version, release)

	forceInstall := false
	update := false
	err = ctx.InstallRPM(rpmname, version, release, forceInstall, update)
	return err
}