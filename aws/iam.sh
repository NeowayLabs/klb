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

fn aws_iam_addrole2profile(profileid, role) {
	aws iam add-role-to-instance-profile --instance-profile-name $profileid --role-name $role >[1=]
}