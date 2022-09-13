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
First of all you will need an empty directory where configs will be put by dcfg.
I hope you will find one! It also may be a git repo.

You need to create a dcfg config file (`dcfg.conf` by default) inside this directory. You can use 
`dcfg init` or `touch dcfg.conf` command. 
Then you will need to define ***bindings***, ***additions***, optionally 
***context directory*** and ***pins***. More information about dcfg config files can be found below.

After that, to copy defined ***additions*** to the ***context directory*** respecting
***bindings*** run:

`dcfg add`

Here you are!

To remove outdated additions (which were defined in the dcfg config file, 
were added using `dcfg add` command, but after a while were removed from it)
from context directory simply run:

`dcfg clean`

This command will not remove `.git` directory, dcfg config file and pinned directories and files.

## dcfg config file
There are 4 entities:
* ***addition*** (`add` directive) - global path to some directory or file which you want to store using dcfg. 
* ***context directory*** (`ctx` directive) - directory where all additions will be put.
* ***binding*** (`bind` directive) - global path to local path (relative to ***context directory***) binding.
* ***pin*** - (`pin` directive) - alien directory or file which you want to store with ***additions*** (e.g. readme file).

The only required things to be defined are ***additions*** and ***bindings***.

Example dcfg config file:
```
# 'ctx' directive - set context directory (can be used only once).
# Syntax: ctx [local path].
ctx ./  # ./ is default value

# Bindings (order makes sense):
# 'bind' directive - bind absolute path to a local one.
# Syntax: bind [absolute path] [local path (relative to the context dir path)].
bind ~ home/  # directories and files from $HOME will be copied to ./home/
bind / root/  # directories and files from / will be copied to ./root/

# Additions:
# 'add' directive - copy directories and files to the destination directory respecting bindings.
# Syntax: add [absolute path]
add ~/.config/i3   # ~/.config/i3  will be copied to ./home/.config/i3
add ~/.Xresources  # ~/.Xresources will be copied to ./home/.Xresources
add /etc/hostname  # /etc/hostname will be copied to ./root/etc/hostname

# Pins:
# 'pin' directive - pin non-addition file or directory so that it will not be removed when running 'dcfg clean'.
# Syntax: pin [local path (relative to the context dir path)]
pin README.md
```
