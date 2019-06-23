# Custom Controller for Copybird Custom Resource


## Generating Code for Custom Controller 
In case of any changes to types for the custom controller, use the following commands to regenerate client and deepcopy files

```
ROOT_PACKAGE="github.com/copybird/copybird/operator"
CUSTOM_RESOURCE_NAME="backup"
CUSTOM_RESOURCE_VERSION="v1"

go get -u k8s.io/code-generator/...
cd $GOPATH/src/k8s.io/code-generator

./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"

```

## Run localy 
First create custom resource definition in your cluster: 
```
kubectl apply -f operator/crd/crd.yaml
```

To run the controller:
``` 
go run main.go operator
```

And then in a separate shell, create custom resource:
```
kubectl create -f operator/example/backup-example.yaml
```

As output you get the following logs when creating, updating or deleting custom resource:
```
INFO[0000] Successfully constructed k8s client          
INFO[0000] Starting Foo controller                      
INFO[0000] Waiting for informer caches to sync          
INFO[0001] Starting workers                             
INFO[0001] Started workers               
```

