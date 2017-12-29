#!groovy

pipeline {
    agent {
        label 'master'
    }
    options {
        buildDiscarder(logRotator(numToKeepStr: '30'))
    }
    parameters {
        string(
            name: 'BRANCH_WEBAPP',
            defaultValue: 'ocmIntegration',
            description: 'Branch of webapp'
        )
        string(
            name: 'BRANCH_FRAMEWORK',
            defaultValue: 'ocmIntegration',
            description: 'Branch of framework, by default, use the same branch as webapp'
        )
        string(
            name: 'BRANCH_DATABASE',
            defaultValue: 'ocmIntegration',
            description: 'Branch of database, by default, use the same branch as webapp'
        )
    }
    environment {
        PATH = "${PATH}"
        MYSQL_CONTAINER_NAME = 'ocm-core-db'
        MYSQL_VERSION = '5.5.53'
        APP_NAME = 'ocmapps'
        DOCKER_DIR = "${WORKSPACE}/pipelines/docker"
    }
    stages {
        stage ('initialize') {
            steps {
                script {
                    //set properties
                    env.branchWebApp = "${params.BRANCH_WEBAPP}"
                    env.branchFramework = "${params.BRANCH_FRAMEWORK}"
                    env.branchDatabase = "${params.BRANCH_DATABASE}"
                }
            }
        }
        stage ('checkout') {
            environment {
                gitUrl = 'git@git.zias.io:csun'
            }
            steps {
                script {
                    // checkout webapp
                    checkoutSCM("${gitUrl}", 'webapp', "${branchWebApp}")
                    // checkout framework
                    checkoutSCM("${gitUrl}", 'framework', "${branchFramework}")
                    // checkout database
                    checkoutSCM("${gitUrl}", 'database', "${branchDatabase}")
                    //checkout pipeline
                    checkoutSCM("${gitUrl}", 'pipelines', 'master')
                }
            }
        }
        stage ('prepare webapp') {
            environment {
                propTemplate = 'test/dev01/build.properties.template'
                catalinaHome = "${env.JENKINS_HOME}/workspace/apache-tomcat-7.0.61"
            }
            steps {
                script {
                    // prepare Catalina
                    def catalinaExists = fileExists "${catalinaHome}"
                    if (catalinaExists) {
                        sh "rm -rf ${catalinaHome}"
                    }
                    dir ("${env.JENKINS_HOME}/workspace") {
                        sh returnStdout: true, script: '''
                            wget -q https://archive.apache.org/dist/tomcat/tomcat-7/v7.0.61/bin/apache-tomcat-7.0.61.tar.gz
                            tar zxf apache-tomcat-7.0.61.tar.gz
                            rm -rf apache-tomcat-7.0.61.tar.gz
                        '''
                    }
                    // check Catalina version
                    sh "${catalinaHome}/bin/catalina.sh version"

                    dir ("${catalinaHome}/bin") {
                        echo "starting catalina"
                        sh returnStdout: true, script: '''
                            ./catalina.sh start
                            sleep 30
                            ./catalina.sh stop
                        '''
                    }

                    dir('webapp/build') {
                        withAnt(installation: 'ant', jdk: 'jdk1.8.0_60') {
                            sh "ant -f initBuildPros.xml initBuildPros -Dtemplate=${propTemplate} -Djobname=${APP_NAME}"
                            updateProperty('my.build', 'my.catalina.home', "${catalinaHome}")
                            updateProperty('my.build', 'repository.type', 'repository.local')
                            updateProperty('my.build', 'database.hostname', "${MYSQL_CONTAINER_NAME}")
                            updateProperty('my.build', 'database.name', "${APP_NAME}")
                            updateProperty('my.build', 'EmailManager.type', 'EmailManager.mock')
                            updateProperty('my.build', 'internal.api.port', '8082')

                        	updateProperty('test', 'project.name', "${APP_NAME}")
                            updateProperty('test', 'product.internal.api.url', 'http://localhost:8082')
                        	updateProperty('test', 'internal.api.url', 'http://localhost:8080')
                        }
                    }
                }
            }
        }
        stage ('compile & build') {
            steps {
                script {
                    dir("${WORKSPACE}/webapp/build") {
                        withAnt(installation: 'ant', jdk: 'jdk1.8.0_60') {
                            sh "ant clean war"
                        }
                    }
                }
            }
        }
        stage ('prepare webapp docker iamge') {
            environment {
                catalinaHome = "${env.JENKINS_HOME}/workspace/apache-tomcat-7.0.61"
            }
            steps {
                script {
                    sh '''
                        unzip -q ${WORKSPACE}/webapp/dist/${APP_NAME}.war -d ${DOCKER_DIR}/catalina/${APP_NAME}
                        cp ${catalinaHome}/conf/Catalina/localhost/${APP_NAME}.xml ${DOCKER_DIR}/catalina/
                    '''
                }
            }
        }
        stage ('prepare database') {
            steps {
                script {
                    dir ("${env.WORKSPACE}/pipelines/docker/mysql") {
                        sh returnStatus: true, returnStdout: false, script: '''
                            sed -i "s/DB_NAME/${APP_NAME}/g zg_shard_configuration.sql"
                            sed -i "s/SHARD_HOST/${MYSQL_CONTAINER_NAME}/g zg_shard_configuration.sql"
                        '''
                        sh returnStatus: true, returnStdout: false, script: '''
                            docker stop ${MYSQL_CONTAINER_NAME} 2>1 > /dev/null || true
                            docker rm ${MYSQL_CONTAINER_NAME} 2>1 > /dev/null || true
                            docker run -d -i \
                                --name ${MYSQL_CONTAINER_NAME} \
                                -p 3306:3306 \
                                -v ${WORKSPACE}/database:/database \
                                -v ${WORKSPACE}/pipelines/docker/mysql/conf.d:/etc/mysql/conf.d \
                                -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
                                mysql:${MYSQL_VERSION}
                            sleep 30
                            docker exec -i ${MYSQL_CONTAINER_NAME} \
                                bash -c "cd /database/scripts; ./refresh-db.sh"
                            mysql --host=localhost \
                                --user=root \
                                --protocol=tcp \
                                < zg_shard_configuration.sql
                            mysqldump --all-databases \
                                --ignore-table=mysql.user \
                                --host=localhost \
                                --user=root \
                                --protocol=tcp \
                                > ${APP_NAME}.sql
                            gzip -f ${APP_NAME}.sql
                        '''
                    }
                }
            }
        }
        stage ('publish database docker image') {
            steps {
                script {
                    publishDockerImage("ocm/core", "${DOCKER_DIR}/mysql", "database", 'us-west-2', '122972921717')
                    publishDockerImage("ocm/core", "${DOCKER_DIR}/webapp", "webapp", 'us-west-2', '122972921717')
                }
            }
        }
    }
}

