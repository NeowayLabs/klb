# Spot request related functions

fn aws_spot_request_instance(price, type, count, instancejson) {
	requestid <= (
		aws ec2 request-spot-instances
					--spot-price $price
					--instance-count $count
					--type $type
					--launch-specification "file://"+$instancejson |
		jq -r ".SpotInstanceRequests[].SpotInstanceRequestId" |
		xargs echo -n
	)

	return $requestid
}

fn aws_spot_request_cancel(requestids) {
	-aws ec2 cancel-spot-instance-requests --spot-instance-request-ids $requestids

	return $status
}

fn aws_spot_request_describe_all() {
	filters = ()

	requests <= aws_spot_request_describe("", $filters)

	return $requests
}

fn aws_spot_request_get(requestid) {
	json <= (
		aws ec2 describe-spot-instance-requests
						--spot-instance-request-ids $requestid |
		jq ".SpotInstanceRequests[0]"
	)

	return $json
}

fn aws_spot_request_describe(requestid, filters) {
	requests  = ()
	filteropt = ()

	if $requestid != "" {
		requests = ("--spot-instance-request-ids" $requestid)
	}

	filterStr = ""

	for f in $filters {
		if $filterStr == "" {
			filterStr = "Name="+$f[0]+",Values="+$f[1]
		} else {
			filterStr = $filterStr+",Name="+$f[0]+",Values="+$f[1]
		}
	}

	if $filterStr != "" {
		filteropt = ("--filters" $filterStr)
	}

	requests <= (
		aws ec2 describe-spot-instance-requests $requests $filteropt |
		jq -j ".SpotInstanceRequests[] | .InstanceId, \" \", .Status.UpdateTime, \" \", .Status.Code, \" \", .Status.State, \"\n\""
	)

	requests <= split($requests, "\n")

	return $requests
}
