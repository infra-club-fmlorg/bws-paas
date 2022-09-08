if [ -z $(docker network list | grep -w application-network)]; then
  docker network create application-network
fi

mkdir -p ~/bws-paas-queue/active
mkdir -p ~/bws-paas-queue/incoming
