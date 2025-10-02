pipeline {
  agent any

  environment {
    DOCKER_CREDENTIALS = credentials('docker-credential')
    GITHUB_CREDENTIALS = credentials('github-credential')
    SSH_KEY_CREDS = 'ssh-key'   
    HOST = credentials('host')
    USERNAME = credentials('username')
    DB_HOST = credentials('db-host')
    DB_PORT = credentials('db-port')
    DB_USER = credentials('db-user')
    DB_PASS = credentials('db-pass')
    DB_NAME_PROD = credentials('db-name-prod')           
    DB_NAME_TESTING = credentials('db-name-test') 
    PORT = credentials('app-port')
    TARGET = ''
    IMAGE_NAME = ''
    IMAGE_FULL = ''
    IMAGE_TAG = ''
  }

  stages {
    stage('Check skip') {
    steps {
        script {
        if (env.CHANGE_ID) {
            def title = env.CHANGE_TITLE ?: ''
            echo "PR title: ${title}"
            if (title.toLowerCase().contains('[skip ci]')) {
            currentBuild.result = 'ABORTED'
            error('Aborted by [skip ci] in PR title')
            }
        }
        }
    }
    }
    stage('Checkout') {
      steps {
        checkout scm
        script {
          echo "BRANCH_NAME = ${env.BRANCH_NAME}"
        }
      }
    }

    stage('Set Target & Image') {
      steps {
        script {
          def b = env.BRANCH_NAME?.toLowerCase()
          if (b == 'develop' || b == 'development') {
            env.TARGET = 'dev'
            env.IMAGE_NAME = 'arthurhozanna/learn_jenkins_develop'
          } else if (b == 'staging' || b == 'stage') {
            env.TARGET = 'staging'
            env.IMAGE_NAME = 'arthurhozanna/learn_jenkins_staging'
          } else if (b == 'master' || b == 'main' || b == 'live') {
            env.TARGET = 'prod'
            env.IMAGE_NAME = 'arthurhozanna/learn_jenkins_prod'
          } else {

            env.TARGET = 'dev'
            env.IMAGE_NAME = 'arthurhozanna/learn_jenkins_develop'
          }

          def shortSha = (env.GIT_COMMIT ?: sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim())
          env.IMAGE_TAG = "${shortSha ?: 'local'}-${env.BUILD_NUMBER ?: '0'}"
          env.IMAGE_FULL = "${env.IMAGE_NAME}:${env.IMAGE_TAG}"

          echo "TARGET=${env.TARGET} IMAGE=${env.IMAGE_FULL}"
        }
      }
    }
    stage('Unit Test') {
      steps {
        // sh '''
        //   docker run --rm \
        //     -v "$WORKSPACE":/app -w /app \
        //     golang:1.24.2-alpine3.20 sh -eux -c "
        //       apk add --no-cache git
        //       go env
        //       go mod tidy
        //       if go test ./...; then
        //         echo '✅ All tests passed'
        //       else
        //         echo '❌ Some tests failed' 1>&2
        //         exit 1
        //       fi
        //     "
        // '''
        sh '''
          # Create .env file with all required variables
          cat > .env << EOF
          DB_USERNAME=${DB_USER}
          DB_PASSWORD=${DB_PASS}  
          DB_HOST=${DB_HOST}
          DB_PORT=${DB_PORT}
          DB_NAME=${DB_NAME}
          DB_NAME_TESTING=${DB_NAME_TESTING}
          PORT=${PORT}
          EOF
          docker run --rm \
            -v "$WORKSPACE":/app -w /app \
            golang:1.24.2-alpine3.20 sh -eux -c "
              apk add --no-cache git
              go env
              go mod tidy
              if go test ./...; then
                echo '✅ All tests passed'
              else
                echo '❌ Some tests failed' 1>&2
                exit 1
              fi
            "
        '''
      }
    }
    stage('Build & Push Image') {
      when {
        expression { 
          def targetBranches = ['develop', 'staging', 'master', 'main', 'live']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
        // expression { return !env.CHANGE_ID && (currentBuild.result == null || currentBuild.result == 'SUCCESS') }
      }
      steps {
        script {
          withCredentials([usernamePassword(credentialsId: 'docker-credential', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
            sh """
              echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
              docker build -t ${IMAGE_FULL} .
              docker push ${IMAGE_FULL}
            """
          }
        }
      }
    }

    stage('Deploy to Remote Host') {
      when {
        expression { 
          def targetBranches = ['develop', 'staging', 'master', 'main', 'live']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
        // expression { return !env.CHANGE_ID && (currentBuild.result == null || currentBuild.result == 'SUCCESS') }
      }
      steps {
        script {
          withCredentials([sshUserPrivateKey(credentialsId: env.SSH_KEY_CREDS, keyFileVariable: 'SSH_KEY_FILE', usernameVariable: 'SSH_USER')]) {
            sh """
              ssh -o StrictHostKeyChecking=no -i ${SSH_KEY_FILE} ${USERNAME}@${HOST} '
                cd /path/to/app || exit 1
                # tarik image terbaru berdasarkan tag
                docker-compose -f docker-compose.${TARGET}.yaml pull || true
                docker-compose -f docker-compose.${TARGET}.yaml up -d --remove-orphans
              '
            """
          }
        }
      }
    }
  }

  post {
    always {
      echo "Finished pipeline for branch=${env.BRANCH_NAME} target=${env.TARGET} image=${env.IMAGE_FULL}"
    }
    failure {
      echo "Build failed for ${env.BRANCH_NAME}"
    }
  }
}