# OpenFaaS Kubernetes info function

Create `view` cluster role binding:

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: openfaas-fn-view
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:serviceaccount:openfaas-fn:default
```

```bash
kubectl apply -f ./deploy/readonly-role.yaml
```

Create a secret named `fninfo-token`:

```bash
kubectl -n openfaas-fn create secret generic fninfo-token --from-literal=token=c1d116c6bfb
```

Deploy (requires OpenFaaS Operator):

```yaml
apiVersion: o6s.io/v1alpha1
kind: Function
metadata:
  name: fninfo
  namespace: openfaas-fn
spec:
  name: fninfo
  image: stefanprodan/fninfo:latest
  environment:
    secrets_path: "/var/openfaas"
  labels:
    release: "ga"
  secrets:
    - fninfo-token
  limits:
    cpu: "2000m"
    memory: "256Mi"
  requests:
    cpu: "100m"
    memory: "64Mi"
```

```bash
kubectl apply -f ./deploy/fninfo.yaml
```

Invoke function:

```bash
echo "test" | faas invoke fninfo | jq .

{
  "Hostname": "fninfo-6c7bd759cd-pr52v",
  "Namespaces": [
    {
      "Name": "kube-system",
      "Pods": 9,
      "Deployments": 6,
      "Services": 5
    },
    {
      "Name": "openfaas",
      "Pods": 7,
      "Deployments": 4,
      "Services": 3
    },
    {
      "Name": "openfaas-fn",
      "Pods": 4,
      "Deployments": 4,
      "Services": 4
    }
  ],
  "Environment": [
    "Http_X_Forwarded_For=10.56.0.160:35656",
    "Http_X_Envoy_Expected_Rq_Timeout_Ms=15000",
    "Http_Content_Type=text/plain",
    "Http_User_Agent=Go-http-client/2.0",
    "Http_X_Call_Id=baf4758d-cf8a-4f4d-af8e-f2a2c6e055dd",
    "Http_X_Request_Id=a1c2b3f8-85c8-4696-8f47-04487421bbfa",
    "Http_X_Start_Time=1527236562012143249",
    "Http_X_Envoy_Internal=true",
    "Http_X_Forwarded_Proto=https",
    "Http_Accept_Encoding=gzip",
    "Http_Method=POST",
    "Http_ContentLength=-1",
    "Http_Path=/function/fninfo"
  ],
  "Request": "test"
}
```

Add a random response delay between 1 to 5 seconds:

```bash
echo "delay" | faas invoke fninfo
```
