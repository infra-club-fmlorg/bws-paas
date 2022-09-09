if test -z "$(docker network list | grep -w paas-network)"; then
  docker network create paas-network
fi

mkdir -p ~/bws-paas-queue/active
mkdir -p ~/bws-paas-queue/incoming
