#!/usr/bin/env nash
import ../../azure/login
import ../../azure/lb

resgroup  = $ARGS[1]
probename = $ARGS[2]
lbname    = $ARGS[3]
port      = $ARGS[4]
protocol  = $ARGS[5]
interval  = $ARGS[6]
count     = $ARGS[7]

azure_login()

probe <= azure_lb_probe_new($probename, $resgroup)
probe <= azure_lb_probe_set_lbname($probe, $lbname)
probe <= azure_lb_probe_set_port($probe, $port)
probe <= azure_lb_probe_set_protocol($probe, $protocol)
probe <= azure_lb_probe_set_interval($probe, $interval)

if len($ARGS) == "9" {
        path = $ARGS[8]
        probe <= azure_lb_probe_set_path($probe, $path)
}

azure_lb_probe_create($probe)
