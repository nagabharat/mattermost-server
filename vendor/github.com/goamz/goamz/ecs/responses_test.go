package ecs

var CreateClusterResponse = `
<CreateClusterResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
	<ResponseMetadata>
		<RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
	</ResponseMetadata>
   <CreateClusterResult>
    <cluster>
      <clusterArn>arn:aws:ecs:region:aws_account_id:cluster/default</clusterArn>
      <clusterName>default</clusterName>
      <status>ACTIVE</status>
    </cluster>
  </CreateClusterResult>
</CreateClusterResponse> 
`
var DeregisterContainerInstanceResponse = `
<DeregisterContainerInstanceResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
   <DeregisterContainerInstanceResult>
    <containerInstance>
     <agentConnected>False</agentConnected>
     <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID</containerInstanceArn>
     <ec2InstanceId>instance_id</ec2InstanceId>
     <status>INACTIVE</status>
     <registeredResources>
        <member>
          <integerValue>2048</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>CPU</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>3955</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>MEMORY</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>0</integerValue>
          <longValue>0</longValue>
          <type>STRINGSET</type>
          <stringSetValue>
            <member>2376</member>
            <member>22</member>
            <member>51678</member>
            <member>2375</member>
          </stringSetValue>
          <name>PORTS</name>
          <doubleValue>0.0</doubleValue>
        </member>
     </registeredResources>
     <remainingResources>
        <member>
          <integerValue>2048</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>CPU</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>3955</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>MEMORY</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>0</integerValue>
          <longValue>0</longValue>
          <type>STRINGSET</type>
          <stringSetValue>
            <member>2376</member>
            <member>22</member>
            <member>51678</member>
            <member>2375</member>
          </stringSetValue>
          <name>PORTS</name>
          <doubleValue>0.0</doubleValue>
        </member>
     </remainingResources>
    </containerInstance>
  </DeregisterContainerInstanceResult>
</DeregisterContainerInstanceResponse> 
`

var DeregisterTaskDefinitionResponse = `
<DeregisterTaskDefinitionResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <DeregisterTaskDefinitionResult>
    <taskDefinition>
      <revision>2</revision>
      <family>sleep360</family>
      <containerDefinitions>
        <member>
          <portMappings/>
          <essential>true</essential>
          <environment>
            <member>
              <name>envVar</name>
              <value>foo</value>
            </member>
          </environment>
          <entryPoint>
            <member>/bin/sh</member>
          </entryPoint>
          <name>sleep</name>
          <command>
            <member>sleep</member>
            <member>360</member>
          </command>
          <cpu>10</cpu>
          <image>busybox</image>
          <memory>10</memory>
        </member>
      </containerDefinitions>
      <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
    </taskDefinition>
  </DeregisterTaskDefinitionResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</DeregisterTaskDefinitionResponse>
`

var DescribeClustersResponse = `
<DescribeClustersResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <DescribeClustersResult>
    <failures/>
    <clusters>
      <member>
        <clusterName>test</clusterName>
        <clusterArn>arn:aws:ecs:us-east-1:aws_account_id:cluster/test</clusterArn>
        <status>ACTIVE</status>
      </member>
      <member>
        <clusterName>default</clusterName>
        <clusterArn>arn:aws:ecs:us-east-1:aws_account_id:cluster/default</clusterArn>
        <status>ACTIVE</status>
      </member>
    </clusters>
  </DescribeClustersResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeClustersResponse>
`

