#!/usr/bin/env nash

import klb/aws/all

igwTags = (
	(Name klb-igw-tests)
        (Env TEST)
)

fn create() {
	igwId <= aws_igw_create($igwTags)
	echo "Internet gateway created: " $igwId
        return $igwId
}

fn delete(igwId) {
	aws_igw_delete($igwId)
}

fn test_igw() {
	igwId <= create()
	delete($igwId)
}

test_igw()