def checkoutSCM(String url, String repoName, String branchName) {
    checkout([
        $class: 'GitSCM',
        branches: [[name: "${branchName}"]],
        doGenerateSubmoduleConfigurations: false,
        extensions: [
            [$class: 'WipeWorkspace'],
            [$class: 'RelativeTargetDirectory', relativeTargetDir: "${repoName}"]],
        submoduleCfg: [],
        userRemoteConfigs: [[
            credentialsId: '599505f7-c287-43f0-aca4-4f26072519ff',
            url: "${url}/${repoName}.git"
        ]]
    ])
}

def updateProperty(String file, String key, String value) {
    sh "ant updateProperty -Dfile=${file} -Dkey=${key} -Dvalue=${value}"
}

def publishDockerImage(String awsRepoName, String dockerFilePath, String imageType, String region, String account) {

    dir("${WORKSPACE}/${imageType}") {
        COMMIT_ID = sh(returnStdout: true, script: 'git rev-parse HEAD').trim().take(7)
    }

    IMAGE_NAME = "${awsRepoName}"
    IMAGE_TAG = "ocm-core-${imageType}.${COMMIT_ID}"

    dir("${WORKSPACE}/${dockerFilePath}") {
        IMAGE = docker.build("${IMAGE_NAME}:${IMAGE_TAG}", "-f ./Dockerfile .")
    
        withAWS(role: 'jenkinsfarm-cross-account', roleAccount: "${account}") {
            sh "eval \$(aws ecr get-login --region ${region} --no-include-email)"
            docker.withRegistry("https://${account}.dkr.ecr.${region}.amazonaws.com") {
                IMAGE.push()
                IMAGE.push('latest')
            }
        }
    }
}

def pushDockerImage(String region, String account) {
    withAWS(role: 'jenkinsfarm-cross-account', roleAccount: "${account}") {
        sh "eval \$(aws ecr get-login --region ${region} --no-include-email)"
        docker.withRegistry("https://${account}.dkr.ecr.${region}.amazonaws.com") {
            IMAGE.push()
            IMAGE.push('latest')
        }
    }
}