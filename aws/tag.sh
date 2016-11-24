# Common aws functions

fn aws_tag(resource, tags) {
	for tag in $tags {
		taglen <= len($tag)

		if $taglen != "2" {
			echo "Invalid tag: " $tag
			abort
		} else {
			(
				aws ec2 create-tags
						--resources $resource
						--tags "Key="+$tag[0]+",Value="+$tag[1]
						>[1=]
			)
		}
	}
}
