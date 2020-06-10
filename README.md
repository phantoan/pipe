# PipeCD

Continuous Delivery for Declarative Kubernetes Application and Infrastructure

## Proposal

https://docs.google.com/document/d/1Z3NqnsxgraD9f55F0TK6e4oLV4296Hd7xUb5GxwdaJQ

## Status

This project is under **PROTOTYPE** development phase.

## Overview

This project aims to explore and develop a unified delivery infrastructure for CA projects.

**Powerful**
- Unifed Deployment System: kubernetes (plain-yaml, helm, kustomize), terraform, lambda, cloudrun...
- Progressive Deployment Strategies: canary, bluegreen, rolling update
- Automated Analysis by Metrics, Log, Smoke Test
- Automated Rollback
- Automated Configuration Drift Detection
- Insights shows Delivery Perfomance

**Easy to Use**
- Operations by Pull Request: scale, rollout, rollback by PR
- Realtime Visualization of application state
- Deployment Pipeline to see what is happenning
- Intuitive UI

**Easy to Operate**
- Just 2 components: `piped` and `control-plane`
- Piped can be run on kubernetes, vm or even local machine
- Easy to operate multi-tenancy, multi-cluster
- Security: your credentials are not exposed outside of your cluster

## License

Apache License 2.0, see [LICENSE](https://github.com/kapetaniosci/pipe/blob/master/README.md).
