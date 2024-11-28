-include Defaults.mk

# Linux distro (try and set to /etc/os-release ID)
OS_REL := $(shell sed -n "s/^ID\s*=\s*['"\""]\(.*\)['"\""]/\1/p" /etc/os-release)
OS ?= $(OS_REL)

ARCH_REL := $(shell uname -p)
ARCH ?= $(ARCH_REL)

# List of variables to save and replace in files
VARLIST := OS

# Project Information
VARLIST += WAREWULF VERSION RELEASE
WAREWULF ?= warewulf
VERSION ?= $(shell scripts/rpm-version.sh $$(scripts/get-version.sh))
RELEASE ?= $(shell scripts/rpm-release.sh $$(scripts/get-version.sh))

# Use LSB-compliant paths if OS is known
ifneq ($(OS),)
  USE_LSB_PATHS := true
endif

# Always default to GNU autotools default paths if PREFIX has been redefined
ifdef PREFIX
  USE_LSB_PATHS := false
endif

# System directory paths
VARLIST += PREFIX BINDIR SYSCONFDIR SRVDIR DATADIR MANDIR DOCDIR LOCALSTATEDIR RELEASE
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
SYSCONFDIR ?= $(PREFIX)/etc
DATADIR ?= $(PREFIX)/share
MANDIR ?= $(DATADIR)/man
DOCDIR ?= $(DATADIR)/doc
CACHEDIR ?= /var/cache

ifeq ($(USE_LSB_PATHS),true)
  SRVDIR ?= /srv
  LOCALSTATEDIR ?= /var/local
else
  SRVDIR ?= $(PREFIX)/srv
  LOCALSTATEDIR ?= $(PREFIX)/var
endif

# OS-Specific Service Locations
VARLIST += TFTPDIR FIREWALLDDIR SYSTEMDDIR BASHCOMPDIR LOGROTATEDIR DRACUTMODDIR
SYSTEMDDIR ?= /usr/lib/systemd/system
BASHCOMPDIR ?= /etc/bash_completion.d
FIREWALLDDIR ?= /usr/lib/firewalld/services
LOGROTATEDIR ?= /etc/logrotate.d
DRACUTMODDIR ?= /usr/lib/dracut/modules.d
ifeq ($(OS),suse)
  TFTPDIR ?= /srv/tftpboot
endif
ifeq ($(OS),ubuntu)
  TFTPDIR ?= /srv/tftp
endif
# Default to Red Hat / Rocky Linux
TFTPDIR ?= /var/lib/tftpboot

# Warewulf directory paths
VARLIST += WWCLIENTDIR WWCONFIGDIR WWPROVISIONDIR WWOVERLAYDIR WWCHROOTDIR WWTFTPDIR WWDOCDIR WWDATADIR IPXESOURCE
WWCONFIGDIR := $(SYSCONFDIR)/$(WAREWULF)
WWPROVISIONDIR := $(LOCALSTATEDIR)/$(WAREWULF)/provision
WWOVERLAYDIR := $(LOCALSTATEDIR)/$(WAREWULF)/overlays
WWCHROOTDIR := $(LOCALSTATEDIR)/$(WAREWULF)/chroots
WWTFTPDIR := $(TFTPDIR)/$(WAREWULF)
WWDOCDIR := $(DOCDIR)/$(WAREWULF)
WWDATADIR := $(DATADIR)/$(WAREWULF)
WWCLIENTDIR ?= /warewulf

CONFIG := $(shell pwd)

IPXESOURCE ?= $(PREFIX)/share/ipxe

# helper functions
godeps=$(shell 2>/dev/null go list -mod vendor -deps -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}' $(1) | sed "s%${PWD}/%%g")

# use GOPROXY for older git clients and speed up downloads
GOPROXY ?= https://proxy.golang.org
export GOPROXY

# built tags needed for wwbuild binary
WW_GO_BUILD_TAGS := containers_image_openpgp containers_image_ostree

.PHONY: defaults
defaults:
	printf " $(foreach V,$(VARLIST),$V := $(strip $($V))\n)" >Defaults.mk
