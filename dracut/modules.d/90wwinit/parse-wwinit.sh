#!/bin/sh
# root=wwinit

[ -z "$root" ] && root=$(getarg root=)

if [ "${root}" = "tmpfs" ] || [ "${root}" = "initramfs" ] || [ "${root}" = "persistent" ] ; then
    info "warewulf: root=${root}"
    uuid=$(dmidecode -s system-uuid)
    assetkey=$(dmidecode -s chassis-asset-tag | sed -E -e 's/(^ +| +$)//g' -e 's/^(Unknown|Not Specified)$//g' -e 's/ /_/g')
    export wwinit_node="$(getarg wwinit.node)"; info "warewulf: wwinit_node=${wwinit_node}"
    export wwinit_ip="$(getarg wwinit.ip)"; info "warewulf: wwinit_ip=${wwinit_ip}"
    export wwinit_port="$(getarg wwinit.port)"; info "warewulf: wwinit_port=${wwinit_port}"
    export wwinit_mac="$(getarg wwid)"; info "warewulf: wwinit_mac=${wwinit_mac}"
    export wwinit_containername="$(getarg wwinit.containername)"; info "warewulf: wwinit_mac=${wwinit_containername}"
    wwinit_uri="http://${wwinit_ip}:${wwinit_port}/provision/${wwinit_mac}?assetkey=${assetkey}&uuid=${uuid}"
    export wwinit_container="${wwinit_uri}&stage=container&compress=gz"; info "warewulf: wwinit_container=${wwinit_container}"
    export wwinit_system="${wwinit_uri}&stage=system&compress=gz"; info "warewulf: wwinit_system=${wwinit_system}"
    export wwinit_runtime="${wwinit_uri}&stage=runtime&compress=gz"; info "warewulf: wwinit_runtime=${wwinit_runtime}"
    if [ -n "$(getarg wwinit.KernelOverride)" ]
    then
        export wwinit_kmods="${wwinit_uri}&stage=kmods&compress=gz"; info "warewulf: wwinit_kmods=${wwinit_kmods}"
    fi
    wwinit_tmpfs_size=$(getarg wwinit.tmpfs.size=)
    if [ -n "$wwinit_tmpfs_size" ]
    then
        info "warewulf: wwinit.tmpfs.size=${wwinit_tmpfs_size}"
        export wwinit_tmpfs_size_option="-o size=${wwinit_tmpfs_size}"
    fi
    if [ "${root}" = "persistent" ] ; then 
        export wwinit_persistent="1"; info "warewulf: wwinit_persistent=$wwinit_persistent"
    fi
    if [ -n "${wwinit_container}" ] ; then
        info "warewulf: Found root=${root} and a Warewulf container image. Will boot from Warewulf."
        rootok=1
    else
        die "warewulf: Found root=${root} but no container image. Cannot boot from Warewulf."
    fi
fi