var DescribeContainerInstancesResponse = `
<DescribeContainerInstancesResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">  
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
   <DescribeContainerInstancesResult>
    <failures/>
    <containerInstances>
      <member>
       <agentConnected>true</agentConnected>
       <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID</containerInstanceArn>
       <ec2InstanceId>instance_id</ec2InstanceId>
       <status>ACTIVE</status>
       <registeredResources>
          <member>
            <integerValue>2048</integerValue>
            <longValue>0</longValue>
            <type>INTEGER</type>
            <name>CPU</name>
            <doubleValue>0.0</doubleValue>
          </member>
          <member>
            <integerValue>3955</integerValue>
            <longValue>0</longValue>
            <type>INTEGER</type>
            <name>MEMORY</name>
            <doubleValue>0.0</doubleValue>
          </member>
          <member>
            <integerValue>0</integerValue>
            <longValue>0</longValue>
            <type>STRINGSET</type>
            <stringSetValue>
              <member>2376</member>
              <member>22</member>
              <member>51678</member>
              <member>2375</member>
            </stringSetValue>
            <name>PORTS</name>
            <doubleValue>0.0</doubleValue>
          </member>
       </registeredResources>
       <remainingResources>
          <member>
            <integerValue>2048</integerValue>
            <longValue>0</longValue>
            <type>INTEGER</type>
            <name>CPU</name>
            <doubleValue>0.0</doubleValue>
          </member>
          <member>
            <integerValue>3955</integerValue>
            <longValue>0</longValue>
            <type>INTEGER</type>
            <name>MEMORY</name>
            <doubleValue>0.0</doubleValue>
          </member>
          <member>
            <integerValue>0</integerValue>
            <longValue>0</longValue>
            <type>STRINGSET</type>
            <stringSetValue>
              <member>2376</member>
              <member>22</member>
              <member>51678</member>
              <member>2375</member>
            </stringSetValue>
            <name>PORTS</name>
            <doubleValue>0.0</doubleValue>
          </member>
       </remainingResources>
      </member>
    </containerInstances>
  </DescribeContainerInstancesResult>
</DescribeContainerInstancesResponse> 
`

var DescribeTaskDefinitionResponse = `
<DescribeTaskDefinitionResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <DescribeTaskDefinitionResult>
    <taskDefinition>
      <revision>2</revision>
      <family>sleep360</family>
      <containerDefinitions>
        <member>
          <portMappings/>
          <essential>true</essential>
          <environment>
            <member>
              <name>envVar</name>
              <value>foo</value>
            </member>
          </environment>
          <entryPoint>
            <member>/bin/sh</member>
          </entryPoint>
          <name>sleep</name>
          <command>
            <member>sleep</member>
            <member>360</member>
          </command>
          <cpu>10</cpu>
          <image>busybox</image>
          <memory>10</memory>
        </member>
      </containerDefinitions>
      <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
    </taskDefinition>
  </DescribeTaskDefinitionResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeTaskDefinitionResponse>
`

var DescribeTasksResponse = ` 
<DescribeTasksResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <DescribeTasksResult>
    <failures/>
    <tasks>
      <member>
        <containers>
          <member>
            <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
            <name>sleep</name>
            <containerArn>arn:aws:ecs:us-east-1:aws_account_id:container/UUID</containerArn>
            <networkBindings/>
            <lastStatus>RUNNING</lastStatus>
          </member>
        </containers>
        <overrides>
          <containerOverrides>
            <member>
              <name>sleep</name>
            </member>
          </containerOverrides>
        </overrides>
        <desiredStatus>RUNNING</desiredStatus>
        <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
        <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID</containerInstanceArn>
        <lastStatus>RUNNING</lastStatus>
        <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
      </member>
    </tasks>
  </DescribeTasksResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeTasksResponse>
`

var DiscoverPollEndpointResponse = `
<DiscoverPollEndpointResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <DiscoverPollEndpointResult>
    <endpoint>https://ecs-x-1.us-east-1.amazonaws.com/</endpoint>
  </DiscoverPollEndpointResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</DiscoverPollEndpointResponse>
`

var ListClustersResponse = `
<ListClustersResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ListClustersResult>
    <clusterArns>
      <member>arn:aws:ecs:us-east-1:aws_account_id:cluster/default</member>
      <member>arn:aws:ecs:us-east-1:aws_account_id:cluster/test</member>
    </clusterArns>
    <nextToken>token_UUID</nextToken>
  </ListClustersResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</ListClustersResponse>
`

