// +build slow

package azure_test

import (
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func TestVMBackup(t *testing.T) {
	t.Parallel()
	vmtesttimeout := 45 * time.Minute
	fixture.Run(t, "VMBackupOsDiskOnly", vmtesttimeout, location, testVMBackupOsDiskOnly)
	fixture.Run(t, "VMBackup", vmtesttimeout, location, testVMBackup)
}

func testVMBackupOsDiskOnly(t *testing.T, f fixture.F) {
	vmSize := "Basic_A2"
	sku := "Standard_LRS"
	bkpprefix := "klb-bkp"

	resources := createVMResources(t, f)
	vm := createVM(t, f, resources.availSet, resources.nic, vmSize, sku)

	defer deleteBackups(t, f, vm, bkpprefix)
	bkpresgroup := backupVM(t, f, vm, bkpprefix)
	assertResourceGroupExists(t, f, bkpresgroup)

	// TODO: call restore procedure
	// TODO: validate VMs have the same osdisk
}

func testVMBackup(t *testing.T, f fixture.F) {

	vmSize := "Basic_A2"
	sku := "Standard_LRS"
	bkpprefix := "klb-tests"

	f.Shell.DisableTryAgain() // TODO: REMOVE THIS

	resources := createVMResources(t, f)
	vm := createVM(t, f, resources.availSet, resources.nic, vmSize, sku)

	disks := []VMDisk{
		// Different sizes is important to validate behavior
		{Name: genUniqName(), Size: 50, Sku: sku},
		{Name: genUniqName(), Size: 20, Sku: sku},
		{Name: genUniqName(), Size: 100, Sku: sku},
	}
	attachDisks(t, f, vm, disks)

	defer deleteBackups(t, f, vm, bkpprefix)

	bkpresgroup := backupVM(t, f, vm, bkpprefix)
	assertResourceGroupExists(t, f, bkpresgroup)

	backups := listBackups(t, f, vm, bkpprefix)
	assertEqualStringSlice(t, []string{bkpresgroup}, backups)

	recoveredVMName := recoverVM(
		t,
		f,
		resources.vnet,
		resources.subnet,
		vmSize,
		sku,
		bkpresgroup,
	)

	assertRecoveredVMDisks(t, f, vm, recoveredVMName)
}

func assertEqualStringSlice(t *testing.T, slice1 []string, slice2 []string) {
	if len(slice2) != len(slice2) {
		t.Fatalf("%+v != %+v", slice1, slice2)
	}

	for _, str1 := range slice1 {
		found := false
		for _, str2 := range slice2 {
			if str1 == str2 {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("%+v != %+v", slice1, slice2)
		}
	}
}

func assertResourceGroupExists(t *testing.T, f fixture.F, resgroup string) {
	fixture.NewResourceGroup(f.Ctx, t, f.Session, f.Logger).AssertExists(t, resgroup)
}

func assertRecoveredVMDisks(t *testing.T, f fixture.F, vmName string, recoveredVMName string) {
	vm := azure.NewVM(f)
	originalOSDisk := vm.OsDisk(t, vmName)
	restoredOSDisk := vm.OsDisk(t, recoveredVMName)

	// WHY: Names cant be equal
	if originalOSDisk.SizeGB != restoredOSDisk.SizeGB ||
		originalOSDisk.OsType != restoredOSDisk.OsType {
		t.Fatalf("os disk: %+v != %+v", originalOSDisk, restoredOSDisk)
	}

	originalDataDisks := vm.DataDisks(t, vmName)
	recoveredDataDisks := vm.DataDisks(t, recoveredVMName)

	fail := func() {
		t.Fatalf("expected data disks:\n%+v\n\ngot:\n%+v\n\n", originalDataDisks, recoveredDataDisks)
	}

	if len(originalDataDisks) != len(recoveredDataDisks) {
		fail()
	}

	for _, dataDisk := range originalDataDisks {
		found := false
		for _, recoveredDataDisk := range recoveredDataDisks {
			// WHY: Disks cant have the same name
			// LUN is fundamental for backup process, must be identical
			if recoveredDataDisk.Lun == dataDisk.Lun {
				if recoveredDataDisk.SizeGB != dataDisk.SizeGB {
					t.Fatalf(
						"expected disks with same LUN to have same size, %+v != %+v",
						dataDisk,
						recoveredDataDisk,
					)
				}
				found = true
			}
		}

		if !found {
			fail()
		}
	}
}

func backupVM(t *testing.T, f fixture.F, vmname string, prefix string) string {

	res := execWithIPC(t, f, func(output string) {
		f.Shell.Run(
			"./testdata/backup_vm.sh",
			vmname,
			f.ResGroupName,
			prefix,
			f.Location,
			output,
		)
	})
	res = strings.TrimSpace(res)
	return strings.Trim(res, "\n")
}

func listBackups(t *testing.T, f fixture.F, vmname string, prefix string) []string {

	res := execWithIPC(t, f, func(outputfile string) {
		f.Shell.Run(
			"./testdata/list_backups.sh",
			vmname,
			prefix,
			outputfile,
		)
	})

	return strings.Split(strings.Trim(strings.TrimSpace(res), "\n"), "\n")
}

func deleteBackups(t *testing.T, f fixture.F, vmname string, bkpprefix string) {
	f.Shell.Run("./testdata/delete_backups.sh", vmname, bkpprefix)
}

func recoverVM(
	t *testing.T,
	f fixture.F,
	vnet string,
	subnet string,
	vmSize string,
	sku string,
	backupResgroup string,
) string {
	vmName := "recoveredVM-" + genUniqName()
	keyFile := "./testdata/key.pub"
	ostype := "linux"

	f.Shell.Run(
		"./testdata/recover_backup.sh",
		vmName,
		f.ResGroupName,
		f.Location,
		vmSize,
		vnet,
		subnet,
		keyFile,
		ostype,
		sku,
		backupResgroup,
	)
	return vmName
}
