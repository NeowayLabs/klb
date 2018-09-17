#!/usr/bin/env nash

import klb/azure/login
import klb/azure/lb

resgroup        = $ARGS[1]
rulename        = $ARGS[2]
lbname          = $ARGS[3]
probename       = $ARGS[4]
frontendipname  = $ARGS[5]
backendpoolname = $ARGS[6]
protocol        = $ARGS[7]
frontendport    = $ARGS[8]
backendport     = $ARGS[9]

azure_login()

rule <= azure_lb_rule_new($rulename, $resgroup)
rule <= azure_lb_rule_set_lbname($rule, $lbname)
rule <= azure_lb_rule_set_probename($rule, $probename)
rule <= azure_lb_rule_set_frontendipname($rule, $frontendipname)
rule <= azure_lb_rule_set_frontendport($rule, $frontendport)
rule <= azure_lb_rule_set_backendport($rule, $backendport)
rule <= azure_lb_rule_set_protocol($rule, $protocol)
rule <= azure_lb_rule_set_backend_pool_name($rule, $backendpoolname)

azure_lb_rule_create($rule)
