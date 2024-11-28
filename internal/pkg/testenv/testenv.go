// Package testenv provides functions and data structures for
// constructing and manipulating a temporary Warewulf environment for
// use during automated testing.
//
// The testenv package should only be used in tests.
package testenv

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/node"
)

const initWarewulfConf = `WW_INTERNAL: 0`
const initDefaultsConf = `WW_INTERNAL: 45`
const initNodesConf = `WW_INTERNAL: 45
nodeprofiles:
  default: {}
nodes:
  node1: {}
`

type TestEnv struct {
	BaseDir string
}

const Sysconfdir = "etc"
const Bindir = "bin"
const Datadir = "share"
const Localstatedir = "var/local"
const Srvdir = "srv"
const Tftpdir = "srv/tftp"
const Firewallddir = "usr/lib/firewalld/services"
const Systemddir = "usr/lib/systemd/system"
const WWOverlaydir = "var/lib/warewulf/overlays"
const WWChrootdir = "var/lib/warewulf/chroots"
const WWProvisiondir = "srv/warewulf"
const WWClientdir = "warewulf"
const Cachedir = "cache"

// New creates a test environment in a temporary directory and configures
// Warewulf to use it.
//
// Caller is responsible to delete env.BaseDir by calling
// env.RemoveAll. Note that this does not restore Warewulf to its
// previous state.
//
// Asserts no errors occur.
func New(t *testing.T) (env *TestEnv) {
	env = new(TestEnv)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "ww4test-*")
	assert.NoError(t, err)
	env.BaseDir = tmpDir

	env.WriteFile(t, path.Join(Sysconfdir, "warewulf/nodes.conf"), initNodesConf)
	env.WriteFile(t, path.Join(Sysconfdir, "warewulf/warewulf.conf"), initWarewulfConf)
	env.WriteFile(t, path.Join(Datadir, "warewulf/defaults.conf"), initDefaultsConf)

	// re-read warewulf.conf
	conf := warewulfconf.New()
	err = conf.Read(env.GetPath(path.Join(Sysconfdir, "warewulf/warewulf.conf")))
	assert.NoError(t, err)

	conf.Paths.Sysconfdir = env.GetPath(Sysconfdir)
	conf.Paths.Bindir = env.GetPath(Bindir)
	conf.Warewulf.DataStore = env.GetPath(Datadir)
	conf.Paths.Localstatedir = env.GetPath(Localstatedir)
	conf.Paths.Srvdir = env.GetPath(Srvdir)
	conf.TFTP.TftpRoot = env.GetPath(Tftpdir)
	conf.Paths.Firewallddir = env.GetPath(Firewallddir)
	conf.Paths.Systemddir = env.GetPath(Systemddir)
	conf.Paths.WWOverlaydir = env.GetPath(WWOverlaydir)
	conf.Paths.WWChrootdir = env.GetPath(WWChrootdir)
	conf.Paths.WWProvisiondir = env.GetPath(WWProvisiondir)
	conf.Paths.WWClientdir = env.GetPath(WWClientdir)
	conf.Paths.Cachedir = env.GetPath(Cachedir)

	for _, confPath := range []string{
		conf.Paths.Sysconfdir,
		conf.Paths.Bindir,
		conf.Warewulf.DataStore,
		conf.Paths.Localstatedir,
		conf.Paths.Srvdir,
		conf.TFTP.TftpRoot,
		conf.Paths.Firewallddir,
		conf.Paths.Systemddir,
		conf.Paths.WWOverlaydir,
		conf.Paths.WWChrootdir,
		conf.Paths.WWProvisiondir,
		conf.Paths.WWClientdir,
	} {
		env.MkdirAll(t, confPath)
	}

	// node.init() has already run, so set the config path again
	node.ConfigFile = env.GetPath(path.Join(Sysconfdir, "warewulf/nodes.conf"))

	return
}

// GetPath returns the absolute path name for fileName specified
// relative to the test environment.
func (env *TestEnv) GetPath(fileName string) string {
	return path.Join(env.BaseDir, fileName)
}

// MkdirAll creates dirName and any intermediate directories relative
// to the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) MkdirAll(t *testing.T, dirName string) {
	err := os.MkdirAll(env.GetPath(dirName), 0755)
	assert.NoError(t, err)
}

// WriteFile writes content to fileName, creating any necessary
// intermediate directories relative to the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) WriteFile(t *testing.T, fileName string, content string) {
	dirName := filepath.Dir(fileName)
	env.MkdirAll(t, dirName)

	f, err := os.Create(env.GetPath(fileName))
	assert.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(content)
	assert.NoError(t, err)
	err = os.Chtimes(env.GetPath(fileName),
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC),
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC))
	assert.NoError(t, err)
}

// ReadFile returns the content of fileName as converted to a
// string.
//
// Asserts no errors occur.
func (env *TestEnv) ReadFile(t *testing.T, fileName string) string {
	buffer, err := os.ReadFile(env.GetPath(fileName))
	assert.NoError(t, err)
	return string(buffer)
}

// RemoveAll deletes the temporary directory, and all its contents,
// for the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) RemoveAll(t *testing.T) {
	err := os.RemoveAll(env.BaseDir)
	assert.NoError(t, err)
}

// Writes to absolute path, but checks if given file name
// is within testenv.
//
// Asserts no errors occur.
func (env *TestEnv) WriteFileAbs(t *testing.T, fileName string, content string) {
	ok := strings.HasPrefix(fileName, env.BaseDir)
	if !ok {
		assert.Fail(t, "given filename is not in testenv")
	}
	dirName := filepath.Dir(fileName)
	err := os.MkdirAll(dirName, 0755)
	assert.NoError(t, err)
	f, err := os.Create(fileName)
	assert.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(content)
	assert.NoError(t, err)
	err = os.Chtimes(fileName,
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC),
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC))
	assert.NoError(t, err)
}