var ListContainerInstancesResponse = `
<ListContainerInstancesResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ListContainerInstancesResult>
    <containerInstanceArns>
      <member>arn:aws:ecs:us-east-1:aws_account_id:container-instance/uuid-1</member>
      <member>arn:aws:ecs:us-east-1:aws_account_id:container-instance/uuid-2</member>
    </containerInstanceArns>
     <nextToken>token_UUID</nextToken>
  </ListContainerInstancesResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</ListContainerInstancesResponse>
`

var ListTaskDefinitionsResponse = `
<ListTaskDefinitionsResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ListTaskDefinitionsResult>
    <taskDefinitionArns>
      <member>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:1</member>
      <member>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</member>
    </taskDefinitionArns>
     <nextToken>token_UUID</nextToken>
  </ListTaskDefinitionsResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</ListTaskDefinitionsResponse>
`

var ListTasksResponse = `
<ListTasksResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ListTasksResult>
    <taskArns>
      <member>arn:aws:ecs:us-east-1:aws_account_id:task/uuid_1</member>
      <member>arn:aws:ecs:us-east-1:aws_account_id:task/uuid_2</member>
    </taskArns>
    <nextToken>token_UUID</nextToken>
  </ListTasksResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</ListTasksResponse>
`

var RegisterContainerInstanceResponse = `
<RegisterContainerInstanceResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
   <RegisterContainerInstanceResult>
    <containerInstance>
     <agentConnected>True</agentConnected>
     <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID</containerInstanceArn>
     <ec2InstanceId>instance_id</ec2InstanceId>
     <status>ACTIVE</status>
     <registeredResources>
        <member>
          <integerValue>2048</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>CPU</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>3955</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>MEMORY</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>0</integerValue>
          <longValue>0</longValue>
          <type>STRINGSET</type>
          <stringSetValue>
            <member>2376</member>
            <member>22</member>
            <member>51678</member>
            <member>2375</member>
          </stringSetValue>
          <name>PORTS</name>
          <doubleValue>0.0</doubleValue>
        </member>
     </registeredResources>
     <remainingResources>
        <member>
          <integerValue>2048</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>CPU</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>3955</integerValue>
          <longValue>0</longValue>
          <type>INTEGER</type>
          <name>MEMORY</name>
          <doubleValue>0.0</doubleValue>
        </member>
        <member>
          <integerValue>0</integerValue>
          <longValue>0</longValue>
          <type>STRINGSET</type>
          <stringSetValue>
            <member>2376</member>
            <member>22</member>
            <member>51678</member>
            <member>2375</member>
          </stringSetValue>
          <name>PORTS</name>
          <doubleValue>0.0</doubleValue>
        </member>
     </remainingResources>
    </containerInstance>
  </RegisterContainerInstanceResult>
</RegisterContainerInstanceResponse> 
`

var RegisterTaskDefinitionResponse = ` 
<RegisterTaskDefinitionResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <RegisterTaskDefinitionResult>
    <taskDefinition>
      <revision>2</revision>
      <family>sleep360</family>
      <containerDefinitions>
        <member>
          <portMappings/>
          <essential>true</essential>
          <environment>
            <member>
              <name>envVar</name>
              <value>foo</value>
            </member>
          </environment>
          <entryPoint>
            <member>/bin/sh</member>
          </entryPoint>
          <name>sleep</name>
          <command>
            <member>sleep</member>
            <member>360</member>
          </command>
          <cpu>10</cpu>
          <image>busybox</image>
          <memory>10</memory>
          <mountPoints>
            <member>
              <containerPath>/tmp/myfile</containerPath>
              <readOnly>false</readOnly>
              <sourceVolume>/srv/myfile</sourceVolume>
            </member>
            <member>
              <containerPath>/tmp/myfile2</containerPath>
              <readOnly>true</readOnly>
              <sourceVolume>/srv/myfile2</sourceVolume>
            </member>
          </mountPoints>
          <volumesFrom>
           <member>
              <readOnly>true</readOnly>
              <sourceContainer>foo</sourceContainer>
            </member>
          </volumesFrom>
        </member>
      </containerDefinitions>
      <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
      <volumes>
        <member>
          <name>/srv/myfile</name>
          <host>
            <sourcePath>/srv/myfile</sourcePath>
          </host>
        </member>
      </volumes>
    </taskDefinition>
  </RegisterTaskDefinitionResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</RegisterTaskDefinitionResponse>
`

