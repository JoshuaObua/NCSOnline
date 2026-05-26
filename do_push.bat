cd /d d:\ncs-online
set GIT_SSH_COMMAND=ssh -i d:\ncs-online\.deploy_key_ncs_online -o IdentitiesOnly=yes
git push -u origin master
