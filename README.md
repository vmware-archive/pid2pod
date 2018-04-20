# pid2pod

This is a small demo/experiment that shows how Linux process IDs (PIDs) can be mapped to Kubernetes pod metadata.

It works by looking up the target process's [cgroup] metadata in `/proc/$PID/cgroup`.
This metadata contains the names of each cgroup assigned to the process.
In the case of Docker containers created using the `docker` CLI or created by the kubelet, these cgroup names contain the Docker container ID.
We can map this container ID to a Kubernetes pod by doing a lookup against the local [kubelet API][kubelet-api].

Thanks to the [SPIRE] implementation of [SPIFFE], which inspired this method ([see `k8s.go`][spiffe-impl]).

## Example Use Cases

- You want to enrich [audit][go-audit] or other host-level events with pod-level metadata.

- You want to control access to a UNIX domain socket based on [the calling pod][uds-pid-lookup].

## Caveats

- Not well tested.
  Please use this as an example but don't run this in production.

- Only works for the Docker runtime.
  Other container runtimes may use a different scheme for cgroup names or may not even have workload PIDs visible from the host.

- Assumes the kubelet read-only API (`--read-only-port`) is running on `http://localhost:10255/`.
  We could add token/certificate authentication to move past this assumption.

- Has some inherent race conditions since a lookup from PID to pod metadata takes some time.
  The PID may correspond to a different process by the time the data is used.

  This may be avoided in some scenarios.
  For example, if you have a server holding an open UNIX domain socket you know the client has had the same PID for the whole connection.

## Usage

The demo lists all processes on the system that it can map back to a pod.
It must be run in privileged mode, in the host PID namespace, and in the host network namespace.

```console
$ docker run --rm --privileged --pid=host --net=host pid2pod:latest
PID 4549 (kube-proxy): &pid2pod.ID{Namespace:"kube-system", PodName:"kube-proxy-mdnzr", PodUID:"6344523c-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"controller-revision-hash":"1193416634", "k8s-app":"kube-proxy", "pod-template-generation":"1"}, ContainerID:"7d6ca8ec4775491b941aeb9d55fe52fe14236575c26b395f36d0afa8717bdaae", ContainerName:"kube-proxy"}
PID 4992 (dashboard): &pid2pod.ID{Namespace:"kube-system", PodName:"kubernetes-dashboard-5498ccf677-xk55s", PodUID:"63fe5f9d-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"addonmanager.kubernetes.io/mode":"Reconcile", "app":"kubernetes-dashboard", "pod-template-hash":"1054779233", "version":"v1.8.1"}, ContainerID:"3d23b4ffce2f10bbcc10ffee5acc5bdb2446e48c9253bb51405d65ece5bf0dd2", ContainerName:"kubernetes-dashboard"}
PID 5079 (kube-dns): &pid2pod.ID{Namespace:"kube-system", PodName:"kube-dns-86f4d74b45-qfs6f", PodUID:"634ab7a3-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"k8s-app":"kube-dns", "pod-template-hash":"4290830601"}, ContainerID:"8ae06d9c2b923cad2a3fd8f1c5ebaad18483d6af369f968209cc940e7f501b1a", ContainerName:"kubedns"}
PID 5173 (storage-provisi): &pid2pod.ID{Namespace:"kube-system", PodName:"storage-provisioner", PodUID:"640b24b0-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"addonmanager.kubernetes.io/mode":"Reconcile", "integration-test":"storage-provisioner"}, ContainerID:"f22bff5ee043e33d952df17f9cf6d79d8b8de3b82d565566b0bef49d2a4fa809", ContainerName:"storage-provisioner"}
PID 5291 (dnsmasq-nanny): &pid2pod.ID{Namespace:"kube-system", PodName:"kube-dns-86f4d74b45-qfs6f", PodUID:"634ab7a3-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"k8s-app":"kube-dns", "pod-template-hash":"4290830601"}, ContainerID:"07567b209503648459d7e5bd5ef1e943355dc479ca68c629a30cc139810a95ae", ContainerName:"dnsmasq"}
PID 5339 (dnsmasq): &pid2pod.ID{Namespace:"kube-system", PodName:"kube-dns-86f4d74b45-qfs6f", PodUID:"634ab7a3-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"k8s-app":"kube-dns", "pod-template-hash":"4290830601"}, ContainerID:"07567b209503648459d7e5bd5ef1e943355dc479ca68c629a30cc139810a95ae", ContainerName:"dnsmasq"}
PID 5378 (sidecar): &pid2pod.ID{Namespace:"kube-system", PodName:"kube-dns-86f4d74b45-qfs6f", PodUID:"634ab7a3-44b7-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"k8s-app":"kube-dns", "pod-template-hash":"4290830601"}, ContainerID:"fa1685e9e27b95556313641bd6b4e7f3db7893a8dcb6d85f5cc369d895be17db", ContainerName:"sidecar"}
PID 16962 (sh): &pid2pod.ID{Namespace:"default", PodName:"testcli-5b8d779dfd-vpwkw", PodUID:"a8c810f7-44d3-11e8-b7fb-b2d1d7d0d428", PodLabels:map[string]string{"pod-template-hash":"1648335898", "run":"testcli"}, ContainerID:"8fe33943cbef5cdc741b265cad4c061e6232ec6187b9e00e2c2e800df7fce4aa", ContainerName:"testcli"}
```

[SPIFFE]: https://spiffe.io/
[SPIRE]: https://github.com/spiffe/spire
[cgroup]: https://en.wikipedia.org/wiki/Cgroups
[kubelet-api]: https://stackoverflow.com/questions/35075195/is-there-api-documentation-for-kublet
[spiffe-impl]: https://github.com/spiffe/spire/blob/719c696036a99c0cd918ff4747d5fd7f15c782e4/pkg/agent/plugin/workloadattestor/k8s/k8s.go#L61
[go-audit]: https://github.com/slackhq/go-audit
[uds-pid-lookup]: https://stackoverflow.com/questions/8104904/identify-program-that-connects-to-a-unix-domain-socket