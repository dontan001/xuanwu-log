pipeline {
    agent none
    stages {

    stage("xuanwu-log") {
            environment {
                EXAMPLE_CREDS = credentials("${credential}")                        //can be used in this stage only
            }
          agent {
        kubernetes {
yaml """
spec:
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  containers:
  - name: "build"
    command:
    - "cat"
    image: "457798666374.dkr.ecr.us-west-2.amazonaws.com/jenkinsslave/xuanwu:alpine"
    imagePullPolicy: "Always"
    resources:
      limits:
        memory: "4096Mi"
        cpu: "4000m"
      requests:
        memory: "4096Mi"
        cpu: "4000m"
    tty: true
    volumeMounts:
    - mountPath: "/jenkins-common"
      name: "volume-0"
      readOnly: false
    - mountPath: "/var/run/docker.sock"
      name: "dockersock"
  - name: "jnlp"
    image: "jenkins/inbound-agent:4.3-4"
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
    volumeMounts:
    - mountPath: "/jenkins-common"
      name: "volume-0"
      readOnly: false
  volumes:
  - name: "volume-0"
    persistentVolumeClaim:
      claimName: "jenkins-common"
      readOnly: false
  - name: "dockersock"
    hostPath:
      path: "/var/run/docker.sock"
"""
        }
    }

        steps {
            container('build') {
                timestamps {
                    checkout([$class: "GitSCM", branches: [[name: "${branch}"]], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: "${credential}", url: "https://github.com/Kyligence/xuanwu-log.git"]]])
                      script {
                          withCredentials([file(credentialsId: 'xuanwueks-awsuser', variable: 'aws')]) {
                              sh "mkdir -p ~/.aws"
                              sh 'cat $aws > ~/.aws/credentials'
                          }

                          sh "pwd && ls -alt"
                          sh "aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 429636537981.dkr.ecr.us-west-2.amazonaws.com"

                          /*
                          sh "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/xuanwu-backup ./cmd/schedule/main.go"
                          sh "docker build -t registry.kyligence.io/xuanwu/xuanwu-log:${tag_schedule} -f docker/schedule/Dockerfile ."
                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                            sh "docker push registry.kyligence.io/xuanwu/xuanwu-log:${tag_schedule}"
                          }

                          sh "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/xuanwu-api ./cmd/api/main.go"
                          sh "docker build -t registry.kyligence.io/xuanwu/xuanwu-log:${tag_api} -f docker/api/Dockerfile ."
                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                            sh "docker push registry.kyligence.io/xuanwu/xuanwu-log:${tag_api}"
                          }
                          */

                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                            sh "docker pull docker.io/dontan001/xuanwu-log:schedule"
                            sh "docker tag docker.io/dontan001/xuanwu-log:schedule registry.kyligence.io/xuanwu/xuanwu-log:schedule"
                            sh "docker push registry.kyligence.io/xuanwu/xuanwu-log:schedule"

                            sh "docker pull docker.io/dontan001/xuanwu-log:api"
                            sh "docker tag docker.io/dontan001/xuanwu-log:api registry.kyligence.io/xuanwu/xuanwu-log:api"
                            sh "docker push registry.kyligence.io/xuanwu/xuanwu-log:api"
                         }

                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                              sh "docker pull docker.io/grafana/loki:2.2.1"
                              sh "docker tag docker.io/grafana/loki:2.2.1 registry.kyligence.io/xuanwu/loki:2.2.1"
                              sh "docker push registry.kyligence.io/xuanwu/loki:2.2.1"
                          }
                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                            sh "docker pull docker.io/grafana/fluent-bit-plugin-loki:2.1.0-amd64"
                            sh "docker tag docker.io/grafana/fluent-bit-plugin-loki:2.1.0-amd64 registry.kyligence.io/xuanwu/fluent-bit-plugin-loki:2.1.0-amd64"
                            sh "docker push registry.kyligence.io/xuanwu/fluent-bit-plugin-loki:2.1.0-amd64"
                          }
                          withDockerRegistry(credentialsId: 'registry-kyligence-io', url: 'https://registry.kyligence.io') {
                            sh "docker pull docker.io/nginxinc/nginx-unprivileged:1.19-alpine"
                            sh "docker tag docker.io/nginxinc/nginx-unprivileged:1.19-alpine registry.kyligence.io/xuanwu/nginx-unprivileged:1.19-alpine"
                            sh "docker push registry.kyligence.io/xuanwu/nginx-unprivileged:1.19-alpine"
                          }

                          sh "docker save -o  xuanwu-log-schedule.tar registry.kyligence.io/xuanwu/xuanwu-log:${tag_schedule}"
                          sh "docker save -o  xuanwu-log-api.tar registry.kyligence.io/xuanwu/xuanwu-log:${tag_api}"
                          sh "docker save -o  loki.tar registry.kyligence.io/xuanwu/loki:2.2.1"
                          sh "docker save -o  fluent-bit.tar registry.kyligence.io/xuanwu/fluent-bit-plugin-loki:2.1.0-amd6"
                          sh "docker save -o  nginx.tar registry.kyligence.io/xuanwu/nginx-unprivileged:1.19-alpine"

                          def awscp = [:]
                          awscp["global_template"] = {
                              withAWS(region:"us-east-1", credentials:'aws_global_s3_cp') {
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-schedule.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-api.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"loki.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"fluent-bit.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"nginx.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                              }
                          }
                          awscp["cn_template"] = {
                              withAWS(region:"cn-north-1", credentials:'aws_cn_s3_cp') {
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-schedule.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-api.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"loki.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"fluent-bit.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                                s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"nginx.tar", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/images/", acl:'PublicRead')
                              }
                          }
                          parallel awscp

                          sh "docker rmi registry.kyligence.io/xuanwu/xuanwu-log:${tag_schedule}"
                          sh "docker rmi registry.kyligence.io/xuanwu/xuanwu-log:${tag_api}"
                          sh "docker rmi registry.kyligence.io/xuanwu/loki:2.2.1"
                          sh "docker rmi registry.kyligence.io/xuanwu/fluent-bit-plugin-loki:2.1.0-amd64"
                          sh "docker rmi registry.kyligence.io/xuanwu/nginx-unprivileged:1.19-alpine"
                        }

                }
            }
        }

    } // stage end


  }
}