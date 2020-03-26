## Provision Server
1. Manually run the contents of `./setup.sh` as `root`, then log in as `deployer`.
1. Set up a local `~/.ssh/config` rule with the following (replacing values in `{}` with their literal values):
```
Host {ENV_NAME}.pubgolf.co
    HostName {REMOTE IP}
    User deployer
    IdentityFile {/absolute/path/to/pubgolf_rsa}
```
1. Set up a local `/etc/hosts` rule to bypass Cloudflare and connect directly to the server IP.
1. Use `docker-machine` from your local computer to install `docker` on the server.

## Give CI/CD the Ability to Remotely Control Docker
1. Run the following command to generate the Docker certs:
  ```
  ./generate-client-certs.sh ~/.docker/machine/certs/ca.pem ~/.docker/machine/certs/ca-key.pem {ENV_NAME}
  ```
1. Copy the generated certs and `~/.docker/machine/certs/ca.pem` into LastPass and GitHub secrets.

## Configure Certbot and Nginx
1. Replace `example.com` in `./bootstrap-nginx.conf` and run the following:
  ```
  scp ./bootstrap-nginx.html {REMOTE_HOST}:/home/deployer/webroot/index.html
  scp ./bootstrap-nginx.conf {REMOTE_HOST}:/home/deployer/nginx-conf/nginx.conf
  ```
1. Open two SSH windows to the server:
  1. Run the contents of `./bootstrap-nginx.sh` and leave it running in the foreground.
  1. In the second window, run the contents of `./download-certbot-certs.sh`.
  1. When the previous step finishes, exit the `./bootstrap-nginx.sh` process.
1. Run the contents of `./download-certbot-config.sh` on the server.

## Set Up the Database
1. Run `mkdir -p /home/deployer/data/pubgolf_database`.
1. Run migrations (the standard CI/CD flow should take care of this for you).

## Set Up Env Config
1. Copy `.env.example` or the decrypted version of `.env.staging.gpg` into `.env.{ENV_NAME}`.
1. Encrypt it by running `gpg --symmetric --cipher-algo AES256 .env.{ENV_NAME}`.
  1. Enter the content of `../keys/gpg.key` as the password.
  1. Check the `.gpg` version of the file into the repo (but not the plaintext version).

