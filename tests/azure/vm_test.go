package azure_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genUniqName() string {
	return fmt.Sprintf("klbvmtests-%d", rand.Intn(9999999))
}

type VMResources struct {
	availSet string
	vnet     string
	subnet   string
	nic      string
}

func createVM(
	t *testing.T,
	f fixture.F,
	availset string,
	nic string,
	vmSize string,
	sku string,
) string {
	vm := genUniqName()
	username := "core"
	osDisk := genUniqName()
	imageUrn := "OpenLogic:CentOS:7.2:7.2.20161026"
	keyFile := "./testdata/key.pub"

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		vmSize,
		username,
		availset,
		nic,
		osDisk,
		imageUrn,
		keyFile,
		sku,
	)

	f.Logger.Println("creating VM")
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm, availset, vmSize, nic)
	f.Logger.Println("created VM with success, attaching a disk")
	return vm
}

func testVMCreation(t *testing.T, f fixture.F, vmSize string, sku string) {

	resources := createVMResources(t, f)
	vm := createVM(t, f, resources.availSet, resources.nic, vmSize, sku)

	diskname := "createVMExtraDisk"
	size := 10

	attachNewDiskOnVM(t, f, vm, diskname, size, sku)

	vms := azure.NewVM(f)
	vms.AssertAttachedDataDisk(t, vm, diskname, size, sku)

	f.Logger.Println("VM with attached disk created with success")
}

func testStandardDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Basic_A0", "Standard_LRS")
}

func testPremiumDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Standard_DS4_v2", "Premium_LRS")
}

func execWithIPC(t *testing.T, f fixture.F, exec func(string)) string {
	outfile, err := ioutil.TempFile("", "klb_vm_script_ipc")
	if err != nil {
		t.Fatalf("error creating output file: %s", err)
	}
	defer os.Remove(outfile.Name()) // clean up

	exec(outfile.Name())

	f.Logger.Println("executed script, reading output")
	out, err := ioutil.ReadAll(outfile)
	if err != nil {
		t.Fatalf("error reading output file: %s", err)
	}
	return string(out)
}

type VMDisk struct {
	Name string
	Sku  string
	Size int
}

func attachDisks(t *testing.T, f fixture.F, vmname string, disks []VMDisk) {
	// TODO: Improve using concurrency.
	vms := azure.NewVM(f)

	for _, disk := range disks {
		attachNewDiskOnVM(t, f, vmname, disk.Name, disk.Size, disk.Sku)
		vms.AssertAttachedDataDisk(t, vmname, disk.Name, disk.Size, disk.Sku)
	}
}

func testVMSnapshot(t *testing.T, f fixture.F, vmSize string, sku string) {
	resources := createVMResources(t, f)
	vm := createVM(t, f, resources.availSet, resources.nic, vmSize, sku)

	disks := []VMDisk{
		// Different sizes is important to validate behavior
		{Name: genUniqName(), Size: 10, Sku: sku},
		{Name: genUniqName(), Size: 20, Sku: sku},
		{Name: genUniqName(), Size: 30, Sku: sku},
	}

	vms := azure.NewVM(f)

	attachDisks(t, f, vm, disks)

	idsraw := execWithIPC(t, f, func(outfile string) {
		f.Shell.Run(
			"./testdata/create_vm_snapshots.sh",
			f.ResGroupName,
			vm,
			outfile,
		)
	})

	f.Logger.Println("created snapshots, retrieving ids")
	ids := strings.Split(strings.Trim(string(idsraw), "\n"), "\n")
	f.Logger.Printf("parsed ids: %s", ids)

	if len(ids) != len(disks) {
		t.Fatalf("expected %d snapshots, got %d", len(disks), len(ids))
	}

	nic := genNicName()
	createVMNIC(f, nic, resources.vnet, resources.subnet)
	vmbackup := createVM(t, f, resources.availSet, nic, vmSize, sku)

	for _, id := range ids {
		attachSnapshotOnVM(t, f, vmbackup, id, sku)
	}

	originaldisks := vms.DataDisks(t, vm)
	backupdisks := vms.DataDisks(t, vmbackup)

	if len(originaldisks) != len(backupdisks) {
		t.Fatalf("expected disks %q == %q", originaldisks, backupdisks)
	}

	for _, originaldisk := range originaldisks {
		size := originaldisk.SizeGB
		got := false
		for _, backupdisk := range backupdisks {
			if backupdisk.SizeGB == size {
				got = true
			}
		}
		if !got {
			t.Fatalf(
				"unable to find disk: %q on backup disks %q",
				originaldisk,
				backupdisks,
			)
		}
	}
}

