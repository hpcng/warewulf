#!/bin/bash

check() {
    # Don't include in hostonly mode
    [[ $hostonly ]] && return 1

    # Don't include by default
    return 255
}

depends() {
    echo network
    return 0
}

install() {
    inst_multiple cpio curl
    inst_hook cmdline 30 "$moddir/parse-wwinit.sh"
    inst_hook pre-mount 30 "$moddir/load-wwinit.sh"
}
