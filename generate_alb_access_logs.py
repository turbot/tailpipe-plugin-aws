import boto3
import time
import requests
import logging
from botocore.exceptions import ClientError

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class AWSInfraSetup:
    def __init__(self, region='us-east-1'):
        self.ec2 = boto3.client('ec2', region_name=region)
        self.elbv2 = boto3.client('elbv2', region_name=region)
        self.s3 = boto3.client('s3', region_name=region)
        self.region = region
        # Store resource IDs for cleanup
        self.resources = {
            'instance_id': None,
            'security_group_id': None,
            'target_group_arn': None,
            'alb_arn': None
        }

    def create_security_group(self):
        try:
            response = self.ec2.create_security_group(
                GroupName='TestALBLoggingSG',
                Description='Security group for web server'
            )
            security_group_id = response['GroupId']
            
            # Allow inbound HTTP traffic
            self.ec2.authorize_security_group_ingress(
                GroupId=security_group_id,
                IpPermissions=[
                    {
                        'IpProtocol': 'tcp',
                        'FromPort': 80,
                        'ToPort': 80,
                        'IpRanges': [{'CidrIp': '0.0.0.0/0'}]
                    },
                    {
                        'IpProtocol': 'tcp',
                        'FromPort': 22,
                        'ToPort': 22,
                        'IpRanges': [{'CidrIp': '0.0.0.0/0'}]
                    }
                ]
            )
            logger.info(f"Created Security Group: {security_group_id}")
            self.resources['security_group_id'] = security_group_id
            return security_group_id
        except ClientError as e:
            logger.error(f"Error creating security group: {e}")
            raise

    def create_ec2_instance(self, security_group_id):
        user_data = '''#!/bin/bash
            yum update -y
            yum install -y httpd
            systemctl start httpd
            systemctl enable httpd
            echo "<h1>Hello from EC2</h1>" > /var/www/html/index.html
        '''

        try:
            response = self.ec2.run_instances(
                ImageId='ami-085ad6ae776d8f09c',  # Amazon Linux 2 AMI ID
                InstanceType='t2.micro',
                MinCount=1,
                MaxCount=1,
                UserData=user_data,
                SecurityGroupIds=[security_group_id],
                TagSpecifications=[
                    {
                        'ResourceType': 'instance',
                        'Tags': [
                            {
                                'Key': 'Name',
                                'Value': 'WebServer'
                            }
                        ]
                    }
                ]
            )
            instance_id = response['Instances'][0]['InstanceId']
            logger.info(f"Created EC2 instance: {instance_id}")
            
            # Wait for instance to be running
            waiter = self.ec2.get_waiter('instance_running')
            waiter.wait(InstanceIds=[instance_id])
            
            self.resources['instance_id'] = instance_id
            return instance_id
        except ClientError as e:
            logger.error(f"Error creating EC2 instance: {e}")
            raise

    def create_alb(self, security_group_id, instance_id):
        try:
            # Create target group
            response = self.elbv2.create_target_group(
                Name='WebServerTG',
                Protocol='HTTP',
                Port=80,
                VpcId='vpc-xxxxxx',  # Replace with your VPC ID
                HealthCheckProtocol='HTTP',
                HealthCheckPort='80',
                HealthCheckPath='/',
                TargetType='instance'
            )
            target_group_arn = response['TargetGroups'][0]['TargetGroupArn']
            self.resources['target_group_arn'] = target_group_arn
            
            # Register EC2 instance with target group
            self.elbv2.register_targets(
                TargetGroupArn=target_group_arn,
                Targets=[{'Id': instance_id}]
            )
            
            # Create ALB
            response = self.elbv2.create_load_balancer(
                Name='WebServerALB',
                Subnets=['subnet-xxxxx', 'subnet-yyyyy'],  # Replace with your subnet IDs
                SecurityGroups=[security_group_id],
                Scheme='internet-facing',
                Tags=[{'Key': 'Name', 'Value': 'WebServerALB'}]
            )
            alb_arn = response['LoadBalancers'][0]['LoadBalancerArn']
            self.resources['alb_arn'] = alb_arn
            
            # Create listener
            self.elbv2.create_listener(
                LoadBalancerArn=alb_arn,
                Protocol='HTTP',
                Port=80,
                DefaultActions=[
                    {
                        'Type': 'forward',
                        'TargetGroupArn': target_group_arn
                    }
                ]
            )
            
            logger.info(f"Created ALB: {alb_arn}")
            return alb_arn
        except ClientError as e:
            logger.error(f"Error creating ALB: {e}")
            raise

    def setup_alb_logging(self, alb_arn, bucket_name):
        try:
            # Create S3 bucket for logs
            self.s3.create_bucket(Bucket=bucket_name)
            
            # Enable access logging
            self.elbv2.modify_load_balancer_attributes(
                LoadBalancerArn=alb_arn,
                Attributes=[
                    {
                        'Key': 'access_logs.s3.enabled',
                        'Value': 'true'
                    },
                    {
                        'Key': 'access_logs.s3.bucket',
                        'Value': bucket_name
                    },
                    {
                        'Key': 'access_logs.s3.prefix',
                        'Value': 'alb-logs'
                    }
                ]
            )
            logger.info(f"Enabled ALB logging to bucket: {bucket_name}")
        except ClientError as e:
            logger.error(f"Error setting up ALB logging: {e}")
            raise

    def cleanup(self):
        """
        Cleanup all AWS resources created except S3 bucket
        """
        logger.info("Starting cleanup of AWS resources...")

        # Delete ALB
        if self.resources['alb_arn']:
            try:
                logger.info("Deleting Application Load Balancer...")
                self.elbv2.delete_load_balancer(LoadBalancerArn=self.resources['alb_arn'])
                # Wait for ALB to be deleted
                while True:
                    try:
                        self.elbv2.describe_load_balancers(LoadBalancerArns=[self.resources['alb_arn']])
                        time.sleep(10)
                    except self.elbv2.exceptions.LoadBalancerNotFoundException:
                        break
                logger.info("ALB deleted successfully")
            except ClientError as e:
                logger.error(f"Error deleting ALB: {e}")

        # Delete Target Group
        if self.resources['target_group_arn']:
            try:
                logger.info("Deleting Target Group...")
                self.elbv2.delete_target_group(TargetGroupArn=self.resources['target_group_arn'])
                logger.info("Target Group deleted successfully")
            except ClientError as e:
                logger.error(f"Error deleting Target Group: {e}")

        # Terminate EC2 instance
        if self.resources['instance_id']:
            try:
                logger.info("Terminating EC2 instance...")
                self.ec2.terminate_instances(InstanceIds=[self.resources['instance_id']])
                waiter = self.ec2.get_waiter('instance_terminated')
                waiter.wait(InstanceIds=[self.resources['instance_id']])
                logger.info("EC2 instance terminated successfully")
            except ClientError as e:
                logger.error(f"Error terminating EC2 instance: {e}")

        # Delete Security Group
        if self.resources['security_group_id']:
            try:
                logger.info("Deleting Security Group...")
                # Wait for instance to be fully terminated before deleting security group
                time.sleep(60)  # Additional wait to ensure all dependencies are cleared
                self.ec2.delete_security_group(GroupId=self.resources['security_group_id'])
                logger.info("Security Group deleted successfully")
            except ClientError as e:
                logger.error(f"Error deleting Security Group: {e}")

        logger.info("Cleanup completed")