func testVMSnapshotStandard(t *testing.T, f fixture.F) {
	testVMSnapshot(t, f, "Basic_A2", "Standard_LRS")
}

func testVMSnapshotPremium(t *testing.T, f fixture.F) {
	testVMSnapshot(t, f, "Standard_DS4_v2", "Premium_LRS")
}

func testDuplicatedAvailabilitySet(t *testing.T, f fixture.F) {
	name := "duplicatedAvSet"

	createAvSet := func() {
		f.Shell.Run(
			"./testdata/create_vm_avail_set.sh",
			f.ResGroupName,
			name,
			f.Location,
			"1",
			"1",
		)
	}

	createAvSet()
	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, name)

	createAvSet()
	availSets.AssertExists(t, name)
}

func attachDiskOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	diskname string,
	diskSizeGB int,
	sku string,
) {
	f.Shell.Run(
		"./testdata/attach_new_disk.sh",
		f.ResGroupName,
		vmname,
		diskname,
		strconv.Itoa(diskSizeGB),
		sku,
	)
}

func attachSnapshotOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	snapshotid string,
	disksku string,
) string {
	diskname := genUniqName()
	f.Shell.Run(
		"./testdata/attach_snapshot.sh",
		f.ResGroupName,
		f.Location,
		vmname,
		diskname,
		disksku,
		snapshotid,
	)
	return diskname
}

func createVMResources(t *testing.T, f fixture.F) VMResources {

	resources := VMResources{}
	resources.availSet = genAvailSetName()
	resources.vnet = genVnetName()
	resources.subnet = genSubnetName()
	resources.nic = genNicName()

	nsg := genNsgName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	updatedomain := "3"
	faultdomain := "3"

	f.Shell.Run(
		"./testdata/create_vm_avail_set.sh",
		f.ResGroupName,
		resources.availSet,
		f.Location,
		updatedomain,
		faultdomain,
	)

	createVNet(t, f, vnetDescription{
		name:     resources.vnet,
		vnetAddr: vnetAddress,
	})

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		resources.subnet,
		f.ResGroupName,
		resources.vnet,
		subnetAddress,
		nsg,
	)

	createVMNIC(f, resources.nic, resources.vnet, resources.subnet)

	return resources
}

func createVMNIC(f fixture.F, nic string, vnet string, subnet string) {
	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,
		subnet,
	)
}

func attachNewDiskOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	diskname string,
	diskSizeGB int,
	sku string,
) {
	f.Shell.Run(
		"./testdata/attach_new_disk.sh",
		f.ResGroupName,
		vmname,
		diskname,
		strconv.Itoa(diskSizeGB),
		sku,
	)
}

func TestVM(t *testing.T) {
	t.Parallel()
	vmtesttimeout := 45 * time.Minute
	fixture.Run(t, "VMCreationStandardDisk", vmtesttimeout, location, testStandardDiskVM)
	fixture.Run(t, "VMCreationPremiumDisk", vmtesttimeout, location, testPremiumDiskVM)
	fixture.Run(t, "VMSnapshotStandard", vmtesttimeout, location, testVMSnapshotStandard)
	fixture.Run(t, "VMSnapshotPremium", vmtesttimeout, location, testVMSnapshotPremium)
	fixture.Run(t, "VMDuplicatedAvSet", vmtesttimeout, location, testDuplicatedAvailabilitySet)
}
