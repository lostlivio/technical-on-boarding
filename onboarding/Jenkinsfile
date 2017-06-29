

podTemplate(label: 'OnboardingWorkflow', containers: [
    containerTemplate(name: 'jnlp', image: 'quay.io/samsung_cnct/custom-jnlp:0.1', args: '${computer.jnlpmac} ${computer.name}'),
    containerTemplate(name: 'technical-on-boarding', image: 'docker.io/library/golang:1.8', ttyEnabled: true, command: 'cat', alwaysPullImage: true, resourceRequestMemory: '128Mi', resourceLimitMemory: '1Gi')\
  ]){

    node('OnboardingWorkflow') {
     

        customContainer ('technical-on-boarding') {
            stage('checkout') {
                checkout scm

                // Ugly hack. The code must live in a path compatible with Golang's layout, and be resolvable within Go apps via its full form ("github.com/samsung-cnct/...").
                // There may be alternatives that allow it to run from local path, especially vendoring-in dependencies via glide.
                // Performing the checkout to an explicit directory DID NOT work because it was somehow pulling `master` instead.
                kubesh "mkdir -p go/src/github.com/samsung-cnct/technical-on-boarding/ && cp -r ./onboarding go/src/github.com/samsung-cnct/technical-on-boarding/onboarding"

            }

            withEnv(["GOPATH=${WORKSPACE}/go/", "CHECKOUT_PATH=./go/src/github.com/samsung-cnct/technical-on-boarding"]){
                stage('docker env setup') {
                    kubesh 'apt-get -qq update && apt-get -qq -y install build-essential'
                }

                stage('dependencies'){
                    kubesh 'make -C ./go/src/github.com/samsung-cnct/technical-on-boarding/onboarding setup'
                }


                stage('test'){ // In the Makefile this provides `lint` and `vet` too.
                    kubesh 'make -C ./go/src/github.com/samsung-cnct/technical-on-boarding/onboarding test'
                }

                stage('build'){
                    kubesh 'make -C ./go/src/github.com/samsung-cnct/technical-on-boarding/onboarding build'
                }
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
  "kubectl exec -i ${env.HOSTNAME} -c ${env.CONTAINER_NAME} -- /bin/sh -c 'cd ${env.WORKSPACE} && export GOPATH=${env.GOPATH} CHECKOUT_PATH=${env.CHECKOUT_PATH} && ${command}'"
}

def customContainer(String name, Closure body) {
  withEnv(["CONTAINER_NAME=$name"]) {
    body()
  }
}
