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
    <a href="https://pkg.go.dev/github.com/Slug-Boi/cocommit"> <img src="https://img.shields.io/badge/_-reference-blue?logo=go&label=%E2%80%8E%20" alt="Go Reference" /></a>
    <a href="https://docs.github.com/en/pull-requests/committing-changes-to-your-project/creating-and-editing-commits/creating-a-commit-with-multiple-authors"> <img src="https://img.shields.io/badge/Github_Docs-Co--Authoring-grey?logo=github&labelColor=black" alt="Co-Author Docs Github Page" /></a>
</p>


## Install:
*Insert blerb about go install if it works*

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

# Syntax for the author file
The syntax for the author file can be found at the top of the template file included in the repo. It should look like this (opt) is optional syntax:  
```
name_short|Name|Username|email (opt: |ex) (opt: ;;group1 or ;;group1|group2|group3...)
```
opt explained:  
ex -> excludes the given author for all and negation selection commands  
group -> groups an author which can then be called as an argument to add all people from that group. An author can be a part of multiple groups 

# Why?
Co-authoring commits is a feature that is supported by github and gitlab and other git hosting services but creating the commits can be a bit of a pain. Co-authoring is extremely useful as teams can be much more transparent in who worked on what and it can be a great way to give credit to people who have helped on projects. This will make git-blame a lot more useful as you can quickly see who to contact or talk to about a specific part of the code. I strongly believe that this feature is underutilized and i attribute it mostly to the fact that is combersome to use. This tool aims to fix and streamline that process. (It even allows for automation of the process with the CLI mode)

# Bugs and Feature Requests
If you find any bugs or have any feature requests please open an issue on the github page and I will look at it, please note this is very much a passion project of mine so things will take time. Feel free to also open a PR if you want to contribute to the project this will be greatly appreciated
