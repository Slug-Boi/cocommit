# cocommit
Currently built in:
- [x] Go
- [ ] C++
- [ ] Python

Currently tested in:
- [x] Go
- [ ] C++
- [ ] Python

# How to run
Download the binary for your OS on the [release page](https://github.com/Slug-Boi/cocommit/releases)  
Once downloaded you need to create an alias for your shell guides found below:  
[Bash Guide](https://linuxize.com/post/how-to-create-bash-aliases/)  
*^please note if you use another shell than bash you will need to look up how to do it in that shell^*  
[Powershell Guide](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/set-alias?view=powershell-7.4)  

Once you've created an alias for the program you need to set an env variable (authors_file) in your shell. This should be the path pointing to your author.txt file   

you can now run it using the alias shorthand you assigned to it 
## Usage:
```zsh
$ cocommit "message" <name1> [name2] [name3]...

$ cocommit "message" <name:email1> [name:email2] [name:email3]...

$ cocommit "message" <name:email1> <name1> [name:email2]

$ cocommit "message" all
*adds all comitters execpt ones tagged with ex*

$ cocommit "message" ^<name>
*adds all users except the negated users and users tagged with ex*

$ cocommit
*prints usage*

$ cocommit users
*prints list of users*
```

# Syntax for the author file
The syntax for the author file can be found at the top of the template file included in the repo. It should look like this (opt) is optional syntax:  
```
name_short|Name|Username|email (opt: |ex)
```
opt explained:  
ex -> excludes the given author for all and negation commands

# Why?
Writing co-authors onto commits can be pretty tedious so automating this process as a simple shell alias is a lot nicer

# Workflows
This repo is sort of a test bed for working with Dagger CI but therefore also should support automated testing and building at some point 
*coming soon*

