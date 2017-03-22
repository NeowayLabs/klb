#!/usr/bin/env nash

resgroup  = $ARGS[1]
probename = $ARGS[2]
lbname    = $ARGS[3]
port      = $ARGS[4]
protocol  = $ARGS[5]
interval  = $ARGS[6]
count     = $ARGS[7]

probe <= azure_lb_probe_new($probename, $resgroup)
probe <= azure_lb_probe_set_lbname($probe, $lbname)
probe <= azure_lb_probe_set_port($probe, $port)
probe <= azure_lb_probe_set_protocol($probe, $protocol)
probe <= azure_lb_probe_set_interval($probe, $interval)

if len($ARGS) == "9" {
        path <= $ARGS[8]
        probe <= azure_lb_probe_set_path($probe, $interval)
}

azure_lb_probe_create($probe)
