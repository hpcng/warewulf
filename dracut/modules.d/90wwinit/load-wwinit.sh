#!/bin/bash

info "Mounting tmpfs at $NEWROOT"
mount -t tmpfs ${wwinit_tmpfs_size_option} tmpfs "$NEWROOT"

for archive in "${wwinit_container}" "${wwinit_kmods}" "${wwinit_system}" "${wwinit_runtime}"
do
    if [ -n "${archive}" ]
    then
        info "Loading ${archive}"
	#Load only runtime overlays from a static privledge port
        [[ "$archive" == *"runtime"* ]] && localport="--local-port 986" || localport=""
        (curl --silent $localport -L "${archive}" | gzip -d | cpio -im --directory="${NEWROOT}") || die "Unable to load ${archive}"
    fi
done
