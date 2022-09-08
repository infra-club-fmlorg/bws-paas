if [ -z $(docker network list | grep -w application-network) ]; then
  docker network create application-network
fi

mkdir -p /tmp/bws-paas-queue/{active,incoming}
