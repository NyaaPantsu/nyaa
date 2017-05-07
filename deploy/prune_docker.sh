# Prune images and volumes
# 
# Docker tends to take a lot of space. This will remove dangling images and
# volumes not used by at least container.
# WARNING: You might not want to run this if you have stuff that are dangling
# that you want to keep.
#
#!/bin/bash

docker image prune -f
docker volume prune -f
