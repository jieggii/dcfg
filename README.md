# dcfg
dcfg - distribute config. 
Simple tool for copying and distributing your config files.

## Installation
**From sources:**
```shell
git clone https://github.com/jieggii/dcfg
cd dcfg/
make
sudo make install
```

## How to use
First of all, you will need an empty directory where config files will be put by dcfg.
I hope you will find one! It also may be a git repo.

You need to create a dcfg config file (`dcfg.conf` by default) in this directory. You can use 
`dcfg init` or `touch dcfg.conf` commands. 
Then you will need to define ***bindings***, ***additions***, optionally 
***context directory*** and ***pins***. More information about dcfg config files can be found below.

After that, to copy defined ***additions*** to the ***context directory*** respecting
***bindings*** run

`dcfg add`

Here you are!


## dcfg config file
There are 4 entities:
* ***addition*** - global path to some directory or file which you want to store using dcfg. 
* ***context directory*** (`./` by default) - directory where all additions will be put.
* ***binding*** - global path to local path (relative to ***context directory***) binding.
* ***pin*** - alien directory or file which you want to store with ***additions*** (e.g. readme file or `.git` directory).
