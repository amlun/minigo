#!/bin/bash

# 修改APP_NAME为云效上的应用名
APP_NAME=minigo-backend

PROG_NAME=$0
ACTION=$1
APP_START_TIMEOUT=20             # 等待应用启动的时间
APP_PORT=8808                    # 应用端口
APP_HOME=/usr/local/boc/backend # 从package.tgz中解压出来的包放到这个目录下
APP_OUT=/var/log/boc/app.log

# 进入应用目录
cd $APP_HOME
# make it executable
chmod +x $APP_NAME
usage() {
  echo "Usage: $PROG_NAME {start|stop|restart}"
  exit 2
}

health_check() {
  exptime=0
  echo "checking health"
  while true; do
    status_code=$(curl -L http://127.0.0.1:${APP_PORT}/status -o /dev/null -w '%{http_code}\n' -s)
    echo "code is $status_code"
    if [ "$status_code" != "200" ]; then
      echo -n -e "\rapplication not started"
    else
      break
    fi
    sleep 1
    ((exptime++))

    echo -e "\rWait app to pass health check: $exptime..."

    if [ $exptime -gt ${APP_START_TIMEOUT} ]; then
      echo 'app start failed'
      exit 1
    fi
  done
  echo "health check success"
}
start_application() {
  echo "starting $APP_NAME"
  nohup ./$APP_NAME >>${APP_OUT} 2>&1 &
  echo "started $APP_NAME"
}

stop_application() {
  checkpid=$(ps -ef | grep ${APP_NAME} | grep -v grep | grep -v 'deploy.sh' | awk '{print$2}')

  if [[ ! $checkpid ]]; then
    echo -e "\rno running $APP_NAME"
    return
  fi

  echo "stop $APP_NAME"
  times=60
  for e in $(seq 60); do
    sleep 1
    COSTTIME=$(($times - $e))
    checkpid=$(ps -ef | grep ${APP_NAME} | grep -v grep | grep -v 'deploy.sh' | awk '{print$2}')
    if [[ $checkpid ]]; then
      kill -9 $checkpid
      echo -e "\r        -- stopping $APP_NAME lasts $(expr $COSTTIME) seconds."
    else
      echo -e "\r$APP_NAME has exited"
      break
    fi
  done
  echo ""
}
start() {
  start_application
  health_check
}
stop() {
  stop_application
}
case "$ACTION" in
start)
  start
  ;;
stop)
  stop
  ;;
restart)
  stop
  start
  ;;
health-check)
  health_check
  ;;
*)
  usage
  ;;
esac
