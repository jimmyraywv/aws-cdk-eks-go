package main

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	kms "github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/jsii-runtime-go"
	Log "github.com/sirupsen/logrus"
	"jimmyray.io/cdk/utils"
	"os"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	// "github.com/aws/jsii-runtime-go"
)

const (
	PropsFile string = "config.properties"
)

var p utils.Properties

type CdkStackProps struct {
	cdk.StackProps
}

func CdkEksStack(scope constructs.Construct, id string, props *CdkStackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := cdk.NewStack(scope, &id, &sprops)

	var tagsLabels = make(map[string]*string)
	tagsLabels[PropsKeys_TagBilling.String()] = jsii.String(p[PropsKeys_TagBilling.String()])
	tagsLabels[PropsKeys_TagEnv.String()] = jsii.String(p[PropsKeys_TagEnv.String()])
	tagsLabels[PropsKeys_TagOwner.String()] = jsii.String(p[PropsKeys_TagOwner.String()])

	//var cfnTags = make([]cdk.CfnTag, int(len(tagsLabels)))
	//for k, v := range tagsLabels {
	//	cfnTags = append(cfnTags, cdk.CfnTag{Key: jsii.String(k), Value: v})
	//}

	// Lookup exiting VPC
	var vcpId = p[PropsKeys_VpcId.String()]
	vpc := ec2.Vpc_FromLookup(stack, jsii.String(vcpId), &ec2.VpcLookupOptions{VpcId: jsii.String(vcpId)})

	// SGs
	var sharedNodeSg = ec2.NewSecurityGroup(stack, jsii.String("SharedNode"),
		&ec2.SecurityGroupProps{SecurityGroupName: jsii.String("cdk-eks-SharedNode"),
			AllowAllOutbound: jsii.Bool(true), Description: jsii.String(""), Vpc: vpc})

	var controlPlaneSg = ec2.NewSecurityGroup(stack, jsii.String("ControlPlane"),
		&ec2.SecurityGroupProps{SecurityGroupName: jsii.String("cdk-eks-ControlPlane"),
			AllowAllOutbound: jsii.Bool(true), Description: jsii.String(""), Vpc: vpc})

	var eniSg = ec2.NewSecurityGroup(stack, jsii.String("ClusterENI"),
		&ec2.SecurityGroupProps{SecurityGroupName: jsii.String("cdk-eks-ClusterENI"),
			AllowAllOutbound: jsii.Bool(true), Description: jsii.String(""), Vpc: vpc})

	sharedNodeSg.AddEgressRule(ec2.Peer_AnyIpv4(), ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscOutEverywhere.String()]), jsii.Bool(false))
	sharedNodeSg.AddIngressRule(sharedNodeSg, ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscInEverywhere.String()]), jsii.Bool(false))
	sharedNodeSg.AddIngressRule(eniSg, ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscInEverywhere.String()]), jsii.Bool(false))

	controlPlaneSg.AddEgressRule(ec2.Peer_AnyIpv4(), ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscOutEverywhere.String()]), jsii.Bool(false))

	eniSg.AddEgressRule(ec2.Peer_AnyIpv4(), ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscOutEverywhere.String()]), jsii.Bool(false))
	eniSg.AddIngressRule(sharedNodeSg, ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscInEverywhere.String()]), jsii.Bool(false))
	eniSg.AddIngressRule(eniSg, ec2.Port_AllTcp(), jsii.String(p[PropsKeys_MiscInEverywhere.String()]), jsii.Bool(false))

	// EKS cluster
	var clusterId = p[PropsKeys_EksClusterId.String()]
	var cluster = eks.NewCluster(stack, jsii.String(clusterId), &eks.ClusterProps{
		Version:           eks.KubernetesVersion_V1_21(),
		ClusterName:       jsii.String(clusterId),
		OutputClusterName: jsii.Bool(true),
		SecurityGroup:     eniSg,
		Vpc:               vpc,
		VpcSubnets: &[]*ec2.SubnetSelection{
			{
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
			},
			{
				SubnetType: ec2.SubnetType_PUBLIC,
			},
		},
		EndpointAccess: eks.EndpointAccess_PUBLIC_AND_PRIVATE(),
		MastersRole: iam.Role_FromRoleArn(stack, jsii.String("master-role"), jsii.String(p[PropsKeys_EksAdminRoleArn.String()]),
			&iam.FromRoleArnOptions{Mutable: jsii.Bool(false)}),
		OutputMastersRoleArn: jsii.Bool(true),
		SecretsEncryptionKey: kms.Key_FromKeyArn(stack, jsii.String("eks-secrets-key"), jsii.String(p[PropsKeys_EksSecretsKeyArn.String()])),
		ClusterLogging: &[]eks.ClusterLoggingTypes{eks.ClusterLoggingTypes_API, eks.ClusterLoggingTypes_AUDIT,
			eks.ClusterLoggingTypes_AUTHENTICATOR, eks.ClusterLoggingTypes_SCHEDULER, eks.ClusterLoggingTypes_CONTROLLER_MANAGER},
		DefaultCapacity: jsii.Number(0),
		Tags:            &tagsLabels,
	})

	// Lookup managed IAM policies
	var managedPolicies []iam.IManagedPolicy
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("kms-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerKmsArn.String()])))
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("s3-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerS3Arn.String()])))
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("node-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerNodeArn.String()])))
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("ecr-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerEcrArn.String()])))
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("ssm-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerSsmArn.String()])))
	managedPolicies = append(managedPolicies, iam.ManagedPolicy_FromManagedPolicyArn(stack, jsii.String("cni-policy"), jsii.String(p[PropsKeys_IamPolicyEksWorkerCniArn.String()])))

	// EKS worker node role
	var workerRole = iam.NewRole(stack, jsii.String("worker-role"), &iam.RoleProps{
		AssumedBy: iam.NewServicePrincipal(jsii.String(p[PropsKeys_Ec2Service.String()]), &iam.ServicePrincipalOpts{
			Conditions: nil,
		}),
		Description:        jsii.String("CDK EKS Worker Role"),
		ManagedPolicies:    &managedPolicies,
		MaxSessionDuration: cdk.Duration_Hours(jsii.Number(1)),
		RoleName:           jsii.String("cdk-eks-workers"),
	})

	//var ami = ec2.MachineImage_Lookup(&ec2.LookupMachineImageProps{
	//	Name: jsii.String(p[PropsKeys_Ec2AmiName.String()]),
	//})

	// EBS volume
	var bd = ec2.NewBlockDeviceVolume(&ec2.EbsDeviceProps{DeleteOnTermination: jsii.Bool(true),
		VolumeSize: jsii.Number(100), VolumeType: ec2.EbsDeviceVolumeType_GP2,
		Encrypted: jsii.Bool(true), KmsKey: kms.Key_FromKeyArn(stack, jsii.String("enc-key"),
			jsii.String(p[PropsKeys_EbsKmsKeyArn.String()])),
	}, jsii.String(p[PropsKeys_EbsRootDeviceName.String()]))

	var lt = ec2.NewLaunchTemplate(stack, jsii.String(p[PropsKeys_EC2LaunchTemplateName.String()]), &ec2.LaunchTemplateProps{
		BlockDevices: &[]*ec2.BlockDevice{{DeviceName: jsii.String(p[PropsKeys_EbsRootDeviceName.String()]),
			Volume: bd, MappingEnabled: jsii.Bool(true)}},
		InstanceType:       ec2.InstanceType_Of(ec2.InstanceClass_STANDARD5_AMD, ec2.InstanceSize_XLARGE),
		LaunchTemplateName: jsii.String(p[PropsKeys_EC2LaunchTemplateName.String()]),
		//MachineImage:                      ami,
		NitroEnclaveEnabled: jsii.Bool(false),
		RequireImdsv2:       jsii.Bool(true),
		SecurityGroup:       eniSg,
	})

	// Add LT tags
	lt.Tags().SetTag(jsii.String(PropsKeys_TagBilling.String()), jsii.String(p[PropsKeys_TagBilling.String()]), jsii.Number(0), jsii.Bool(true))
	lt.Tags().SetTag(jsii.String(PropsKeys_TagEnv.String()), jsii.String(p[PropsKeys_TagEnv.String()]), jsii.Number(0), jsii.Bool(true))
	lt.Tags().SetTag(jsii.String(PropsKeys_TagOwner.String()), jsii.String(p[PropsKeys_TagOwner.String()]), jsii.Number(0), jsii.Bool(true))

	// Escape hatch to override IMDS properties
	var cfnLt = lt.Node().DefaultChild().(ec2.CfnLaunchTemplate)
	cfnLt.AddPropertyOverride(jsii.String("LaunchTemplateData.MetadataOptions"), map[string]interface{}{
		"HttpPutResponseHopLimit": jsii.Number(1),
	})

	eks.NewNodegroup(stack, jsii.String("ng-al2"), &eks.NodegroupProps{
		AmiType:            eks.NodegroupAmiType_AL2_X86_64,
		CapacityType:       eks.CapacityType_ON_DEMAND,
		Cluster:            cluster,
		DesiredSize:        jsii.Number(3),
		ForceUpdate:        jsii.Bool(true),
		Labels:             &tagsLabels,
		LaunchTemplateSpec: &eks.LaunchTemplateSpec{Id: lt.LaunchTemplateId(), Version: lt.VersionNumber()},
		MaxSize:            jsii.Number(5),
		MinSize:            jsii.Number(3),
		NodegroupName:      jsii.String("al2"),
		NodeRole:           workerRole,
		Subnets:            &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT},
		Tags:               &tagsLabels,
	})

	return stack
}

func main() {
	utils.InitLogs(nil, Log.DebugLevel)

	var err error
	p, err = utils.ReadProperties(PropsFile)
	if err != nil {
		errorData := utils.ErrorLog{Skip: 1, Event: ErrorCode_NoDataFound.String(), Message: ErrorCode_NoProps.String(), ErrorData: err.Error()}
		utils.LogErrors(errorData)
		os.Exit(1)
	}

	switch p[PropsKeys_LoggerLevel.String()] {
	case "debug":
		utils.Logger.SetLevel(Log.DebugLevel)
	case "error":
		utils.Logger.SetLevel(Log.ErrorLevel)
	case "fatal":
		utils.Logger.SetLevel(Log.FatalLevel)
	case "info":
		utils.Logger.SetLevel(Log.InfoLevel)
	case "warn":
		utils.Logger.SetLevel(Log.WarnLevel)
	default:
		utils.Logger.SetLevel(Log.DebugLevel)
	}

	utils.Logger.WithFields(Log.Fields{"properties": p.String()}).Debug("properties")

	app := cdk.NewApp(nil)
	e := env()

	utils.Logger.WithFields(Log.Fields{"env": e}).Debug("env")

	CdkEksStack(app, p[PropsKeys_StackName.String()], &CdkStackProps{
		cdk.StackProps{
			Env: e,
		},
	})

	app.Synth(nil)
}

func env() *cdk.Environment {
	return &cdk.Environment{
		Account: jsii.String(p[PropsKeys_StackAccount.String()]),
		Region:  jsii.String(p[PropsKeys_StackRegion.String()]),
	}
}
