bosh create-release --force --tarball api-v3.tgz --name api-v3 --version 1.0.9

bosh ur api-v3.tgz

bosh -d api-v3 -n deploy test.yml
