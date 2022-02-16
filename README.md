



operator-sdk init --domain data.konghq.com --repo github.com/MadeByMakers/kong-operator-for-k8s


operator-sdk create api --group data --version v1alpha1 --kind Service --resource --controller
operator-sdk create api --group data --version v1alpha1 --kind Route --resource --controller
operator-sdk create api --group data --version v1alpha1 --kind Plugin --resource --controller

operator-sdk create api --group declarative --version v1alpha1 --kind RestService --resource --controller
operator-sdk create api --group declarative --version v1alpha1 --kind SoapService --resource --controller
operator-sdk create api --group declarative --version v1alpha1 --kind TcpProxy --resource --controller
operator-sdk create api --group declarative --version v1alpha1 --kind UdpProxy --resource --controller



operator-sdk create webhook --group data --version v1alpha1 --kind Service --defaulting --programmatic-validation




REFERENCIAS:

https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/

https://cloud.redhat.com/blog/kubernetes-operators-best-practices#error

https://github.com/redhat-cop/operator-utils
