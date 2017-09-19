podTemplate(label: "technical-on-boarding", containers: [
    containerTemplate(name: 'jnlp', image: "quay.io/samsung_cnct/custom-jnlp:0.1", args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'docker', image: 'docker', command: 'cat', ttyEnabled: true)
  ], volumes: [
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
    hostPathVolume(hostPath: '/var/lib/docker/scratch', mountPath: '/mnt/scratch'),
    secretVolume(mountPath: '/home/jenkins/.docker/', secretName: 'samsung-cnct-quay-robot-dockercfg')
  ]) {
    node("technical-on-boarding") {
      customContainer('docker') {
        stage('Checkout') {
          checkout scm
          // retrieve the URI used for checking out the source
          // this assumes one branch with one uri
          git_uri = scm.getRepositories()[0].getURIs()[0].toString()
          git_branch = scm.getBranches()[0].toString()
        }
        stage('Setup') {
          kubesh 'apk add --update --no-cache build-base git'
        }
        stage('Build') {
          kubesh "make -f Makefile.docker IMAGE_DEVL_TAG=${env.JOB_BASE_NAME}.${env.BUILD_ID} build"
        }
        stage('Test') {
          kubesh "make -f Makefile.docker IMAGE_DEVL_TAG=${env.JOB_BASE_NAME}.${env.BUILD_ID} test"
        }
        stage('Publish') {
          kubesh "make -f Makefile.docker GITHUB_BRANCH=${git_branch} GITHUB_URI=${git_uri} IMAGE_DEVL_TAG=${env.JOB_BASE_NAME}.${env.BUILD_ID} IMAGE_PROD_TAG=${env.RELEASE_VERSION} publish"
        }
      }
    }
  }

def kubesh(command) {
  if (env.CONTAINER_NAME) {
    if ((command instanceof String) || (command instanceof GString)) {
      command = kubectl(command)
    }

    if (command instanceof LinkedHashMap) {
      command["script"] = kubectl(command["script"])
    }
  }
  sh(command)
}

def kubectl(command) {
  "kubectl exec -i ${env.HOSTNAME} -c ${env.CONTAINER_NAME} -- /bin/sh -c 'cd ${env.WORKSPACE} && ${command}'"
}

def customContainer(String name, Closure body) {
  withEnv(["CONTAINER_NAME=$name"]) {
    body()
  }
}

// vi: ft=groovy
