#!/bin/bash

# Technical On Boarding script for Samsung SDSA, Inc 

step_counter=1

onramp_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


echo "********************************************************************************"
echo "*                                                                              *"
echo "* Installing the Samsung Stack for OSX.                                        *"
echo "*                                                                              *"
echo "* Before you start:                                                            *"
echo "*                                                                              *"
echo "********************************************************************************"
echo

default_install_dir="${HOME}/Development/Samsung"

# Ask where to install everything.
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
read -ep "Choose an install directory [press enter to use \`${default_install_dir}\`]: " -r install_dir
if test -z "${install_dir}";
then
    # If we just read a 0-length string, use the default install directory.
    install_dir="${default_install_dir}"
fi
echo


# If the chosen install directory does not exist, create it.
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Preparing install directory..."
echo
if [ ! -d "${install_dir}" ];
then
    mkdir -p "${install_dir}"
fi
echo "Installing to: \`${install_dir}\`."
echo


# Check if Homebrew is installed, and if it is not, install it.
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Ensuring Homebrew is installed..."
echo
if ! type "brew" > /dev/null;
then
    echo "Homebrew is required. Installing now."
    ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
    brew doctor
    brew update
    echo "********************************************************************************"
    echo "*                                                                              *"
    echo "* Please take a moment to fix any issues that \`brew doctor\` may have           *"
    echo "* complained about.                                                            *"
    echo "*                                                                              *"
    echo "* Press any key when finished (or if there were no issues).                    *"
    echo "*                                                                              *"
    echo "********************************************************************************"
    read -n 1 -s -r
    echo
    echo "Continuing..."
    echo
fi

# Check for required prerequisites
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Check for required prerequisites..."
echo


# Install required Homebrew packages.
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Installing Homebrew packages..."
echo

# wget
brew install wget
# silver searcher
brew install ag
# git
brew install git
# python and pip
brew install python && sudo easy_install pip
# ruby
brew install rbenv ruby-build
# go
brew install go

# install virtualbox + vagrant
brew cask install virtualbox
brew cask install vagrant
brew install boot2docker

echo
echo "Installing Python packages..."
echo
pip install cql   # cql is used to connect to cassandra
pip install ansible==2.0.0.2
pip install awscli==1.10.4
pip install boto==2.39.0
pip install paramiko==2.0.0
pip install cffi==1.5.0


brew doctor  # The installs can expose issues that weren't originally exposed in the original check. Run it again.
new_brew_errors=$(brew doctor 2>&1 | grep -c "Warning")

if [ "$new_brew_errors" -gt 0 ]
then
    echo "********************************************************************************"
    echo "*                                                                              *"
    echo "* Brew doctor is not happy after running the installs.                         *"
    echo "*                                                                              *"
    echo "* Please open a new tab and fix any issues that \`brew doctor\` may have         *"
    echo "* complained about, particularly around permissions and symlinks.              *"
    echo "*                                                                              *"
    echo "* Failure to do so may result in failures later in the script                  *"
    echo "*                                                                              *"
    echo "* Press any key when finished.                                                 *"
    echo "*                                                                              *"
    echo "********************************************************************************"
    read -n 1 -s -r
    echo
    echo "Continuing..."
    echo
fi

# Update PATH and adding exports
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Updating PATH..."
echo
PATH=/usr/local/bin:${PATH}

bash_profile=$HOME/.bash_profile
touch $bash_profile

# color on your terminal
echo "# terminal colors"  >> $bash_profile
echo "export CLICOLOR=1"  >> $bash_profile
echo "export LSCOLORS=GxFxCxDxBxegedabagaced"  >> $bash_profile

echo "please add environment variables when complete to the .bash_profile"


echo
# Terraform goodness (may change once our plugins are updated for version 0.7.0)
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo " Installing terraform..manually.....(step should be removed in the near future and handled by brew)..."
echo

mkdir -p $HOME/bin && cd $HOME/bin
wget -O terraform.zip https://releases.hashicorp.com/terraform/0.6.16/terraform_0.6.16_darwin_amd64.zip
wget -O terraform-provider-coreosbox.tar.gz https://github.com/samsung-cnct/terraform-provider-coreosbox/releases/download/v0.0.1/terraform-provider-coreosbox_darwin_amd64.tar.gz
unzip terraform.zip && rm terraform.zip
tar xzf terraform-provider-coreosbox.tar.gz && rm terraform-provider-coreosbox.tar.gz

echo 'export PATH=$HOME/bin:PATH' >> $bash_profile


# Initialize boot2docker
echo
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Step $((step_counter++)) <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo
echo "Initializing boot2docker..."
echo
boot2docker config | sed 's/DockerPort = 0/DockerPort = 2376/' > ~/.boot2docker/profile
boot2docker init
boot2docker up

# Finished

echo
echo '********************************************************************************'
echo '*                                                                              *'
echo '* FINISHED!                                                                    *'
echo '*                                                                              *'
echo '* The Samsung Stack should now be installed                                    *'
echo '*                                                                              *'
echo '* please read github page on contributing, here we fork a projec               *'
echo '*  and not contribute directly via a branch                                     *'
echo '********************************************************************************'


open https://guides.github.com/activities/contributing-to-open-source/