var RunTaskResponse = ` 
<RunTaskResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <RunTaskResult>
    <failures/>
    <tasks>
      <member>
        <containers>
          <member>
            <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
            <name>sleep</name>
            <containerArn>arn:aws:ecs:us-east-1:aws_account_id:container/UUID</containerArn>
            <networkBindings/>
            <lastStatus>RUNNING</lastStatus>
          </member>
        </containers>
        <overrides>
          <containerOverrides>
            <member>
              <name>sleep</name>
            </member>
          </containerOverrides>
        </overrides>
        <desiredStatus>RUNNING</desiredStatus>
        <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
        <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID</containerInstanceArn>
        <lastStatus>PENDING</lastStatus>
        <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
      </member>
    </tasks>
  </RunTaskResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</RunTaskResponse>
`

var StartTaskResponse = ` 
<StartTaskResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <StartTaskResult>
    <failures/>
    <tasks>
      <member>
        <containers>
          <member>
            <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
            <name>sleep</name>
            <containerArn>arn:aws:ecs:us-east-1:aws_account_id:container/UUID</containerArn>
            <networkBindings/>
            <lastStatus>RUNNING</lastStatus>
          </member>
        </containers>
        <overrides>
          <containerOverrides>
            <member>
              <name>sleep</name>
            </member>
          </containerOverrides>
        </overrides>
        <desiredStatus>RUNNING</desiredStatus>
        <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
        <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID</containerInstanceArn>
        <lastStatus>PENDING</lastStatus>
        <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
      </member>
    </tasks>
  </StartTaskResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</StartTaskResponse>
`

var StopTaskResponse = ` 
<StopTaskResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <StopTaskResult>
    <task>
      <containers>
        <member>
          <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
          <name>sleep</name>
          <containerArn>arn:aws:ecs:us-east-1:aws_account_id:container/UUID</containerArn>
          <networkBindings/>
          <lastStatus>RUNNING</lastStatus>
        </member>
      </containers>
      <overrides>
        <containerOverrides>
          <member>
            <name>sleep</name>
          </member>
        </containerOverrides>
      </overrides>
      <desiredStatus>STOPPED</desiredStatus>
      <taskArn>arn:aws:ecs:us-east-1:aws_account_id:task/UUID</taskArn>
      <containerInstanceArn>arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID</containerInstanceArn>
      <lastStatus>RUNNING</lastStatus>
      <taskDefinitionArn>arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2</taskDefinitionArn>
    </task>
  </StopTaskResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</StopTaskResponse>
`

var SubmitContainerStateChangeResponse = `
<SubmitContainerStateChangeResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <SubmitContainerStateChangeResult>
   <acknowledgment>ACK</acknowledgment>
  </SubmitContainerStateChangeResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</SubmitContainerStateChangeResponse>
`

var SubmitTaskStateChangeResponse = `
<SubmitTaskStateChangeResponse xmlns="http://ecs.amazonaws.com/doc/2014-11-13/">
  <SubmitTaskStateChangeResult>
   <acknowledgment>ACK</acknowledgment>
  </SubmitTaskStateChangeResult>
  <ResponseMetadata>
    <RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
  </ResponseMetadata>
</SubmitTaskStateChangeResponse>
`
