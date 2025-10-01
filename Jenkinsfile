pipeline {
  agent any

  environment {
    IMAGE_NAME = 'Learn_Jenkins'   
    DOCKER_CREDENTIALS = credentials('docker-credential')
    GITHUB_CREDENTIALS = credentials('github-credential')
    SSH_KEY = credentials('ssh-key')
    HOST = credentials('host')
    USERNAME = credentials('username')   
  }

  stages {
    stage('Check Commit Message') {
      steps {
        script{
          echo "Commit Message: ${env.COMMIT_MESSAGE}"
          echo "IMAGE_NAME: ${IMAGE_NAME}"
          echo "DOCKER_CREDENTIALS: ${DOCKER_CREDENTIALS}"
          echo "GITHUB_CREDENTIALS: ${GITHUB_CREDENTIALS}"
          echo "SSH_KEY: ${SSH_KEY}"
          echo "HOST: ${HOST}"
          echo "USERNAME: ${USERNAME}"
        }
      }
    }
}