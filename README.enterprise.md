# Cilium Enterprise

## (**EXPERIMENTAL**) `main-ce`: Development branch for Cilium Enterprise

[`main-ce`](https://github.com/isovalent/cilium/tree/main-ce) is an experimental
development branch for Cilium Enterprise that is regularly synchronized with the
OSS main branch. Please see the following resources for more details:

- [CFP](https://docs.google.com/document/d/1zcAttwopuhA_BzpfjH0yjlxScijSv2ZqGkphSHpnrRE/edit)
- [GitHub Issue](https://github.com/isovalent/roadmap/issues/616)

We'll share the proposal with the engineering team after Cilium Enterprise v1.13
is released. In the meantime, we are currently experimenting with `main-ce`
branch using a small number of enterprise-only features:

- [Egress Gateway HA](https://github.com/isovalent/cilium/pull/850)
- [Buried Treasure](https://github.com/isovalent/cilium/pull/869)
- FQDN Ingress

Please post a message in [#prj-cee-oss-codebase Slack channel](https://isovalent.slack.com/archives/C05188W95GT)
if you are interested in participating in the experiment.

## Resolving conflicts in autogenerated files

There are some files that are autogenerated:

- [Documentation/helm-values.rst](Documentation/helm-values.rst)
- [install/kubernetes/cilium/README.md](install/kubernetes/cilium/README.md)
- [install/kubernetes/cilium/values.yaml](install/kubernetes/cilium/values.yaml)

If any of these files conflict with upstream, remove them and regenerate:

    rm -f Documentation/helm-values.rst install/kubernetes/cilium/README.md install/kubernetes/cilium/values.yaml
    make -C install/kubernetes
    make -C Documentation update-helm-values