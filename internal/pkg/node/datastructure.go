package node

import (
	"net"

	"github.com/warewulf/warewulf/internal/pkg/wwtype"
)

const undef string = "UNDEF"

/******
 * YAML data representations
 ******/
/*
Structure of which goes to disk
*/
type NodesYaml struct {
	NodeProfiles map[string]*Profile `yaml:"nodeprofiles"`
	Nodes        map[string]*Node    `yaml:"nodes"`
}

/*
Node is the datastructure describing a node and a profile which in disk format.
*/
type Node struct {
	id    string
	valid bool // Is set true, if called by the constructor
	// exported values
	Discoverable wwtype.WWbool     `yaml:"discoverable,omitempty" lopt:"discoverable" sopt:"e" comment:"Make discoverable in given network (true/false)"`
	AssetKey     string            `yaml:"asset key,omitempty" lopt:"asset" comment:"Set the node's Asset tag (key)"`
	Profile      `yaml:"-,inline"` // include all values set in the profile, but inline them in yaml output if these are part of Node
}

/*
Holds the data which can be set for profiles and nodes.
*/
type Profile struct {
	id string
	// exported values
	Profiles       []string               `yaml:"profiles,omitempty" lopt:"profile" sopt:"P" comment:"Set the node's profile members (comma separated)"`
	Comment        string                 `yaml:"comment,omitempty" lopt:"comment" comment:"Set arbitrary string comment"`
	ClusterName    string                 `yaml:"cluster name,omitempty" lopt:"cluster" sopt:"c" comment:"Set cluster group"`
	ContainerName  string                 `yaml:"container name,omitempty" lopt:"container" sopt:"C" comment:"Set container name"`
	Ipxe           string                 `yaml:"ipxe template,omitempty" lopt:"ipxe" comment:"Set the iPXE template name"`
	RuntimeOverlay []string               `yaml:"runtime overlay,omitempty" lopt:"runtime" sopt:"R" comment:"Set the runtime overlay"`
	SystemOverlay  []string               `yaml:"system overlay,omitempty" lopt:"wwinit" sopt:"O" comment:"Set the system overlay"`
	Kernel         *KernelConf            `yaml:"kernel,omitempty"`
	Ipmi           *IpmiConf              `yaml:"ipmi,omitempty"`
	Init           string                 `yaml:"init,omitempty" lopt:"init" sopt:"i" comment:"Define the init process to boot the container"`
	Root           string                 `yaml:"root,omitempty" lopt:"root" comment:"Define the rootfs" `
	NetDevs        map[string]*NetDev     `yaml:"network devices,omitempty"`
	Tags           map[string]string      `yaml:"tags,omitempty"`
	PrimaryNetDev  string                 `yaml:"primary network,omitempty" lopt:"primarynet" sopt:"p" comment:"Set the primary network interface"`
	Disks          map[string]*Disk       `yaml:"disks,omitempty"`
	FileSystems    map[string]*FileSystem `yaml:"filesystems,omitempty"`
	Resources      map[string]Resource    `yaml:"resources,omitempty"`
}

type IpmiConf struct {
	UserName   string            `yaml:"username,omitempty" lopt:"ipmiuser" comment:"Set the IPMI username"`
	Password   string            `yaml:"password,omitempty" lopt:"ipmipass" comment:"Set the IPMI password"`
	Ipaddr     net.IP            `yaml:"ipaddr,omitempty" lopt:"ipmiaddr" comment:"Set the IPMI IP address" type:"IP"`
	Gateway    net.IP            `yaml:"gateway,omitempty" lopt:"ipmigateway" comment:"Set the IPMI gateway" type:"IP"`
	Netmask    net.IP            `yaml:"netmask,omitempty" lopt:"ipminetmask" comment:"Set the IPMI netmask" type:"IP"`
	Port       string            `yaml:"port,omitempty" lopt:"ipmiport" comment:"Set the IPMI port"`
	Interface  string            `yaml:"interface,omitempty" lopt:"ipmiinterface" comment:"Set the node's IPMI interface (defaults: 'lan')"`
	EscapeChar string            `yaml:"escapechar,omitempty" lopt:"ipmiescapechar" comment:"Set the IPMI escape character (defaults: '~')"`
	Write      wwtype.WWbool     `yaml:"write,omitempty" lopt:"ipmiwrite" comment:"Enable the write of impi configuration (true/false)"`
	Template   string            `yaml:"template,omitempty" lopt:"ipmitemplate" comment:"template used for ipmi command"`
	Tags       map[string]string `yaml:"tags,omitempty"`
}

