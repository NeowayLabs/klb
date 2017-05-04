import klb/azure/vm

cases = (
	(
		"year"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2018.05.04.1449-vm"
		"klb-ex-bkp-2018.05.04.1449-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"month"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.06.04.1449-vm"
		"klb-ex-bkp-2017.06.04.1449-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"day"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.05.30.1449-vm"
		"klb-ex-bkp-2017.05.30.1449-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"hour"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.05.04.1549-vm"
		"klb-ex-bkp-2017.05.04.1549-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"minute"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.05.04.1450-vm"
		"klb-ex-bkp-2017.05.04.1450-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"one"
		"klb-ex-bkp-2017.05.04.1449-vm"
		"klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"threeResgroups"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.05.04.1451-vm\nklb-ex-bkp-2017.05.04.1450-vm"
		"klb-ex-bkp-2017.05.04.1451-vm klb-ex-bkp-2017.05.04.1450-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
	(
		"ignoreTrailingNewline"
		"klb-ex-bkp-2017.05.04.1449-vm\nklb-ex-bkp-2017.05.04.1450-vm\n"
		"klb-ex-bkp-2017.05.04.1450-vm klb-ex-bkp-2017.05.04.1449-vm"
	)
)

for c in $cases {
	name = $c[0]
	input = $c[1]
	expectation <= echo $c[2]
	ordered <= _azure_vm_backup_order_list($input)
	input_list <= split($input, "\n")
	# WHY: easier to compare :-)
	res <= echo $ordered

	if $expectation != $res {
		print("test %q expected %q got %q\n", $name, $expectation, $res)
		exit("1")
	}
}

print("success")
