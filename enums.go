package main

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

const (
	ErrorCode_NoProps     ErrorCode = "CouldNotReadProperties"
	ErrorCode_NoDataFound ErrorCode = "NoDataFound"
)

type PropsKeys string

func (p PropsKeys) String() string {
	return string(p)
}

const (
	PropsKeys_Ec2AmiName            PropsKeys = "ec2.ami.name"
	PropsKeys_Ec2Service            PropsKeys = "ec2.service"
	PropsKeys_EC2LaunchTemplateName PropsKeys = "ec2.launch.template.name"

	PropsKeys_EbsKmsKeyArn      PropsKeys = "ebs.kms.key.arn"
	PropsKeys_EbsRootDeviceName PropsKeys = "ebs.root.device.name"

	PropsKeys_EksAdminRoleArn  PropsKeys = "eks.admin.role.arn"
	PropsKeys_EksSecretsKeyArn PropsKeys = "eks.secrets.key.arn"
	PropsKeys_EksClusterId     PropsKeys = "eks.cluster.id"

	PropsKeys_IamPolicyEksWorkerKmsArn  PropsKeys = "iam.policy.eks.worker.kms.arn"
	PropsKeys_IamPolicyEksWorkerS3Arn   PropsKeys = "iam.policy.eks.worker.s3.arn"
	PropsKeys_IamPolicyEksWorkerNodeArn PropsKeys = "iam.policy.eks.worker.node.arn"
	PropsKeys_IamPolicyEksWorkerEcrArn  PropsKeys = "iam.policy.eks.worker.ecr.arn"
	PropsKeys_IamPolicyEksWorkerSsmArn  PropsKeys = "iam.policy.eks.worker.ssm.arn"
	PropsKeys_IamPolicyEksWorkerCniArn  PropsKeys = "iam.policy.eks.worker.cni.arn"

	PropsKeys_LoggerLevel PropsKeys = "logger.level"

	PropsKeys_MiscInEverywhere  PropsKeys = "in.everywhere"
	PropsKeys_MiscOutEverywhere PropsKeys = "out.everywhere"

	PropsKeys_StackAccount PropsKeys = "stack.account"
	PropsKeys_StackName    PropsKeys = "stack.name"
	PropsKeys_StackRegion  PropsKeys = "stack.region"

	PropsKeys_TagBilling PropsKeys = "billing"
	PropsKeys_TagEnv     PropsKeys = "env"
	PropsKeys_TagOwner   PropsKeys = "owner"

	PropsKeys_VpcId PropsKeys = "vpc.id"
)