type KernelConf struct {
	Version string `yaml:"version,omitempty" lopt:"kernelversion" comment:"Set kernel version" json:"version,omitempty"`
	Args    string `yaml:"args,omitempty" lopt:"kernelargs" sopt:"A" comment:"Set kernel arguments" json:"args,omitempty"`
}

type NetDev struct {
	Type    string            `yaml:"type,omitempty" lopt:"type" sopt:"T" comment:"Set device type of given network"`
	OnBoot  wwtype.WWbool     `yaml:"onboot,omitempty" lopt:"onboot" comment:"Enable/disable network device (true/false)"`
	Device  string            `yaml:"device,omitempty" lopt:"netdev" sopt:"N" comment:"Set the device for given network"`
	Hwaddr  string            `yaml:"hwaddr,omitempty" lopt:"hwaddr" sopt:"H" comment:"Set the device's HW address for given network" type:"MAC"`
	Ipaddr  net.IP            `yaml:"ipaddr,omitempty" comment:"IPv4 address in given network" sopt:"I" lopt:"ipaddr" type:"IP"`
	Ipaddr6 net.IP            `yaml:"ip6addr,omitempty" lopt:"ipaddr6" comment:"IPv6 address" type:"IP"`
	Prefix  net.IP            `yaml:"prefix,omitempty"`
	Netmask net.IP            `yaml:"netmask,omitempty" lopt:"netmask" sopt:"M" comment:"Set the networks netmask" type:"IP"`
	Gateway net.IP            `yaml:"gateway,omitempty" lopt:"gateway" sopt:"G" comment:"Set the node's network device gateway" type:"IP"`
	MTU     string            `yaml:"mtu,omitempty" lopt:"mtu" comment:"Set the mtu" type:"uint"`
	Tags    map[string]string `yaml:"tags,omitempty"`
	primary bool
}

/*
Holds the disks of a node
*/
type Disk struct {
	WipeTable  bool                  `yaml:"wipe_table,omitempty" lopt:"diskwipe" comment:"whether or not the partition tables shall be wiped"`
	Partitions map[string]*Partition `yaml:"partitions,omitempty"`
}

/*
partition definition, the label must be uniq so its used as the key in the
Partitions map
*/
type Partition struct {
	Number             string `yaml:"number,omitempty" lopt:"partnumber" comment:"set the partition number, if not set next free slot is used" type:"uint"`
	SizeMiB            string `yaml:"size_mib,omitempty" lopt:"partsize" comment:"set the size of the partition, if not set maximal possible size is used"  type:"uint"`
	StartMiB           string `yaml:"start_mib,omitempty" comment:"the start of the partition" type:"uint"`
	TypeGuid           string `yaml:"type_guid,omitempty" comment:"Linux filesystem data will be used if empty"`
	Guid               string `yaml:"guid,omitempty" comment:"the GPT unique partition GUID"`
	WipePartitionEntry bool   `yaml:"wipe_partition_entry,omitempty" comment:"if true, Ignition will clobber an existing partition if it does not match the config"`
	ShouldExist        bool   `yaml:"should_exist,omitempty" lopt:"partcreate" comment:"create partition if not exist"`
	Resize             bool   `yaml:"resize,omitempty" comment:" whether or not the existing partition should be resize"`
}

/*
Definition of a filesystem. The device is uniq so its used as key
*/
type FileSystem struct {
	Format         string   `yaml:"format,omitempty" lopt:"fsformat" comment:"format of the file system"`
	Path           string   `yaml:"path,omitempty" lopt:"fspath" comment:"the mount point of the file system"`
	WipeFileSystem bool     `yaml:"wipe_filesystem,omitempty" lopt:"fswipe" comment:"wipe file system at boot"`
	Label          string   `yaml:"label,omitempty" comment:"the label of the filesystem"`
	Uuid           string   `yaml:"uuid,omitempty" comment:"the uuid of the filesystem"`
	Options        []string `yaml:"options,omitempty" comment:"any additional options to be passed to the format-specific mkfs utility"`
	MountOptions   string   `yaml:"mount_options,omitempty" comment:"any special options to be passed to the mount command"`
}

type Resource interface{}
