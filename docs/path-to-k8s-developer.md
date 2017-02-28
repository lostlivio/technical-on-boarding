#Path to Becoming a Kubernetes Developer

Document intends to describe a successful path to becoming a Kubernetes (k8s) developer.
 
You can find helpful links to a lot of the information mentioned here at [k8s-cheatsheet](k8s-developer-cheatsheet.md).

Table of Contents
==================
  * [Code of Conduct](#code-of-conduct)
  * [The Kubernetes Community](#the-kubernetes-community)
    * [Google Groups](#google-groups)
    * [SIGS](#sigs)
    * [SLACK](#slack)  
    * [Zoom](#zoom)  
  * [Assumptions](#assumptions)
    * [Working Knowledge](#working-knowledge)
    * [Initial Setup](#initial-setup)
  * [Kubernetes Repos](#kubernetes-repos)
  * [Local Fork](#local-fork)
  * [Building Kubernetes](#building-kubernetes)
    * [Clean](#clean)
    * [Make](#make)
  * [Testing](testing)
  * [Local Cluster](#local-cluster)
    * [Development Process: Linux](development-process-linux)
    * [Development Process: OSX](development-process-osx)
  * [First PR](#first-pr)
  * [Submitting PR](#submitting-pr)
    * [Running Updates](#running-updates)

## Code of Conduct
Review the Cloud Native Computing Foundation (CNCF) document on the
[Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md) as it makes it 
clear what is expected from all individuals. In summary, play nice, be polite, and keep an atmosphere open for
collaboration among a viariety of viewpoints.  

## The Kubernetes Community
The Kubernetes community provides a variety of channels for communication. It can be confusing to understand which 
channel is appropriate for which situation. These are the channels we have found most useful:

### Google Groups
- [kubernetes-dev](https://groups.google.com/forum/#!forum/kubernetes-dev)
- [kubernetes-announce](https://groups.google.com/forum/#!forum/kubernetes-dev-announce)
- [kubernetes-pm](https://groups.google.com/forum/#!forum/kubernetes-pm)
- [kubernetes-users](https://groups.google.com/forum/#!forum/kubernetes-users)


### SIGS
The Kubernetes community is presently organized in a non-hierarchical set of Special Interest Groups (SIG's) focused 
on specific domains, eg: testing, scalability, storage, etc.  The full list of SIG's is available at 
[link](https://github.com/kubernetes/community/blob/master/README.md#special-interest-groups-sig)).

SIG's are expected to maintain a number of communication channels including:
- google groups / mailing lists for discussion
- slack channels for discussion
- github teams for notification purposes
- video meetings held at least once every 3 weeks

In general, if you have a question or feature in a specific area, you start with a SIG, and they will route you to the appropriate person.

### SLACK
Kubernetes has a publicly available slack instance at [link], with searchable archives at  [k8s slack channel](http://slack.k8s.io/).

Not all members of the community are immediately available here, but it is often the quickest way to get an answer.

## Assumptions
We expect you to be familiar with:

- git
  - Cloning
  - Forking
  - Rebasing
  - 2nd Auth
  - ssh key generation
  
- github
  - ssh key generation
  - 2nd Auth
  - Pull Requests
  - Review process

- golang
  - GOPATH
  - go get
  - go run
  - go test
  - go install
  - cross compilation

- docker
  - pull
  - run
  - exec
  - info
  - images
  

  
### Initial Setup

- git; Optionally with [git completion](https://github.com/git/git/blob/master/contrib/completion/git-completion.bash)
- docker, latest; Optionally [docker version manager](https://github.com/getcarina/dvm)
  - [linux box](https://docs.docker.com/engine/installation/) (do not use sudo-apt-get)
  - darwin: `brew install docker`
  - once installed, make sure you give docker at least 4GB (recommend 6GB) of ram and at least 4cpus
- Virtual-Box; Optionally if you are on osx

## Kubernetes Repos
You can find all github Kubernetes project repos to work from here

[Kubernetes Repos](https://github.com/kubernetes)
The Kubernetes project is composed of a variety of repos spread across two GitHub organizations:

- [kubernetes](https://github.com/kubernetes/kubernetes) - K8s codebase, main repo to contribute to.
- [kubernetes.gitHub.io](https://github.com/kubernetes.github.io) - K8s online website documents and guides
- [community](https://github.com/community) - K8s community site, including all SIG details.
- [kubernetes incubator](https://github.com/kubernetes-incubator)


## Local Fork
To contribute, the project must be forked to your own repo. Read on Github how the process works, it is fairly simple 
to do, but ask for help if you need to. Additional info can be found here: [kubernetes fork](https://github.com/kubernetes/kubernetes/blob/master/docs/devel/development.md).

Once you create a fork on github, do the following locally before making a clone:

- make the following directory:

  `mkdir  ${GOPATH}/go/src/k8s.io`

- then clone your fork into the k8s.io 

  ```
  cd ${GOPATH}/go/src/k8s.io /
  git clone https://githube.com/username/kubernetes.git
  ```

- set up upstream correctly

  ```
  cd kubernetes \
  git remote add upstream https://github.com/kubernetes/kubernetes.git
  ```

- Make sure you also install `godep`, `go-bindata`, `golint`, and `cfssl/cmd`:

  ```
  go get -u github.com/tools/godep \
  go get -u github.com/jteeuwen/go-bindata/go-bindata \
  go get -u github.com/golang/lint/golint \
  go get -u github.com/cloudflare/cfssl/cmd/cfssl \
  go get -u github.com/cloudflare/cfssl/cmd/...
  ```

- (optional) Follow the pre-commit instructions (ensures checks will pass and documentation is up-to-date before you can commit a change)

   ```
   cd .git/hooks/ \
   ln -s ../../hooks/pre-commit . \ 
   cd ../../
   ```

## Building Kubernetes
There are several flags to select from when running `make` to build k8s, regardless, all output 
should be in `_output/` folder. Here we list the most common ways to build.
 
- `make`
  - creates a build based on your current ARCH
- `make cross`
  - creates builds based on all available ARCHs
- `make quick-release` or `make release-skip-tests`
  - creates build based on current ARCH but skips tests
  
Failure to successfully build will cause results in a failure message that usually gives you a good idea why a build
was not possible. In addition building takes places on container images. It is suggested that you give these
containers at least 6gb of ram. You could give them at 11gb to run in parallel, but it has not been seen to give any
advantage in how long a build takes.

To change the memory limit on osx: 
- make sure you are using docker for mac.
- click on the whaley docker icon on the top right
- select `preferences`
- select the `Advanced` tab
- move slider to adjust the memory value you desire
- click on `Apply & Restart`

  
### Clean
Before building, make sure you run 

`make clean`

or if you need to

`rm -rf _output/ `

in order to remove any remnants of the previous build.

## Testing
For a good guide regarding testing you can follow
[E2E Testing: E2E, Local, Conformance, CI, etc...](https://github.com/kubernetes/community/blob/master/contributors/devel/e2e-tests.md)
which contains in good detail how to run tests.
For the most part you would run

`go run hack/e2e.go -v --build --up --test --down`

that builds, runs a cluster, runs testing, and then tears the tests down. 
If you wish to test your code changes more precisely, assuming unit tests,  you can go to the path containing your 
local test and run 

`go test <flags>`

and include any additional arguments/flags you need, details found [here](https://golang.org/pkg/testing/).

## Local Cluster
Tip: Minikube is a great tool to understand k8s, in particular running commands using `kubectl` to get services running
pods up, get an idea of what a namespace is, etc... However, what it is not is a dev tool for k8s. 


#### Development Process: Linux
Usually when you have made K8s changes you wish to test them, and sometimes you want to be able to build a release 
with your changes. The suggested way to make this happen is to

make quick-release`

and then bring up a local stack

`./hack/local-up-cluster.sh`

there may be other tools available to help you get started developing/building/testing but they are usually 
bloated and unnecessary. Minikube, for example runs on a pre-built version of K8s, and you will have to build it 
pointing to your local version of K8s.

#### Development Process: Darwin (OSX)

Current state:
- Kubelet will not run, so you can bring up a local cluster, but do not expect to talk to kubelet.
- All other components will run, no gurantee about how well they will run, nor how stable.
- E2E testing is possible if deploying a cluster to gke via kube-up.

##### Local Development Path

Current suggestion is to use a guide like

[Developer Guide, IBM](https://developer.ibm.com/opentech/2016/06/15/kubernetes-developer-guide-part-1/). 
The alternative to doing so would be to switch over to a linux laptop, or use a linux box to
use for developing, such as a nuc. Recommended you get at least an i5 with 16gb of ram and a linux distro such as
Ubuntu LTS.

Based on the current status of k8s development on osx, and as mentioned in the guide above, you will have to use a 
VM, such a virtual-box. Currently [Veertu-Desktop](https://veertu.com/forums/topic/veertu-desktop/) is the suggested 
route as it works with the native hypervisor and has been found to perform very well and stable. You are free
to use what you wish.

## First PR
Several approaches to selecting which repo you wish to start contributing to, and generally depends on the level of
confidence you have with the k8s codebase. Several suggestions:

- If you are starting with someone senior already in a project, consult them to get an idea how you can help and
contribute there.
- Go over known issues, see if there is something there that you can take up. There are some help-wanted label items 
but so far they seem to be very specific and oddly not very beginner friendly. 
- Go through the manuals, things you reviewed in the kubernetes website find errors, issues in grammar, etc... and fix
these. The community will only be better with better documentation. Not the most attractive place to contribute, but
it is a start, and an important one also.

## Submitting PR

If you have followed the [Forking](#forking) section carefully you will be prepared to create your first commmit and
have automatic checks run. That will save you time in the end of having tests fail upstream and then having to go back
and make changes or updates.

If you opted NOT to do so, then you can manually run the checks by running:

`make verify` 

which can take 30min+ on mac with 16gb ram and 2.5GHZ i7, either way you will get a list of errors and potentially ways 
to rectify the issue. The next section will cover a few of the most common.

Please revisit [Life of a Pull Request (very important)](https://github.com/kubernetes/kubernetes/blob/master/docs/devel/pull-requests.md)
It will save you a lot of time and pain by following that make verify etc.. process.


### Running Updates
As you make changes to code, there will on ocassion be a need to run updates to the build, documentations, and ownerships.

The list below is not exhaustive, but there are some of the updates you will, from time to time, be asked to perform 
(they are sorted in order they would need to be run at if asked, e.g., update-generated-swagger-docs.sh would be run before
you run update-swagger-spec.sh):
- `./hack/update-bazel.sh`
- `./hack/update-generated-swagger-docs.sh`
- `./hack/update-swagger-spec.sh`
- `./hack/update-federation-openapi-spec.sh`
- `./hack/update-openapi-spec.sh`

Lastly, if you do change GO code, it will be a good idea to run the following (the linter will ask you if you forget)

`gofmt -s -w pkg/kubectl/cmd/version.go`
