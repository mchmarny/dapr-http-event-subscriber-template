# Kubernetes Deployment

This section overviews deployment into Kubernetes. For illustration purposes, all commands in this document will be based on Microsoft Azure Kubernetes Service. Also, Dapr supports a wide array of state and pubsub backing services across multiple Cloud and on-prem deployments. This demo uses:

* [state backing service options](https://github.com/dapr/docs/tree/master/howto/setup-state-store)
* [pubsub backing service options](https://github.com/dapr/docs/tree/master/howto/setup-pub-sub-message-broker) 

## Prerequisite

* [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)

Also, to simplify all the scripts in this doc, set a few `az` CLI defaults:

```shell
az account set --subscription <id or name>
az configure --defaults location=<preferred location> group=<your resource group>
```

## Cluster (optional)

If you don't already have Kubernates cluster, you can create it on Azure with all the necessary add-ons for this demo using this command:

```shell
az aks create --name demo \
              --kubernetes-version 1.15.10 \
              --enable-managed-identity \
              --vm-set-type VirtualMachineScaleSets \
              --node-vm-size Standard_F4s_v2 \
              --enable-addons monitoring,http_application_routing \
              --generate-ssh-keys
```

## Install dapr

See instructions how to install dapr into your Kubernetes cluster [here](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md#installing-dapr-on-a-kubernetes-cluster)

## Component-backing services 

Assuming you have a Kubernates cluster and `kubectl` CLI configure to connect you are ready to setup the `dapr` components and their backing services:

> This demo installs into the `default` namespace in your cluster. When installing into a different namespace, make sure to append the `-n <your namespace name>` to all commands below (secret, component, and deployment) 

### State

To configure `dapr` state component in this demo I will use Azure Table Storage. To set it up, follow [these instructions](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal). Once finished, you will need to configure also the Kubernates secrets to hold the Azure Table Storage account information:

```shell
kubectl create secret generic demo-state-secret \
  --from-literal=account-name='' \
  --from-literal=account-key=''
```

Once the secret is configured, deploy the `dapr` state component:

```shell
kubectl apply -f deploy/component/state.yaml
```

### PubSub

To configure `dapr` pubsub component in this demo I will use Azure Service Bus. To set it up, follow [these instructions](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal). Once finished, you will need to configure the Kubernates secret to hold Azure Service Bus connection string information. 


```shell
kubectl create secret generic demo-pubsub-secret \
  --from-literal=connection-string=''
```

Once the secret is configured, deploy the `dapr` pubsub topic components:

```shell
kubectl apply -f deploy/component/pubsub.yaml
```


## Deploying Service

> Note, the following commands use `starter` as the name of the service which is defined in the `deploy/deployment.yaml` file. Change commands accordingly if you renamed your service. 

Once the necessary components are created, you can deploy the service with a single command:

```shell
kubectl apply -f deploy/deployment.yaml
```

You can check on the status of your deployment like this: 

```shell
kubectl get pods -l env=starter
```

The results should look similar to this (make sure each pod has READY status 2/2)

```shell
NAME                     READY   STATUS    RESTARTS   AGE
starter-89666d54b-hkd5t  2/2     Running   0          18s
```


## Observability 

In addition, you can also observe Dapr metrics, logs, and traces for this demo. 

### Metrics in Grafana dashboard 

> Instructions on how to setup Grafana for Dapr are [here](https://github.com/dapr/docs/blob/master/howto/setup-monitoring-tools/setup-prometheus-grafana.md)

Forward port

http://localhost:8080/

#### Logs in Kibana dashboard 

> Instructions on how to setup Kibana for Dapr are [here](https://github.com/dapr/docs/blob/master/howto/setup-monitoring-tools/setup-fluentd-es-kibana.md)

Forward port

http://localhost:5601/

#### Traces in Zipkin dashboard 

> Instructions on how to setup Zipkin for Dapr are [here](https://github.com/dapr/docs/blob/master/howto/diagnose-with-tracing/zipkin.md)

http://localhost:9411/zipkin/

> Note, if your Zipkin isn't deployed in the `default` namespace you will have to edit the `exporterAddress` in [deployment/tracing/zipkin.yaml](deployment/tracing/zipkin.yaml)

Then just restart all the deployments 

```shell
kubectl rollout restart deployment starter
```