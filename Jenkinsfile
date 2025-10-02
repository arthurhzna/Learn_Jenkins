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
          def target = ''
          def imageName = ''
          
          if (b == 'develop' || b == 'development') {
            target = 'dev'
            imageName = 'arthurhozanna/learn_jenkins_develop'
          } else if (b == 'staging' || b == 'stage') {
            target = 'staging'
            imageName = 'arthurhozanna/learn_jenkins_staging'
          } else if (b == 'master' || b == 'main' || b == 'live') {
            target = 'prod'
            imageName = 'arthurhozanna/learn_jenkins_prod'
          } else {
            target = 'dev'
            imageName = 'arthurhozanna/learn_jenkins_develop'
          }

          def shortSha = (env.GIT_COMMIT ?: sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim())
          def imageTag = "${shortSha ?: 'local'}-${env.BUILD_NUMBER ?: '0'}"
          def imageFull = "${imageName}:${imageTag}"

          // Set environment variables
          env.TARGET = target
          env.IMAGE_NAME = imageName
          env.IMAGE_TAG = imageTag
          env.IMAGE_FULL = imageFull

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
          set -x  # Enable debug mode
          
          # Debug: print all variables
          echo "DB_USER: ${DB_USER}"
          echo "DB_PASS: ${DB_PASS}"
          echo "DB_HOST: ${DB_HOST}" 
          echo "DB_PORT: ${DB_PORT}"
          echo "DB_NAME_PROD: ${DB_NAME_PROD}"
          echo "DB_NAME_TESTING: ${DB_NAME_TESTING}"
          echo "PORT: ${PORT}"
          
          # Create .env file with all required variables - SIMPLIFIED VERSION
          echo "DB_USERNAME=${DB_USER}" > .env
          echo "DB_PASSWORD=${DB_PASS}" >> .env  
          echo "DB_HOST=${DB_HOST}" >> .env
          echo "DB_PORT=${DB_PORT}" >> .env
          echo "DB_NAME=${DB_NAME_PROD}" >> .env
          echo "DB_NAME_TESTING=${DB_NAME_TESTING}" >> .env
          echo "PORT=${PORT}" >> .env

          # Show content of .env file
          echo "=== .env file content ==="
          cat .env
          echo "========================="

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