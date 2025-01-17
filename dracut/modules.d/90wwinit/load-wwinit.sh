#!/bin/bash

function setup_disks() {
	# format and prepare the disk(s)
        /usr/bin/ignition --root=/sysroot --platform=metal --stage=disks || die "warewulf: failed to partition/format disk"
	#mount /dev/disk/by-partlabel/root "$NEWROOT" || die "warewulf: failed to mount new root"
        /usr/bin/ignition --root=/sysroot --platform=metal --stage=mount || die "warewulf: failed to mount disk"

}


info "warewulf: Running warewulf v4 dracut init ${wwinit_persistent}"
archives="${wwinit_container} ${wwinit_kmods} ${wwinit_system} ${wwinit_runtime}"
if [ -z ${wwinit_persistent} ] ; then 
	info "warewulf: Mounting tmpfs at $NEWROOT"
	mount -t tmpfs -o mpol=interleave ${wwinit_tmpfs_size_option} tmpfs "$NEWROOT"
else
	info "warewulf: Using persistent setup"
	# get the igintion config and store in /run/igintion.ign, don't use igntion as so we can
	# construct the download link
	curl --silent ${localport}  -L "http://${wwinit_ip}:${wwinit_port}/overlay/persistent/ignition.json?node=${wwinit_node}" | jq '.storage.filesystems |= map(select(.device=="/dev/disk/by-partlabel/rootfs").path="/")' > /run/ignition.json || die "warewulf: failed to fetch ignition configuration"
	setup_disks
	if [ -e ${NEWROOT}/warewulf/ww_${wwinit_containername} ] ; then 
		info "warewulf: found  ${NEWROOT}/warewulf/ww_${wwinit_container} running rsync for update"
		# we need only the overlays as we update the root 
		# with rsync when the same container is used
		archives="${wwinit_system} ${wwinit_runtime}"
		rsync -aux --delete --exclude=/proc/ --exclude=/sys/ --exclude=/dev rsync://${wwinit_ip}/${wwinit_containername} $NEWROOT || die "warewulf: failed to run rsync"
	else
		if [  -e ${NEWROOT}/warewulf/ww_* ] ; then
			info "warewulf: found $(ls ${NEWROOT}/warewulf/ww_*) but container is ${wwinit_containername}, wiping $NEWROOT"
			# we didn't have the wanted container but provisioned with warewulf, wipe the rootfs
			/usr/bin/ignition --root=/sysroot --platform=metal --stage=umount || die "warewulf: failed to mount disk"
			jq '.storage.filesystems |= map(select(.device=="/dev/disk/by-partlabel/rootfs").wipeFilesystem=true)' /run/ignition.json > /run/ignition.json.mod || die "couldn't reformation ignition configuration"
			mv /run/ignition.json.mod /run/ignition.json || die "warewulf: couldn't mv ignition configuration back"
			setup_disks
		fi
	fi
fi
for archive in  $archives ; do
    if [ -n "${archive}" ] ; then
	info "warewulf: Loading ${archive}"
	# Load runtime overlay from a static privledged port.
	# Others use default settings.
	localport=""
	if [[ "${archive}" == "${wwinit_runtime}" ]] ; then
	    localport="--local-port 1-1023"
	fi
	(curl --silent ${localport} -L "${archive}" | gzip -d | cpio --quiet -u -im --directory="${NEWROOT}") || die "warewulf: Unable to load ${archive}"
    fi
done
if [ ! -z ${wwinit_persistent} ] ; then
	# this avoids that ignition runs a second time
	info "warewulf: removing ignition config from wwinit image"
	rm -rf ${NEWROOT}/warewulf/ignition.json
	echo "Container name of persistent install: ${wwinit_containername}" > ${NEWROOT}/warewulf/ww_${wwinit_containername}
	echo "# created from http://${wwinit_ip}:${wwinit_port}/overlay/persistent/fstab.ww?node=${wwinit_node}" > ${NEWROOT}/etc/fstab
	curl --silent ${localport}  -L "http://${wwinit_ip}:${wwinit_port}/overlay/persistent/fstab.ww?node=${wwinit_node}" >> ${NEWROOT}/etc/fstab || die "warewulf: failed to write correct fstab"
fi
info "warewulf: Finished warewulf v4 dracut"

