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
          checkout scm
          if (env.CHANGE_ID) {
            def title = env.CHANGE_TITLE ?: ''
            echo "PR title: ${title}"
            if (title.toLowerCase().contains('[skip ci]')) {
              currentBuild.result = 'ABORTED'
              error('Aborted by [skip ci] in PR title')
            }
          }
          def lastMsg = sh(script: 'git log -1 --pretty=%B', returnStdout: true).trim().toLowerCase()
          echo "Last commit message: ${lastMsg}"
          if (lastMsg.contains('[skip ci]') || lastMsg.contains('[ci skip]')) {
            currentBuild.result = 'ABORTED'
            error('Aborted by skip token in commit message')
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

          TARGET = target
          IMAGE_NAME = imageName
          IMAGE_TAG = imageTag
          IMAGE_FULL = imageFull

          echo "TARGET=${TARGET} IMAGE=${IMAGE_FULL}"
        }
      }
    }
    stage('Unit Test') {
      steps {
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
          def targetBranches = ['develop', 'staging', 'master', 'main']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
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
    stage('Update docker-compose.yaml') {
      when {
        expression { 
          def targetBranches = ['develop', 'staging', 'master', 'main']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
      }
      steps {
        script {
          sh """
            sed -i 's|image: ${IMAGE_NAME}:[^\\n]*|image: ${IMAGE_FULL}|' docker-compose.${TARGET}.yaml
          """
        }
      }
    }

    stage('Commit and Push Changes') {
      when {
        expression { 
          def targetBranches = ['develop', 'staging', 'master', 'main']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
      }
      steps {
        script {
          withCredentials([usernamePassword(credentialsId: 'github-credential', usernameVariable: 'GITHUB_USERNAME', passwordVariable: 'GITHUB_PASSWORD')]) {
            sh """
              git config --global user.name 'Jenkins CI'
              git config --global user.email 'jenkins@example.com'
              git remote set-url origin https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/arthurhzna/Learn_Jenkins.git
              git add docker-compose.${TARGET}.yaml
              git commit -m 'Update image version to ${IMAGE_TAG} [skip ci]' || echo 'No changes to commit'
              git pull origin ${env.BRANCH_NAME} --rebase
              git push origin HEAD:${env.BRANCH_NAME}
            """
          }
        }
      }
    }
    stage('Deploy to Remote Host') {
      when {
        expression {
          def targetBranches = ['develop', 'staging', 'master', 'main']
          return !env.CHANGE_ID && targetBranches.contains(env.BRANCH_NAME)
        }
      }
      steps {
        script {
          def deployEnv = (params.DEPLOY_ENV ?: env.DEPLOY_ENV ?: env.BRANCH_NAME ?: 'dev').toLowerCase()
          def appBase = "/home/arthurhozana123/go"
          def appName = "Learn_Jenkins"
          def suffix = (deployEnv == 'prod' || deployEnv == 'production') ? '_prod' : (deployEnv == 'staging' ? '_staging' : '_dev')
          def APP_DIR = "${appBase}/${appName}${suffix}"
          withCredentials([
            sshUserPrivateKey(credentialsId: env.SSH_KEY_CREDS, keyFileVariable: 'SSH_KEY_FILE', usernameVariable: 'SSH_USER'),
            usernamePassword(credentialsId: 'github-credential', usernameVariable: 'GITHUB_USERNAME', passwordVariable: 'GITHUB_PASSWORD')
          ]) {
            sh """
              ssh -o StrictHostKeyChecking=no -i ${SSH_KEY_FILE} ${USERNAME}@${HOST} '
                # Define app directory
                APP_DIR="${APP_DIR}"

                # ensure parent exists
                mkdir -p "$(dirname "\$APP_DIR")"

                if [ -d "\$APP_DIR/.git" ]; then
                  echo "Directory exists. Pulling latest changes."
                  cd "\$APP_DIR"
                  git pull origin ${env.BRANCH_NAME}
                else
                  echo "Directory does not exist. Cloning repository using HTTPS (PAT)."
                  git clone -b ${env.BRANCH_NAME} https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/arthurhzna/Learn_Jenkins.git "\$APP_DIR"
                  cd "\$APP_DIR"
                fi

                # Pull latest image
                docker-compose -f docker-compose.${TARGET}.yaml pull || true

                # Copy .env.example to .env and update with sed
                cp .env.example .env
                sed -i "s/^DB_USERNAME=.*/DB_USERNAME=${DB_USER}/" .env
                sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=${DB_PASS}/" .env
                sed -i "s/^DB_HOST=.*/DB_HOST=${DB_HOST}/" .env
                sed -i "s/^DB_PORT=.*/DB_PORT=${DB_PORT}/" .env
                sed -i "s/^DB_NAME=.*/DB_NAME=${DB_NAME_PROD}/" .env
                sed -i "s/^PORT=.*/PORT=${PORT}/" .env

                # Show .env content for debugging
                echo "=== .env file content ==="
                cat .env
                echo "========================="

                # Deploy
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