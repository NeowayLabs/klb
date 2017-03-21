#!/usr/bin/env nash

import ../../azure/login

resgroup        = $ARGS[1]
rulename        = $ARGS[2]
lbname          = $ARGS[3]
probename       = $ARGS[4]
frontendipname  = $ARGS[5]
frontendport    = $ARGS[6]
backendport     = $ARGS[7]
protocol        = $ARGS[8]
addresspoolname = $ARGS[9]

rule <= azure_lb_rule_new($rulename, $resgroup)
rule <= azure_lb_rule_set_lbname($rule, $lbname)
rule <= azure_lb_rule_set_probename($rule, $probename)
rule <= azure_lb_rule_set_frontendipname($rule, $frontendipname)
rule <= azure_lb_rule_set_frontendport($rule, $frontendport)
rule <= azure_lb_rule_set_backendport($rule, $backendport)
rule <= azure_lb_rule_set_protocol($rule, $protocol)
rule <= azure_lb_rule_set_addresspoolname($rule, $addresspoolname)

azure_lb_rule_create($rule)
