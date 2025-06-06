<h1 align="center"> 
    cocommit
</h1>
<p align="center">
   <img src="imgs/cocommit_logo.png" alt="cocommit logo" width="300" height="300"/>
</p>
<p align="center">
    <a href="https://github.com/Slug-Boi/cocommit/releases/latest"> <img src="https://img.shields.io/github/v/release/Slug-Boi/cocommit?logo=github" alt="Github Latest Release" /></a>
    <a href="https://github.com/Slug-Boi/cocommit/actions/workflows/test_push.yml"> <img src="https://img.shields.io/github/actions/workflow/status/Slug-Boi/cocommit/test_push.yml?label=tests" alt="Github Actions Tests" /></a>
    <a href="https://github.com/Slug-Boi/cocommit/actions/workflows/build_test_release.yml"> <img src="https://img.shields.io/github/actions/workflow/status/Slug-Boi/cocommit/build_test_release.yml" alt="Github Actions Build" /></a>  
    <a href="https://codecov.io/gh/Slug-Boi/cocommit" > <img src="https://codecov.io/gh/Slug-Boi/cocommit/branch/main/graph/badge.svg?token=EVGSNJGCYJ"/></a>  
    <a href="https://pkg.go.dev/github.com/Slug-Boi/cocommit"> <img src="https://img.shields.io/badge/_-reference-blue?logo=go&label=%E2%80%8E%20" alt="Go Reference" /></a>
    <a href="https://docs.github.com/en/pull-requests/committing-changes-to-your-project/creating-and-editing-commits/creating-a-commit-with-multiple-authors"> <img src="https://img.shields.io/badge/Github_Docs-Co--Authoring-grey?logo=github&labelColor=black" alt="Co-Author Docs Github Page" /></a>
</p>


## Install:
### Install Script
The repo contains two different install scripts for unix based system and windows. You can copy paste one of the commands below to use it.  
**Unix:**
```console
curl -L -o install.sh https://github.com/Slug-Boi/cocommit/raw/refs/heads/main/installer/install.sh && chmod +x install.sh && ./install.sh
```
**Windows:**  
*The windows script is untested so please let me know if it works and use at your own risk*
```console 
Invoke-WebRequest -Uri https://github.com/Slug-Boi/cocommit/raw/refs/heads/main/installer/install.ps1 -OutFile install.ps1; .\install.ps1
```
### Go Install
It should be possible to install the program using the go install command:
```
go install github.com/Slug-Boi/cocommit
```
You will more than likely have to add the binary to your PATH after the fact if your go bin directory is not in your PATH already. Or you can create an alias to it using the instructions below.