def ping_alb(alb_dns, num_pings=10, interval=5):
    """
    Ping the ALB endpoint for specified number of times
    """
    for i in range(num_pings):
        try:
            response = requests.get(f"http://{alb_dns}")
            logger.info(f"Ping {i+1}: Status {response.status_code}")
        except requests.RequestException as e:
            logger.error(f"Ping {i+1} failed: {e}")
        time.sleep(interval)

def main():
    # Initialize setup
    setup = AWSInfraSetup()
    
    try:
        # Create infrastructure
        security_group_id = setup.create_security_group()
        instance_id = setup.create_ec2_instance(security_group_id)
        alb_arn = setup.create_alb(security_group_id, instance_id)
        
        # Setup logging
        bucket_name = 'aws-waf-logs-test-s3-bucket'  # Replace with your desired bucket name
        setup.setup_alb_logging(alb_arn, bucket_name)
        
        # Get ALB DNS name
        alb_info = setup.elbv2.describe_load_balancers(LoadBalancerArns=[alb_arn])
        alb_dns = alb_info['LoadBalancers'][0]['DNSName']
        
        # Start pinging ALB
        ping_alb(alb_dns, num_pings=20, interval=10)

    except Exception as e:
        logger.error(f"An error occurred: {e}")
    finally:
        # Cleanup all resources except S3 bucket
        user_input = input("Do you want to cleanup all resources? (yes/no): ")
        if user_input.lower() == 'yes':
            setup.cleanup()

if __name__ == "__main__":
    main()