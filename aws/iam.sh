# IAM related functions

fn aws_iam_create(name, policyFile) {
	roleid <= (
		aws iam create-role
				--role-name $name
				--assume-role-policy-document "file://"+$policyFile |
		jq ".Role.RoleId" |
		xargs echo -n
	)

	return $roleid
}

fn aws_iam_deleterole(name) {
	-aws iam delete-role --role-name $name
}

fn aws_iam_createpolicy(name, path) {
	arn <= (
		aws iam create-policy
				--policy-name $name
				--policy-document "file://"+$path |
		jq -r ".Policy.Arn" |
		xargs echo -n
	)

	return $arn
}

fn aws_iam_deletepolicy(policyArn) {
	-aws iam delete-policy --policy-arn $policyArn
}

fn aws_iam_attachpolicy(roleName, policyArn) {
	aws iam attach-role-policy --policy-arn $policyArn --role-name $roleName
}

fn aws_iam_putpolicy(roleName, policyName, policyFile) {
	aws iam put-role-policy --role-name $roleName --policy-name $policyName --policy-document "file://"+$policyFile >[1=]
}

fn aws_iam_profile(name) {
	profileid <= (
		aws iam create-instance-profile
					--instance-profile-name $name |
		jq ".InstanceProfile.InstanceProfileId" |
		xargs echo -n
	)

	return $profileid
}

fn aws_iam_deleteprofile(name) {
	-aws iam delete-instance-profile --instance-profile-name $name
}

fn aws_iam_addrole2profile(profile, role) {
	aws iam add-role-to-instance-profile --instance-profile-name $profile --role-name $role >[1=]
}

fn aws_iam_delete_rolefromprofile(role, profile) {
	-aws iam remove-role-from-instance-profile --instance-profile-name $profile --role-name $role
}
