<p align="center">
<a href="https://mostafahussein.github.io/projects/conceal/#gh-light-mode-only">
<img width="300" src="https://raw.githubusercontent.com/mostafahussein/conceal/develop/dark.svg#gh-light-mode-only">
</a>
<a href="https://mostafahussein.github.io/projects/conceal/#gh-dark-mode-only">
<img width="300" src="https://raw.githubusercontent.com/mostafahussein/conceal/develop/light.svg#gh-dark-mode-only">
</a>
</p>

<p align="center">
  <a href="https://opensource.org/licenses/Apache-2.0">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  </a>
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/mostafahussein/conceal?color=purple"/>
</p>

<div align="center">

  </div>

![1  Conceal main](https://user-images.githubusercontent.com/4104127/205219226-b73b462b-7354-48f1-9912-b6b1f3796cbc.png)

[Conceal](https://mostafahussein.github.io/projects/conceal/) is an open‑source command line utility. It provides a secure method to get your secrets from your existing password manager.

# Features

Conceal provides the following features:

- Configured Session:
  - You can configure for how long your main password will be valid.
  - The password you entered will be saved locally and encrypted with OpenPGP.
- Integration with different password managers:
  - Currently, we support [Enpass](https://www.enpass.io), but more will come soon.

# Getting Started

## Installation

### Linux (AMD64)

```shell
curl -Lo conceal https://github.com/mostafahussein/conceal/releases/download/$(curl -s https://api.github.com/repos/mostafahussein/conceal/releases/latest | grep tag_name | cut -d '"' -f 4)/conceal-linux-amd64
chmod +x conceal
sudo mv conceal /usr/local/bin/conceal
```

### Linux (ARM64)

```shell
curl -Lo conceal https://github.com/mostafahussein/conceal/releases/download/$(curl -s https://api.github.com/repos/mostafahussein/conceal/releases/latest | grep tag_name | cut -d '"' -f 4)/conceal-linux-arm64
chmod +x conceal
sudo mv conceal /usr/local/bin/conceal
```

### Linux (ARM7)

```shell
curl -Lo conceal https://github.com/mostafahussein/conceal/releases/download/$(curl -s https://api.github.com/repos/mostafahussein/conceal/releases/latest | grep tag_name | cut -d '"' -f 4)/conceal-linux-arm
chmod +x conceal
sudo mv conceal /usr/local/bin/conceal
```

### macOS (AMD64)

```shell
curl -Lo conceal https://github.com/mostafahussein/conceal/releases/download/$(curl -s https://api.github.com/repos/mostafahussein/conceal/releases/latest | grep tag_name | cut -d '"' -f 4)/conceal-darwin-amd64
chmod +x conceal
sudo mv conceal /usr/local/bin/conceal
```
### macOS (ARM64)

```shell
curl -Lo conceal https://github.com/mostafahussein/conceal/releases/download/$(curl -s https://api.github.com/repos/mostafahussein/conceal/releases/latest | grep tag_name | cut -d '"' -f 4)/conceal-darwin-arm64
chmod +x conceal
sudo mv conceal /usr/local/bin/conceal
```

## How it works

Conceal will start connecting to your password manager and fetches the needed secrets (e.g. username and password) and set these values as environment variables based on what you have defined in the configuration file.
<p align="center">
<img width="600" src="https://user-images.githubusercontent.com/4104127/205493754-77d8ee0b-cc2f-4b97-ab71-c74c6fe6551f.png#gh-light-mode-only">
<img width="600" src="https://user-images.githubusercontent.com/4104127/205493816-22081f9e-d076-4297-8697-a0f07673da3e.png#gh-dark-mode-only">
</p>

## Roadmap

- [ ] Integrate with Secret managers
  - [x] [Enpass](https://www.enpass.io)
  - [ ] [Vault](https://www.vaultproject.io/)
  - [ ] [Lastpass](https://www.lastpass.com/)
- [x] Ability to execute different command-line utilities (e.g. kubectl, oc, aws)
- [x] Support for different environments (e.g. dev, prod)
- [ ] CI/CD integration
  - [ ] Github Actions
  - [ ] Gitlab CI

## Usage

```
A cli utility that provides a secure method to get your secrets from your existing password manager.

Usage:
  conceal [command]

Available Commands:
  exec        Execute commands for a given profile
  gen         Generate command alias
  init        Initialize Conceal Configuration
  version     Print the version number of conceal

Flags:
  -h, --help   help for conceal
```

## Configuration

After adding the binary to your system, you need to create a local directory and add the configuration file, this can be done by executing:
```
$ conceal init
```

`conceal init` offers 3 flags:
- `--secret-manager` (`-s` for short) a flag that is utilized for the secret manager that you are going to use, by default it will be `enpass`
- `--timeout` (`-t` for short) a flag that is utilized for keeping the password valid for defined number of minutes, by default it will be `15` minutes
- `--vault-location` (`-l` for short) a flag that is utilized for defining the secret manager location, by default it will use a default location for enpass which will be `~/Documents/Enpass/Vaults/primary`

### Adding command configuration

In order to add a new command configuration, you need to add `resource` block to `~/.conceal/config.hcl`.

#### Scenario 1

Let's see how `kubectl` will work with `conceal`:

Assuming that, you have `kubectl` installed and the authentication to your kubernetes cluster is being done through [AWS IAM](https://aws.amazon.com/iam/)

1. In Enpass, create a new login with the name `AWS_ACCESS_KEY`
2. Add the value of `AWS_ACCESS_KEY_ID` as the username and the value of `AWS_SECRET_ACCESS_KEY` as the password
3. Add the following block to `~/.conceal/config.hcl`
```
resource "profile" "k8s" {
  environment "default" {
    command = "kubectl"
    env = {
      id = "AWS_ACCESS_KEY"
      login = "AWS_ACCESS_KEY_ID"
      password = "AWS_SECRET_ACCESS_KEY"
    }
  }
}
```

The block above tells `conceal` the following:

- We need to define two environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` where the values can be found under `AWS_ACCESS_KEY` inside Enpass.
- Execute `kubectl` command

#### Scenario 2

Another example in case you are using OpenShift and you want to avoid using `oc login` every time you access your cluster.

1. In Enpass, create a new login with the name `openshift_login`
2. Add your openshift username in the username field value and the password in the password field value
3. Add the following block to `~/.conceal/config.hcl`
```
resource "profile" "openshift" {
  environment "default" {
    command = "oc"
    args = "login -u $OC_USERNAME -p $OC_PASSWORD https://localhost:8443"
    env = {
      id = "openshift_login"
      login = "OC_USERNAME"
      password = "OC_PASSWORD"
    }
  }
}

```
The above block tells `conceal` the following:

- Define a system environment variables named `OC_USERNAME` and `OC_PASSWORD` based on the values that we have added to `openshift_login` inside Enpass.
- We need to execute `oc` command as follows:

```
oc login -u $OC_USERNAME -p $OC_PASSWORD https://localhost:8443
```

### Supporting Multiple environments

If `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` on development environment are different from the production environment, and you want to make conceal handle both environments, you can update your config like below

```
resource "profile" "amazon" {
  environment "default" {
    command = "aws"
    env = {
      id = "AWS_ACCESS_KEY_DEV"
      login = "AWS_ACCESS_KEY_ID"
      password = "AWS_SECRET_ACCESS_KEY"
    }
  }
  environment "prod" {
    command = "aws"
    env = {
      id = "AWS_ACCESS_KEY_PROD"
      login = "AWS_ACCESS_KEY_ID"
      password = "AWS_SECRET_ACCESS_KEY"
    }
  }
}
```

Note: You can switch between environments by passing `-p default` or `-p prod`

### Adding command aliases

Once the configuration step is done, you can execute kubectl commands using the following command
```
conceal -p k8s "get pods"
```
But as this will look different than the normal kubectl commands and might be harder to type, you can generate an alias for `kubectl` and then add it to your `.bashrc` file or its equivalent depends on which shell you use.

In order to generate an alias you can execute the following command

```
conceal gen -a kubectl -p k8s -e default
```
The above command will generate an alias for `kubectl` so you can start typing normal `kubectl` commands and it will be handled by conceal.

## Contributing

We'd love for you to contribute to this tool. You can request new features by creating an [issue](https://github.com/mostafahussein/conceal/issues), or submit a [pull request](https://github.com/mostafahussein/conceal/pulls) with your contribution.

## License

Copyright &copy; 2022 Mostafa Hussein

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

