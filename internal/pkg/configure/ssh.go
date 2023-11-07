package configure

import (
	"os"
	"path"

	warewulfconf "github.com/hpcng/warewulf/internal/pkg/config"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/pkg/errors"
)

func SSH() error {
	if os.Getuid() == 0 {
		wwlog.Info("Updating system keys")
		conf := warewulfconf.Get()
		wwkeydir := path.Join(conf.Paths.Sysconfdir, "warewulf/keys") + "/"

		err := os.MkdirAll(path.Join(conf.Paths.Sysconfdir, "warewulf/keys"), 0755)
		if err != nil {
			wwlog.Error("Could not create base directory: %s", err)
			os.Exit(1)
		}

		for _, k := range [4]string{"rsa", "dsa", "ecdsa", "ed25519"} {
			keytype := "ssh_host_" + k + "_key"
			if !util.IsFile(path.Join(wwkeydir, keytype)) {
				wwlog.Info("Setting up key: %s\n", keytype)
				wwlog.Debug("Creating new %s key", keytype)
				err = util.ExecInteractive("ssh-keygen", "-q", "-t", k, "-f", path.Join(wwkeydir, keytype), "-C", "", "-N", "")
				if err != nil {
					wwlog.Error("Failed to exec ssh-keygen: %s", err)
					return errors.Wrap(err, "failed to exec ssh-keygen command")
				}
			} else {
				wwlog.Info("Skipping, key already exists: %s", keytype)
			}
		}
	} else {
		wwlog.Info("Updating user's keys")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		wwlog.Error("Could not obtain the user's home directory: %s", err)
		os.Exit(1)
	}

	authorizedKeys := path.Join(homeDir, "/.ssh/authorized_keys")
	rsaPriv := path.Join(homeDir, "/.ssh/id_rsa")
	rsaPub := path.Join(homeDir, "/.ssh/id_rsa.pub")

	if !util.IsFile(authorizedKeys) {
		wwlog.Info("Setting up: %s", authorizedKeys)
		err = util.ExecInteractive("ssh-keygen", "-q", "-t", "rsa", "-f", rsaPriv, "-C", "", "-N", "")
		if err != nil {
			return errors.Wrap(err, "failed to exec ssh-keygen command")
		}
		err := util.CopyFile(rsaPub, authorizedKeys)
		if err != nil {
			return errors.Wrap(err, "failed to copy keys")
		}
	} else {
		wwlog.Info("Skipping, authorized_keys already exists: %s", authorizedKeys)
	}

	return nil
}
