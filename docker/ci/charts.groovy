pipeline {
    agent {
        kubernetes {
            yaml """
spec:
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  containers:
  - name: "xuanwu"
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
    stages {
        stage("build&push helm") {
          steps {
              container('xuanwu') {
                  timestamps {
                      checkout([$class: "GitSCM", branches: [[name: "${xuanwu_charts_branch}"]], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: "${credential}", url: "https://github.com/Kyligence/xuanwu-charts.git"]]])

                      script {
                            sh "ls -lh"
                            sh "helm package ./xuanwu-log --app-version=${version} --version=${version}"
                            def awscp = [:]
                            awscp["global_template"] = {
                                withAWS(region:"us-east-1", credentials:'aws_global_s3_cp') {
                                  s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-${version}.tgz", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/charts/", acl:'PublicRead')
                                }
                            }
                            awscp["cn_template"] = {
                                withAWS(region:"cn-north-1", credentials:'aws_cn_s3_cp') {
                                  s3Upload(pathStyleAccessEnabled: true, payloadSigningEnabled: true, file:"xuanwu-log-${version}.tgz", bucket:"public.kyligence.io", path:"xuanwu/$release_type/$tag/charts/", acl:'PublicRead')
                                }
                            }
                            parallel awscp
                            withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: '${helm-credential}', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                                sh "curl -v -F file=@xuanwu-log-${version}.tgz -u ${env.USERNAME}:${env.PASSWORD} http://devops-nexus:8081/service/rest/v1/components?repository=kyligence-helm"
                            }
                      }
                  }
              }
          }
        } // stage end

  }
}