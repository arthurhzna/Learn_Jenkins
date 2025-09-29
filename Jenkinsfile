pipeline {
  agent any

  stages {
    stage('Print Info') {
      steps {
        script {
          echo "Branch Name: ${env.BRANCH_NAME}"
          echo "Is Pull Request: ${env.CHANGE_ID}" // kalau PR ada isinya
          
          if (env.CHANGE_ID) {
            echo "🔵 Ini job PR (hasil merge simulasi)!"
          } else {
            echo "🟢 Ini job branch murni (kode asli)!"   
          }
        }
      }
    }
  }
}
