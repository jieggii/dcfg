<p align="center">
    <img alt="dcfg logo" src="https://imgur.com/pRCVSAo.jpg" height=225 />
</p>
<p align="center">
    Minimalist tool for copying, storing and distributing your system-wide and user config files.
</p>

## Conception
Long story short - dcfg copies config files from different places to the directory you've 
chosen. The program has its own config file in this directory and provides you some entities 
to do it in a more convenient way.

Comfy documentation and detailed explanation can be found [here](https://jieggii.github.io/dcfg).

## Installation
At first install **diff** and **go** packages, they are the only dcfg dependencies.
**diff** is needed to show differences between files while using `dcfg extract` command
and **go** is used to compile dcfg itself. It is most likely, that they are already
installed on your system but be aware!

**On Debian-based distributions:**
```shell
$ sudo apt install diff go
```

**On Arch-based distributions:**
```shell
$ sudo pacman -S diff go
```
After installing dependencies you can proceed to installation of dcfg.
There are some interchangeable options to install it:
### A. From sources
```shell
$ git clone https://github.com/jieggii/dcfg
$ cd dcfg/
$ make
$ sudo make install
$ dcfg --version  # check dcfg version
```
> dcfg binary will be installed to `/usr/bin/dcfg`.

### B. From AUR (not yet)
```shell
$ yay -S dcfg
$ dcfg --version  # check dcfg version
```
> dcfg binary will be installed to `/usr/bin/dcfg`.

### C. Using go (not yet)
```shell
$ go install github.com/jieggii/dcfg@latest
$ dcfg --version  # check dcfg version
```
> dcfg binary will be installed to `$GOPATH/bin`. Ensure that `$GOPATH/bin` is in your `$PATH`.

## Simple use case
First: why you will ever want to store your config files in one place if they are already
stored on your computer? If you ask this question, most probably that you really don't need it.
But I will answer anyway just in case. I see two main reasons here:
1. To share your funky so-called **dotfiles** with the outer world and your nerdy internet friends
2. To back up your necessary configuration files, so they don't get lost if you blow your machine
3. To use same configuration on other computer

It seems that's all.

So! You've decided to store your config files to distribute them. 
How do you do that using dcfg?

### Step 1. Choose directory and initialize dcfg
At first, you need directory where all your config files will be stored!
It is most likely that this directory will also be a git repository.

```shell
$ mkdir ~/dotfiles  # create `dotfiles` directory in $HOME
$ cd ~/dotfiles  # go into it
$ git init  # initialize git repository
```
And then run
```shell
$ dcfg init
```
It will simply create `dcfg.json` file - dcfg config file.

### Step 2. Choose what to store and where to store
When dcfg is initialized we can define **bindings** and **targets**, so we can later
collect targets according to the bindings.

What is target:
> Target is simply a target file or directory which we want to be stored using dcfg. 

What is binding:
> Binding is like mount point or alias. 
> They map **source** path to **destination** path.
> For example, binding `/home/user -> ./user-home` means that all _targets_ 
> from `/home/user` will be copied to `./user-home`

#### Step 2.1. Define bindings
In usual simple case you will want to define only one or two bindings.
The most common will be `/home/<username> -> ./user-home`.

It can be simply defined using this command:
```shell
$ dcfg bind /home/<username> ./user-home
```

Probably you will also need to create another binding, in case if you wish to store
global config files from `/etc` directory:
```shell
$ dcfg bind /etc ./etc
```

I also would like to point out, that you can add `--remove` flag to remove binding. For example:
```shell
$ dcfg bind /usr/etc/ ./usr-etc
$ dcfg bind --remove /usr/etc/
```

#### Step 2.2. Define targets
After bindings are defined, we can define our targets. Remember, that each target needs
suitable binding at first, so dcfg would know where to put it.
```shell
$ dcfg add ~/.xinitrc ~/.config/i3
$ dcfg add /etc/issue
```

Target can be removed using `dcfg remove` command:
```shell
$ dcfg remove ~/.xinitrc
```

### Step 3. Collect targets
Now it's time to finally collect all defined targets! Just type:
```shell
$ dcfg collect
```

Then, if you run `tree -a -I .git` command to see contents of the current directory,
you will see the similar picture: 
```
.
├── dcfg.json
├── etc
│   └── issue
└── user-home
    ├── .config
    │   └── i3
    │       ├── config
    │       └── themes
    │           ├── catppuccin-frappe
    │           ├── catppuccin-latte
    │           ├── catppuccin-macchiato
    │           └── catppuccin-mocha
    └── .xinitrc
```
Now you can commit & push updates to the remote repository.

> You should run `dcfg collect` command every time you change your target files to
> keep the storage up-to-date.

### Extracting targets to their sources
Now from other side: you want to install (extract) collected targets to the machine.
All you need is to use `dcfg extract` command. It will extract all targets to their destinations.
It will show diff(s) and kindly ask before each copy operation by default.

#### Overwriting target source prefixes
One important thing to mention about extracting is that you can _overwrite target source prefixes_.
Sounds a bit complicated... What does it mean?
It will be simpler to explain using example:

User Bob with home directory located at `/home/bob` collected some config files from 
his home directory to `./user-home`.

Then Alice downloaded his repository containing these configuration files and dcfg config file
on her own computer.
Alice has home directory at `/home/alice` and once she tries to extract Bob's targets she gets
an error, because dcfg tries to extract targets to `/home/bob`, which does not exist.
But fortunately there is way to fix this. She needs to use `--overwrite-source-prefix` option when extracting.
The command will look like this:
```shell
$ dcfg extract --overwrite-source-prefix /home/bob:/home/alice
```
And now Bob's dotfiles will be put into Alice's home directory!
