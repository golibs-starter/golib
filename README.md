# Golib

Common core for Golang project.

### ⚠️ **Notice** ⚠️
Our modules are in private repo, so you need to config something bellow before start your develop.
#### 1. Setup `GOPRIVATE`

Run the command `export GOPRIVATE="gitlab.id.vin"` to add `gitlab.id.vin` as private repo.

For future usage, you might add above command to `.bashrc` or `.zshrc`.

#### 2. Add credentials to private host
Run the following command line to load `https://gitlab.id.vin/` using SSH:
```shell
git config --global url."git@gitlab.id.vin:".insteadOf "https://gitlab.id.vin/"
```

Or with access token in URL:
```shell
git config --global url."https://oath2:{your_access_token}@gitlab.id.vin/".insteadOf https://gitlab.id.vin/
```

> TODO Add more instructions
