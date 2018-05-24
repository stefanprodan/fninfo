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

Deploy (requires OpenFaaS Operator):

```yaml
apiVersion: o6s.io/v1alpha1
kind: Function
metadata:
  name: fninfo
  namespace: openfaas-fn
spec:
  name: gofast
  image: stefanprodan/fninfo
  limits:
    cpu: "1000m"
    memory: "128Mi"
  requests:
    cpu: "10m"
    memory: "64Mi"
```

```bash
kubectl apply -f ./deploy/fninfo.yaml
```

Invoke function:

```bash
echo "test" | faas invoke fninfo | jq .

{
  "Hostname": "fninfo-6498fbd77c-x7vqh",
  "Pods": [
    "certinfo-7874f9c8f5-5gjrr",
    "fninfo-6498fbd77c-x7vqh",
    "nodeinfo-b8fdcd9d4-rgdzq",
    "sentimentanalysis-76b4968d64-58zjn"
  ],
  "Services": [
    "certinfo",
    "fninfo",
    "nodeinfo",
    "sentimentanalysis"
  ],
  "Deployments": [
    "certinfo",
    "fninfo",
    "nodeinfo",
    "sentimentanalysis"
  ],
  "Environment": [
    "fprocess=./handler",
    "HOME=/home/app",
    "Http_User_Agent=Go-http-client/2.0",
    "Http_Authorization=Basic YWRtaW46YWRtaW4=",
    "Http_X_Forwarded_Proto=https",
    "Http_X_Request_Id=0937af80-869e-468f-867b-c950bb0ddcd3",
    "Http_X_Forwarded_For=10.56.0.160:39364",
    "Http_Content_Type=text/plain",
    "Http_X_Call_Id=795b7945-af8b-4a0e-943b-b296d19dbc69",
    "Http_X_Envoy_Internal=true",
    "Http_X_Envoy_Expected_Rq_Timeout_Ms=15000",
    "Http_X_Start_Time=1527200696113845058",
    "Http_Accept_Encoding=gzip",
    "Http_Method=POST",
    "Http_ContentLength=-1",
    "Http_Path=/function/fninfo"
  ],
  "Request": "test\n"
}
```
