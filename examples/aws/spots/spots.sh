#!/usr/bin/env nash

import klb/aws/all

fn create_spot_requests() {
	requestid <= aws_spot_request_instance("0.004", "one-time", "1", "spot1.json")

	echo "Spot request: " $requestid

	# describe request created
	filters = ()

	requests <= aws_spot_request_describe($requestid, $filters)

	for r in $requests {
		printf "%s%s%s\n" $NASH_GREEN $r $NASH_RESET
	}

	aws_spot_request_cancel($requestid)
}

create_spot_requests()