### Manual Binary Install
Download the binary for your OS on the [release page](https://github.com/Slug-Boi/cocommit/releases)  
Once downloaded you need to create an alias for your shell guides found below:  
[Bash Guide for alias](https://linuxize.com/post/how-to-create-bash-aliases/)  
*^please note if you use another shell than bash you will need to look up how to do it in that shell^*  
[Powershell Guide for alias](https://stackoverflow.com/questions/24914589/how-to-create-permanent-powershell-aliases)  

Once you've created an alias for the program you need to create an author file. You can run the program and it will prompt you to create an author file if one is not found in the default location.

Optionally you can set an env variable (author_file) in your shell. This should be the path pointing to your author file  
For bash you just need to add this to your .bashrc file:
```
export author_file='path/to/your/aurhors'
```
[Powershell guide for env variable](https://stackoverflow.com/a/714918)

*Please note that the syntax line at the top of the authors file should not be deleted and that you must add at least one author to the file to run the program*

## Usage:
The CLI has two modes of operation, the CLI commands and the TUI below are guides on how to use both

### TUI
To launch the TUI run the program with no args  
```
$ cocommit
```
From here you will be asked to write a commit message and then select authors from a list. The TUI has lots of keybinds that can be seen on the list view by pressing `?` (You can get to the list view by either creating a commit message or launching the TUI with the -a flag (stands for author_list)). You can create authors, add temp authors or do all of the usual selections like group selections, negated selections or select all. Below is a small gif/video showing a run through of the TUI.

<a href="https://asciinema.org/a/OFRu5t0A2cSugw49VV6GvumyF" target="_blank"><img src="imgs/cocommit.gif" /></a>
^click me for video with timestamps^

### CLI

```
$ cocommit "message" <name1> [name2] [name3]...

$ cocommit "message" <name:email1> [name:email2] [name:email3]...

$ cocommit "message" <name:email1> <name1> [name:email2]

$ cocommit "message" all
*adds all comitters execpt ones tagged with ex*

$ cocommit "message" ^<name1> ^[name2]
*adds all users except the negated users and users tagged with ex*

$ cocommit "message" <group_name>
*adds all users that has that group tag in author file*

$ cocommit "message"
*Runs git commit -m "message"*

$ cocommit users
*prints list of users*

$ cocommit update
*updates the cocommit binary to the newest version hosted on github*
```

### Commitizen Integration
**Note:** *This feature requires the commitizen package to be installed on your system.*  
If you want to write the commit message using commitizen you can use the cz command:
```
$ cocommit cz 
```
This will open the commitizen prompt and then you can select authors after as usual. By default it runs using the TUI mode of cocommit but you can use the -c flag to run it in CLI mode:
```
$ cocommit cz -c
```
This should allow for all the same utility as the normal CLI mode in terms of author selection

### Lazygit Config
If you use lazygit you can add the following to your lazygit config file to use cocommit:
```
customCommands:
    - key: '<c-A>'
      context: 'global'
      showOutput: true
      prompts:
        - type: 'input'
          title: 'Commit message'
          key: 'message'
          initialValue: ''
        - type: 'input'
          title: 'Authors'
          key: 'authors'
          initialValue: ''
      command: 'cocommit -p "{{.Form.message}}" {{.Form.authors}}'
    - key: '<c-a>'
      context: 'global'
      subprocess: true
      command: 'cocommit -p; exit'
```
A sample lazygit config file can be found [here](https://github.com/Slug-Boi/cocommit/blob/main/lazygit_config/config.yml)

# Syntax for the author file
The syntax for the author file is json below is a small example with fake information to show what it looks like. The author file can be edited safely from the tool so there is no real need to edit this manually. Whilst this format is a little heavier than the old custom CSV format it is much easier to work with and handle json so rest assured this is best way forward
```json
{
   "Authors":{
      "Morgan Rivers":{
         "shortname":"mr",
         "longname":"Morgan Rivers",
         "username":"morgan-rivers",
         "email":"mrivers@example.com",
         "ex":false,
         "groups":[
            "dev",
            "qa",
            "design"
         ]
      },
      "Taylor Chen":{
         "shortname":"tc",
         "longname":"Taylor Chen",
         "username":"tchen",
         "email":"tchen@example.org",
         "ex":true,
         "groups":[
            "dev",
            "admin",
            "support"
         ]
      },
      "Jordan Smithfield":{
         "shortname":"js",
         "longname":"Jordan Smithfield",
         "username":"jsmith",
         "email":"j.smithfield@test.net",
         "ex":false,
         "groups":[
            "marketing",
            "content",
            "social"
         ]
      }
   }
}
```

# Why?
Co-authoring commits is a feature that is supported by github and gitlab and other git hosting services but creating the commits can be a bit of a pain. Co-authoring is extremely useful as teams can be much more transparent in who worked on what and it can be a great way to give credit to people who have helped on projects. This will make git-blame a lot more useful as you can quickly see who to contact or talk to about a specific part of the code. I strongly believe that this feature is underutilized and i attribute it mostly to the fact that is combersome to use. This tool aims to fix and streamline that process. (It even allows for automation of the process with the CLI mode)

# Bugs and Feature Requests
If you find any bugs or have any feature requests please open an issue on the github page and I will look at it, please note this is very much a passion project of mine so things will take time. Feel free to also open a PR if you want to contribute to the project this will be greatly appreciated
