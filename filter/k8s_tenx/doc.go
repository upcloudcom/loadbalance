//create: 2018/01/18 17:07:44 change: 2018/01/24 11:06:35 upcloudcom@foxmail.com
/*
$ cat /etc/tenx/extention.conf
{
    "groups": [
        {
            "address": "10.39.0.116",
            "domain": "otest.enncloud.cn",
            "id": "group-darso",
            "is_default": true,
            "name": " \u5185\u7f51\u7f51\u7edc",
            "type": "private"
        },
        {
            "address": "123.58.243.23",
            "domain": "itest.enncloud.cn",
            "id": "group-default",
            "is_default": true,
            "name": "\u516c\u7f51\u7f51\u7edc",
            "type": "public"
        }
    ],
    "nodes": [
        {
            "address": "10.39.0.116",
            "group": "group-darso",
            "host": "test-slave-116"
        },
        {
            "address": "10.39.0.117",
            "group": "group-default",
            "host": "test-slave-117"
        }
    ]
}

$ cat domain.json |python -m json.tool
{
    "domain": "",
    "externalip": "10.39.0.113"
}

annotation:
	binding_domains: www.greatgas.cn,pay.greatgas.cn
	binding_port: "80"
	system/lbgroup: group-thudj
	tenxcloud.com/https: "true"
	tenxcloud.com/schemaPortname: greatgas-nginx-1/HTTP

crt /etc/sslkeys/certs/wangxuz.tenxsep.greatgas-nginx

annotations:
	system/lbgroup: group-default
	tenxcloud.com/schemaPortname: sshproxy-1/TCP/51582,sshproxy-2/HTTP,sshproxy-3/TCP/41678

ports:
- name: sshproxy-1
port: 22
protocol: TCP
targetPort: 22
- name: sshproxy-2
port: 23
protocol: TCP
targetPort: 23
- name: sshproxy-3
port: 24
protocol: TCP
targetPort: 24

annotations:
	binding_domains: www.xx.com
	binding_port: "22"
	system/lbgroup: group-gulre
	tenxcloud.com/schemaPortname: sshproxy-pub-0/TCP/23558

- name: sshproxy-pub-0
  port: 22
	  protocol: TCP
		  targetPort: 22

POD:

apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    annotations:
      kubernetes.io/created-by: |
        {"kind":"SerializedReference","apiVersion":"v1","reference":{"kind":"ReplicaSet","namespace":"lijiaob","name":"webshell-1259653225","uid":"7b68f816-fce1-11e7-9d36-5254b24cbf5e","apiVersion":"extensions","resourceVersion":"413982376"}}
    creationTimestamp: 2018-01-19T06:25:02Z
    generateName: webshell-1259653225-
    labels:
      ClusterID: CID-516874818ed4
      UserID: "8"
      name: webshell
      pod-template-hash: "1259653225"
      tenxcloud.com/appName: webshell
      tenxcloud.com/svcName: webshell
    name: webshell-1259653225-crmss
    namespace: lijiaob
    ownerReferences:
    - apiVersion: extensions/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: ReplicaSet
      name: webshell-1259653225
      uid: 7b68f816-fce1-11e7-9d36-5254b24cbf5e
    resourceVersion: "413982454"
    selfLink: /api/v1/namespaces/lijiaob/pods/webshell-1259653225-crmss
    uid: 7b6a7135-fce1-11e7-9d36-5254b24cbf5e
  spec:
    affinity:
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
          - matchExpressions:
            - key: system/build-node
              operator: NotIn
              values:
              - "true"
    containers:
    - env:
      - name: PATH
        value: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
      - name: APP_NAME
        value: webshell
      - name: SERVICE_NAME
        value: webshell
      - name: CLUSTER_ID
        value: CID-516874818ed4
      - name: USER_ID
        value: "8"
      image: reg.enncloud.cn/lijiaob/webshell:master
      imagePullPolicy: Always
      name: webshell
      ports:
      - containerPort: 80
        protocol: TCP
      resources:
        limits:
          cpu: 500m
          memory: 100Mi
        requests:
          cpu: 500m
          memory: 100Mi
      terminationMessagePath: /dev/termination-log
      terminationMessagePolicy: File
      volumeMounts:
      - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
        name: default-token-a3xc3
        readOnly: true
    dnsPolicy: ClusterFirst
    imagePullSecrets:
    - name: registrysecret
    nodeName: slave-186
    restartPolicy: Always
    schedulerName: default-scheduler
    securityContext: {}
    serviceAccount: default
    serviceAccountName: default
    terminationGracePeriodSeconds: 30
    volumes:
    - name: default-token-a3xc3
      secret:
        defaultMode: 420
        secretName: default-token-a3xc3
  status:
    conditions:
    - lastProbeTime: null
      lastTransitionTime: 2018-01-19T06:25:02Z
      status: "True"
      type: Initialized
    - lastProbeTime: null
      lastTransitionTime: 2018-01-19T06:25:05Z
      status: "True"
      type: Ready
    - lastProbeTime: null
      lastTransitionTime: 2018-01-19T06:25:02Z
      status: "True"
      type: PodScheduled
    containerStatuses:
    - containerID: docker://b2200465bdceaf11afcbc15348b522103a5aa6eecb5a42164eadc803f29c894f
      image: reg.enncloud.cn/lijiaob/webshell:master
      imageID: docker-pullable://reg.enncloud.cn/lijiaob/webshell@sha256:a062b89b468b4864d1ea959c04f2b1df0815c0dad35180b2c6443f8963b3573a
      lastState: {}
      name: webshell
      ready: true
      restartCount: 0
      state:
        running:
          startedAt: 2018-01-19T06:25:05Z
    hostIP: 10.39.1.186
    phase: Running
    podIP: 192.168.55.105
    qosClass: Guaranteed
    startTime: 2018-01-19T06:25:02Z
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""

*/
package k8s_tenx
