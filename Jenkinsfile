#!/usr/bin/env groovy
pipeline {
  agent { label 'mesos-sous' }
  options {
    // Version 3
		// Set to 1 day to allow people to input whether they want to go to Prod on the Master branch build/deploys
  	timeout(time: 1, unit: 'DAYS')
  }
  stages {
    stage('Inited Values') {
      steps {
        echo "BUILD_NUMBER=$BUILD_NUMBER"
        echo "BRANCH_NAME=$BRANCH_NAME"
        echo "NODE_NAME=$NODE_NAME"
        echo "NODE_LABELS=$NODE_LABELS"
        echo "BUILD_URL=$BUILD_URL"
        script {
          def notifier = new org.gradiant.jenkins.slack.SlackNotifier()

          env.SLACK_CHANNEL = '#team-eng-sous-bots, #tech-deploy'

          notifier.notifyStart()
        }
      }
    }
    stage('Test'){
      parallel {
        stage('Unit') {
  				agent { label 'mesos-sous' }
          steps {
            echo "unit test step"
            echo $PATH
            sh '''#!/usr/bin/env bash
            echo $PATH
            make test-unit
            '''
          }
        }
        stage('Smoke') {
   				agent { label 'mesos-sous' }
					steps {
            echo "smoke test step"
            sh '''#!/usr/bin/env bash
            make test-smoke
            '''
          }
        }
        stage('Integration') {
   				agent { label 'mesos-sous' }
          steps {
            echo "integration test"
            sh '''#!/usr/bin/env bash
            make test-integration
            '''
          }
        }
      }
    }
    stage('Build') {
		options {
			timeout(time: 10, unit: 'MINUTES')
		}
    steps {
        echo 'Build in Jenkinsfile'
        sh 'make release'
        echo 'leaving Jenkinsfile stage build'
      }
      environment {
        SOUS_CMD_TAG = "latest"
      }
    }
    stage('master-branch-deploy'){
      when{
        branch 'master'
      }
      parallel {
        stage('Deploy ci-sf') {
  				agent { label 'mesos-sous' }
          options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
            retry(4) {
							script {
            		def deploy = new com.opentable.sous.Deploy()
            		deploy.execute()
            	}
						}
          }
          environment {
            SOUS_CLUSTER = "ci-sf"
          }
        }
        stage('Deploy pp-sf') {
   				agent { label 'mesos-sous' }
          options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
          	retry(4) {
							script {
              	def deploy = new com.opentable.sous.Deploy()
								deploy.execute()
							}
						}
          }
          environment {
            SOUS_CLUSTER = "pp-sf"
          }
        }
        stage('Deploy ci-uswest2') {
   				agent { label 'mesos-sous' }
					options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
						retry(4) {
							script {
            		def deploy = new com.opentable.sous.Deploy()
              	deploy.execute()
							}
						}
          }
          environment {
            SOUS_CLUSTER = "ci-uswest2"
          }
        }
        stage('Deploy pp-uswest2') {
   				agent { label 'mesos-sous' }
 					options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
						retry(4) {
							script {
            		def deploy = new com.opentable.sous.Deploy()
              	deploy.execute()
							}
            }
          }
          environment {
            SOUS_CLUSTER = "pp-uswest2"
          }
        }
      }
    }
		stage('Input Prod Deployment') {
      agent none
      when{
        branch 'master'
      }
      steps {
        script {
        	env.DEPLOY_TO_PROD = input message: 'User input required',
          parameters: [choice(name: 'Deploy to Prod?', choices: 'no\nyes', description: 'Choose "yes" if you want to deploy this build to production.')]
        }
      }
		}
    stage('Deploy to Production') {
      when {
      	environment name: 'DEPLOY_TO_PROD', value: 'yes'
      }
      parallel {
        stage('Deploy prod-sf') {
  				agent { label 'mesos-qa-uswest2' }
					options {
						timeout(time: 5, unit: 'MINUTES')
					}
          steps {
						retry(4) {
            	script {
            		def deploy = new com.opentable.sous.Deploy()
              	deploy.execute()
            	}
						}
          }
          environment {
            SOUS_CLUSTER = "prod-sc"
          }
        }
        stage('Deploy prod-ln') {
   				agent { label 'mesos-qa-uswest2' }
					options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
          	retry(4) {
							script {
              	def deploy = new com.opentable.sous.Deploy()
            		deploy.execute()
            	}
						}
          }
          environment {
            SOUS_CLUSTER = "prod-ln"
          }
        }
        stage('Deploy prod-euwest1') {
   				agent { label 'mesos-qa-uswest2' }
					options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
						retry(4) {
							script {
              	deploy = new com.opentable.sous.Deploy()
            		deploy.execute()
							}
            }
          }
          environment {
            SOUS_CLUSTER = "prod-euwest1"
          }
        }
        stage('Deploy prod-uswest2') {
   				agent { label 'mesos-qa-uswest2' }
					options {
						timeout(time: 5, unit: 'MINUTES')
					}
					steps {
						retry(4) {
						script {
              def deploy = new com.opentable.sous.Deploy()
            	deploy.execute()
							}
            }
          }
          environment {
            SOUS_CLUSTER = "prod-uswest2"
          }
        }
      }
    }
	}
  post {
    always {
      echo 'This will always run'

      script {
        def notifier = new org.gradiant.jenkins.slack.SlackNotifier()

        env.SLACK_CHANNEL = '#team-eng-sous-bots, #tech-deploy'
        env.CHANGE_LIST = 'true'
        env.NOTIFY_SUCCESS = 'true'

        notifier.notifyResult()
      }

      //slackSend color: 'good', message: 'Message from Jenkins Pipeline'
      //script {
      //  def slack = new com.opentable.sous.Slack()
      //  slack.call(currentBuild.currentResult, '#team-eng-sous-bots')
      //}
    }
    success {
      echo 'This will run only if successful'
    }
    failure {
      echo 'This will run only if failed'
    }
    unstable {
      echo 'This will run only if the run was marked as unstable'
    }
    changed {
      echo 'This will run only if the state of the Pipeline has changed'
      echo 'For example, if the Pipeline was previously failing but is now successful'
    }
  }
}
